package main

import (
	"log"
	"net/http"

	"mobileordering/internal/config"
	"mobileordering/internal/database"
)

func main() {
	// 1. 初始化配置
	config.Init()

	// 2. 连接数据库
	database.InitDB(&config.Cfg.DB)

	// 3. 启动服务
	addr := config.Cfg.Server.Port
	log.Printf("Starting server on %s...", addr)
	log.Printf("AppID configured: %s...", maskString(config.Cfg.Wx.AppID))

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
