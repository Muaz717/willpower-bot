package pullupsService

import (
	"context"
	"time"

	"github.com/Muaz717/willpower-bot/internal/app/storage"
	"github.com/Muaz717/willpower-bot/internal/lib/gym"
)

type PullUpsService struct {
	store storage.PullUpsGym
}

func New(store storage.PullUpsGym) *PullUpsService {
	return &PullUpsService{
		store: store,
	}
}

func (srv *PullUpsService) Create(ctx context.Context, chatId int, pullUps gym.PullUps) (int, error) {
	now := time.Now()
	date := now.Format("02-01-2006")

	pullUps.Date = date

	return srv.store.Create(ctx, chatId, pullUps)
}

func (srv *PullUpsService) GetAll(ctx context.Context, chatId int) ([]gym.PullUps, error) {
	return srv.store.GetAll(ctx, chatId)
}

func (srv *PullUpsService) Delete(ctx context.Context, chatId, pullUpsId int) error {
	return srv.store.Delete(ctx, chatId, pullUpsId)
}
