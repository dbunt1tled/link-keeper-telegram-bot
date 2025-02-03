package storage

import (
	"crypto/sha1"
	"fmt"
	"io"
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

func (p *Page) Hash() (string, error) {
	var err error
	h := sha1.New()
	_, err = io.WriteString(h, p.URL)
	if err != nil {
		return "", errors.Wrap(err, "failed to write url to hash")
	}
	_, err = io.WriteString(h, p.UserName)
	if err != nil {
		return "", errors.Wrap(err, "failed to write user name to hash")
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
