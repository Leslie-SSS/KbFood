-- 添加 platform 字段到 master_product 表
ALTER TABLE master_product ADD COLUMN platform TEXT DEFAULT '探探糖';

-- 创建索引
CREATE INDEX idx_master_product_platform ON master_product(platform);
