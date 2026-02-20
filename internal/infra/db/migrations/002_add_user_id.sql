-- Migration: Add user_id for multi-user support
-- Binds blocked products and notifications to user's Bark Key

-- 屏蔽商品表：添加 user_id，复合主键
CREATE TABLE IF NOT EXISTS blocked_product_new (
    activity_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    create_time TEXT NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (activity_id, user_id)
);

-- 通知配置表：添加 user_id，复合主键
CREATE TABLE IF NOT EXISTS notification_config_new (
    activity_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    target_price REAL NOT NULL,
    last_notify_time TEXT,
    create_time TEXT NOT NULL DEFAULT (datetime('now')),
    update_time TEXT NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (activity_id, user_id)
);

-- 迁移现有数据到指定用户 (Bark Key: 3eHKCA7aL6fY9Raipx3fEP)
INSERT INTO blocked_product_new SELECT activity_id, '3eHKCA7aL6fY9Raipx3fEP', create_time FROM blocked_product;
INSERT INTO notification_config_new SELECT activity_id, '3eHKCA7aL6fY9Raipx3fEP', target_price, last_notify_time, create_time, update_time FROM notification_config;

-- 替换原表
DROP TABLE blocked_product;
DROP TABLE notification_config;
ALTER TABLE blocked_product_new RENAME TO blocked_product;
ALTER TABLE notification_config_new RENAME TO notification_config;

-- 索引
CREATE INDEX IF NOT EXISTS idx_blocked_product_user ON blocked_product(user_id);
CREATE INDEX IF NOT EXISTS idx_notification_user ON notification_config(user_id);

-- 用户设置表：存储用户的 Bark Key（用于后台推送）
CREATE TABLE IF NOT EXISTS user_settings (
    user_id TEXT PRIMARY KEY,
    bark_key TEXT NOT NULL,
    create_time TEXT NOT NULL DEFAULT (datetime('now')),
    update_time TEXT NOT NULL DEFAULT (datetime('now'))
);
