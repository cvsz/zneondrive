package store

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/cvsz/zneondrive/services/game-api/internal/core"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/001_init.sql
var migrationSQL string

var starterParts = []string{
	"part_chassis_starter_prototype",
	"part_engine_ice_street_i",
	"part_ecu_legacy_zero",
}

type Postgres struct {
	pool *pgxpool.Pool
}

func OpenPostgres(ctx context.Context, databaseURL string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}
	p := &Postgres{pool: pool}
	if err := p.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return p, nil
}

func (p *Postgres) Close() {
	p.pool.Close()
}

func (p *Postgres) Ping(ctx context.Context) error {
	if err := p.pool.Ping(ctx); err != nil {
		return fmt.Errorf("postgres ping: %w", err)
	}
	return nil
}

func (p *Postgres) EnsureSchema(ctx context.Context) error {
	migrations := []struct {
		name string
		sql  string
	}{
		{name: "001_init", sql: migrationSQL},
		{name: "002_game_tickets", sql: ticketMigrationSQL},
	}
	for _, migration := range migrations {
		if _, err := p.pool.Exec(ctx, migration.sql); err != nil {
			return fmt.Errorf("apply schema %s: %w", migration.name, err)
		}
	}
	return nil
}

func (p *Postgres) Bootstrap(ctx context.Context, resumeKeyHash string) (core.Snapshot, bool, error) {
	var accountID string
	err := p.pool.QueryRow(ctx, "SELECT id FROM accounts WHERE resume_key_hash=$1", resumeKeyHash).Scan(&accountID)
	if err == nil {
		snapshot, snapErr := p.snapshotByAccount(ctx, accountID)
		return snapshot, false, snapErr
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return core.Snapshot{}, false, fmt.Errorf("lookup resume key: %w", err)
	}

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return core.Snapshot{}, false, fmt.Errorf("begin bootstrap: %w", err)
	}
	defer tx.Rollback(ctx)

	accountID, err = core.NewID("acct")
	if err != nil {
		return core.Snapshot{}, false, err
	}
	characterID, err := core.NewID("char")
	if err != nil {
		return core.Snapshot{}, false, err
	}
	vehicleID, err := core.NewID("veh")
	if err != nil {
		return core.Snapshot{}, false, err
	}
	partIDs, err := core.NormalizeParts(starterParts)
	if err != nil {
		return core.Snapshot{}, false, err
	}
	buildHash, err := core.BuildHash(partIDs)
	if err != nil {
		return core.Snapshot{}, false, err
	}
	partJSON, err := json.Marshal(partIDs)
	if err != nil {
		return core.Snapshot{}, false, fmt.Errorf("marshal starter build: %w", err)
	}

	if _, err = tx.Exec(ctx, "INSERT INTO accounts(id,resume_key_hash) VALUES($1,$2)", accountID, resumeKeyHash); err != nil {
		return core.Snapshot{}, false, fmt.Errorf("insert account: %w", err)
	}
	if _, err = tx.Exec(ctx, "INSERT INTO characters(id,account_id) VALUES($1,$2)", characterID, accountID); err != nil {
		return core.Snapshot{}, false, fmt.Errorf("insert character: %w", err)
	}
	if _, err = tx.Exec(ctx, "INSERT INTO vehicles(id,owner_character_id,starter_lineage) VALUES($1,$2,true)", vehicleID, characterID); err != nil {
		return core.Snapshot{}, false, fmt.Errorf("insert starter vehicle: %w", err)
	}
	if _, err = tx.Exec(ctx, "INSERT INTO vehicle_builds(vehicle_id,revision,part_ids,validation_hash) VALUES($1,1,$2,$3)", vehicleID, partJSON, buildHash); err != nil {
		return core.Snapshot{}, false, fmt.Errorf("insert starter build: %w", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return core.Snapshot{}, false, fmt.Errorf("commit bootstrap: %w", err)
	}

	snapshot, err := p.snapshotByAccount(ctx, accountID)
	return snapshot, true, err
}

