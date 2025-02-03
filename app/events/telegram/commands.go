package telegram

import (
	"log"
	"net/url"
	"strings"
	"tBot/internal/storage"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

const (
	cmdHelp  = "/help"
	cmdStart = "/start"
	cmdRnd   = "/rnd"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)
	log.Printf("doCmd: command: %s, from %s", text, username)
	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}
	switch text {
	case cmdHelp:
		return p.sendHelp(chatID)
	case cmdStart:
		return p.sendStart(chatID)
	case cmdRnd:
		return p.sendRnd(chatID, username)
	default:
		return newSendMessage(chatID, p.tg)(MsgUnknown)
	}
}

func (p *Processor) sendHelp(chatID int) error {
	return newSendMessage(chatID, p.tg)(MsgHelp)
}

func (p *Processor) sendStart(chatID int) error {
	return newSendMessage(chatID, p.tg)(MsgStart)
}

func (p *Processor) sendRnd(chatID int, username string) error {
	send := newSendMessage(chatID, p.tg)
	page, err := p.storage.PickRandom(username)
	if err != nil {
		if errors.Is(err, storage.ErrPageNotFound) {
			return send(MsgNoSavedPages)
		}
		return errors.Wrap(err, "failed to pick random page")
	}
	if err = send(page.URL); err != nil {
		return errors.Wrap(err, "failed to send page")
	}

	return p.storage.Remove(page)
}

func (p *Processor) savePage(chatID int, pageURL string, username string) error {
	send := newSendMessage(chatID, p.tg)
	page := &storage.Page{
		URL:       pageURL,
		UserName:  username,
		CreatedAt: time.Time{},
	}
	isExist, err := p.storage.Exists(page)
	if err != nil {
		return err
	}
	if isExist {
		return send(msgAlreadySavedPage)
	}
	if err = p.storage.Save(page); err != nil {
		return err
	}

	return send(MsgSavedPage)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}

func newSendMessage(chatID int, tg *tgbotapi.BotAPI) func(string) error {
	return func(msg string) error {
		_, err := tg.Send(tgbotapi.NewMessage(int64(chatID), msg))
		return err
	}
}
