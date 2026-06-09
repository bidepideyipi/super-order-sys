package service

import (
	"super-order-web/internal/model"

	"gorm.io/gorm"
)

// FinancialTransactionService 财务流水服务
type FinancialTransactionService struct {
	db *gorm.DB
}

// NewFinancialTransactionService 创建财务流水服务
func NewFinancialTransactionService(db *gorm.DB) *FinancialTransactionService {
	return &FinancialTransactionService{db: db}
}

// List 获取财务流水列表
func (s *FinancialTransactionService) List(page, pageSize int, category string) ([]model.FinancialTransaction, int64, error) {
	var transactions []model.FinancialTransaction
	var total int64

	query := s.db.Model(&model.FinancialTransaction{})

	if category != "" && category != "all" {
		query = query.Where("category = ?", category)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&transactions).Error

	return transactions, total, err
}

// Create 创建财务流水
func (s *FinancialTransactionService) Create(transaction *model.FinancialTransaction) error {
	return s.db.Create(transaction).Error
}

// GetLatestBalance 获取最新余额
func (s *FinancialTransactionService) GetLatestBalance() (float64, error) {
	var transaction model.FinancialTransaction
	err := s.db.Order("created_at DESC").First(&transaction).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}
	return transaction.Balance, nil
}
