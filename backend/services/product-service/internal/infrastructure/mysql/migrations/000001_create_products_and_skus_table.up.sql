CREATE TABLE IF NOT EXISTS products (
    id BIGINT NOT NULL PRIMARY KEY,
    shop_id BIGINT NOT NULL,
    shop_name VARCHAR(128) NOT NULL,
    title VARCHAR(255) NOT NULL,
    sub_title VARCHAR(255) NOT NULL,
    main_image_url VARCHAR(512) NOT NULL,
    description TEXT NOT NULL,
    review_status VARCHAR(32) NOT NULL,
    sale_status VARCHAR(32) NOT NULL,
    created_at DATETIME(3) NOT NULL,
    updated_at DATETIME(3) NOT NULL,
    KEY idx_products_shop_id (shop_id),
    KEY idx_products_public (review_status, sale_status, updated_at)
);

CREATE TABLE IF NOT EXISTS product_skus (
    id BIGINT NOT NULL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    price_amount BIGINT NOT NULL,
    currency VARCHAR(8) NOT NULL,
    stock INT NOT NULL,
    sale_status VARCHAR(32) NOT NULL,
    created_at DATETIME(3) NOT NULL,
    updated_at DATETIME(3) NOT NULL,
    KEY idx_product_skus_product_id (product_id),
    CONSTRAINT fk_product_skus_product_id FOREIGN KEY (product_id) REFERENCES products (id)
);

INSERT INTO products (id, shop_id, shop_name, title, sub_title, main_image_url, description, review_status, sale_status, created_at, updated_at)
VALUES
    (3100000000001, 4100000000001, '极客数码旗舰店', '机械键盘 K87', '热插拔三模办公游戏机械键盘', 'https://example.com/images/k87.jpg', '支持蓝牙、2.4G、有线三模连接，适合办公和游戏。', 'approved', 'on_sale', '2026-06-23 00:00:00.000', '2026-06-23 00:00:00.000'),
    (3100000000002, 4100000000002, '山系生活方式店', '户外徒步背包 28L', '轻量化通勤露营两用背包', 'https://example.com/images/bag-28l.jpg', '带电脑仓、防泼水面料，适合轻徒步和城市通勤。', 'approved', 'on_sale', '2026-06-23 00:00:00.000', '2026-06-23 00:00:00.000');

INSERT INTO product_skus (id, product_id, name, price_amount, currency, stock, sale_status, created_at, updated_at)
VALUES
    (3200000000001, 3100000000001, '白光红轴', 39900, 'CNY', 120, 'on_sale', '2026-06-23 00:00:00.000', '2026-06-23 00:00:00.000'),
    (3200000000002, 3100000000001, 'RGB茶轴', 45900, 'CNY', 80, 'on_sale', '2026-06-23 00:00:00.000', '2026-06-23 00:00:00.000'),
    (3200000000003, 3100000000002, '曜石黑', 29900, 'CNY', 56, 'on_sale', '2026-06-23 00:00:00.000', '2026-06-23 00:00:00.000'),
    (3200000000004, 3100000000002, '沙岩黄', 31900, 'CNY', 34, 'on_sale', '2026-06-23 00:00:00.000', '2026-06-23 00:00:00.000');
