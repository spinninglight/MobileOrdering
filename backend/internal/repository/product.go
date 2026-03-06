// internal/repository/product_repository.go
package repository

import (
	"mobileordering/internal/model"

	"gorm.io/gorm"
	"errors"
	"fmt"
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

// internal/repository/product_repository.go

func (r *ProductRepository) GetSkusByIDs(skuIDs []int64) ([]model.ProductSKU, error) {
    var skus []model.ProductSKU
    
    // 直接查询，不需要 Transaction 包装
    err := r.DB.
        Where("id IN ? AND stock != 0", skuIDs). // 过滤掉库存明确为0的
        Find(&skus).Error
        
    return skus, err
}

// GetProductNameByID 根据商品ID获取商品名称（仅查询 name 字段，高效）
func (r *ProductRepository) GetProductNameByID(productID int64) (string, error) {
	var name string

	err := r.DB.
		Model(&model.Product{}).
		Select("name").
		Where("id = ? AND status = 1", productID).
		Scan(&name).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("product not found or not on shelf")
		}
		return "", err
	}

	return name, nil
}