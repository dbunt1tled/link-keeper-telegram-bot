package telegram

import (
	"tBot/app/events"
	"tBot/internal/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

const BatchSize = 2

type Meta struct {
	ChatID   int
	Username string
}

type Processor struct {
	tg      *tgbotapi.BotAPI
	offset  int
	storage storage.Storage
}

func NewProcessor(tg *tgbotapi.BotAPI, storage storage.Storage) *Processor {
	return &Processor{
		tg:      tg,
		offset:  0,
		storage: storage,
	}
}
func (p *Processor) Fetch(limit int) (events.ChEvent, error) {

	u := tgbotapi.NewUpdate(p.offset)
	u.Timeout = 60
	u.Limit = limit

	updates := p.tg.GetUpdatesChan(u)
	ch := make(chan events.Event, 10)
	go func() {
		defer close(ch)
		for update := range updates {
			ch <- event(update)
			p.offset = update.UpdateID + 1
		}
	}()
	return ch, nil
}
func (p *Processor) Process(e events.Event) error {
	switch e.Type {
	case events.Message:
		return p.processMessage(e)
	default:
		return errors.New("unsupported event type")
	}
}

func (p *Processor) processMessage(e events.Event) error {
	m, err := meta(e)
	if err != nil {
		return errors.Wrap(err, "failed to get meta")
	}
	if err = p.doCmd(e.Text, m.ChatID, m.Username); err != nil {
		return errors.Wrap(err, "failed to do command")
	}

	return nil
}

func meta(e events.Event) (Meta, error) {
	res, ok := e.Meta.(Meta)
	if !ok {
		return Meta{}, errors.New("invalid meta")
	}
	return res, nil
}

func event(update tgbotapi.Update) events.Event {
	uType := fetchType(update)
	ev := events.Event{
		Type: uType,
		Text: fetchText(update),
	}
	if uType == events.Message {
		ev.Meta = Meta{
			ChatID:   int(update.Message.Chat.ID),
			Username: update.Message.From.UserName,
		}
	}

	return ev
}

func fetchText(update tgbotapi.Update) string {
	if update.Message == nil {
		return ""
	}
	return update.Message.Text
}

func fetchType(update tgbotapi.Update) events.Type {
	if update.Message == nil {
		return events.Unknown
	}
	return events.Message
}
