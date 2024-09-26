package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/Muaz717/willpower-bot/internal/app/processor"
	"github.com/Muaz717/willpower-bot/internal/app/service"
	"github.com/Muaz717/willpower-bot/internal/app/storage"
	"github.com/Muaz717/willpower-bot/internal/app/storage/postgres"
	"github.com/Muaz717/willpower-bot/internal/config"
	"github.com/Muaz717/willpower-bot/internal/lib/logger/sl"
	"github.com/Muaz717/willpower-bot/internal/lib/logger/slogpretty"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	ctx   context.Context
	cfg   *config.Config
	log   *slog.Logger
	store *storage.Storage
	srv   *service.Service
	proc  *processor.Processor
	bot   *tgbotapi.BotAPI
}

func New() (*App, error) {
	a := &App{}

	a.ctx = context.Background()

	a.cfg = config.New()

	a.log = a.SetupLogger()

	db, err := postgres.New(a.ctx, *a.cfg)
	if err != nil {
		a.log.Error("failed to init database: ", sl.Err(err))

		return nil, fmt.Errorf("failed to init database: %w", err)
	}

	a.store = storage.New(db)

	a.srv = service.New(a.store)

	a.proc = processor.New(a.srv, a.log, a.ctx)

	return a, nil
}

func (a *App) StartBot() error {

	bot, err := a.proc.InitBot()
	if err != nil {
		a.log.Error("failed to init bot: ", sl.Err(err))

		return fmt.Errorf("failed to init bot: %w", err)
	}
	a.bot = bot

	return nil
}

func (a *App) SetupLogger() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	log := slog.New(handler)

	return log
}
