package telegram

import (
	"github.com/MoJIoToK/Go_Telebot_pocket/pkg/config"
	"github.com/MoJIoToK/Go_Telebot_pocket/pkg/storage"
	pocket "github.com/MoJIoToK/go_Pocket_SDK"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot         *tgbotapi.BotAPI
	client      *pocket.Client
	redirectURL string
	storage     storage.TokenStorage
	messages    config.Messages
}

func NewBot(bot *tgbotapi.BotAPI, client *pocket.Client, redirectURL string, storage storage.TokenStorage, messages config.Messages) *Bot {
	return &Bot{
		bot:         bot,
		client:      client,
		redirectURL: redirectURL,
		storage:     storage,
		messages:    messages,
	}
}

// Start() -
func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)
	// if err != nil {
	// 	return err
	// }

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() { // проверка команда ли это или нет
			if err := b.handleCommand(update.Message); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}

			continue
		}

		if err := b.handleMessage(update.Message); err != nil {
			b.handleError(update.Message.Chat.ID, err)
		}
	}

	return nil
}
