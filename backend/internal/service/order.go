package service

import (
	"super-order-web/internal/model"
	"super-order-web/pkg/util"

	"gorm.io/gorm"
)

// OrderService 订单服务
type OrderService struct {
	db *gorm.DB
}

// NewOrderService 创建订单服务
func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

// List 获取订单列表
func (s *OrderService) List(page, pageSize int, status string) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := s.db.Model(&model.Order{})

	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Customer").Preload("Items").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders).Error

	return orders, total, err
}

// GetByID 根据ID获取订单
func (s *OrderService) GetByID(id int64) (*model.Order, error) {
	var order model.Order
	err := s.db.Preload("Customer").Preload("Items.SKU").Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByOrderNo 根据订单号获取订单
func (s *OrderService) GetByOrderNo(orderNo string) (*model.Order, error) {
	var order model.Order
	err := s.db.Preload("Customer").Preload("Items.SKU").Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// Create 创建订单
func (s *OrderService) Create(order *model.Order, items []model.OrderItem) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		order.OrderNo = util.GenerateOrderNo()

		if err := tx.Create(order).Error; err != nil {
			return err
		}

		for i := range items {
			items[i].OrderID = order.ID
		}

		if err := tx.Create(&items).Error; err != nil {
			return err
		}

		return nil
	})
}

// Update 更新订单
func (s *OrderService) Update(order *model.Order) error {
	return s.db.Model(&model.Order{}).
		Where("id = ?", order.ID).
		Updates(map[string]interface{}{
			"customer_id":      order.CustomerID,
			"order_date":      order.OrderDate,
			"status":          order.Status,
			"total_cost_amount": order.TotalCostAmount,
			"total_sale_amount": order.TotalSaleAmount,
			"remarks":         order.Remarks,
		}).Error
}

// UpdateStatus 更新订单状态
func (s *OrderService) UpdateStatus(id int64, status string) error {
	return s.db.Model(&model.Order{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Delete 删除订单
func (s *OrderService) Delete(id int64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("order_id = ?", id).Delete(&model.OrderItem{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", id).Delete(&model.Order{}).Error
	})
}

// Settle 订单结算
func (s *OrderService) Settle(id int64) error {
	return s.db.Model(&model.Order{}).
		Where("id = ?", id).
		Update("is_settled", true).Error
}

// GetByStatus 根据状态获取订单列表
func (s *OrderService) GetByStatus(statuses []string) ([]model.Order, error) {
	var orders []model.Order
	err := s.db.Preload("Customer").Preload("Items.SKU").
		Where("status IN ?", statuses).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

// GetUnsettled 获取未结算的订单
func (s *OrderService) GetUnsettled() ([]model.Order, error) {
	var orders []model.Order
	err := s.db.Preload("Customer").Preload("Items.SKU").
		Where("is_settled = ?", false).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

// GetItems 获取订单明细
func (s *OrderService) GetItems(orderID int64) ([]model.OrderItem, error) {
	var items []model.OrderItem
	err := s.db.Preload("SKU").Where("order_id = ?", orderID).Find(&items).Error
	return items, err
}
