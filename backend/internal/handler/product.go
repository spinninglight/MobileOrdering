// internal/handler/product_handler.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mobileordering/internal/service"
)

type ProductHandler struct {
	Service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{Service: service}
}

// GetProductDetails 获取特定商品详情的接口
func (h *ProductHandler) GetProductDetails(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的商品 ID"})
		return
	}

	data, err := h.Service.GetProductDetail(productID)
	if err != nil {
		// 区分商品不存在和其他数据库报错
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在或已下架"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询商品详情失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": data,
	})
}