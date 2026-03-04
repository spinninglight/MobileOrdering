package model

import (
"gorm.io/gorm"
"time"
)

// Category 商品分类表
type Category struct {
	ID         int64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	MerchantID int64          `gorm:"column:merchant_id;not null;index:idx_category_merchant_visible;index:idx_category_sort" json:"merchant_id"`
	Name       string         `gorm:"column:name;size:50;not null" json:"name"`
	SortOrder  int            `gorm:"column:sort_order;default:0" json:"sort_order"` // 数值越大越靠前
	IsVisible  int8           `gorm:"column:is_visible;default:1" json:"is_visible"`
	Products   []Product      `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty"`
}

// TableName 设置表名为 categories
func (Category) TableName() string {
	return "t_categories"
}

// Product 商品基础信息表
type Product struct {
	ID           int64               `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CategoryID   int64               `gorm:"column:category_id;not null;index:idx_product_category_status" json:"category_id"`
	MerchantID   int64               `gorm:"column:merchant_id;not null;index:idx_product_merchant_status;index:idx_product_created" json:"merchant_id"`
	Name         string              `gorm:"column:name;size:100;not null" json:"name"`
	Description  string              `gorm:"column:description" json:"description,omitempty"`
	ImageURL     string              `gorm:"column:image_url;size:255" json:"image_url,omitempty"`
	SalesCount   int                 `gorm:"column:sales_count;default:0" json:"sales_count"`
	LowestPrice  float64             `gorm:"column:lowest_price;type:decimal(10,2);not null" json:"lowest_price"`
	Status       int8                `gorm:"column:status;default:1;index:idx_product_category_status;index:idx_product_merchant_status" json:"status"` // 0=下架, 1=上架
	CreatedAt    time.Time           `gorm:"column:created_at;index:idx_product_created" json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
	DeletedAt    gorm.DeletedAt      `json:"deleted_at,omitempty"`
	SKUs         []ProductSKU        `gorm:"foreignKey:ProductID" json:"skus,omitempty"`
	Attributes   []ProductAttribute  `gorm:"foreignKey:ProductID" json:"attributes,omitempty"`
}

// TableName 设置表名为 products
func (Product) TableName() string {
	return "t_products"
}

// ProductSKU 商品规格表 (Stock Keeping Unit)
type ProductSKU struct {
	ID        int64   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProductID int64   `gorm:"column:product_id;not null;index:idx_sku_product;index:idx_sku_product_stock" json:"product_id"`
	SpecName  string  `gorm:"column:spec_name;size:50;not null" json:"spec_name"`  // 如: 容量
	SpecValue string  `gorm:"column:spec_value;size:50;not null" json:"spec_value"` // 如: 大杯
	Price     float64 `gorm:"column:price;type:decimal(10,2);not null" json:"price"`
	Stock     int     `gorm:"column:stock;default:-1;index:idx_sku_product_stock" json:"stock"` // -1 表示无限
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

// TableName 设置表名为 product_skus
func (ProductSKU) TableName() string {
return "t_product_skus"
}

// ProductAttribute 商品属性表 (如甜度、冰度)
type ProductAttribute struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProductID  int64  `gorm:"column:product_id;not null;index:idx_attr_product" json:"product_id"`
	AttrName   string `gorm:"column:attr_name;size:50;not null" json:"attr_name"`     // 如: 甜度
	AttrValues string `gorm:"column:attr_values;size:255;not null" json:"attr_values"` // 用逗号隔开: 正常,半糖,无糖
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty"`
}

// TableName 设置表名为 product_attributes
func (ProductAttribute) TableName() string {
	return "t_product_attributes"
}