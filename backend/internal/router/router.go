package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"mobileordering/internal/handler"
	"mobileordering/internal/service"
)

func NewRouter() *gin.Engine {

	r := gin.Default()

	r.Use(cors.Default())

	menuservice := service.NewMenuService()

	loginHandler := handler.NewLoginHandler()
	menuHandler := handler.NewMenuHandler(menuservice)

	api := r.Group("/api")

	v1 := api.Group("v1")
	{
		v1.POST("/login", loginHandler.Login)
	}

	merchant := v1.Group("merchant")
	{
		merchant.GET("/merchant/:id/menu", menuHandler.GetMerchantMenu)
	}

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	return r
}