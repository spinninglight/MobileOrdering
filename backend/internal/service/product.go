package service

import (
	"mobileordering/internal/model/VO"
	"mobileordering/internal/repository"
	"strings"
)



type ProductService struct {
	Repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{Repo:repo}
}

func (s *ProductService) GetProductDetail(productID int64) (*VO.ProductDetailVO, error) {
	product, err := s.Repo.GetProductDetails(productID)
	if err != nil {
		return nil, err
	}

	// 核心逻辑：转换属性格式
	attrVOs := make([]VO.ProductAttributeVO, 0)
	for _, attr := range product.Attributes {
		attrVOs = append(attrVOs, VO.ProductAttributeVO{
			ID:         attr.ID,
			AttrName:   attr.AttrName,
			// 通过逗号分割字符串转为数组
			AttrValues: strings.Split(attr.AttrValues, ","),
		})
	}

	// 组装返回对象
	res := &VO.ProductDetailVO{
		Product:    *product,
		Attributes: attrVOs,
	}

	return res, nil
}