func (p *Postgres) CreateSession(ctx context.Context, accountID, tokenHash string, expiresAt time.Time) error {
	_, err := p.pool.Exec(ctx, `
		INSERT INTO sessions(token_hash,account_id,expires_at)
		VALUES($1,$2,$3)
		ON CONFLICT(token_hash) DO UPDATE SET account_id=excluded.account_id, expires_at=excluded.expires_at
	`, tokenHash, accountID, expiresAt.UTC())
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

func (p *Postgres) SnapshotBySession(ctx context.Context, tokenHash string) (core.Snapshot, error) {
	var accountID string
	err := p.pool.QueryRow(ctx, `
		SELECT account_id
		FROM sessions
		WHERE token_hash=$1 AND expires_at > now()
	`, tokenHash).Scan(&accountID)
	if errors.Is(err, pgx.ErrNoRows) {
		return core.Snapshot{}, ErrUnauthorized
	}
	if err != nil {
		return core.Snapshot{}, fmt.Errorf("resolve session: %w", err)
	}
	return p.snapshotByAccount(ctx, accountID)
}

func (p *Postgres) CompleteQuest(ctx context.Context, tokenHash, questID, operationID string) (core.Snapshot, core.RewardReceipt, error) {
	if operationID == "" {
		return core.Snapshot{}, core.RewardReceipt{}, ErrOperationKey
	}
	questNumber, err := core.ParseQuestID(questID)
	if err != nil {
		return core.Snapshot{}, core.RewardReceipt{}, err
	}
	previousQuest, hasPrevious, err := core.PreviousQuestID(questID)
	if err != nil {
		return core.Snapshot{}, core.RewardReceipt{}, err
	}

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return core.Snapshot{}, core.RewardReceipt{}, fmt.Errorf("begin quest mutation: %w", err)
	}
	defer tx.Rollback(ctx)

	var accountID, characterID string
	err = tx.QueryRow(ctx, `
		SELECT a.id, c.id
		FROM sessions s
		JOIN accounts a ON a.id=s.account_id
		JOIN characters c ON c.account_id=a.id
		WHERE s.token_hash=$1 AND s.expires_at > now()
	`, tokenHash).Scan(&accountID, &characterID)
	if errors.Is(err, pgx.ErrNoRows) {
		return core.Snapshot{}, core.RewardReceipt{}, ErrUnauthorized
	}
	if err != nil {
		return core.Snapshot{}, core.RewardReceipt{}, fmt.Errorf("authorize quest mutation: %w", err)
	}

	var existingQuest string
	err = tx.QueryRow(ctx, "SELECT quest_id FROM quest_completions WHERE operation_id=$1", operationID).Scan(&existingQuest)
	if err == nil {
		if existingQuest != questID {
			return core.Snapshot{}, core.RewardReceipt{}, ErrOperationKey
		}
		if err := tx.Commit(ctx); err != nil {
			return core.Snapshot{}, core.RewardReceipt{}, err
		}
		snapshot, snapErr := p.snapshotByAccount(ctx, accountID)
		return snapshot, core.RewardReceipt{QuestID: questID, Applied: false}, snapErr
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return core.Snapshot{}, core.RewardReceipt{}, fmt.Errorf("check quest operation: %w", err)
	}

	if hasPrevious {
		var exists bool
		err = tx.QueryRow(ctx, `
			SELECT EXISTS(
				SELECT 1 FROM quest_completions WHERE character_id=$1 AND quest_id=$2
			)
		`, characterID, previousQuest).Scan(&exists)
		if err != nil {
			return core.Snapshot{}, core.RewardReceipt{}, fmt.Errorf("check quest prerequisite: %w", err)
		}
		if !exists {
			return core.Snapshot{}, core.RewardReceipt{}, ErrOutOfOrder
		}
	}

	money := int64(100 + questNumber*10)
	xp := int64(50 + questNumber*5)
	reputation := int64(1)

	var inserted string
	err = tx.QueryRow(ctx, `
		INSERT INTO quest_completions(
			character_id,quest_id,operation_id,reward_money,reward_xp,reward_reputation
		)
		VALUES($1,$2,$3,$4,$5,$6)
		ON CONFLICT(character_id,quest_id) DO NOTHING
		RETURNING quest_id
	`, characterID, questID, operationID, money, xp, reputation).Scan(&inserted)

	applied := true
	if errors.Is(err, pgx.ErrNoRows) {
		applied = false
		money, xp, reputation = 0, 0, 0
	} else if err != nil {
		return core.Snapshot{}, core.RewardReceipt{}, fmt.Errorf("record quest completion: %w", err)
	}

	if applied {
		if _, err = tx.Exec(ctx, `
			UPDATE characters
			SET money=money+$2, xp=xp+$3, reputation=reputation+$4
			WHERE id=$1
		`, characterID, money, xp, reputation); err != nil {
			return core.Snapshot{}, core.RewardReceipt{}, fmt.Errorf("apply quest reward: %w", err)
		}
		if questID == "MQ012" {
			if _, err = tx.Exec(ctx, "UPDATE vehicles SET roadworthy=true WHERE owner_character_id=$1 AND starter_lineage=true", characterID); err != nil {
				return core.Snapshot{}, core.RewardReceipt{}, fmt.Errorf("mark starter roadworthy: %w", err)
			}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return core.Snapshot{}, core.RewardReceipt{}, fmt.Errorf("commit quest mutation: %w", err)
	}

	snapshot, err := p.snapshotByAccount(ctx, accountID)
	receipt := core.RewardReceipt{
		QuestID: questID, Applied: applied, Money: money, XP: xp, Reputation: reputation,
	}
	return snapshot, receipt, err
}

func (p *Postgres) ReviseBuild(ctx context.Context, tokenHash, vehicleID string, expectedRevision int, partIDs []string, operationID string) (core.Snapshot, error) {
	if operationID == "" {
		return core.Snapshot{}, ErrOperationKey
	}
	normalized, err := core.NormalizeParts(partIDs)
	if err != nil {
		return core.Snapshot{}, err
	}
	buildHash, err := core.BuildHash(normalized)
	if err != nil {
		return core.Snapshot{}, err
	}
	partJSON, err := json.Marshal(normalized)
	if err != nil {
		return core.Snapshot{}, fmt.Errorf("marshal build: %w", err)
	}

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return core.Snapshot{}, fmt.Errorf("begin build mutation: %w", err)
	}
	defer tx.Rollback(ctx)

	var accountID, characterID string
	err = tx.QueryRow(ctx, `
		SELECT a.id, c.id
		FROM sessions s
		JOIN accounts a ON a.id=s.account_id
		JOIN characters c ON c.account_id=a.id
		WHERE s.token_hash=$1 AND s.expires_at > now()
	`, tokenHash).Scan(&accountID, &characterID)
	if errors.Is(err, pgx.ErrNoRows) {
		return core.Snapshot{}, ErrUnauthorized
	}
	if err != nil {
		return core.Snapshot{}, fmt.Errorf("authorize build mutation: %w", err)
	}

	var existingVehicle string
	var existingRevision int
	err = tx.QueryRow(ctx, "SELECT vehicle_id,revision FROM vehicle_builds WHERE operation_id=$1", operationID).Scan(&existingVehicle, &existingRevision)
	if err == nil {
		if existingVehicle != vehicleID {
			return core.Snapshot{}, ErrOperationKey
		}
		if err := tx.Commit(ctx); err != nil {
			return core.Snapshot{}, err
		}
		return p.snapshotByAccount(ctx, accountID)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return core.Snapshot{}, fmt.Errorf("check build operation: %w", err)
	}

	var currentRevision int
	err = tx.QueryRow(ctx, `
		SELECT active_revision
		FROM vehicles
		WHERE id=$1 AND owner_character_id=$2
		FOR UPDATE
	`, vehicleID, characterID).Scan(&currentRevision)
	if errors.Is(err, pgx.ErrNoRows) {
		return core.Snapshot{}, ErrNotFound
	}
	if err != nil {
		return core.Snapshot{}, fmt.Errorf("lock vehicle: %w", err)
	}
	if currentRevision != expectedRevision {
		return core.Snapshot{}, ErrConflict
	}

	nextRevision := currentRevision + 1
	if _, err = tx.Exec(ctx, `
		INSERT INTO vehicle_builds(vehicle_id,revision,part_ids,validation_hash,operation_id)
		VALUES($1,$2,$3,$4,$5)
	`, vehicleID, nextRevision, partJSON, buildHash, operationID); err != nil {
		return core.Snapshot{}, fmt.Errorf("insert build revision: %w", err)
	}
	if _, err = tx.Exec(ctx, "UPDATE vehicles SET active_revision=$2 WHERE id=$1", vehicleID, nextRevision); err != nil {
		return core.Snapshot{}, fmt.Errorf("activate build revision: %w", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return core.Snapshot{}, fmt.Errorf("commit build mutation: %w", err)
	}
	return p.snapshotByAccount(ctx, accountID)
}

func (p *Postgres) snapshotByAccount(ctx context.Context, accountID string) (core.Snapshot, error) {
	var snapshot core.Snapshot
	err := p.pool.QueryRow(ctx, `
		SELECT
			a.id,
			c.id,
			c.money,
			c.xp,
			c.reputation,
			v.id,
			v.active_revision,
			v.starter_lineage,
			v.roadworthy
		FROM accounts a
		JOIN characters c ON c.account_id=a.id
		JOIN vehicles v ON v.owner_character_id=c.id AND v.starter_lineage=true
		WHERE a.id=$1
		LIMIT 1
	`, accountID).Scan(
		&snapshot.AccountID,
		&snapshot.CharacterID,
		&snapshot.Money,
		&snapshot.XP,
		&snapshot.Reputation,
		&snapshot.VehicleID,
		&snapshot.ActiveBuildRevision,
		&snapshot.StarterLineage,
		&snapshot.Roadworthy,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return core.Snapshot{}, ErrNotFound
	}
	if err != nil {
		return core.Snapshot{}, fmt.Errorf("load snapshot: %w", err)
	}

	var partJSON []byte
	if err = p.pool.QueryRow(ctx, `
		SELECT part_ids
		FROM vehicle_builds
		WHERE vehicle_id=$1 AND revision=$2
	`, snapshot.VehicleID, snapshot.ActiveBuildRevision).Scan(&partJSON); err != nil {
		return core.Snapshot{}, fmt.Errorf("load active build: %w", err)
	}
	if err = json.Unmarshal(partJSON, &snapshot.ActivePartIDs); err != nil {
		return core.Snapshot{}, fmt.Errorf("decode active build: %w", err)
	}

	rows, err := p.pool.Query(ctx, "SELECT quest_id FROM quest_completions WHERE character_id=$1 ORDER BY quest_id", snapshot.CharacterID)
	if err != nil {
		return core.Snapshot{}, fmt.Errorf("load quest completions: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var questID string
		if err = rows.Scan(&questID); err != nil {
			return core.Snapshot{}, fmt.Errorf("scan quest completion: %w", err)
		}
		snapshot.CompletedQuests = append(snapshot.CompletedQuests, questID)
	}
	if err = rows.Err(); err != nil {
		return core.Snapshot{}, fmt.Errorf("iterate quest completions: %w", err)
	}

	return snapshot, nil
}
