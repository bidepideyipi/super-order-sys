package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

const baseURL = "http://localhost:8173/api"

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Helper: 发送 HTTP 请求
func makeRequest(method, url string, body interface{}, token string) (*Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ==================== SKU 接口测试 ====================

func TestSKUList(t *testing.T) {
	url := fmt.Sprintf("%s/sku/list", baseURL)
	resp, err := makeRequest("GET", url, nil, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 前端期望: res.data.data 是数组
	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}

	// 验证返回数据是数组
	data, ok := resp.Data.([]interface{})
	if !ok {
		t.Errorf("期望返回数组，实际: %T", resp.Data)
	} else {
		t.Logf("✓ SKU 列表测试通过，返回 %d 条数据", len(data))
	}
}

func TestSKUListPaginated(t *testing.T) {
	url := fmt.Sprintf("%s/sku/list-paginated?page=1&page_size=10", baseURL)
	resp, err := makeRequest("GET", url, nil, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}

	// 验证返回数据包含分页信息
	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Errorf("期望返回分页对象，实际: %T", resp.Data)
	} else {
		if _, hasTotal := data["total"]; !hasTotal {
			t.Error("分页数据缺少 total 字段")
		}
		if _, hasData := data["data"]; !hasData {
			t.Error("分页数据缺少 data 字段")
		}
		t.Logf("✓ SKU 分页列表测试通过")
	}
}

func TestSKUSearch(t *testing.T) {
	url := fmt.Sprintf("%s/sku/search?keyword=test", baseURL)
	resp, err := makeRequest("GET", url, nil, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}
	t.Logf("✓ SKU 搜索测试通过")
}

func TestSKUSearchPaginated(t *testing.T) {
	url := fmt.Sprintf("%s/sku/search-paginated?keyword=test&page=1&page_size=10", baseURL)
	resp, err := makeRequest("GET", url, nil, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}
	t.Logf("✓ SKU 分页搜索测试通过")
}

func TestSKUSearchWithCategory(t *testing.T) {
	url := fmt.Sprintf("%s/sku/search-with-category?keyword=test&category_id=CAT001", baseURL)
	resp, err := makeRequest("GET", url, nil, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}
	t.Logf("✓ SKU 按分类搜索测试通过")
}

func TestSKUCreate(t *testing.T) {
	// 使用时间戳生成唯一的 SKU code
	uniqueCode := fmt.Sprintf("TEST_SKU_%d", time.Now().Unix())
	skuData := map[string]interface{}{
		"sku_code":     uniqueCode,
		"name":         "测试SKU",
		"description":  "测试描述",
		"category_id":  "CAT001",
		"unit":         "个",
		"box_spec":     "10个/盒",
		"box_quantity": 10,
		"cost_price":   100.0,
		"sale_price":   150.0,
	}

	url := fmt.Sprintf("%s/sku", baseURL)
	resp, err := makeRequest("POST", url, skuData, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}
	t.Logf("✓ SKU 创建测试通过")
}

func TestSKUUpdate(t *testing.T) {
	skuData := map[string]interface{}{
		"name":        "更新后的SKU名称",
		"description": "更新后的描述",
		"category_id": "CAT001",
		"unit":        "个",
		"cost_price":  120.0,
		"sale_price":  180.0,
	}

	url := fmt.Sprintf("%s/sku/1", baseURL)
	resp, err := makeRequest("PUT", url, skuData, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}
	t.Logf("✓ SKU 更新测试通过")
}

