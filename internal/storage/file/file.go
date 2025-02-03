package file

import (
	"encoding/gob"
	"math/rand"
	"os"
	"path/filepath"
	"tBot/internal/storage"
	"time"

	"github.com/pkg/errors"
)

type Storage struct {
	basePath string
}

func NewStorage(basePath string) Storage {
	return Storage{
		basePath: basePath,
	}
}
func (s Storage) Exists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, errors.Wrap(err, "failed to get file name")
	}
	filePath := filepath.Join(s.basePath, p.UserName, fileName)
	_, err = os.Stat(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, errors.Wrap(err, "failed to stat file")
	}
	return true, nil
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return errors.Wrap(err, "failed to get file name")
	}
	filePath := filepath.Join(s.basePath, p.UserName, fileName)
	if err = os.Remove(filePath); err != nil {
		return errors.Wrap(err, "failed to remove file"+filePath)
	}
	return nil
}

func (s Storage) PickRandom(userName string) (*storage.Page, error) {
	var err error

	path := filepath.Join(s.basePath, userName)
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read directory")
	}
	if len(files) == 0 {
		return nil, storage.ErrPageNotFound
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Intn(len(files))
	filePath := filepath.Join(path, files[index].Name())
	return s.decodePage(filePath)
}

func (s Storage) Save(p *storage.Page) error {
	filePath := filepath.Join(s.basePath, p.UserName)
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to create directory")
	}
	fName, err := fileName(p)
	if err != nil {
		return errors.Wrap(err, "failed to get file name")
	}
	filePath = filepath.Join(filePath, fName)
	if err := os.WriteFile(filePath, []byte(p.URL), os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to write file")
	}
	file, err := os.Create(filePath)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer func() {
		_ = file.Close()
	}()

	if err := gob.NewEncoder(file).Encode(p); err != nil {
		return errors.Wrap(err, "failed to encode")
	}
	return nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}
	defer func() {
		_ = file.Close()
	}()

	var page storage.Page
	if err := gob.NewDecoder(file).Decode(&page); err != nil {
		return nil, errors.Wrap(err, "failed to decode")
	}
	return &page, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
