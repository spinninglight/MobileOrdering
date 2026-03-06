package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"mobileordering/internal/handler"
	"mobileordering/internal/service"
	"mobileordering/internal/repository"
	"mobileordering/internal/database"
)

func NewRouter() *gin.Engine {

	r := gin.Default()

	r.Use(cors.Default())

	menuRepository := repository.NewMenuRepository(database.DB)
	productRepository := repository.NewProductRepository(database.DB)

	menuService := service.NewMenuService(menuRepository)
	productService := service.NewProductService(productRepository)

	loginHandler := handler.NewLoginHandler()
	menuHandler := handler.NewMenuHandler(menuService)
	productHandler := handler.NewProductHandler(productService)

	orderRepo := repository.NewOrderRepository(database.DB)
	orderSvc := service.NewOrderService(orderRepo,productRepository)
	orderHandler := handler.NewOrderHandler(orderSvc)

	api := r.Group("/api")

	v1 := api.Group("v1")
	{
		v1.POST("/login", loginHandler.Login)
	}

	merchant := v1.Group("merchant")
	{
		merchant.GET("/:id/menu", menuHandler.GetMerchantMenu)
	}

	products := v1.Group("products")
	{
		products.GET("/:product_id/details", productHandler.GetProductDetails)
	}

	order := v1.Group("order")
	{
		order.POST("/create",orderHandler.CreateOrder)
	}

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	return r
}