func TestSKUDelete(t *testing.T) {
	url := fmt.Sprintf("%s/sku/999", baseURL) // 使用不存在的ID测试
	resp, err := makeRequest("DELETE", url, nil, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	t.Logf("✓ SKU 删除测试通过，响应: code=%d", resp.Code)
}

// ==================== Category 接口测试 ====================

func TestCategoryList(t *testing.T) {
	url := fmt.Sprintf("%s/category/list", baseURL)
	resp, err := makeRequest("GET", url, nil, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 前端期望: res.data.data 是数组
	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}

	data, ok := resp.Data.([]interface{})
	if !ok {
		t.Errorf("期望返回数组，实际: %T", resp.Data)
	} else {
		t.Logf("✓ Category 列表测试通过，返回 %d 条数据", len(data))
	}
}

func TestCategoryCreate(t *testing.T) {
	categoryData := map[string]interface{}{
		"category_name": "测试分类",
	}

	url := fmt.Sprintf("%s/category", baseURL)
	resp, err := makeRequest("POST", url, categoryData, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}
	t.Logf("✓ Category 创建测试通过")
}

func TestCategoryUpdate(t *testing.T) {
	categoryData := map[string]interface{}{
		"category_name": "更新后的分类名称",
	}

	url := fmt.Sprintf("%s/category/CAT001", baseURL)
	resp, err := makeRequest("PUT", url, categoryData, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}
	t.Logf("✓ Category 更新测试通过")
}

func TestCategoryDelete(t *testing.T) {
	url := fmt.Sprintf("%s/category/CAT999", baseURL)
	resp, err := makeRequest("DELETE", url, nil, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	t.Logf("✓ Category 删除测试通过，响应: code=%d", resp.Code)
}

// ==================== Customer 接口测试 ====================

func TestCustomerList(t *testing.T) {
	url := fmt.Sprintf("%s/customer/list", baseURL)
	resp, err := makeRequest("GET", url, nil, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}

	data, ok := resp.Data.([]interface{})
	if !ok {
		t.Errorf("期望返回数组，实际: %T", resp.Data)
	} else {
		t.Logf("✓ Customer 列表测试通过，返回 %d 条数据", len(data))
	}
}

func TestCustomerCreate(t *testing.T) {
	customerData := map[string]interface{}{
		"customer_name": "测试客户",
	}

	url := fmt.Sprintf("%s/customer", baseURL)
	resp, err := makeRequest("POST", url, customerData, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}
	t.Logf("✓ Customer 创建测试通过")
}

func TestCustomerUpdate(t *testing.T) {
	customerData := map[string]interface{}{
		"customer_name": "更新后的客户名称",
	}

	url := fmt.Sprintf("%s/customer/CUST001", baseURL)
	resp, err := makeRequest("PUT", url, customerData, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}
	t.Logf("✓ Customer 更新测试通过")
}

func TestCustomerDelete(t *testing.T) {
	url := fmt.Sprintf("%s/customer/CUST999", baseURL)
	resp, err := makeRequest("DELETE", url, nil, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	t.Logf("✓ Customer 删除测试通过，响应: code=%d", resp.Code)
}

// ==================== Order 接口测试 ====================

func TestOrderList(t *testing.T) {
	url := fmt.Sprintf("%s/order/list", baseURL)
	resp, err := makeRequest("GET", url, nil, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}
	t.Logf("✓ Order 列表测试通过")
}

func TestOrderGet(t *testing.T) {
	url := fmt.Sprintf("%s/order/1", baseURL)
	resp, err := makeRequest("GET", url, nil, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	t.Logf("✓ Order 详情测试通过，响应: code=%d", resp.Code)
}

func TestOrderGetProcessing(t *testing.T) {
	url := fmt.Sprintf("%s/order/processing", baseURL)
	resp, err := makeRequest("GET", url, nil, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}

	data, ok := resp.Data.([]interface{})
	if !ok {
		t.Errorf("期望返回数组，实际: %T", resp.Data)
	} else {
		t.Logf("✓ Order 进行中订单测试通过，返回 %d 条数据", len(data))
	}
}

func TestOrderGetUnsettled(t *testing.T) {
	url := fmt.Sprintf("%s/order/unsettled", baseURL)
	resp, err := makeRequest("GET", url, nil, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}

	data, ok := resp.Data.([]interface{})
	if !ok {
		t.Errorf("期望返回数组，实际: %T", resp.Data)
	} else {
		t.Logf("✓ Order 未结算订单测试通过，返回 %d 条数据", len(data))
	}
}

func TestOrderGetItems(t *testing.T) {
	url := fmt.Sprintf("%s/order/1/items", baseURL)
	resp, err := makeRequest("GET", url, nil, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}

	data, ok := resp.Data.([]interface{})
	if !ok {
		t.Errorf("期望返回数组，实际: %T", resp.Data)
	} else {
		t.Logf("✓ Order 明细测试通过，返回 %d 条数据", len(data))
	}
}

func TestOrderCreate(t *testing.T) {
	orderData := map[string]interface{}{
		"customer_id": "CUST001",
		"order_date":  "2026-06-12",
		"remarks":     "测试订单",
		"items": []map[string]interface{}{
			{
				"sku_code":          "SKU001",
				"product_name":      "测试商品",
				"quantity":          10,
				"cost_price":        100.0,
				"sale_price":        150.0,
				"total_cost_amount": 1000.0,
				"total_sale_amount": 1500.0,
			},
		},
	}

	url := fmt.Sprintf("%s/order", baseURL)
	resp, err := makeRequest("POST", url, orderData, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}
	t.Logf("✓ Order 创建测试通过")
}

func TestOrderUpdate(t *testing.T) {
	orderData := map[string]interface{}{
		"customer_id":       "CUST001",
		"order_date":        "2026-06-12",
		"total_cost_amount": 2000.0,
		"total_sale_amount": 3000.0,
		"remarks":           "更新后的订单",
	}

	url := fmt.Sprintf("%s/order/1", baseURL)
	resp, err := makeRequest("PUT", url, orderData, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}
	t.Logf("✓ Order 更新测试通过")
}

func TestOrderDelete(t *testing.T) {
	url := fmt.Sprintf("%s/order/999", baseURL)
	resp, err := makeRequest("DELETE", url, nil, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	t.Logf("✓ Order 删除测试通过，响应: code=%d", resp.Code)
}

// ==================== OrderItem 接口测试 ====================

func TestOrderItemCreate(t *testing.T) {
	orderItemData := map[string]interface{}{
		"order_id":           1,
		"sku_code":           "SKU001",
		"product_name":       "测试商品",
		"quantity":           5,
		"cost_price":         100.0,
		"sale_price":         150.0,
		"total_cost_amount":  500.0,
		"total_sale_amount":  750.0,
	}

	url := fmt.Sprintf("%s/order-item", baseURL)
	resp, err := makeRequest("POST", url, orderItemData, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}
	t.Logf("✓ OrderItem 创建测试通过")
}

func TestOrderItemUpdate(t *testing.T) {
	orderItemData := map[string]interface{}{
		"sku_code":           "SKU001",
		"product_name":       "测试商品",
		"quantity":           10,
		"cost_price":         100.0,
		"sale_price":         150.0,
		"total_cost_amount":  1000.0,
		"total_sale_amount":  1500.0,
	}

	url := fmt.Sprintf("%s/order-item/1", baseURL)
	resp, err := makeRequest("PUT", url, orderItemData, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}
	t.Logf("✓ OrderItem 更新测试通过")
}

func TestOrderItemDelete(t *testing.T) {
	url := fmt.Sprintf("%s/order-item/999", baseURL)
	resp, err := makeRequest("DELETE", url, nil, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	t.Logf("✓ OrderItem 删除测试通过，响应: code=%d", resp.Code)
}

// ==================== Financial 接口测试 ====================

func TestFinancialList(t *testing.T) {
	url := fmt.Sprintf("%s/financial/list", baseURL)
	resp, err := makeRequest("GET", url, nil, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}
	t.Logf("✓ Financial 列表测试通过")
}

func TestFinancialGetBalance(t *testing.T) {
	url := fmt.Sprintf("%s/financial/balance", baseURL)
	resp, err := makeRequest("GET", url, nil, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}
	t.Logf("✓ Financial 余额测试通过")
}

func TestFinancialGet(t *testing.T) {
	url := fmt.Sprintf("%s/financial/1", baseURL)
	resp, err := makeRequest("GET", url, nil, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	t.Logf("✓ Financial 详情测试通过，响应: code=%d", resp.Code)
}

func TestFinancialCreate(t *testing.T) {
	financialData := map[string]interface{}{
		"category":      "销售收入",
		"description":   "测试收入",
		"amount_change": 5000.0,
	}

	url := fmt.Sprintf("%s/financial", baseURL)
	resp, err := makeRequest("POST", url, financialData, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}
	t.Logf("✓ Financial 创建测试通过")
}

func TestFinancialUpdate(t *testing.T) {
	financialData := map[string]interface{}{
		"category":      "销售收入",
		"description":   "更新后的收入",
		"amount_change": 6000.0,
		"balance":       6000.0,
	}

	url := fmt.Sprintf("%s/financial/1", baseURL)
	resp, err := makeRequest("PUT", url, financialData, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("期望 code=0, 实际: %d, message: %s", resp.Code, resp.Message)
	}
	t.Logf("✓ Financial 更新测试通过")
}

func TestFinancialDelete(t *testing.T) {
	url := fmt.Sprintf("%s/financial/999", baseURL)
	resp, err := makeRequest("DELETE", url, nil, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	t.Logf("✓ Financial 删除测试通过，响应: code=%d", resp.Code)
}

// ==================== Common 接口测试 ====================

func TestCommonGetImage(t *testing.T) {
	url := fmt.Sprintf("%s/common/image/SKU001", baseURL)
	resp, err := makeRequest("GET", url, nil, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	t.Logf("✓ Common 获取图片测试通过，响应: code=%d", resp.Code)
}
