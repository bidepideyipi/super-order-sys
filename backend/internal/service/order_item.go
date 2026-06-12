package service

import (
	"super-order-web/internal/model"

	"gorm.io/gorm"
)

// OrderItemService 订单明细服务
type OrderItemService struct {
	db *gorm.DB
}

// NewOrderItemService 创建订单明细服务
func NewOrderItemService(db *gorm.DB) *OrderItemService {
	return &OrderItemService{db: db}
}

// GetByID 根据ID获取订单明细
func (s *OrderItemService) GetByID(id int64) (*model.OrderItem, error) {
	var item model.OrderItem
	err := s.db.Preload("SKU").Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Create 创建订单明细
func (s *OrderItemService) Create(item *model.OrderItem) error {
	return s.db.Create(item).Error
}

// Update 更新订单明细
func (s *OrderItemService) Update(item *model.OrderItem) error {
	return s.db.Model(&model.OrderItem{}).
		Where("id = ?", item.ID).
		Updates(map[string]interface{}{
			"sku_id":            item.SKUID,
			"sku_code":          item.SKUCode,
			"product_name":      item.ProductName,
			"quantity":          item.Quantity,
			"cost_price":        item.CostPrice,
			"sale_price":        item.SalePrice,
			"total_cost_amount": item.TotalCostAmount,
			"total_sale_amount": item.TotalSaleAmount,
			"settled_amount":    item.SettledAmount,
		}).Error
}

// Delete 删除订单明细
func (s *OrderItemService) Delete(id int64) error {
	return s.db.Where("id = ?", id).Delete(&model.OrderItem{}).Error
}

// GetByOrderID 获取订单的所有明细
func (s *OrderItemService) GetByOrderID(orderID int64) ([]model.OrderItem, error) {
	var items []model.OrderItem
	err := s.db.Preload("SKU").Where("order_id = ?", orderID).Find(&items).Error
	return items, err
}
