package event_consumer

import (
	"log"
	"tBot/app/events"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}
func (c *Consumer) Start() error {
	for {
		ev, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERROR] Failed to fetch events: %v", err)
			continue
		}
		if len(ev) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		err = c.handleEvents(ev)
		if err != nil {
			log.Printf("[ERROR] Failed to handle events: %v", err)
		}
	}
}

func (c *Consumer) handleEvents(ev []events.Event) error {
	for _, e := range ev {
		log.Printf("[INFO] new event: %v", e)
		if err := c.processor.Process(e); err != nil {
			log.Printf("[ERROR] Failed handle event: %v", err)
			continue
		}
	}
	return nil
}
