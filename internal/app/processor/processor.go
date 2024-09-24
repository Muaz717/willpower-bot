package processor

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/Muaz717/willpower-bot/internal/app/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/looplab/fsm"
)

const (
	commandStart = "start"
)

type Processor struct {
	srv *service.Service
	log *slog.Logger
	ctx context.Context
}

func New(srv *service.Service, log *slog.Logger, ctx context.Context) *Processor {
	return &Processor{
		srv: srv,
		log: log,
		ctx: ctx,
	}
}

var GymKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("🏋️‍♂️ Добавить тренировку"),
		tgbotapi.NewKeyboardButton("📈 Статистика тренировок"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("💪 Добавить подтягивания"),
		tgbotapi.NewKeyboardButton("📈 Статистика подтягиваний"),
	),
)

var (
	addWorkout  = GymKeyboard.Keyboard[0][0].Text
	workoutStat = GymKeyboard.Keyboard[0][1].Text
	addPullups  = GymKeyboard.Keyboard[1][0].Text
	pullupsStat = GymKeyboard.Keyboard[1][1].Text
)

func (p *Processor) InitBot() (*tgbotapi.BotAPI, error) {
	const op = "app.bot.InitBot"

	log := p.log

	newFSM := fsm.NewFSM(
		"canceled",
		fsm.Events{
			{Name: "addWorkout", Src: []string{"canceled"}, Dst: "addingWorkout"},
			{Name: "addPullups", Src: []string{"canceled"}, Dst: "addingPullups"},
			{Name: "cancel", Src: []string{"addingWorkout", "addingPullups"}, Dst: "canceled"},
		},
		fsm.Callbacks{},
	)

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message != nil {

			switch update.Message.Command() {
			case commandStart:
				err := p.StartCommand(bot, update)
				if err != nil {
					log.Info(err.Error())
				}
			}

			switch update.Message.Text {
			case addWorkout:
				err := p.NewWorkout(bot, update, newFSM)
				if err != nil {
					log.Info(err.Error())
				}
			case workoutStat:
				err := p.WorkoutStat(bot, update)
				if err != nil {
					log.Info(err.Error())
				}
			case addPullups:
				err := p.NewPullups(bot, update, newFSM)
				if err != nil {
					log.Info(err.Error())
				}
			case pullupsStat:
				err := p.PullupsStat(bot, update)
				if err != nil {
					log.Info(err.Error())
				}
			default:
				switch newFSM.Current() {
				case addingWorkoutFSM:
					err := p.AddWorkout(bot, update, newFSM)
					if err != nil {
						log.Info(err.Error())
					}
				case addingPullupsFSM:
					err := p.AddPullups(bot, update, newFSM)
					if err != nil {
						log.Info(err.Error())
					}
				}
			}

		} else if update.CallbackQuery != nil {

			switch update.CallbackQuery.Data {
			case "Отмена 🚫":
				err := p.NewWorkoutCallback(bot, update, newFSM)
				if err != nil {
					log.Info(err.Error())
				}
			case "Отмена ⛔":
				err := p.NewPullupsCallback(bot, update, newFSM)
				if err != nil {
					log.Info(err.Error())
				}
			}
		}

	}

	return bot, err
}
