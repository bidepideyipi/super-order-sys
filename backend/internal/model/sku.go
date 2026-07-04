package model

import (
	customtime "super-order-web/pkg/time"
)

// SKU 商品
type SKU struct {
	ID          int64        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	SKUCode     string       `gorm:"column:sku_code;uniqueIndex;not null;type:TEXT" json:"sku_code"`
	Name        string       `gorm:"column:name;not null;type:TEXT" json:"name"`
	Description string       `gorm:"column:description;type:TEXT" json:"description"`
	Spec        string       `gorm:"column:spec;type:TEXT" json:"spec"`
	Unit        string       `gorm:"column:unit;type:TEXT;default:'个'" json:"unit"`
	CategoryID  string       `gorm:"column:category_id;not null;type:TEXT" json:"category_id"`
	BoxSpec     string       `gorm:"column:box_spec;type:TEXT" json:"box_spec"`
	BoxQuantity int          `gorm:"column:box_quantity;default:1" json:"box_quantity"`
	CostPrice   float64      `gorm:"column:cost_price;default:0" json:"cost_price"`
	SalePrice   float64      `gorm:"column:sale_price;default:0" json:"sale_price"`
	IsDeleted   bool         `gorm:"column:is_deleted;default:0" json:"is_deleted"`
	CreatedAt   customtime.CustomTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   customtime.CustomTime `gorm:"column:updated_at" json:"updated_at"`
	Category    *SKUCategory `gorm:"foreignKey:CategoryID;references:CategoryID" json:"category,omitempty"`
}

// TableName 指定表名
func (SKU) TableName() string {
	return "sku"
}
