package main

import (
	"log"
	"net/http"

	"mobileordering/internal/config"
	"mobileordering/internal/handler"
)

func main() {
	// 1. 初始化配置
	config.Init()

	// 2. 注册路由
	loginHandler := handler.NewLoginHandler()

	http.HandleFunc("/api/login", loginHandler.ServeHTTP)

	// 健康检查接口
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// 3. 启动服务
	addr := config.Cfg.ServerPort
	log.Printf("Starting server on %s...", addr)
	log.Printf("AppID configured: %s...", maskString(config.Cfg.WxAppID))

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func maskString(s string) string {
	if len(s) <= 4 {
		return "***"
	}
	return s[:2] + "***" + s[len(s)-2:]
}
