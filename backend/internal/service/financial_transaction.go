package service

import (
	"time"
	"super-order-web/internal/model"
	customtime "super-order-web/pkg/time"

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

// ListAll 获取所有财务流水
func (s *FinancialTransactionService) ListAll() ([]model.FinancialTransaction, error) {
	var transactions []model.FinancialTransaction
	err := s.db.Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}

// Create 创建财务流水
func (s *FinancialTransactionService) Create(transaction *model.FinancialTransaction) error {
	transaction.CreatedAt = customtime.NewCustomTime(time.Now())
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

// GetByID 根据ID获取财务流水
func (s *FinancialTransactionService) GetByID(id int64) (*model.FinancialTransaction, error) {
	var transaction model.FinancialTransaction
	err := s.db.Where("id = ?", id).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// Update 更新财务流水
func (s *FinancialTransactionService) Update(transaction *model.FinancialTransaction) error {
	return s.db.Model(&model.FinancialTransaction{}).
		Where("id = ?", transaction.ID).
		Updates(map[string]interface{}{
			"category":      transaction.Category,
			"description":   transaction.Description,
			"amount_change": transaction.AmountChange,
			"balance":       transaction.Balance,
			"is_settled":    transaction.IsSettled,
		}).Error
}

// Delete 删除财务流水
func (s *FinancialTransactionService) Delete(id int64) error {
	return s.db.Where("id = ?", id).Delete(&model.FinancialTransaction{}).Error
}
