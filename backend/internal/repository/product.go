// internal/repository/product_repository.go
package repository

import (
	"mobileordering/internal/model"

	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

// GetProductDetails 获取商品及其关联的 SKU 和属性信息
func (r *ProductRepository) GetProductDetails(productID int64) (*model.Product, error) {
	var product model.Product

	// 预加载 SKUs 和 Attributes，且只查询处于上架状态(status=1)的商品
	err := r.DB.
		Preload("SKUs").
		Preload("Attributes").
		Where("id = ? AND status = 1", productID).
		First(&product).Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}