-- =============================================
-- Super Order 数据库创建脚本
-- 数据库: SQLite3
-- 生成时间: 2024-06-09
-- =============================================

-- =============================================
-- 1. SKU分类表
-- =============================================
CREATE TABLE sku_category (
    category_id TEXT PRIMARY KEY,
    category_name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- 2. 客户表
-- =============================================
CREATE TABLE customer (
    customer_id TEXT PRIMARY KEY,
    customer_name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- 3. 财务流水表
-- =============================================
CREATE TABLE financial_transaction (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category TEXT NOT NULL,
    description TEXT,
    amount_change REAL NOT NULL,
    balance REAL NOT NULL,
    is_settled INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- 4. SKU商品表
-- =============================================
CREATE TABLE sku (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sku_code TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    spec TEXT,
    unit TEXT DEFAULT '个',
    category_id TEXT NOT NULL,
    box_spec TEXT,
    box_quantity INTEGER DEFAULT 1,
    cost_price REAL DEFAULT 0,
    sale_price REAL DEFAULT 0,
    is_deleted INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- 5. 订单表
-- =============================================
CREATE TABLE `order` (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_no TEXT UNIQUE NOT NULL,
    customer_id TEXT NOT NULL,
    order_date TEXT NOT NULL,
    status TEXT DEFAULT 'pending',
    is_settled INTEGER DEFAULT 0,
    total_cost_amount REAL DEFAULT 0,
    total_sale_amount REAL DEFAULT 0,
    remarks TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- 6. 订单明细表
-- =============================================
CREATE TABLE order_item (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_id INTEGER NOT NULL,
    sku_id INTEGER,
    sku_code TEXT NOT NULL,
    product_name TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    cost_price REAL NOT NULL,
    sale_price REAL NOT NULL,
    total_cost_amount REAL NOT NULL,
    total_sale_amount REAL NOT NULL,
    settled_amount REAL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- 索引
-- =============================================
CREATE INDEX idx_sku_category_id ON sku(category_id);
CREATE INDEX idx_order_customer_id ON `order`(customer_id);
CREATE INDEX idx_order_order_no ON `order`(order_no);
CREATE INDEX idx_order_order_date ON `order`(order_date);
CREATE INDEX idx_order_status ON `order`(status);
CREATE INDEX idx_order_item_order_id ON order_item(order_id);
CREATE INDEX idx_order_item_sku_id ON order_item(sku_id);
CREATE INDEX idx_financial_transaction_category ON financial_transaction(category);
CREATE INDEX idx_financial_transaction_created_at ON financial_transaction(created_at);
