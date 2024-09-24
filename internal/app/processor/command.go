package processor

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processor) StartCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	text := fmt.Sprintf("Добро пожаловать, %s.\nЭто бот тренажерного зала WillpowerGym.", update.Message.From.FirstName)
	msg.Text = text

	msg.ReplyMarkup = GymKeyboard

	_, err := bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
