//go:build integration

package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/cvsz/zneondrive/services/game-api/internal/core"
)

func TestPostgresPromotionRollsBackInterruptedQuestTransaction(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		if os.Getenv("CI") == "true" {
			t.Fatalf("docker is required for the PostgreSQL transaction-failover evidence gate: %v", err)
		}
		t.Skip("docker is unavailable")
	}

	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	network := "zneondrive-tx-failover-" + suffix
	primaryName := "zneondrive-tx-primary-" + suffix
	replicaName := "zneondrive-tx-replica-" + suffix
	replicaVolume := "zneondrive-tx-replica-data-" + suffix

	mustDocker(t, "network", "create", network)
	mustDocker(t, "volume", "create", replicaVolume)
	t.Cleanup(func() {
		_ = docker("rm", "-f", replicaName)
		_ = docker("rm", "-f", primaryName)
		_ = docker("volume", "rm", "-f", replicaVolume)
		_ = docker("network", "rm", network)
	})

	mustDocker(t,
		"run", "-d", "--name", primaryName,
		"--network", network,
		"-p", "127.0.0.1::5432",
		"-e", "POSTGRES_DB="+promotionDatabase,
		"-e", "POSTGRES_USER="+promotionUser,
		"-e", "POSTGRES_PASSWORD="+promotionPassword,
		"postgres:17-alpine",
		"postgres",
		"-c", "wal_level=replica",
		"-c", "max_wal_senders=5",
		"-c", "max_replication_slots=5",
		"-c", "hot_standby=on",
	)
	waitDockerPostgres(t, primaryName, promotionDatabase, promotionUser)

	mustDocker(t, "exec", "-e", "PGPASSWORD="+promotionPassword, primaryName,
		"psql", "-v", "ON_ERROR_STOP=1", "-U", promotionUser, "-d", promotionDatabase,
		"-c", "CREATE ROLE zneondrive_replica WITH REPLICATION LOGIN")
	mustDocker(t, "exec", primaryName, "sh", "-ceu",
		"printf '%s\\n' 'host replication zneondrive_replica 0.0.0.0/0 trust' >> \"$PGDATA/pg_hba.conf\"")
	mustDocker(t, "exec", "-e", "PGPASSWORD="+promotionPassword, primaryName,
		"psql", "-v", "ON_ERROR_STOP=1", "-U", promotionUser, "-d", promotionDatabase,
		"-c", "SELECT pg_reload_conf()")

	ctx := context.Background()
	primaryPort := dockerMappedPort(t, primaryName)
	primaryURL := fmt.Sprintf("postgres://%s:%s@127.0.0.1:%s/%s?sslmode=disable", promotionUser, promotionPassword, primaryPort, promotionDatabase)
	primary, err := OpenPostgres(ctx, primaryURL)
	if err != nil {
		t.Fatal(err)
	}
	if err = primary.EnsureSchema(ctx); err != nil {
		primary.Close()
		t.Fatal(err)
	}

	mustDocker(t,
		"run", "--rm",
		"-v", replicaVolume+":/var/lib/postgresql/data",
		"postgres:17-alpine", "sh", "-ceu",
		"rm -rf /var/lib/postgresql/data/*; chown postgres:postgres /var/lib/postgresql/data",
	)
	mustDocker(t,
		"run", "--rm", "--user", "postgres", "--network", network,
		"-v", replicaVolume+":/var/lib/postgresql/data",
		"postgres:17-alpine",
		"pg_basebackup", "-h", primaryName, "-U", "zneondrive_replica",
		"-D", "/var/lib/postgresql/data", "-Fp", "-Xs", "-R", "-P",
	)
	mustDocker(t,
		"run", "-d", "--name", replicaName,
		"--network", network,
		"-p", "127.0.0.1::5432",
		"-v", replicaVolume+":/var/lib/postgresql/data",
		"postgres:17-alpine",
		"postgres", "-c", "hot_standby=on",
	)
	waitDockerPostgres(t, replicaName, promotionDatabase, promotionUser)
	waitSQLTruth(t, primaryName, "SELECT count(*) = 1 FROM pg_stat_replication WHERE state = 'streaming'")

	resumeKey, err := core.NewSecret()
	if err != nil {
		primary.Close()
		t.Fatal(err)
	}
	resumeHash := core.HashSecret(resumeKey)
	before, created, err := primary.Bootstrap(ctx, resumeHash)
	if err != nil {
		primary.Close()
		t.Fatal(err)
	}
	if !created {
		primary.Close()
		t.Fatal("expected fresh bootstrap on isolated primary")
	}

	sessionKey, err := core.NewSecret()
	if err != nil {
		primary.Close()
		t.Fatal(err)
	}
	sessionHash := core.HashSecret(sessionKey)
	if err = primary.CreateSession(ctx, before.AccountID, sessionHash, time.Now().UTC().Add(time.Hour)); err != nil {
		primary.Close()
		t.Fatal(err)
	}

	operationID, ok := core.ScopeQuestOperationID(before.CharacterID, "tx-failover-mq001-"+before.AccountID)
	if !ok {
		primary.Close()
		t.Fatal("expected character-scoped durable quest operation id")
	}

	mustDocker(t, "exec", "-e", "PGPASSWORD="+promotionPassword, primaryName,
		"psql", "-v", "ON_ERROR_STOP=1", "-U", promotionUser, "-d", promotionDatabase,
		"-c", `
CREATE OR REPLACE FUNCTION zneondrive_test_pause_quest_insert()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
  IF NEW.quest_id = 'MQ001' THEN
    PERFORM pg_sleep(30);
  END IF;
  RETURN NEW;
END
$$;
DROP TRIGGER IF EXISTS zneondrive_test_pause_quest_insert ON quest_completions;
CREATE TRIGGER zneondrive_test_pause_quest_insert
AFTER INSERT ON quest_completions
FOR EACH ROW EXECUTE FUNCTION zneondrive_test_pause_quest_insert();`)

	triggerLSN := strings.TrimSpace(mustDockerOutput(t, "exec", "-e", "PGPASSWORD="+promotionPassword, primaryName,
		"psql", "-At", "-U", promotionUser, "-d", promotionDatabase, "-c", "SELECT pg_current_wal_lsn()"))
	waitReplicaLSN(t, replicaName, triggerLSN)

	type questResult struct {
		snapshot core.Snapshot
		receipt  core.RewardReceipt
		err      error
	}
	resultCh := make(chan questResult, 1)
	go func() {
		snapshot, receipt, callErr := primary.CompleteQuest(ctx, sessionHash, "MQ001", operationID)
		resultCh <- questResult{snapshot: snapshot, receipt: receipt, err: callErr}
	}()

	waitSQLTruth(t, primaryName, `SELECT EXISTS(
SELECT 1
FROM pg_stat_activity
WHERE datname = current_database()
  AND state = 'active'
  AND wait_event = 'PgSleep'
  AND query ILIKE '%INSERT INTO quest_completions%'
)`)

	mustDocker(t, "stop", "-t", "0", primaryName)
	primary.Close()

	select {
	case interrupted := <-resultCh:
		if interrupted.err == nil {
			t.Fatalf("quest mutation unexpectedly succeeded despite primary termination: snapshot=%+v receipt=%+v", interrupted.snapshot, interrupted.receipt)
		}
	case <-time.After(15 * time.Second):
		t.Fatal("interrupted quest mutation did not return after primary termination")
	}

	mustDocker(t, "exec", "-u", "postgres", replicaName, "pg_ctl", "promote", "-D", "/var/lib/postgresql/data", "-w")
	waitSQLTruth(t, replicaName, "SELECT NOT pg_is_in_recovery()")

	replicaPort := dockerMappedPort(t, replicaName)
	promotedURL := fmt.Sprintf("postgres://%s:%s@127.0.0.1:%s/%s?sslmode=disable", promotionUser, promotionPassword, replicaPort, promotionDatabase)
	promoted, err := OpenPostgres(ctx, promotedURL)
	if err != nil {
		t.Fatal(err)
	}
	defer promoted.Close()

	mustDocker(t, "exec", "-e", "PGPASSWORD="+promotionPassword, replicaName,
		"psql", "-v", "ON_ERROR_STOP=1", "-U", promotionUser, "-d", promotionDatabase,
		"-c", `DROP TRIGGER IF EXISTS zneondrive_test_pause_quest_insert ON quest_completions;
DROP FUNCTION IF EXISTS zneondrive_test_pause_quest_insert();`)

	resumed, created, err := promoted.Bootstrap(ctx, resumeHash)
	if err != nil {
		t.Fatal(err)
	}
	if created {
		t.Fatal("promoted standby must resume the durable account")
	}
	if resumed.AccountID != before.AccountID || resumed.CharacterID != before.CharacterID || resumed.VehicleID != before.VehicleID {
		t.Fatalf("durable identity changed after interrupted transaction promotion: before=%+v after=%+v", before, resumed)
	}
	if resumed.Money != 0 || resumed.XP != 0 || resumed.Reputation != 0 || len(resumed.CompletedQuests) != 0 {
		t.Fatalf("uncommitted MQ001 state became visible after promotion: %+v", resumed)
	}

	var durableRows int
	if err = promoted.pool.QueryRow(ctx, "SELECT count(*) FROM quest_completions WHERE operation_id=$1", operationID).Scan(&durableRows); err != nil {
		t.Fatalf("count interrupted quest rows: %v", err)
	}
	if durableRows != 0 {
		t.Fatalf("expected no committed row for interrupted operation, got %d", durableRows)
	}

	newSessionKey, err := core.NewSecret()
	if err != nil {
		t.Fatal(err)
	}
	newSessionHash := core.HashSecret(newSessionKey)
	if err = promoted.CreateSession(ctx, resumed.AccountID, newSessionHash, time.Now().UTC().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}

	afterRetry, retryReceipt, err := promoted.CompleteQuest(ctx, newSessionHash, "MQ001", operationID)
	if err != nil {
		t.Fatalf("retry MQ001 on promoted primary: %v", err)
	}
	if !retryReceipt.Applied {
		t.Fatal("retry after rollback/promotion must apply MQ001 exactly once")
	}
	if afterRetry.Money != 110 || afterRetry.XP != 55 || afterRetry.Reputation != 1 {
		t.Fatalf("unexpected progression after promoted retry: %+v", afterRetry)
	}
	if len(afterRetry.CompletedQuests) != 1 || afterRetry.CompletedQuests[0] != "MQ001" {
		t.Fatalf("expected exactly one MQ001 after retry, got %v", afterRetry.CompletedQuests)
	}

	replayed, replayReceipt, err := promoted.CompleteQuest(ctx, newSessionHash, "MQ001", operationID)
	if err != nil {
		t.Fatal(err)
	}
	if replayReceipt.Applied || replayReceipt.Money != 0 || replayReceipt.XP != 0 || replayReceipt.Reputation != 0 {
		t.Fatalf("post-retry replay must remain idempotent: %+v", replayReceipt)
	}
	if replayed.Money != afterRetry.Money || replayed.XP != afterRetry.XP || replayed.Reputation != afterRetry.Reputation {
		t.Fatalf("post-retry replay changed authoritative progression: before=%+v after=%+v", afterRetry, replayed)
	}

	if errors.Is(err, ErrOperationKey) {
		t.Fatal("same authoritative character/quest operation must not become an operation-key conflict after rollback")
	}

	t.Log("PostgreSQL transaction-boundary failover evidence passed: primary terminated while MQ001 row was uncommitted, promoted standby exposed no partial reward/completion, and retry applied exactly once")
}
