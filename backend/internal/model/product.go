package model

type Product struct {
	ID          int64   `json:"id"`
	CategoryID  int64   `json:"category_id"`
	MerchantID  int64   `json:"merchant_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	SalesCount  int     `json:"sales_count"`
	LowestPrice float64 `json:"lowest_price"`
	Status      int8    `json:"status"`
}

type ProductSku struct {
	ID        int64   `json:"id"`
	ProductID int64   `json:"product_id"`
	SpecName  string  `json:"spec_name"`
	SpecValue string  `json:"spec_value"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
}

type ProductAttribute struct {
	ID         int64  `json:"id"`
	ProductID  int64  `json:"product_id"`
	AttrName   string `json:"attr_name"`
	AttrValues string `json:"attr_values"`
}

type Category struct {
	ID         int64  `json:"id"`
	MerchantID int64  `json:"merchant_id"`
	Name       string `json:"name"`
	SortOrder  int    `json:"sort_order"`
	IsVisible  int8   `json:"is_visible"`

	Products []Product `gorm:"foreignKey:CategoryID"`
}