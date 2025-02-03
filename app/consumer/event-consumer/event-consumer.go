package event_consumer

import (
	"log"
	"tBot/app/events"
)

const (
	workerCount = 10
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
	ev, err := c.fetcher.Fetch(c.batchSize)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch events: %v", err)
		return err
	}

	for i := 0; i < workerCount; i++ {
		c.worker(ev)
	}
	return nil
}

func (c *Consumer) handleEvents(e events.Event) error {
	return c.processor.Process(e)
}

func (c *Consumer) worker(ev events.ChEvent) {
	for e := range ev {
		err := c.handleEvents(e)
		if err != nil {
			log.Printf("[ERROR] Failed to handle events: %v", err)
		}
	}
}
