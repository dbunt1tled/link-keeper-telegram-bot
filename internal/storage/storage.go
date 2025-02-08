package storage

import (
	"time"

	"github.com/pkg/errors"
)

var (
	ErrPageNotFound = errors.New("page not found")
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	Exists(p *Page) (bool, error)
}

type Page struct {
	URL       string
	UserName  string
	CreatedAt time.Time
}
