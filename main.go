package main

import (
	"log"
	tgClient "tBot/app/clients/telegram"
	consumer "tBot/app/consumer/event-consumer"
	processor "tBot/app/events/telegram"
	"tBot/internal/env"
	"tBot/internal/storage/file"
)

func main() {
	cfg := env.MustLoadConfig()

	tg := tgClient.New(cfg.Telegram.Host, cfg.Telegram.Token)
	eventProcessor := processor.NewProcessor(tg, file.NewStorage(cfg.Storage.BasePath))
	eventConsumer := consumer.New(eventProcessor, eventProcessor, processor.BatchSize)
	log.Println("[INFO] Start event consumer")
	if err := eventConsumer.Start(); err != nil {
		log.Fatal("[ERROR] Failed to start event consumer", err)
	}
}
