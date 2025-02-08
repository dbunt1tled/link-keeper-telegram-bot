package database

import (
	"sync"
	"tBot/internal/env"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	m  sync.Once //nolint:gochecknoglobals // singleton
	db *gorm.DB  //nolint:gochecknoglobals // singleton
)

func InitDBConnection(cfg *env.Config) *gorm.DB {
	m.Do(func() {
		var err error
		db, err = gorm.Open(sqlite.Open(cfg.Storage.DB), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
	})
	return DB()
}

func DB() *gorm.DB {
	if db == nil {
		panic("Connection is not initialized. Call InitDBConnection first.")
	}
	return db
}
