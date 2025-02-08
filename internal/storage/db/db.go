package db

import (
	"crypto/rand"
	"math/big"
	"tBot/app/page"
	"tBot/internal/database"
	"tBot/internal/storage"

	"github.com/pkg/errors"
)

type Storage struct {
	repository page.Repository
}

func NewStorage() Storage {
	return Storage{
		repository: page.Repository{},
	}
}
func (s Storage) Exists(p *storage.Page) (bool, error) {
	status := database.PageActive
	_, err := s.repository.First(page.DTO{
		UserName: &p.UserName,
		URL:      &p.URL,
		Status:   &status,
	})
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (s Storage) Remove(p *storage.Page) error {
	pg, err := s.repository.First(page.DTO{
		UserName: &p.UserName,
		URL:      &p.URL,
	})

	if err != nil {
		return errors.Wrap(err, "failed to get page")
	}
	if pg == nil {
		return nil
	}
	_, err = s.repository.Delete(pg.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s Storage) PickRandom(userName string) (*storage.Page, error) {
	var (
		err    error
		count  *int64
		index  int
		idxBig *big.Int
		pages  []page.Pages
	)
	status := database.PageActive
	count, err = s.repository.Count(page.DTO{
		UserName: &userName,
		Status:   &status,
	})
	if err != nil {
		return nil, err
	}
	if *count == 0 {
		return nil, storage.ErrPageNotFound
	}

	idxBig, err = rand.Int(rand.Reader, big.NewInt(*count))
	if err != nil {
		return nil, err
	}
	index = int(idxBig.Int64())
	limit := 1
	pages, err = s.repository.List(page.DTO{
		UserName: &userName,
		Status:   &status,
		Offset:   &index,
		Limit:    &limit,
	})
	if err != nil {
		return nil, err
	}
	if len(pages) == 0 {
		return nil, storage.ErrPageNotFound
	}
	return &storage.Page{
		URL:       pages[0].URL,
		UserName:  pages[0].UserName,
		CreatedAt: pages[0].CreatedAt,
	}, nil
}

func (s Storage) Save(p *storage.Page) error {
	_, err := s.repository.Save(
		p.URL,
		p.UserName,
		database.PageActive,
	)
	if err != nil {
		return err
	}
	return nil
}
