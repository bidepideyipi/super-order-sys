package service

import (
	"super-order-web/internal/model"
	"super-order-web/pkg/util"

	"gorm.io/gorm"
)

// CustomerService 客户服务
type CustomerService struct {
	db *gorm.DB
}

// NewCustomerService 创建客户服务
func NewCustomerService(db *gorm.DB) *CustomerService {
	return &CustomerService{db: db}
}

// List 获取客户列表
func (s *CustomerService) List() ([]model.Customer, error) {
	var customers []model.Customer
	err := s.db.Find(&customers).Error
	return customers, err
}

// GetByID 根据ID获取客户
func (s *CustomerService) GetByID(id string) (*model.Customer, error) {
	var customer model.Customer
	err := s.db.Where("customer_id = ?", id).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// Create 创建客户
func (s *CustomerService) Create(customer *model.Customer) error {
	customer.CustomerID = util.GenerateCustomerID()
	return s.db.Create(customer).Error
}

// Update 更新客户
func (s *CustomerService) Update(customer *model.Customer) error {
	return s.db.Model(&model.Customer{}).
		Where("customer_id = ?", customer.CustomerID).
		Updates(map[string]interface{}{
			"customer_name": customer.CustomerName,
		}).Error
}

// Delete 删除客户
func (s *CustomerService) Delete(id string) error {
	return s.db.Where("customer_id = ?", id).Delete(&model.Customer{}).Error
}
