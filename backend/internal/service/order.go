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