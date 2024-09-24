package service

import (
	"context"

	pullupsService "github.com/Muaz717/willpower-bot/internal/app/service/pull-ups"
	workoutService "github.com/Muaz717/willpower-bot/internal/app/service/workout"
	"github.com/Muaz717/willpower-bot/internal/app/storage"
	"github.com/Muaz717/willpower-bot/internal/lib/gym"
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

type Service struct {
	WorkoutGym
	PullUpsGym
}

func New(store *storage.Storage) *Service {
	return &Service{
		WorkoutGym: workoutService.New(store.WorkoutGym),
		PullUpsGym: pullupsService.New(store.PullUpsGym),
	}
}
