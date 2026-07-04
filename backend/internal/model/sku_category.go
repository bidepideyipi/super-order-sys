package model

import (
	customtime "super-order-web/pkg/time"
)

// SKUCategory SKU分类
type SKUCategory struct {
	CategoryID   string    `gorm:"column:category_id;primaryKey;type:TEXT" json:"category_id"`
	CategoryName string    `gorm:"column:category_name;not null;type:TEXT" json:"category_name"`
	CreatedAt    customtime.CustomTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    customtime.CustomTime `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (SKUCategory) TableName() string {
	return "sku_category"
}
