CREATE TABLE IF NOT EXISTS payment_orders (
    id BIGINT NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    order_group_id BIGINT NOT NULL,
    status VARCHAR(32) NOT NULL,
    channel VARCHAR(32) NOT NULL,
    pay_amount BIGINT NOT NULL,
    currency VARCHAR(8) NOT NULL,
    expire_at DATETIME(3) NOT NULL,
    paid_at DATETIME(3) NULL,
    created_at DATETIME(3) NOT NULL,
    updated_at DATETIME(3) NOT NULL,
    UNIQUE KEY uk_payment_orders_order_group_id (order_group_id),
    KEY idx_payment_orders_user_created (user_id, created_at),
    KEY idx_payment_orders_status_expire (status, expire_at)
);
