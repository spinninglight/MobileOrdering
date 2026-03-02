package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	WxAppID     string
	WxAppSecret string
	JWTSecret   string
	ServerPort  string
}

var Cfg *Config

func Init() {
	// 加载 .env 文件 (如果存在)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	Cfg = &Config{
		WxAppID:     os.Getenv("WX_APP_ID"),
		WxAppSecret: os.Getenv("WX_APP_SECRET"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		ServerPort:  os.Getenv("SERVER_PORT"),
	}

	if Cfg.ServerPort == "" {
		Cfg.ServerPort = ":8080"
	}

	if Cfg.WxAppID == "" || Cfg.WxAppSecret == "" {
		log.Fatal("FATAL: WX_APP_ID and WX_APP_SECRET must be set in environment variables")
	}
}