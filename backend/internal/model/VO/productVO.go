package VO

import (
	"mobileordering/internal/model"
)

// ProductDetailVO 用于前端展示的视图对象
type ProductAttributeVO struct {
	ID         int64    `json:"id"`
	AttrName   string   `json:"attr_name"`
	AttrValues []string `json:"attr_values"` // 数组格式
}

type ProductDetailVO struct {
	model.Product
	Attributes []ProductAttributeVO `json:"attributes"` // 覆盖原有的 Attributes
}
