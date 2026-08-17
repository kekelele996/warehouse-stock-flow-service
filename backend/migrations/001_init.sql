-- wmsflow 供应链仓储出入库协同管理系统 - 初始化迁移
-- 说明：运行时由 GORM AutoMigrate 保持模型与表结构同步，本脚本用于手动初始化与文档参考。

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(64) NOT NULL,
    role VARCHAR(32) NOT NULL,
    owner_id BIGINT,
    status VARCHAR(16) DEFAULT 'Active',
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_owner ON users(owner_id);

CREATE TABLE IF NOT EXISTS owners (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL UNIQUE,
    contact_name VARCHAR(64) NOT NULL,
    phone VARCHAR(32) NOT NULL,
    email VARCHAR(128),
    address VARCHAR(255),
    settlement_method VARCHAR(16) NOT NULL,
    credit_limit NUMERIC(14,2) DEFAULT 0,
    current_debt NUMERIC(14,2) DEFAULT 0,
    status VARCHAR(16) DEFAULT 'Active',
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    owner_id BIGINT NOT NULL,
    name VARCHAR(128) NOT NULL,
    sku VARCHAR(64) NOT NULL UNIQUE,
    barcode VARCHAR(64),
    category VARCHAR(64),
    spec VARCHAR(128),
    unit VARCHAR(16) NOT NULL,
    shelf_life_days INT,
    storage_requirement VARCHAR(16) DEFAULT 'Normal',
    volume NUMERIC(12,3) DEFAULT 0,
    weight NUMERIC(12,3) DEFAULT 0,
    price NUMERIC(14,2) DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_products_owner ON products(owner_id);

CREATE TABLE IF NOT EXISTS bin_locations (
    id BIGSERIAL PRIMARY KEY,
    area VARCHAR(4) NOT NULL,
    rack_no VARCHAR(16) NOT NULL,
    layer_no INT NOT NULL,
    column_no INT NOT NULL,
    capacity NUMERIC(10,2) NOT NULL,
    occupancy_rate NUMERIC(6,2) DEFAULT 0,
    storage_requirement VARCHAR(16) DEFAULT 'Normal',
    status VARCHAR(16) DEFAULT 'Available',
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_bins_area ON bin_locations(area);

CREATE TABLE IF NOT EXISTS inventories (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    owner_id BIGINT NOT NULL,
    bin_location_id BIGINT NOT NULL,
    batch_no VARCHAR(64),
    quantity INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE (product_id, bin_location_id, batch_no)
);

CREATE TABLE IF NOT EXISTS inbound_orders (
    id BIGSERIAL PRIMARY KEY,
    order_no VARCHAR(32) NOT NULL UNIQUE,
    owner_id BIGINT NOT NULL,
    supplier_name VARCHAR(128),
    expected_arrival_date TIMESTAMPTZ,
    actual_arrival_date TIMESTAMPTZ,
    status VARCHAR(16) NOT NULL,
    qc_inspector_id BIGINT,
    keeper_id BIGINT,
    remark VARCHAR(500),
    created_by BIGINT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS inbound_items (
    id BIGSERIAL PRIMARY KEY,
    inbound_order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    batch_no VARCHAR(64),
    expected_qty INT NOT NULL,
    actual_qty INT DEFAULT 0,
    qc_result VARCHAR(16),
    bin_location_id BIGINT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_inbound_items_order ON inbound_items(inbound_order_id);

CREATE TABLE IF NOT EXISTS outbound_orders (
    id BIGSERIAL PRIMARY KEY,
    order_no VARCHAR(32) NOT NULL UNIQUE,
    owner_id BIGINT NOT NULL,
    receiver_name VARCHAR(128) NOT NULL,
    receiver_address VARCHAR(255),
    required_ship_date TIMESTAMPTZ,
    actual_ship_date TIMESTAMPTZ,
    status VARCHAR(16) NOT NULL,
    picker_id BIGINT,
    checker_id BIGINT,
    tracking_no VARCHAR(64),
    created_by BIGINT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS outbound_items (
    id BIGSERIAL PRIMARY KEY,
    outbound_order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    bin_location_id BIGINT,
    expected_qty INT NOT NULL,
    actual_qty INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_outbound_items_order ON outbound_items(outbound_order_id);

CREATE TABLE IF NOT EXISTS operation_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    username VARCHAR(64),
    role VARCHAR(32),
    module VARCHAR(32),
    action VARCHAR(64),
    entity_type VARCHAR(32),
    entity_id VARCHAR(64),
    detail VARCHAR(500),
    ip VARCHAR(64),
    created_at TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_operation_logs_module ON operation_logs(module);
CREATE INDEX IF NOT EXISTS idx_operation_logs_action ON operation_logs(action);
