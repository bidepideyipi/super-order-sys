package model

import "time"

// FinancialTransaction 财务流水
type FinancialTransaction struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Category     string    `gorm:"column:category;not null;type:TEXT" json:"category"`
	Description  string    `gorm:"column:description;type:TEXT" json:"description"`
	AmountChange float64   `gorm:"column:amount_change;not null" json:"amount_change"`
	Balance      float64   `gorm:"column:balance;not null" json:"balance"`
	IsSettled    bool      `gorm:"column:is_settled;default:0" json:"is_settled"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

// TableName 指定表名
func (FinancialTransaction) TableName() string {
	return "financial_transaction"
}
