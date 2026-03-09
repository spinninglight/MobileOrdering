package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient 声明一个全局变量，方便外部调用（或者通过依赖注入传递）
var RedisClient *redis.Client

// RedisConfig 用于定义初始化配置
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

// InitRedis 初始化 Redis 连接
func InitRedis(cfg *RedisConfig) error {
	// 1. 创建 Redis 选项
	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password, // 如果没有密码则为空
		DB:       cfg.DB,       // 默认数据库，通常是 0
		
		// 连接池配置
		PoolSize:     cfg.PoolSize,     // 最大连接数
		MinIdleConns: 5,                // 最小空闲连接数
		DialTimeout:  5 * time.Second,  // 连接超时
		ReadTimeout:  3 * time.Second,  // 读取超时
		WriteTimeout: 3 * time.Second,  // 写入超时
		PoolTimeout:  4 * time.Second,  // 当池中无连接时的等待超时
	}

	// 2. 实例化客户端
	client := redis.NewClient(opts)

	// 3. 使用 Context 测试连接是否通畅 (Ping)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("无法连接到 Redis: %v", err)
	}

	log.Println("✅ Redis 连接初始化成功")
	RedisClient = client
	return nil
}

// CloseRedis 优雅关闭连接（通常在 main.go 的 defer 中调用）
func CloseRedis() {
	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			log.Printf("❌ 关闭 Redis 连接时出错: %v", err)
		}
	}
}