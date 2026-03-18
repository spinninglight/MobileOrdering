package handler

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"mobileordering/internal/service"
	"mobileordering/internal/model/response"
)

type OrderHandler struct {
	svc *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: service}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 假设从中间件获取的当前登录用户ID
	userID := c.GetInt64("user_id")

	order, err := h.svc.PlaceOrder(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.Success(c,order)
}


// GetPendingOrders 获取商家待接单列表
// GET /merchant/orders/pending?page=1&page_size=10
func (h *OrderHandler) GetPendingOrders(c *gin.Context) {
    // 1. 从中间件获取当前登录商家的 shop_id
    shopID := c.GetInt64("shop_id")
    if shopID == 0 {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到商户信息"})
        return
    }

    // 2. 获取分页参数
    page := c.DefaultQuery("page", "1")
    pageSize := c.DefaultQuery("page_size", "10")
    
    // 转换为 int (此处省略具体转换错误处理，实际开发中建议封装工具函数)
    p := 1
    ps := 10
    fmt.Sscanf(page, "%d", &p)
    fmt.Sscanf(pageSize, "%d", &ps)

    // 3. 调用 Service
    orders, total, err := h.svc.GetPendingOrders(c.Request.Context(), shopID, p, ps)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "获取列表失败: " + err.Error()})
        return
    }

    // 4. 返回分页数据
    response.Success(c, gin.H{
        "list":  orders,
        "total": total,
    })
}

// AcceptOrder 商家接单
// POST /merchant/orders/accept
func (h *OrderHandler) AcceptOrder(c *gin.Context) {
    var req struct {
        OrderID int64 `json:"order_id" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
        return
    }

    shopID := c.GetInt64("shop_id")

    err := h.svc.AcceptOrder(c.Request.Context(), shopID, req.OrderID)
    if err != nil {
        // 这里的错误可能是“订单已被接单”，可以根据 err 类型细化状态码
        c.JSON(http.StatusOK, gin.H{"code": 400, "error": err.Error()})
        return
    }

    response.Success(c, "接单成功")
}

// RejectOrder 商家拒单
// POST /merchant/orders/reject
func (h *OrderHandler) RejectOrder(c *gin.Context) {
    var req struct {
        OrderID int64  `json:"order_id" binding:"required"`
        Reason  string `json:"reason" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
        return
    }

    shopID := c.GetInt64("shop_id")

    err := h.svc.RejectOrder(c.Request.Context(), shopID, req.OrderID, req.Reason)
    if err != nil {
        c.JSON(http.StatusOK, gin.H{"code": 400, "error": err.Error()})
        return
    }

    response.Success(c, "拒单成功，已发起退款")
}