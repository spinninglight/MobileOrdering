package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"mobileordering/internal/handler"
	"mobileordering/internal/service"
	"mobileordering/internal/repository"
	"mobileordering/internal/database"
	"mobileordering/internal/middleware"
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
	///api/v1/merchant/
	merchant := v1.Group("merchant")
	{
		merchant.GET("/:id/menu", menuHandler.GetMerchantMenu)
		merchant.POST("/login", loginHandler.MerchantLogin)
	}

	// 私有接口：需要验证 Token
    private := merchant.Group("")
    private.Use(middleware.MerchantAuth()) // 挂载我们写的中间件
    {
        // 当请求到达这里时，中间件已经把 shop_id 塞进 context 了！
        private.GET("/orders/pending", orderHandler.GetPendingOrders)
        private.POST("/orders/accept", orderHandler.AcceptOrder)
        private.POST("/orders/reject", orderHandler.RejectOrder)
    }

	products := v1.Group("products")
	{
		products.GET("/:product_id/details", productHandler.GetProductDetails)
	}

	order := v1.Group("order")
	{
		order.POST("/create",orderHandler.CreateOrder)
		order.POST("/pending",orderHandler.GetPendingOrders)
		order.POST("/accept",orderHandler.AcceptOrder)
		order.POST("/reject",orderHandler.RejectOrder)
	}

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	return r
}