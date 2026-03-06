package repository

import (
	"context"
	"mobileordering/internal/model"
	"gorm.io/gorm"
)

type OrderRepository interface {
	// CreateOrder 事务创建订单
	CreateOrder(ctx context.Context, order *model.Order, items []*model.OrderItem) error
}

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepo{db: db}
}

func (r *orderRepo) CreateOrder(ctx context.Context, order *model.Order, items []*model.OrderItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 插入订单主表
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		// 2. 批量插入订单详情表（设置OrderID）
		for _, item := range items {
			item.OrderID = order.ID
		}
		if err := tx.CreateInBatches(items, 100).Error; err != nil {
			return err
		}
		// 3. 扣减库存 (简单示例：扣减非无限库存的SKU)
		for _, item := range items {
			res := tx.Table("t_product_skus").
				Where("id = ? AND stock > 0", item.SkuID).
				UpdateColumn("stock", gorm.Expr("stock - ?", item.Quantity))
			if res.Error != nil {
				return res.Error
			}
		}
		return nil
	})
}