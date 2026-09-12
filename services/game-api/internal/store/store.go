package store

import (
	"context"
	"errors"
	"time"

	"github.com/cvsz/zneondrive/services/game-api/internal/core"
)

var (
	ErrUnauthorized = errors.New("unauthorized session")
	ErrConflict     = errors.New("state conflict")
	ErrNotFound     = errors.New("not found")
	ErrOutOfOrder   = errors.New("quest prerequisite not complete")
	ErrOperationKey = errors.New("operation id reused for different mutation")
)

type Store interface {
	Ping(context.Context) error
	Bootstrap(context.Context, string) (core.Snapshot, bool, error)
	CreateSession(context.Context, string, string, time.Time) error
	SnapshotBySession(context.Context, string) (core.Snapshot, error)
	CompleteQuest(context.Context, string, string, string) (core.Snapshot, core.RewardReceipt, error)
	ReviseBuild(context.Context, string, string, int, []string, string) (core.Snapshot, error)
	Close()
}
