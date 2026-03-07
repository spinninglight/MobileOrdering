package model

import "time"

// Order 订单主模型
type Order struct {
	ID          int64    `gorm:"primaryKey;column:id" json:"id"`
	OrderSn     string    `gorm:"column:order_sn;type:varchar(64);not null;uniqueIndex" json:"order_sn"`
	ShopID      int64    `gorm:"column:shop_id;index" json:"shop_id"` // 新增：店铺ID
	UserID      int64    `gorm:"column:user_id;index" json:"user_id"`
	TotalAmount float64   `gorm:"column:total_amount" json:"total_amount"`
	PayAmount   float64   `gorm:"column:pay_amount" json:"pay_amount"`
	Status      int8      `gorm:"column:status;index" json:"status"`
	PayType     int8      `gorm:"column:pay_type" json:"pay_type"`
	Remark      string    `gorm:"column:remark" json:"remark"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`

	Items []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
}

func (Order) TableName() string {
	return "t_order"
}

// OrderItem 订单详情模型
type OrderItem struct {
	ID          int64    `gorm:"primaryKey;column:id" json:"id"`
	OrderID     int64    `gorm:"column:order_id;index" json:"order_id"`
	ShopID      int64    `gorm:"column:shop_id;index" json:"shop_id"` // 新增：冗余店铺ID
	ProductID   int64    `gorm:"column:product_id" json:"product_id"`
	ProductName string    `gorm:"column:product_name" json:"product_name"`
	SkuID       int64    `gorm:"column:sku_id" json:"sku_id"`
	SkuInfo     string    `gorm:"column:sku_info" json:"sku_info"`
	Price       float64   `gorm:"column:price" json:"price"`
	Quantity    int       `gorm:"column:quantity" json:"quantity"`
	TotalPrice  float64   `gorm:"column:total_price" json:"total_price"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
}

func (OrderItem) TableName() string {
	return "t_order_item"
}