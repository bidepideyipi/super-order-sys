# Super Order Web - Backend

订单管理系统后端服务

## 技术栈

- Go 1.23
- Gin Web Framework
- GORM ORM
- SQLite 数据库
- Aliyun OSS

## 项目结构

```
backend/
├── cmd/
│   └── server/          # 主程序入口
│       └── main.go
├── internal/
│   ├── config/          # 配置管理
│   ├── handler/         # HTTP处理器
│   ├── middleware/      # 中间件
│   ├── model/           # 数据模型
│   ├── repository/      # 数据访问层
│   ├── service/         # 业务逻辑层
│   └── routes/          # 路由配置
├── pkg/
│   ├── database/        # 数据库初始化
│   ├── oss/            # OSS工具
│   ├── response/       # 统一响应
│   └── util/           # 工具函数
├── config.json         # 配置文件
├── go.mod
└── go.sum
```

## 快速开始

### 1. 安装依赖

```bash
cd backend
go mod download
```

### 2. 配置

编辑 `config.json` 文件：

```json
{
  "server": {
    "port": 8080,
    "mode": "debug"
  },
  "database": {
    "type": "sqlite",
    "sqlite": "../data/super_order.db"
  },
  "oss": {
    "endpoint": "",
    "access_key_id": "",
    "access_key_secret": "",
    "bucket_name": ""
  }
}
```

### 3. 运行

```bash
go run cmd/server/main.go
```

或指定配置文件路径：

```bash
CONFIG_PATH=/path/to/config.json go run cmd/server/main.go
```

### 4. 构建

```bash
go build -o bin/server cmd/server/main.go
```

## API 接口

### 客户管理
- `GET /api/customers` - 获取客户列表
- `GET /api/customers/:id` - 获取客户详情
- `POST /api/customers` - 创建客户
- `PUT /api/customers/:id` - 更新客户
- `DELETE /api/customers/:id` - 删除客户

### SKU分类
- `GET /api/sku-categories` - 获取分类列表
- `GET /api/sku-categories/:id` - 获取分类详情
- `POST /api/sku-categories` - 创建分类
- `PUT /api/sku-categories/:id` - 更新分类
- `DELETE /api/sku-categories/:id` - 删除分类

### SKU管理
- `GET /api/skus` - 获取SKU列表
- `GET /api/skus/all` - 获取所有SKU
- `GET /api/skus/:id` - 获取SKU详情
- `POST /api/skus` - 创建SKU
- `PUT /api/skus/:id` - 更新SKU
- `DELETE /api/skus/:id` - 删除SKU

### 订单管理
- `GET /api/orders` - 获取订单列表
- `GET /api/orders/:id` - 获取订单详情
- `POST /api/orders` - 创建订单
- `PUT /api/orders/:id` - 更新订单
- `PUT /api/orders/:id/status` - 更新订单状态
- `DELETE /api/orders/:id` - 删除订单
- `POST /api/orders/:id/settle` - 订单结算

### 财务流水
- `GET /api/financial-transactions` - 获取流水列表
- `GET /api/financial-transactions/balance` - 获取余额
- `POST /api/financial-transactions` - 创建流水
