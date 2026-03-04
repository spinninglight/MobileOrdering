// internal/repository/menu_repository.go
package repository

import (
	"mobileordering/internal/model"

	"gorm.io/gorm"
)

type MenuRepository struct {
	DB *gorm.DB
}

func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{DB: db}
}

func (r *MenuRepository) GetMerchantMenu(merchantID int64) ([]model.Category, error) {
	var categories []model.Category

	err := r.DB.
		Where("merchant_id = ? AND is_visible = 1", merchantID).
		Order("sort_order DESC").
		Preload("Products", func(db *gorm.DB) *gorm.DB {
			return db.
				Select("id", "category_id", "merchant_id", "name", "image_url", "sales_count", "lowest_price").
				Where("status = 1").
				Order("id DESC")
		}).
		Find(&categories).Error

	return categories, err
}