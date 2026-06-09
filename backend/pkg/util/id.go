package util

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateCustomerID 生成客户ID
func GenerateCustomerID() string {
	return fmt.Sprintf("C%09d", time.Now().Unix()%1000000000)
}

// GenerateOrderNo 生成订单号
func GenerateOrderNo() string {
	return fmt.Sprintf("SO%s%04d", time.Now().Format("20060102"), rand.Intn(10000))
}

// GenerateCategoryID 生成分类ID
func GenerateCategoryID() string {
	return fmt.Sprintf("CAT%09d", time.Now().Unix()%1000000000)
}
