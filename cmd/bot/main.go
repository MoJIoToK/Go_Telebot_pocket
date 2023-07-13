package main

import (
	"log"

	"github.com/MoJIoToK/Go_Telebot_pocket/pkg/config"
	"github.com/MoJIoToK/Go_Telebot_pocket/pkg/server"
	"github.com/MoJIoToK/Go_Telebot_pocket/pkg/storage"
	"github.com/MoJIoToK/Go_Telebot_pocket/pkg/storage/boltdb"
	telegram "github.com/MoJIoToK/Go_Telebot_pocket/pkg/telegram"
	pocket "github.com/MoJIoToK/go_Pocket_SDK"
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	botApi, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}
	botApi.Debug = true

	//pocketClient, err := pocket.new_Client(cfg.PocketConsumerKey)
	pocketClient, err := pocket.new_Client(cfg.PocketConsumerKey)

	db, err := initBolt()
	if err != nil {
		log.Fatal(err)
	}
	storage := boltdb.NewTokenStorage(db)

	bot := telegram.NewBot(botApi, pocketClient, cfg.AuthServerURL, storage, cfg.Messages)

	redirectServer := server.NewAuthServer(cfg.BotURL, storage, pocketClient)

	go func() {
		if err := redirectServer.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := bot.Start(); err != nil {
		log.Fatal(err)
	}
}

func initBolt() (*bolt.DB, error) {
	db, err := bolt.Open("bot.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Batch(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(storage.AccessTokens))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(storage.RequestTokens))
		return err
	}); err != nil {
		return nil, err
	}

	return db, nil
}
