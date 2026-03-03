package database

import (
	"mobileordering/internal/config"
	"testing"
)

func TestInitDB(t *testing.T) {
	// 1. 手动构造一个测试用的配置（或者调用 config.Init() 加载 .env）
	// 注意：确保这里的账号密码在你本地 MySQL 是真实存在的
	testCfg := &config.DatabaseConfig{
		Host:     "127.0.0.1",
		Port:     "3306",
		User:     "root",
		Password: "123456", // 记得替换
		Name:     "mysql",    // 先连系统自带的 mysql 库测试连通性
	}

	// 2. 执行初始化
	// 使用 defer 来处理可能的 panic，让测试优雅退出
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("InitDB 发生 Panic: %v", r)
		}
	}()

	InitDB(testCfg)

	// 3. 验证全局变量 DB 是否被赋值
	if DB == nil {
		t.Fatal("数据库连接成功但 DB 实例为 nil")
	}

	// 4. 尝试执行一个简单的查询
	sqlDB, err := DB.DB()
	if err != nil {
		t.Fatalf("无法获取底层的 sql.DB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("数据库 Ping 失败: %v", err)
	}

	t.Log("✅ 数据库集成测试通过！")
}