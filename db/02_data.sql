-- 开始事务以确保数据一致性
START TRANSACTION;

-- 清空现有数据（如果表非空）
DELETE FROM t_product_attributes;
DELETE FROM t_product_skus;
DELETE FROM t_products;
DELETE FROM t_categories;

-- 重置自增ID
ALTER TABLE t_categories AUTO_INCREMENT = 1;
ALTER TABLE t_products AUTO_INCREMENT = 1;
ALTER TABLE t_product_skus AUTO_INCREMENT = 1;
ALTER TABLE t_product_attributes AUTO_INCREMENT = 1;

-- ==========================================
-- 插入商户ID为1001的测试数据
-- ==========================================

-- 1. 插入商品分类
-- 饮品分类 (排序权重5)
INSERT INTO t_categories (merchant_id, name, sort_order, is_visible) VALUES 
(1001, '精选咖啡', 5, 1);
-- 食物分类 (排序权重3)
INSERT INTO t_categories (merchant_id, name, sort_order, is_visible) VALUES 
(1001, '轻食简餐', 3, 1);

-- 2. 插入商品基础信息
-- 咖啡类商品
INSERT INTO t_products (category_id, merchant_id, name, description, image_url, sales_count, lowest_price, status) VALUES 
(1, 1001, '美式咖啡', '经典浓缩咖啡与热水的完美融合', 'https://example.com/americano.jpg', 150, 25.00, 1),
(1, 1001, '拿铁咖啡', '浓郁奶香与丝滑浓缩的绝妙搭配', 'https://example.com/latte.jpg', 200, 28.00, 1);

-- 食物类商品
INSERT INTO t_products (category_id, merchant_id, name, description, image_url, sales_count, lowest_price, status) VALUES 
(2, 1001, '金枪鱼三明治', '新鲜金枪鱼沙拉搭配全麦面包', 'https://example.com/tuna_sandwich.jpg', 80, 22.00, 1),
(2, 1001, '水果沙拉', '当季新鲜水果搭配自制酸奶酱', 'https://example.com/fruit_salad.jpg', 60, 18.00, 1);

-- 3. 插入商品规格 (SKUs)
-- 美式咖啡规格
INSERT INTO t_product_skus (product_id, spec_name, spec_value, price, stock) VALUES 
(1, '容量', '中杯', 25.00, 100),
(1, '容量', '大杯', 28.00, 50);

-- 拿铁咖啡规格
INSERT INTO t_product_skus (product_id, spec_name, spec_value, price, stock) VALUES 
(2, '容量', '中杯', 28.00, 80),
(2, '容量', '大杯', 32.00, 30);

-- 金枪鱼三明治规格 (不同套餐)
INSERT INTO t_product_skus (product_id, spec_name, spec_value, price, stock) VALUES 
(3, '套餐', '单点', 22.00, 40),
(3, '套餐', '套餐A(含饮料)', 29.00, 20);

-- 4. 插入商品属性 (Attributes)
-- 咖啡通用属性
INSERT INTO t_product_attributes (product_id, attr_name, attr_values) VALUES 
(1, '甜度', '标准糖,半糖,无糖'),
(1, '冰度', '正常冰,少冰,去冰'),
(2, '甜度', '标准糖,半糖,无糖'),
(2, '冰度', '正常冰,少冰,去冰');

-- ==========================================
-- 插入商户ID为1002的测试数据
-- ==========================================

-- 1. 插入商品分类
INSERT INTO t_categories (merchant_id, name, sort_order, is_visible) VALUES 
(1002, '招牌奶茶', 5, 1),
(1002, '甜品', 4, 1);

-- 2. 插入商品基础信息
INSERT INTO t_products (category_id, merchant_id, name, description, image_url, sales_count, lowest_price, status) VALUES 
(3, 1002, '珍珠奶茶', '经典台式风味奶茶', 'https://example.com/milk_tea.jpg', 300, 15.00, 1),
(4, 1002, '抹茶大福', '软糯外皮包裹香甜抹茶馅', 'https://example.com/mochi.jpg', 50, 12.00, 1);

-- 3. 插入商品规格 (SKUs)
INSERT INTO t_product_skus (product_id, spec_name, spec_value, price, stock) VALUES 
(5, '糖度', '正常糖', 15.00, 200),
(5, '糖度', '少糖', 15.00, 150),
(6, '口味', '原味', 12.00, 30);

-- 4. 插入商品属性 (Attributes)
INSERT INTO t_product_attributes (product_id, attr_name, attr_values) VALUES 
(5, '冰度', '正常冰,少冰,去冰');

-- 提交事务
COMMIT;