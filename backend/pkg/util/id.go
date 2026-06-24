package util

import (
	"fmt"
	"time"
)

// GenerateCustomerID 生成客户ID
func GenerateCustomerID() string {
	return fmt.Sprintf("C%09d", time.Now().Unix()%1000000000)
}

// GenerateOrderNo 生成订单号
// 规则：customer_id + yyyyMMdd
func GenerateOrderNo(customerID string) string {
	return customerID + time.Now().Format("20060102")
}

// GenerateCategoryID 生成分类ID
func GenerateCategoryID() string {
	return fmt.Sprintf("CAT%09d", time.Now().Unix()%1000000000)
}
