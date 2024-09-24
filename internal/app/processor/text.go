package processor

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/Muaz717/willpower-bot/internal/lib/gym"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/looplab/fsm"
)

var (
	addWorkoutFSM    = "addWorkout"
	addingWorkoutFSM = "addingWorkout"

	addPullupsFSM    = "addPullups"
	addingPullupsFSM = "addingPullups"

	canceledFSM = "canceled"
)

func (p *Processor) NewWorkout(bot *tgbotapi.BotAPI, update tgbotapi.Update, newFSM *fsm.FSM) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch newFSM.Current() {
	case canceledFSM:
		text := "Добавляем тренировку. Введите ваш вес.\n\nФормат: 87.6"
		msg.Text = text

		_, err := bot.Send(msg)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		err = newFSM.Event(p.ctx, addWorkoutFSM)
		if err != nil {
			return fmt.Errorf("failed to get new event: %w", err)
		}
	case addingWorkoutFSM:
		msg.Text = "Вы начали добавлять тренировку. Введите ваш вес.\n\nФормат: 87.6"
		_, err := bot.Send(msg)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
	}

	return nil
}

func (p *Processor) WorkoutStat(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	workouts, err := p.srv.WorkoutGym.GetAll(p.ctx, int(update.Message.Chat.ID))
	if err != nil {
		p.log.Info(err.Error())

		return fmt.Errorf("failed to get all workouts: %w", err)
	}

	var result string

	for _, workout := range workouts {
		result += fmt.Sprintf("Номер: %s\n", strconv.Itoa(workout.RowNum))
		result += fmt.Sprintf("Дата: %s\n", workout.Date)
		result += fmt.Sprintf("Вес: %s\n\n", strconv.FormatFloat(workout.Weight, 'f', 1, 64))
	}

	msg.Text = result
	_, err = bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

var cancelWorkoutButton = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена 🚫", "Отмена 🚫"),
	),
)

func (p *Processor) AddWorkout(bot *tgbotapi.BotAPI, update tgbotapi.Update, newFSM *fsm.FSM) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	var workout gym.Workout
	input := update.Message.Text

	weight, err := p.isFloat(input, msg, bot)
	if err != nil {
		return err
	}

	workout.Weight = weight

	userId, err := p.srv.WorkoutGym.SaveUser(p.ctx, int(msg.ChatID), update.Message.From.UserName)
	if err != nil {
		p.log.Info("User already exists")
	} else {
		p.log.Info("User saved", slog.Int("user_id", userId))
	}

	workoutId, err := p.srv.WorkoutGym.Create(p.ctx, int(update.Message.Chat.ID), workout)
	if err != nil {
		p.log.Info("failed to add workout")

		return err
	}
	p.log.Info("Workout added", slog.Int("workout_id", workoutId))

	msg.Text = "Тренировка добавлена"
	_, err = bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	err = newFSM.Event(p.ctx, "cancel")
	if err != nil {
		return fmt.Errorf("failed to get new event: %w", err)
	}

	return nil
}

func (p *Processor) NewPullups(bot *tgbotapi.BotAPI, update tgbotapi.Update, newFSM *fsm.FSM) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch newFSM.Current() {
	case canceledFSM:
		text := "Добавляем подтягивания. Введите количество за 3 подхода.\n\nФормат: 10"
		msg.Text = text

		_, err := bot.Send(msg)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		err = newFSM.Event(p.ctx, addPullupsFSM)
		if err != nil {
			return fmt.Errorf("failed to get new event: %w", err)
		}
	case addingPullupsFSM:
		msg.Text = "Вы начали добавлять подтягивания. Введите количество.\n\nФормат: 10"
		_, err := bot.Send(msg)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
	}

	return nil
}

func (p *Processor) PullupsStat(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	pullupsStat, err := p.srv.PullUpsGym.GetAll(p.ctx, int(update.Message.Chat.ID))
	if err != nil {
		p.log.Info(err.Error())

		return fmt.Errorf("failed to get all workouts: %w", err)
	}

	var result string

	for _, pullups := range pullupsStat {
		result += fmt.Sprintf("Номер: %s\n", strconv.Itoa(pullups.RowNum))
		result += fmt.Sprintf("Дата: %s\n", pullups.Date)
		result += fmt.Sprintf("Количество: %s\n\n", strconv.Itoa(pullups.Quantity))
	}

	msg.Text = result
	_, err = bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

var cancelPullupsButton = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена ⛔", "Отмена ⛔"),
	),
)

func (p *Processor) AddPullups(bot *tgbotapi.BotAPI, update tgbotapi.Update, newFSM *fsm.FSM) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	var pullups gym.PullUps
	input := update.Message.Text

	quantity, err := p.isInt(input, msg, bot)
	if err != nil {
		return err
	}
	pullups.Quantity = quantity

	pullupsId, err := p.srv.PullUpsGym.Create(p.ctx, int(msg.ChatID), pullups)
	if err != nil {
		p.log.Info("failed to add pullups")

		return err
	}
	p.log.Info("Pullups added", slog.Int("pullups_id", pullupsId))

	msg.Text = "Подтягивания добавлены"
	_, err = bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	err = newFSM.Event(p.ctx, "cancel")
	if err != nil {
		return fmt.Errorf("failed to get new event: %w", err)
	}

	return nil
}

func (p *Processor) isInt(input string, msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) (int, error) {
	quantity, err := strconv.Atoi(input)

	if err != nil {
		msg.Text = "Неправильный формат!\n\nНажмите ОТМЕНА, если передумали 👇"
		msg.ReplyMarkup = cancelPullupsButton

		_, err = bot.Send(msg)
		if err != nil {
			return 0, fmt.Errorf("failed to send message: %w", err)
		}

		p.log.Info("Неправильный формат")

		return 0, fmt.Errorf("wrong quantity type")
	}

	return quantity, nil
}

func (p *Processor) isFloat(input string, msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) (float64, error) {
	weight, err := strconv.ParseFloat(input, 64)

	if err != nil {
		msg.Text = "Неправильный формат!\n\nНажмите ОТМЕНА, если передумали 👇"
		msg.ReplyMarkup = cancelWorkoutButton

		_, err = bot.Send(msg)
		if err != nil {
			return 0, fmt.Errorf("failed to send message: %w", err)
		}

		p.log.Info("Неправильный формат")

		return 0, fmt.Errorf("wrong weight type")
	}

	return weight, nil
}
