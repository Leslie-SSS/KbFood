-- 商品主表
CREATE TABLE IF NOT EXISTS product (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    activity_id TEXT NOT NULL UNIQUE,
    platform TEXT,
    region TEXT,
    title TEXT,
    shop_name TEXT,
    original_price REAL,
    current_price REAL,
    sales_status INTEGER,
    activity_create_time TEXT,  -- SQLite stores timestamps as TEXT or INTEGER
    create_time TEXT NOT NULL DEFAULT (datetime('now')),
    update_time TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_product_activity_id ON product(activity_id);
CREATE INDEX IF NOT EXISTS idx_product_create_time ON product(activity_create_time);
CREATE INDEX IF NOT EXISTS idx_product_platform ON product(platform);
CREATE INDEX IF NOT EXISTS idx_product_region ON product(region);

-- 标准商品库 (DT专用)
CREATE TABLE IF NOT EXISTS master_product (
    id TEXT PRIMARY KEY,
    region TEXT NOT NULL,
    standard_title TEXT NOT NULL,
    price REAL,
    status INTEGER,
    trust_score INTEGER DEFAULT 0,
    create_time TEXT NOT NULL DEFAULT (datetime('now')),
    update_time TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_master_region ON master_product(region);

-- 候选商品池
CREATE TABLE IF NOT EXISTS candidate_item (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_key TEXT NOT NULL,
    region TEXT NOT NULL,
    title_votes TEXT NOT NULL DEFAULT '{}',  -- JSON stored as TEXT
    total_occurrences INTEGER DEFAULT 0,
    last_price REAL,
    last_status INTEGER,
    first_seen_time TEXT NOT NULL DEFAULT (datetime('now')),
    last_seen_time TEXT NOT NULL DEFAULT (datetime('now')),
    create_time TEXT NOT NULL DEFAULT (datetime('now')),
    update_time TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_candidate_region ON candidate_item(region);
CREATE INDEX IF NOT EXISTS idx_candidate_group_key ON candidate_item(group_key);

-- 屏蔽商品
CREATE TABLE IF NOT EXISTS blocked_product (
    activity_id TEXT PRIMARY KEY,
    create_time TEXT NOT NULL DEFAULT (datetime('now'))
);

-- 通知配置
CREATE TABLE IF NOT EXISTS notification_config (
    activity_id TEXT PRIMARY KEY,
    target_price REAL NOT NULL,
    last_notify_time TEXT,
    create_time TEXT NOT NULL DEFAULT (datetime('now')),
    update_time TEXT NOT NULL DEFAULT (datetime('now'))
);

-- 价格趋势
CREATE TABLE IF NOT EXISTS product_price_trend (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    activity_id TEXT NOT NULL,
    price REAL NOT NULL,
    record_date TEXT NOT NULL,  -- DATE stored as TEXT in format YYYY-MM-DD
    create_time TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE(activity_id, record_date)
);

CREATE INDEX IF NOT EXISTS idx_trend_activity_date ON product_price_trend(activity_id, record_date);

-- 同步状态表
CREATE TABLE IF NOT EXISTS sync_status (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    job_name TEXT NOT NULL UNIQUE,
    last_run_time TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    product_count INTEGER DEFAULT 0,
    error_message TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);
