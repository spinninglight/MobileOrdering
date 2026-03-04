// internal/handler/menu_handler.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"mobileordering/internal/service"
)

type MenuHandler struct {
	Service *service.MenuService
}

func NewMenuHandler(service *service.MenuService) *MenuHandler {
	return &MenuHandler{Service: service}
}

func (h *MenuHandler) GetMerchantMenu(c *gin.Context) {

	merchantIDStr := c.Param("id")
	merchantID, err := strconv.ParseInt(merchantIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid merchant id"})
		return
	}

	data, err := h.Service.GetMerchantMenu(merchantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"merchant_id": merchantID,
		"categories":  data,
	})
}