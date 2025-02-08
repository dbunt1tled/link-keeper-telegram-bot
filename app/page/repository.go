package page

import (
	"tBot/internal/database"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
}

type DTO struct {
	Status   *database.PageStatus
	UserName *string
	URL      *string
	Limit    *int
	Offset   *int
}

func (p Repository) Save(url string, userName string, status database.PageStatus) (*Pages, error) {
	page := &Pages{
		URL:       url,
		UserName:  userName,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := database.DB().Create(page).Error; err != nil {
		return nil, err
	}
	return page, nil
}

func (p Repository) Update(id int, updateDTO DTO) (*Pages, error) {
	var page Pages
	if err := database.DB().Where("id = ?", id).First(&page).Error; err != nil {
		return nil, err
	}
	if updateDTO.Status != nil {
		page.Status = *updateDTO.Status
	}
	if updateDTO.UserName != nil {
		page.UserName = *updateDTO.UserName
	}
	if updateDTO.URL != nil {
		page.URL = *updateDTO.URL
	}
	if err := database.DB().Save(&page).Error; err != nil {
		return nil, err
	}
	return &page, nil
}

func (p Repository) Delete(id int) (*Pages, error) {
	var page Pages
	if err := database.DB().Where("id = ?", id).First(&page).Error; err != nil {
		return nil, err
	}
	page.Status = database.PageDeleted
	if err := database.DB().Save(&page).Error; err != nil {
		return nil, err
	}
	return &page, nil
}

func (p Repository) List(filter DTO) ([]Pages, error) {
	var pages []Pages
	builder := buildFilter(filter)
	if err := builder.Find(&pages).Error; err != nil {
		return nil, err
	}
	return pages, nil
}
func (p Repository) Count(filter DTO) (*int64, error) {
	var count int64
	builder := buildFilter(filter)
	if err := builder.Count(&count).Error; err != nil {
		return nil, err
	}
	return &count, nil
}

func (p Repository) First(filter DTO) (*Pages, error) {
	var page Pages
	builder := buildFilter(filter)
	if err := builder.First(&page).Error; err != nil {
		return nil, err
	}

	return &page, nil
}

func buildFilter(filter DTO) *gorm.DB {
	builder := database.DB().Model(&Pages{})
	if filter.Status != nil {
		builder = builder.Where("status = ?", *filter.Status)
	}
	if filter.UserName != nil {
		builder = builder.Where("user_name = ?", *filter.UserName)
	}
	if filter.URL != nil {
		builder = builder.Where("url = ?", *filter.URL)
	}
	if filter.Limit != nil {
		builder = builder.Limit(*filter.Limit)
	}
	if filter.Offset != nil {
		builder = builder.Offset(*filter.Offset)
	}

	return builder
}
