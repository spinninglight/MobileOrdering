// internal/service/menu_service.go
package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"mobileordering/internal/database"
	"mobileordering/internal/model"
	"mobileordering/internal/repository"

	"github.com/fatih/color"
)

type MenuService struct {
	Repo *repository.MenuRepository
}

func NewMenuService(repo *repository.MenuRepository) *MenuService {
	return &MenuService{Repo: repo}
}

func (s *MenuService) GetMerchantMenu(merchantID int64) ([]model.Category, error) {

	// 未来可以加：
	// - Redis缓存
	// - 数据加工
	// - 权限校验
	// - 商户状态校验

	ctx := context.Background()
    // 1. 定义 Key
    cacheKey := fmt.Sprintf("merchant:menu:%d", merchantID)
    var categories []model.Category

    // 2. 尝试从 Redis 获取
    found, err := database.GetJSON(ctx, cacheKey, &categories)
    if err != nil {
        log.Printf("Redis error: %v", err) // 缓存报错通常不影响业务，记录日志即可
    }
    if found {
		//log.Printf("[Redis Hit] MerchantID: %d, Key: %s", merchantID, cacheKey)
		color.Green("[Redis Hit] MerchantID: %d, Key: %s", merchantID, cacheKey)
        return categories, nil
    }

	// ❌ 未命中，打印日志提示正在查询数据库
    //log.Printf("[Redis Miss] Fetching from DB for MerchantID: %d", merchantID)
	color.Yellow("[Redis Miss] Fetching from DB for MerchantID: %d", merchantID)

    // 3. 缓存穿透，查询数据库
    categories, err = s.Repo.GetMerchantMenu(merchantID)
    if err != nil {
        return nil, err
    }

    // 4. 异步或同步回写 Redis，设置过期时间（比如 10 分钟）
    _ = database.SetJSON(ctx, cacheKey, categories, 10*time.Minute)

    return categories, nil
}