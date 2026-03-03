package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
    Server ServerConfig   `envconfig:"SERVER"`
    DB     DatabaseConfig `envconfig:"DB"`
    Wx     WechatConfig   `envconfig:"WX"`
    Jwt    JwtConfig      `envconfig:"JWT"`
}

type WechatConfig struct {
    // 最终匹配: APP_WX_APP_ID
    AppID     string `envconfig:"APP_ID" required:"true"`
    // 最终匹配: APP_WX_APP_SECRET
    AppSecret string `envconfig:"APP_SECRET" required:"true"`  
}

type JwtConfig struct {
    // 最终匹配: APP_JWT_SECRET
    Secret    string `envconfig:"SECRET" required:"true"`
}

type ServerConfig struct {
    // 最终匹配: APP_SERVER_PORT
    Port      string `envconfig:"PORT" default:":8080"`
}

type DatabaseConfig struct {
    // 最终环境变量名: APP_DB_HOST
    Host     string `envconfig:"HOST" default:"127.0.0.1"`
    // 最终环境变量名: APP_DB_PORT
    Port     string `envconfig:"PORT" default:"3306"`
    // 最终环境变量名: APP_DB_USER
    User     string `envconfig:"USER" required:"true"`
    // 最终环境变量名: APP_DB_PASSWORD
    Password string `envconfig:"PASSWORD" required:"true"`
    // 最终环境变量名: APP_DB_NAME
    Name     string `envconfig:"NAME" required:"true"`
}

var Cfg *AppConfig

func Init() {
	// 1. 尝试加载 .env 文件（本地开发利器）
	// 不直接检查 error，因为在生产环境（如 Docker）通常不需要 .env 文件
	_ = godotenv.Load()

	// 2. 实例化并解析环境变量
	var c AppConfig
	// 注意：如果你的变量名没有统一前缀（如 APP_），第一个参数传空字符串 ""
	if err := envconfig.Process("APP", &c); err != nil {
		log.Fatalf("❌ 无法加载配置: %v", err)
	}

	Cfg = &c
	log.Println("✅ 配置加载成功")
}