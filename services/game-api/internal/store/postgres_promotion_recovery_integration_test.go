//go:build integration

package store

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/cvsz/zneondrive/services/game-api/internal/core"
)

const (
	promotionDatabase = "zneondrive_failover_test"
	promotionUser     = "zneondrive"
	promotionPassword = "zneondrive-test"
)

func TestPostgresPhysicalStandbyPromotionPreservesQuestIdempotency(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		if os.Getenv("CI") == "true" {
			t.Fatalf("docker is required for the PostgreSQL promotion evidence gate: %v", err)
		}
		t.Skip("docker is unavailable")
	}

	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	network := "zneondrive-promotion-" + suffix
	primaryName := "zneondrive-primary-" + suffix
	replicaName := "zneondrive-replica-" + suffix
	replicaVolume := "zneondrive-replica-data-" + suffix

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

	mustDocker(t,
		"run", "--rm", "--network", network,
		"-v", replicaVolume+":/var/lib/postgresql/data",
		"postgres:17-alpine", "sh", "-ceu",
		"rm -rf /var/lib/postgresql/data/*; chown postgres:postgres /var/lib/postgresql/data; exec su-exec postgres pg_basebackup -h "+primaryName+" -U zneondrive_replica -D /var/lib/postgresql/data -Fp -Xs -R -P",
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

	primaryPort := dockerMappedPort(t, primaryName)
	primaryURL := fmt.Sprintf("postgres://%s:%s@127.0.0.1:%s/%s?sslmode=disable", promotionUser, promotionPassword, primaryPort, promotionDatabase)

	ctx := context.Background()
	primary, err := OpenPostgres(ctx, primaryURL)
	if err != nil {
		t.Fatal(err)
	}
	if err = primary.EnsureSchema(ctx); err != nil {
		primary.Close()
		t.Fatal(err)
	}

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

	sessionOne, err := core.NewSecret()
	if err != nil {
		primary.Close()
		t.Fatal(err)
	}
	sessionOneHash := core.HashSecret(sessionOne)
	if err = primary.CreateSession(ctx, before.AccountID, sessionOneHash, time.Now().UTC().Add(time.Hour)); err != nil {
		primary.Close()
		t.Fatal(err)
	}

	operationID, ok := core.ScopeQuestOperationID(before.CharacterID, "promotion-mq001-"+before.AccountID)
	if !ok {
		primary.Close()
		t.Fatal("expected character-scoped durable quest operation id")
	}
	mutated, receipt, err := primary.CompleteQuest(ctx, sessionOneHash, "MQ001", operationID)
	if err != nil {
		primary.Close()
		t.Fatal(err)
	}
	if !receipt.Applied || mutated.Money != 110 || mutated.XP != 55 || mutated.Reputation != 1 {
		primary.Close()
		t.Fatalf("unexpected canonical MQ001 mutation before promotion: snapshot=%+v receipt=%+v", mutated, receipt)
	}
	primary.Close()

	primaryLSN := strings.TrimSpace(mustDockerOutput(t, "exec", "-e", "PGPASSWORD="+promotionPassword, primaryName,
		"psql", "-At", "-U", promotionUser, "-d", promotionDatabase, "-c", "SELECT pg_current_wal_lsn()"))
	if primaryLSN == "" {
		t.Fatal("primary returned an empty WAL LSN")
	}
	waitReplicaLSN(t, replicaName, primaryLSN)

	mustDocker(t, "stop", "-t", "5", primaryName)
	mustDocker(t, "exec", "-u", "postgres", replicaName, "pg_ctl", "promote", "-D", "/var/lib/postgresql/data", "-w")
	waitSQLTruth(t, replicaName, "SELECT NOT pg_is_in_recovery()")

	replicaPort := dockerMappedPort(t, replicaName)
	promotedURL := fmt.Sprintf("postgres://%s:%s@127.0.0.1:%s/%s?sslmode=disable", promotionUser, promotionPassword, replicaPort, promotionDatabase)
	promoted, err := OpenPostgres(ctx, promotedURL)
	if err != nil {
		t.Fatal(err)
	}
	defer promoted.Close()
	if err = promoted.EnsureSchema(ctx); err != nil {
		t.Fatal(err)
	}

	resumed, created, err := promoted.Bootstrap(ctx, resumeHash)
	if err != nil {
		t.Fatal(err)
	}
	if created {
		t.Fatal("promoted standby must resume the existing durable account")
	}
	if resumed.AccountID != before.AccountID || resumed.CharacterID != before.CharacterID || resumed.VehicleID != before.VehicleID {
		t.Fatalf("durable identity changed across PostgreSQL promotion: before=%+v after=%+v", before, resumed)
	}
	if resumed.Money != mutated.Money || resumed.XP != mutated.XP || resumed.Reputation != mutated.Reputation {
		t.Fatalf("durable progression changed across PostgreSQL promotion: before=%+v after=%+v", mutated, resumed)
	}

	sessionTwo, err := core.NewSecret()
	if err != nil {
		t.Fatal(err)
	}
	sessionTwoHash := core.HashSecret(sessionTwo)
	if err = promoted.CreateSession(ctx, resumed.AccountID, sessionTwoHash, time.Now().UTC().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}

	replayed, replayReceipt, err := promoted.CompleteQuest(ctx, sessionTwoHash, "MQ001", operationID)
	if err != nil {
		t.Fatal(err)
	}
	if replayReceipt.Applied {
		t.Fatal("MQ001 replay after standby promotion must not re-apply rewards")
	}
	if replayReceipt.Money != 0 || replayReceipt.XP != 0 || replayReceipt.Reputation != 0 {
		t.Fatalf("promotion replay reissued rewards: %+v", replayReceipt)
	}
	if replayed.Money != 110 || replayed.XP != 55 || replayed.Reputation != 1 {
		t.Fatalf("promotion replay changed authoritative progression: %+v", replayed)
	}
	if len(replayed.CompletedQuests) != 1 || replayed.CompletedQuests[0] != "MQ001" {
		t.Fatalf("expected exactly one durable MQ001 completion after promotion, got %v", replayed.CompletedQuests)
	}

	mq002Operation, ok := core.ScopeQuestOperationID(before.CharacterID, "promotion-mq002-"+before.AccountID)
	if !ok {
		t.Fatal("expected scoped MQ002 operation id")
	}
	afterMQ002, mq002Receipt, err := promoted.CompleteQuest(ctx, sessionTwoHash, "MQ002", mq002Operation)
	if err != nil {
		t.Fatalf("complete MQ002 on promoted primary: %v", err)
	}
	if !mq002Receipt.Applied || afterMQ002.Money <= replayed.Money {
		t.Fatalf("promoted primary did not accept the next authoritative mutation: before=%+v after=%+v receipt=%+v", replayed, afterMQ002, mq002Receipt)
	}

	t.Logf("PostgreSQL physical-standby promotion evidence passed at replayed WAL LSN %s; durable MQ001 remained exactly-once and MQ002 advanced after promotion", primaryLSN)
}

