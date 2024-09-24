package processor

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/looplab/fsm"
)

func (p *Processor) NewWorkoutCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update, newFSM *fsm.FSM) error {

	err := newFSM.Event(p.ctx, "cancel")
	if err != nil {
		return err
	}

	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	if _, err := bot.Request(callback); err != nil {
		return fmt.Errorf("failed to request callback: %w", err)
	}

	delmsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	if _, err := bot.Request(delmsg); err != nil {
		return fmt.Errorf("failed to request delete message: %w", err)
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Операция отменена 🚫")
	msg.ReplyMarkup = GymKeyboard

	if _, err := bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (p *Processor) NewPullupsCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update, newFSM *fsm.FSM) error {

	err := newFSM.Event(p.ctx, "cancel")
	if err != nil {
		return err
	}

	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	if _, err := bot.Request(callback); err != nil {
		return fmt.Errorf("failed to request callback: %w", err)
	}

	delmsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	if _, err := bot.Request(delmsg); err != nil {
		return fmt.Errorf("failed to request delete message: %w", err)
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Операция отменена 🚫")
	msg.ReplyMarkup = GymKeyboard

	if _, err := bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
