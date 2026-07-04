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

// updateOrderTotals 更新订单的总金额
func (s *OrderItemService) updateOrderTotals(orderID int64) error {
	var items []model.OrderItem
	err := s.db.Where("order_id = ?", orderID).Find(&items).Error
	if err != nil {
		return err
	}

	var totalCostAmount float64
	var totalSaleAmount float64
	for _, item := range items {
		totalCostAmount += item.TotalCostAmount
		totalSaleAmount += item.TotalSaleAmount
	}

	return s.db.Model(&model.Order{}).
		Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"total_cost_amount": totalCostAmount,
			"total_sale_amount": totalSaleAmount,
		}).Error
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
	err := s.db.Create(item).Error
	if err != nil {
		return err
	}
	return s.updateOrderTotals(item.OrderID)
}

// Update 更新订单明细
func (s *OrderItemService) Update(item *model.OrderItem) error {
	err := s.db.Model(&model.OrderItem{}).
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
	if err != nil {
		return err
	}

	// 获取 order_id 并更新订单总额
	var orderItem model.OrderItem
	err = s.db.Where("id = ?", item.ID).First(&orderItem).Error
	if err != nil {
		return err
	}
	return s.updateOrderTotals(orderItem.OrderID)
}

// Delete 删除订单明细
func (s *OrderItemService) Delete(id int64) error {
	// 获取 order_id
	var orderItem model.OrderItem
	err := s.db.Where("id = ?", id).First(&orderItem).Error
	if err != nil {
		return err
	}

	// 删除明细
	err = s.db.Where("id = ?", id).Delete(&model.OrderItem{}).Error
	if err != nil {
		return err
	}

	// 更新订单总额
	return s.updateOrderTotals(orderItem.OrderID)
}

// GetByOrderID 获取订单的所有明细
func (s *OrderItemService) GetByOrderID(orderID int64) ([]model.OrderItem, error) {
	var items []model.OrderItem
	err := s.db.Preload("SKU").Where("order_id = ?", orderID).Find(&items).Error
	return items, err
}
