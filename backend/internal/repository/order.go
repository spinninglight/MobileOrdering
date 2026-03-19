package repository

import (
	"context"
	"mobileordering/internal/model"
	"gorm.io/gorm"
)

type OrderRepository interface {
    // 原有的创建订单
    CreateOrder(ctx context.Context, order *model.Order, items []*model.OrderItem) error
    
    // 1. 获取商家的订单列表（支持按状态过滤，如待接单、制作中等）
    GetMerchantOrders(ctx context.Context, shopID uint64, status int8, page, pageSize int) ([]*model.Order, int64, error)
    
    // 2. 获取订单详情（包含 Items）
    GetOrderWithItems(ctx context.Context, orderID string) (*model.Order, []*model.OrderItem, error)
    
    // 3. 更新订单状态（带前置状态检查，防止重复接单）
    UpdateOrderStatus(ctx context.Context, orderID int64, oldStatus, newStatus int8) (bool, error)
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

// GetMerchantOrders 分页获取商家订单列表
func (r *orderRepo) GetMerchantOrders(ctx context.Context, shopID uint64, status int8, page, pageSize int) ([]*model.Order, int64, error) {
    var orders []*model.Order
    var total int64
    offset := (page - 1) * pageSize

    db := r.db.WithContext(ctx).Model(&model.Order{}).Where("shop_id = ?", shopID)
    
    // status 为 -1 或其他逻辑可以代表查询全部
    if status >= 0 {
        db = db.Where("status = ?", status)
    }

    if err := db.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    if err := db.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&orders).Error; err != nil {
        return nil, 0, err
    }

    return orders, total, nil
}

// UpdateOrderStatus 更新状态（CAS 思想：乐观锁防止并发冲突）
func (r *orderRepo) UpdateOrderStatus(ctx context.Context, orderID int64, oldStatus, newStatus int8) (bool, error) {
    // 只有当状态匹配 oldStatus 时才更新，防止两个店员同时点接单
    res := r.db.WithContext(ctx).Model(&model.Order{}).
        Where("id = ? AND status = ?", orderID, oldStatus).
        Update("status", newStatus)

    if res.Error != nil {
        return false, res.Error
    }
    
    // RowsAffected 为 0 说明状态已经被别人改过了
    return res.RowsAffected > 0, nil
}

// GetOrderWithItems 获取详情（商家接单后需要知道做哪几个菜）
func (r *orderRepo) GetOrderWithItems(ctx context.Context, orderSN string) (*model.Order, []*model.OrderItem, error) {
    var order model.Order
    var items []*model.OrderItem

    // 1. 使用 Where 显式指定根据 order_sn 查询，避免 GORM 误判为主键 ID
    // First 会自动加上 LIMIT 1
    err := r.db.WithContext(ctx).
        Where("order_sn = ?", orderSN). 
        First(&order).Error
    
    if err != nil {
        return nil, nil, err // 可能是 record not found
    }

    // 2. 根据查出来的订单自增 ID 去关联查询 items
    // 注意：这里用 order.ID (int64) 比用字符串 order_sn 性能更高
    err = r.db.WithContext(ctx).
        Where("order_id = ?", order.ID). 
        Find(&items).Error

    return &order, items, err
}