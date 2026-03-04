package router

import (
	"net/http"
	"mobileordering/internal/handler"
)

// NewRouter 返回配置好的 http.Handler
func NewRouter() http.Handler {
	// 1. 创建 ServeMux (路由器)
	mux := http.NewServeMux()

	// 2. 实例化 Handler
	loginHandler := handler.NewLoginHandler()

	// 3. 注册业务路由
	mux.HandleFunc("/api/login", loginHandler.ServeHTTP)

	// 4. 注册系统路由 (健康检查)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// 注册商品展示路由

	return mux
}