package main

import (
	"log"
	consumer "tBot/app/consumer/event-consumer"
	processor "tBot/app/events/telegram"
	"tBot/internal/database"
	"tBot/internal/env"
	"tBot/internal/storage/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg := env.MustLoadConfig()
	database.InitDBConnection(cfg)
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = cfg.App.Debug
	eventProcessor := processor.NewProcessor(bot, db.NewStorage())
	eventConsumer := consumer.New(eventProcessor, eventProcessor, processor.BatchSize)
	log.Println("[INFO] Start event consumer")
	if err = eventConsumer.Start(); err != nil {
		log.Fatal("[ERROR] Failed to start event consumer", err)
	}
}
