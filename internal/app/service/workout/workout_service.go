package workoutService

import (
	"context"
	"time"

	"github.com/Muaz717/willpower-bot/internal/app/storage"
	"github.com/Muaz717/willpower-bot/internal/lib/gym"
)

type WorkoutService struct {
	store storage.WorkoutGym
}

func New(store storage.WorkoutGym) *WorkoutService {
	return &WorkoutService{store: store}
}

func (srv *WorkoutService) SaveUser(ctx context.Context, chatId int, userName string) (int, error) {
	return srv.store.SaveUser(ctx, chatId, userName)
}

func (srv *WorkoutService) Create(ctx context.Context, chatId int, workout gym.Workout) (int, error) {
	now := time.Now()
	date := now.Format("02-01-2006")

	workout.Date = date

	return srv.store.Create(ctx, chatId, workout)
}

func (srv *WorkoutService) GetAll(ctx context.Context, chatId int) ([]gym.Workout, error) {
	return srv.store.GetAll(ctx, chatId)
}

func (srv *WorkoutService) Delete(ctx context.Context, chatId, workoutId int) error {
	return srv.store.Delete(ctx, chatId, workoutId)
}
