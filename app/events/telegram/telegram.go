package telegram

import (
	"tBot/app/clients/telegram"
	"tBot/app/events"
	"tBot/internal/storage"

	"github.com/pkg/errors"
)

const BatchSize = 100

type Meta struct {
	ChatID   int
	Username string
}

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

func NewProcessor(tg *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      tg,
		offset:  0,
		storage: storage,
	}
}
func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, err
	}
	if len(updates) == 0 {
		return nil, nil
	}
	evnt := make([]events.Event, 0, len(updates))
	for _, update := range updates {
		evnt = append(evnt, event(update))
	}
	p.offset = updates[len(updates)-1].ID + 1
	return evnt, nil
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
	if err := p.doCmd(e.Text, m.ChatID, m.Username); err != nil {
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

func event(update telegram.Update) events.Event {
	uType := fetchType(update)
	ev := events.Event{
		Type: uType,
		Text: fetchText(update),
	}
	if uType == events.Message {
		ev.Meta = Meta{
			ChatID:   update.Message.Chat.ID,
			Username: update.Message.From.UserName,
		}
	}

	return ev
}

func fetchText(update telegram.Update) string {
	if update.Message == nil {
		return ""
	}
	return update.Message.Text
}

func fetchType(update telegram.Update) events.Type {
	if update.Message == nil {
		return events.Unknown
	}
	return events.Message
}
