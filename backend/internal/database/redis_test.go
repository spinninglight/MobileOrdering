package database

import (
    "context"
    "mobileordering/internal/config"
    "testing"
    "time"
)

func TestInitRedis(t *testing.T) {
    // 1. 模拟全局配置环境
    config.Cfg = &config.AppConfig{
        Env: "test",
    }

    // 2. 构造测试用的 Redis 配置
    // 确保这里的地址和密码与你本地运行的 Redis 一致
    testRedisCfg := &config.RedisConfig{
        Host:     "127.0.0.1",
        Port:     6379,
        Password: "123456", // 如果没设密码就留空
        DB:       0,  // 默认数据库
    }

    // 3. 执行初始化并处理 Panic
    defer func() {
        if r := recover(); r != nil {
            t.Fatalf("InitRedis 发生 Panic: %v", r)
        }
    }()

    InitRedis(testRedisCfg)

    // 4. 验证全局变量 RedisClient (或你定义的 Redis 变量名) 是否被赋值
    if RedisClient == nil {
        t.Fatal("Redis 连接初始化成功但 RedisClient 实例为 nil")
    }

    // 5. 尝试执行 Ping 操作验证连通性
    // 使用带超时的 Context 防止测试死锁
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := RedisClient.Ping(ctx).Result()
    if err != nil {
        t.Fatalf("Redis Ping 失败: %v", err)
    }

    // 6. (可选) 测试基础读写
    err = RedisClient.Set(ctx, "test_key", "hello_redis", 10*time.Second).Err()
    if err != nil {
        t.Fatalf("Redis Set 失败: %v", err)
    }

    val, err := RedisClient.Get(ctx, "test_key").Result()
    if err != nil || val != "hello_redis" {
        t.Fatalf("Redis Get 验证失败: 得到 %v, 错误: %v", val, err)
    }

    t.Log("✅ Redis 集成测试通过！")
}