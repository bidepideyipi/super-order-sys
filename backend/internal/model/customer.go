package model

import "time"

// Customer 客户
type Customer struct {
	CustomerID   string    `gorm:"column:customer_id;primaryKey;type:TEXT" json:"customer_id"`
	CustomerName string    `gorm:"column:customer_name;not null;type:TEXT" json:"customer_name"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
	Orders       []Order   `gorm:"foreignKey:CustomerID" json:"orders,omitempty"`
}

// TableName 指定表名
func (Customer) TableName() string {
	return "customer"
}
