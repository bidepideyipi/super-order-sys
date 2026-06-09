package util

// PageRequest 分页请求
type PageRequest struct {
	Page     int `form:"page" binding:"min=1"`
	PageSize int `form:"page_size" binding:"min=1,max=100"`
}

// GetOffset 获取偏移量
func (p *PageRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetLimit 获取限制数量
func (p *PageRequest) GetLimit() int {
	return p.PageSize
}

// SetDefault 设置默认值
func (p *PageRequest) SetDefault() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
}

// PageResponse 分页响应
type PageResponse struct {
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	PageSize int        `json:"page_size"`
	Data    interface{} `json:"data"`
}
