// internal/service/menu_service.go
package service

import (
	"mobileordering/internal/model"
	"mobileordering/internal/repository"
)

type MenuService struct {
	Repo *repository.MenuRepository
}

func NewMenuService() *MenuService {
	return &MenuService{}
}

func (s *MenuService) GetMerchantMenu(merchantID int64) ([]model.Category, error) {

	// 未来可以加：
	// - Redis缓存
	// - 数据加工
	// - 权限校验
	// - 商户状态校验

	return s.Repo.GetMerchantMenu(merchantID)
}