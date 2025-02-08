package page

import (
	"tBot/internal/database"
	"time"
)

type Pages struct {
	ID        int                 `json:"id" gorm:"column:id"`
	URL       string              `json:"url"  gorm:"column:url"`
	UserName  string              `json:"username" gorm:"column:user_name"`
	Status    database.PageStatus `json:"status" gorm:"column:status"`
	CreatedAt time.Time           `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time           `json:"updatedAt" gorm:"column:updated_at"`
}

func (Pages) TableName() string {
	return "pages"
}