func docker(args ...string) error {
	cmd := exec.Command("docker", args...)
	return cmd.Run()
}

func mustDocker(t *testing.T, args ...string) {
	t.Helper()
	cmd := exec.Command("docker", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("docker command failed: %v: %s", err, strings.TrimSpace(string(output)))
	}
}

func mustDockerOutput(t *testing.T, args ...string) string {
	t.Helper()
	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("docker command failed: %v: %s", err, strings.TrimSpace(string(output)))
	}
	return string(output)
}

func waitDockerPostgres(t *testing.T, container, database, user string) {
	t.Helper()
	deadline := time.Now().Add(45 * time.Second)
	for time.Now().Before(deadline) {
		output, err := exec.Command("docker", "exec", "-e", "PGPASSWORD="+promotionPassword, container,
			"psql", "-At", "-U", user, "-d", database, "-c", "SELECT 1").CombinedOutput()
		if err == nil && strings.TrimSpace(string(output)) == "1" {
			return
		}
		time.Sleep(time.Second)
	}
	t.Fatalf("PostgreSQL container %s target database %s did not become query-ready", container, database)
}

func dockerMappedPort(t *testing.T, container string) string {
	t.Helper()
	mapping := strings.TrimSpace(mustDockerOutput(t, "port", container, "5432/tcp"))
	if mapping == "" {
		t.Fatalf("container %s has no mapped PostgreSQL port", container)
	}
	line := strings.Split(mapping, "\n")[0]
	idx := strings.LastIndex(line, ":")
	if idx < 0 || idx == len(line)-1 {
		t.Fatalf("unexpected PostgreSQL port mapping for %s: %q", container, line)
	}
	return line[idx+1:]
}

func waitSQLTruth(t *testing.T, container, query string) {
	t.Helper()
	deadline := time.Now().Add(45 * time.Second)
	for time.Now().Before(deadline) {
		output, err := exec.Command("docker", "exec", "-e", "PGPASSWORD="+promotionPassword, container,
			"psql", "-At", "-U", promotionUser, "-d", promotionDatabase, "-c", query).CombinedOutput()
		if err == nil && strings.TrimSpace(string(output)) == "t" {
			return
		}
		time.Sleep(time.Second)
	}
	t.Fatalf("PostgreSQL assertion did not become true in container %s", container)
}

func waitReplicaLSN(t *testing.T, replica, requiredLSN string) {
	t.Helper()
	query := fmt.Sprintf("SELECT COALESCE(pg_last_wal_replay_lsn() >= '%s'::pg_lsn, false)", strings.ReplaceAll(requiredLSN, "'", "''"))
	waitSQLTruth(t, replica, query)
}
