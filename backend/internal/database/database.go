package database

import (
	"fmt"
	"log"
	"time"

	"mobileordering/internal/config"
	"mobileordering/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(cfg *config.DatabaseConfig) {
	// 拼接 DSN (Data Source Name)
	// 格式: user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	var err error
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("❌ 数据库连接失败: %v", err)
	}

	// 获取底层 sql.DB 以配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("❌ 获取底层 SQL 对象失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接可复用的最大时间

	if config.Cfg.Env == "dev" {
    	DB.AutoMigrate(&model.Category{},&model.Product{},&model.ProductSKU{},&model.ProductAttribute{},&model.Order{},&model.OrderItem{})
		log.Println("✅ 数据库自动匹配model")
	}

	log.Println("✅ 数据库连接成功")
}