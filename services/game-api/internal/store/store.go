package store

import (
	"context"
	"errors"
	"time"

	"github.com/cvsz/zneondrive/services/game-api/internal/core"
)

var (
	ErrUnauthorized          = errors.New("unauthorized session")
	ErrConflict              = errors.New("state conflict")
	ErrNotFound              = errors.New("not found")
	ErrOutOfOrder            = errors.New("quest prerequisite not complete")
	ErrOperationKey          = errors.New("operation id reused for different mutation")
	ErrInsufficientInventory = errors.New("insufficient inventory")
	ErrBlueprintRequired     = errors.New("required blueprint not unlocked")
	ErrRaceNotReady          = errors.New("vehicle is not eligible to start race")
	ErrRaceOrder             = errors.New("race event is out of order")
)

type Store interface {
	Ping(context.Context) error
	Bootstrap(context.Context, string) (core.Snapshot, bool, error)
	CreateSession(context.Context, string, string, time.Time) error
	SnapshotBySession(context.Context, string) (core.Snapshot, error)
	IssueGameTicket(context.Context, string, string, time.Time) error
	RedeemGameTicket(context.Context, string) (core.Snapshot, error)
	CompleteQuest(context.Context, string, string, string) (core.Snapshot, core.RewardReceipt, error)
	ReviseBuild(context.Context, string, string, int, []string, string) (core.Snapshot, error)
	StartRace(context.Context, string, string, string, string) (core.RaceInstance, error)
	RecordRaceCheckpoint(context.Context, string, int, int64, string) (core.RaceInstance, error)
	FinishRace(context.Context, string, int, int64, string) (core.RaceResult, error)
	Close()
}
