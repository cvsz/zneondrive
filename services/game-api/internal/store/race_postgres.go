package store

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/cvsz/zneondrive/services/game-api/internal/core"
	"github.com/jackc/pgx/v5"
)

//go:embed migrations/004_race_runtime.sql
var raceMigrationSQL string

func (p *Postgres) EnsureRaceSchema(ctx context.Context) error {
	if _, err := p.pool.Exec(ctx, raceMigrationSQL); err != nil {
		return fmt.Errorf("apply race schema 004_race_runtime: %w", err)
	}
	return nil
}

func (p *Postgres) StartRace(ctx context.Context, accountID, vehicleID, raceID, operationID string) (core.RaceInstance, error) {
	if operationID == "" {
		return core.RaceInstance{}, ErrOperationKey
	}
	normalizedRaceID, err := core.NormalizeRaceID(raceID)
	if err != nil {
		return core.RaceInstance{}, err
	}

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return core.RaceInstance{}, fmt.Errorf("begin race start: %w", err)
	}
	defer tx.Rollback(ctx)

	var existing core.RaceInstance
	err = tx.QueryRow(ctx, `
		SELECT id,race_id,account_id,character_id,vehicle_id,build_revision,build_validation_hash,state,next_checkpoint,last_elapsed_ms
		FROM race_instances WHERE operation_id=$1
	`, operationID).Scan(
		&existing.RaceInstanceID, &existing.RaceID, &existing.AccountID, &existing.CharacterID,
		&existing.VehicleID, &existing.BuildRevision, &existing.BuildValidationHash, &existing.State,
		&existing.NextCheckpoint, &existing.LastElapsedMS,
	)
	if err == nil {
		if existing.AccountID != accountID || existing.VehicleID != vehicleID || existing.RaceID != normalizedRaceID {
			return core.RaceInstance{}, ErrOperationKey
		}
		if err = tx.Commit(ctx); err != nil {
			return core.RaceInstance{}, err
		}
		return existing, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return core.RaceInstance{}, fmt.Errorf("check race start operation: %w", err)
	}

	var characterID, buildHash string
	var buildRevision int
	var roadworthy bool
	err = tx.QueryRow(ctx, `
		SELECT c.id,v.active_revision,v.roadworthy,vb.validation_hash
		FROM accounts a
		JOIN characters c ON c.account_id=a.id
		JOIN vehicles v ON v.owner_character_id=c.id
		JOIN vehicle_builds vb ON vb.vehicle_id=v.id AND vb.revision=v.active_revision
		WHERE a.id=$1 AND v.id=$2
		FOR UPDATE OF v
	`, accountID, vehicleID).Scan(&characterID, &buildRevision, &roadworthy, &buildHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return core.RaceInstance{}, ErrNotFound
	}
	if err != nil {
		return core.RaceInstance{}, fmt.Errorf("bind race vehicle: %w", err)
	}
	if !roadworthy {
		return core.RaceInstance{}, ErrRaceNotReady
	}

	raceInstanceID, err := core.NewID("raceinst")
	if err != nil {
		return core.RaceInstance{}, err
	}
	instance := core.RaceInstance{
		RaceInstanceID: raceInstanceID,
		RaceID: normalizedRaceID,
		AccountID: accountID,
		CharacterID: characterID,
		VehicleID: vehicleID,
		BuildRevision: buildRevision,
		BuildValidationHash: buildHash,
		State: "active",
		NextCheckpoint: 0,
		LastElapsedMS: 0,
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO race_instances(
			id,race_id,account_id,character_id,vehicle_id,build_revision,build_validation_hash,state,operation_id
		) VALUES($1,$2,$3,$4,$5,$6,$7,'active',$8)
	`, instance.RaceInstanceID, instance.RaceID, instance.AccountID, instance.CharacterID,
		instance.VehicleID, instance.BuildRevision, instance.BuildValidationHash, operationID)
	if err != nil {
		return core.RaceInstance{}, fmt.Errorf("insert race instance: %w", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return core.RaceInstance{}, fmt.Errorf("commit race start: %w", err)
	}
	return instance, nil
}

func (p *Postgres) RecordRaceCheckpoint(ctx context.Context, raceInstanceID string, checkpointIndex int, elapsedMS int64, operationID string) (core.RaceInstance, error) {
	if operationID == "" {
		return core.RaceInstance{}, ErrOperationKey
	}
	if err := core.ValidateRaceCheckpoint(checkpointIndex, elapsedMS); err != nil {
		return core.RaceInstance{}, err
	}
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return core.RaceInstance{}, fmt.Errorf("begin race checkpoint: %w", err)
	}
	defer tx.Rollback(ctx)

	var existingInstanceID string
	var existingIndex int
	var existingElapsed int64
	err = tx.QueryRow(ctx, `SELECT race_instance_id,checkpoint_index,elapsed_ms FROM race_checkpoints WHERE operation_id=$1`, operationID).
		Scan(&existingInstanceID, &existingIndex, &existingElapsed)
	if err == nil {
		if existingInstanceID != raceInstanceID || existingIndex != checkpointIndex || existingElapsed != elapsedMS {
			return core.RaceInstance{}, ErrOperationKey
		}
		if err = tx.Commit(ctx); err != nil {
			return core.RaceInstance{}, err
		}
		return p.raceInstanceByID(ctx, raceInstanceID)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return core.RaceInstance{}, fmt.Errorf("check checkpoint operation: %w", err)
	}

	instance, err := scanRaceInstance(tx.QueryRow(ctx, `
		SELECT id,race_id,account_id,character_id,vehicle_id,build_revision,build_validation_hash,state,next_checkpoint,last_elapsed_ms
		FROM race_instances WHERE id=$1 FOR UPDATE
	`, raceInstanceID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core.RaceInstance{}, ErrNotFound
		}
		return core.RaceInstance{}, fmt.Errorf("lock race instance: %w", err)
	}
	if instance.State != "active" || checkpointIndex != instance.NextCheckpoint || elapsedMS <= instance.LastElapsedMS {
		return core.RaceInstance{}, ErrRaceOrder
	}
	if _, err = tx.Exec(ctx, `
		INSERT INTO race_checkpoints(race_instance_id,checkpoint_index,elapsed_ms,operation_id)
		VALUES($1,$2,$3,$4)
	`, raceInstanceID, checkpointIndex, elapsedMS, operationID); err != nil {
		return core.RaceInstance{}, fmt.Errorf("insert checkpoint: %w", err)
	}
	if _, err = tx.Exec(ctx, `
		UPDATE race_instances SET next_checkpoint=next_checkpoint+1,last_elapsed_ms=$2 WHERE id=$1
	`, raceInstanceID, elapsedMS); err != nil {
		return core.RaceInstance{}, fmt.Errorf("advance checkpoint cursor: %w", err)
	}
	instance.NextCheckpoint++
	instance.LastElapsedMS = elapsedMS
	if err = tx.Commit(ctx); err != nil {
		return core.RaceInstance{}, fmt.Errorf("commit checkpoint: %w", err)
	}
	return instance, nil
}

func (p *Postgres) FinishRace(ctx context.Context, raceInstanceID string, checkpointCount int, finishElapsedMS int64, operationID string) (core.RaceResult, error) {
	if operationID == "" {
		return core.RaceResult{}, ErrOperationKey
	}
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return core.RaceResult{}, fmt.Errorf("begin race finish: %w", err)
	}
	defer tx.Rollback(ctx)

	var existing core.RaceResult
	err = tx.QueryRow(ctx, `
		SELECT rr.race_instance_id,ri.race_id,ri.vehicle_id,ri.build_revision,rr.checkpoint_count,rr.finish_elapsed_ms,rr.result_hash
		FROM race_results rr JOIN race_instances ri ON ri.id=rr.race_instance_id
		WHERE rr.operation_id=$1
	`, operationID).Scan(&existing.RaceInstanceID, &existing.RaceID, &existing.VehicleID, &existing.BuildRevision,
		&existing.CheckpointCount, &existing.FinishElapsedMS, &existing.ResultHash)
	if err == nil {
		if existing.RaceInstanceID != raceInstanceID || existing.CheckpointCount != checkpointCount || existing.FinishElapsedMS != finishElapsedMS {
			return core.RaceResult{}, ErrOperationKey
		}
		if err = tx.Commit(ctx); err != nil {
			return core.RaceResult{}, err
		}
		return existing, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return core.RaceResult{}, fmt.Errorf("check race finish operation: %w", err)
	}

	instance, err := scanRaceInstance(tx.QueryRow(ctx, `
		SELECT id,race_id,account_id,character_id,vehicle_id,build_revision,build_validation_hash,state,next_checkpoint,last_elapsed_ms
		FROM race_instances WHERE id=$1 FOR UPDATE
	`, raceInstanceID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core.RaceResult{}, ErrNotFound
		}
		return core.RaceResult{}, fmt.Errorf("lock race finish: %w", err)
	}
	if instance.State != "active" || checkpointCount != instance.NextCheckpoint || checkpointCount < 1 || finishElapsedMS <= instance.LastElapsedMS {
		return core.RaceResult{}, ErrRaceOrder
	}
	resultHash, err := core.RaceResultHash(instance, checkpointCount, finishElapsedMS)
	if err != nil {
		return core.RaceResult{}, err
	}
	result := core.RaceResult{
		RaceInstanceID: raceInstanceID,
		RaceID: instance.RaceID,
		VehicleID: instance.VehicleID,
		BuildRevision: instance.BuildRevision,
		CheckpointCount: checkpointCount,
		FinishElapsedMS: finishElapsedMS,
		ResultHash: resultHash,
	}
	if _, err = tx.Exec(ctx, `
		INSERT INTO race_results(race_instance_id,checkpoint_count,finish_elapsed_ms,result_hash,operation_id)
		VALUES($1,$2,$3,$4,$5)
	`, raceInstanceID, checkpointCount, finishElapsedMS, resultHash, operationID); err != nil {
		return core.RaceResult{}, fmt.Errorf("insert race result: %w", err)
	}
	if _, err = tx.Exec(ctx, `UPDATE race_instances SET state='finished',finished_at=now() WHERE id=$1`, raceInstanceID); err != nil {
		return core.RaceResult{}, fmt.Errorf("finish race instance: %w", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return core.RaceResult{}, fmt.Errorf("commit race finish: %w", err)
	}
	return result, nil
}

func (p *Postgres) raceInstanceByID(ctx context.Context, raceInstanceID string) (core.RaceInstance, error) {
	instance, err := scanRaceInstance(p.pool.QueryRow(ctx, `
		SELECT id,race_id,account_id,character_id,vehicle_id,build_revision,build_validation_hash,state,next_checkpoint,last_elapsed_ms
		FROM race_instances WHERE id=$1
	`, raceInstanceID))
	if errors.Is(err, pgx.ErrNoRows) {
		return core.RaceInstance{}, ErrNotFound
	}
	return instance, err
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanRaceInstance(row rowScanner) (core.RaceInstance, error) {
	var instance core.RaceInstance
	err := row.Scan(
		&instance.RaceInstanceID, &instance.RaceID, &instance.AccountID, &instance.CharacterID,
		&instance.VehicleID, &instance.BuildRevision, &instance.BuildValidationHash, &instance.State,
		&instance.NextCheckpoint, &instance.LastElapsedMS,
	)
	return instance, err
}
