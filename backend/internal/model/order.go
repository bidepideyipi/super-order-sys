package model

import "time"

// Order 订单
type Order struct {
	ID              int64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderNo         string         `gorm:"column:order_no;uniqueIndex;not null;type:TEXT" json:"order_no"`
	CustomerID      string         `gorm:"column:customer_id;not null;type:TEXT" json:"customer_id"`
	OrderDate       string         `gorm:"column:order_date;not null;type:TEXT" json:"order_date"`
	Status          string         `gorm:"column:status;type:TEXT;default:'pending'" json:"status"`
	IsSettled       bool           `gorm:"column:is_settled;default:0" json:"is_settled"`
	TotalCostAmount float64        `gorm:"column:total_cost_amount;default:0" json:"total_cost_amount"`
	TotalSaleAmount float64        `gorm:"column:total_sale_amount;default:0" json:"total_sale_amount"`
	Remarks         string         `gorm:"column:remarks;type:TEXT" json:"remarks"`
	CreatedAt       time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at" json:"updated_at"`
	Customer        *Customer      `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Items           []OrderItem   `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "order"
}

// OrderStatus 订单状态
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)
