package service

import (
	"context"
	"fmt"
	"mobileordering/internal/model"
	"mobileordering/internal/repository"
	"time"
)

type OrderService struct {
	orderRepo repository.OrderRepository
	prodcutRepo *repository.ProductRepository
}

func NewOrderService(order repository.OrderRepository, product *repository.ProductRepository) *OrderService {
	return &OrderService{orderRepo:order,prodcutRepo: product}
}

// CreateOrderRequest 下单请求参数
type CreateOrderRequest struct {
	MerchantID int64 `json:"merchant_id" binding:"required"`
	Remark     string `json:"remark"`
	Items      []struct {
		SkuID    int64 `json:"sku_id" binding:"required"`
		Quantity int    `json:"quantity" binding:"required,gt=0"`
		AttrInfo string `json:"attr_info"` // 例如: "半糖,加冰"
	} `json:"items" binding:"required"`
}

func (s *OrderService) PlaceOrder(ctx context.Context, userID int64, req CreateOrderRequest) (*model.Order, error) {
	var skuIDs []int64
	for _, item := range req.Items {
		skuIDs = append(skuIDs, item.SkuID)
	}

	// 1. 获取SKU详情并校验
	skus, err := s.prodcutRepo.GetSkusByIDs(skuIDs)
	if err != nil || len(skus) != len(req.Items) {
		return nil, fmt.Errorf("部分商品已下架或不存在")
	}

	skuMap := make(map[int64]model.ProductSKU)
	for _, sku := range skus {
		skuMap[sku.ID] = sku
	}

	// 2. 计算金额 & 构建订单项
	var totalAmount float64
	var orderItems []*model.OrderItem

	for _, reqItem := range req.Items {
		sku := skuMap[reqItem.SkuID]

		itemTotal := sku.Price * float64(reqItem.Quantity)
		totalAmount += itemTotal

		productName, err := s.prodcutRepo.GetProductNameByID(sku.ProductID)
		if err != nil{
			return nil, fmt.Errorf("访问数据库出错，%s",err)
		}

		orderItems = append(orderItems, &model.OrderItem{
			ShopID:      req.MerchantID,
			ProductID:   sku.ProductID,
			ProductName: productName, // 快照
			SkuID:       sku.ID,
			SkuInfo:     fmt.Sprintf("%s:%s %s", sku.SpecName, sku.SpecValue, reqItem.AttrInfo),
			Price:       sku.Price,
			Quantity:    reqItem.Quantity,
			TotalPrice:  itemTotal,
		})
	}

	// 3. 构建主订单
	order := &model.Order{
		OrderSn:     fmt.Sprintf("SN%d%d", time.Now().Unix(), userID), // 简单订单号生成
		ShopID:      req.MerchantID,
		UserID:      userID,
		TotalAmount: totalAmount,
		PayAmount:   totalAmount, // 实际开发中此处需减去优惠券金额
		Status:      0,           // 待支付
		Remark:      req.Remark,
	}

	// 4. 调用持久化层
	if err := s.orderRepo.CreateOrder(ctx, order, orderItems); err != nil {
		return nil, err
	}

	return order, nil
}

// GetPendingOrders 获取商家的待接单列表 (状态: 1-已支付)
func (s *OrderService) GetPendingOrders(ctx context.Context, shopID int64, page, pageSize int) ([]*model.Order, int64, error) {
    // 状态 1 代表用户已支付且通过回调，正等待商家确认接单
    var pendingStatus int8 = 1
    return s.orderRepo.GetMerchantOrders(ctx, uint64(shopID), pendingStatus, page, pageSize)
}

// AcceptOrder 商家接单
func (s *OrderService) AcceptOrder(ctx context.Context, shopID int64, orderID string) error {
    // 1. 获取订单详情并校验归属权
    order, _, err := s.orderRepo.GetOrderWithItems(ctx, orderID)
    if err != nil {
        return fmt.Errorf("订单不存在")
    }
    
    // 安全检查：防止 A 商家操作 B 商家的订单
    if int64(order.ShopID) != shopID {
        return fmt.Errorf("无权处理该订单")
    }

    // 2. 状态流转校验：只有已支付(1)的订单才能变为制作中(2)
    // 使用 Repository 层的 CAS 更新，防止多终端并发操作
    success, err := s.orderRepo.UpdateOrderStatus(ctx, order.ID, 1, 2)
    if err != nil {
        return fmt.Errorf("更新订单状态失败: %w", err)
    }
    if !success {
        return fmt.Errorf("接单失败：订单已被处理或状态已变更")
    }

    // 3. 异步操作（可选）
    // go s.triggerKitchenPrint(order) // 通知厨房打印机
    
    return nil
}

// RejectOrder 商家拒单（包含状态变更与退款逻辑）
func (s *OrderService) RejectOrder(ctx context.Context, shopID int64, orderID string, reason string) error {
    // 1. 获取并校验订单
    order, _, err := s.orderRepo.GetOrderWithItems(ctx, orderID)
    if err != nil {
        return fmt.Errorf("订单不存在")
    }
    if int64(order.ShopID) != shopID {
        return fmt.Errorf("无权处理该订单")
    }

    // 2. 状态更新：从 已支付(1) 更新为 已取消(4)
    success, err := s.orderRepo.UpdateOrderStatus(ctx, order.ID, 1, 4)
    if err != nil {
        return fmt.Errorf("取消订单失败: %w", err)
    }
    if !success {
        return fmt.Errorf("拒单失败：订单当前状态无法取消")
    }

    // 3. 执行退款流程
    // 在实际生产中，退款可能是异步的，或者需要调用专门的支付服务
    err = s.doRefund(ctx, order, reason)
    if err != nil {
        // 特别注意：状态已经改为已取消，但退款失败了，这里需要记录异常日志并人工介入
        fmt.Printf("[CRITICAL] 订单 %s 取消成功但退款失败: %v\n", order.OrderSn, err)
        return fmt.Errorf("订单已取消，但系统退款异常，请联系客服处理")
    }

    return nil
}

// doRefund 内部退款逻辑封装
func (s *OrderService) doRefund(ctx context.Context, order *model.Order, reason string) error {
    // 逻辑：判断 PayType (1-微信, 2-余额)
    // 调用对应的支付 SDK 执行 Refund 操作
    return nil 
}