package database

import "database/sql/driver"

type PageStatus int64

const (
	PageDeleted PageStatus = iota
	PageActive
)

func (p *PageStatus) Scan(value interface{}) error {
	*p = PageStatus(value.(int64))
	return nil
}

func (p PageStatus) Value() (driver.Value, error) {
	return int64(p), nil
}
