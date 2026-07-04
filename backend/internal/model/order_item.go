package model

import (
	customtime "super-order-web/pkg/time"
)

// OrderItem 订单明细
type OrderItem struct {
	ID               int64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderID          int64    `gorm:"column:order_id;not null;index" json:"order_id"`
	SKUID            int64    `gorm:"column:sku_id;index" json:"sku_id"`
	SKUCode          string   `gorm:"column:sku_code;not null;type:TEXT" json:"sku_code"`
	ProductName      string   `gorm:"column:product_name;not null;type:TEXT" json:"product_name"`
	Quantity         int      `gorm:"column:quantity;not null" json:"quantity"`
	CostPrice        float64  `gorm:"column:cost_price;not null" json:"cost_price"`
	SalePrice        float64  `gorm:"column:sale_price;not null" json:"sale_price"`
	TotalCostAmount  float64  `gorm:"column:total_cost_amount;not null" json:"total_cost_amount"`
	TotalSaleAmount  float64  `gorm:"column:total_sale_amount;not null" json:"total_sale_amount"`
	SettledAmount    float64  `gorm:"column:settled_amount;default:0" json:"settled_amount"`
	CreatedAt        customtime.CustomTime `gorm:"column:created_at" json:"created_at"`
	Order            *Order   `gorm:"foreignKey:OrderID" json:"-"`
	SKU              *SKU     `gorm:"foreignKey:SKUID" json:"sku,omitempty"`
}

// TableName 指定表名
func (OrderItem) TableName() string {
	return "order_item"
}
