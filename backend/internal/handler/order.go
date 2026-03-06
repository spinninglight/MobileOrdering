package handler

import (
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