package main

import (
	"log"

	"mobileordering/internal/config"
	"mobileordering/internal/database"
	"mobileordering/internal/router"
)

func main() {

	// 1️⃣ 初始化配置
	config.Init()

	// 2️⃣ 初始化数据库
	database.InitDB(&config.Cfg.DB)

	// 3️⃣ 初始化路由
	r := router.NewRouter()

	addr := config.Cfg.Server.Port

	log.Printf("Starting server on %s...", addr)
	log.Printf("AppID configured: %s...", maskString(config.Cfg.Wx.AppID))

	// 4️⃣ 启动 Gin 服务
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func maskString(s string) string {
	if len(s) <= 4 {
		return "***"
	}
	return s[:2] + "***" + s[len(s)-2:]
}