package storage

import (
	"context"

	pullupsStorage "github.com/Muaz717/willpower-bot/internal/app/storage/pull-ups"
	workoutStorage "github.com/Muaz717/willpower-bot/internal/app/storage/workout"
	"github.com/Muaz717/willpower-bot/internal/lib/gym"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkoutGym interface {
	SaveUser(ctx context.Context, chatId int, userName string) (int, error)
	Create(ctx context.Context, chatId int, workout gym.Workout) (int, error)
	GetAll(ctx context.Context, chatId int) ([]gym.Workout, error)
	Delete(ctx context.Context, chatId, workoutId int) error
}

type PullUpsGym interface {
	Create(ctx context.Context, chatId int, pullUps gym.PullUps) (int, error)
	GetAll(ctx context.Context, chatId int) ([]gym.PullUps, error)
	Delete(ctx context.Context, chatId, pullUpsId int) error
}

type Storage struct {
	WorkoutGym
	PullUpsGym
}

func New(db *pgxpool.Pool) *Storage {
	return &Storage{
		WorkoutGym: workoutStorage.New(db),
		PullUpsGym: pullupsStorage.New(db),
	}
}
