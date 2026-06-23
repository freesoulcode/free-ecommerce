ALTER TABLE order_groups
    ADD COLUMN payment_deadline_at DATETIME(3) NOT NULL DEFAULT '1970-01-01 00:00:00.000' AFTER item_count,
    ADD COLUMN paid_at DATETIME(3) NULL AFTER payment_deadline_at,
    ADD KEY idx_order_groups_status_payment_deadline (status, payment_deadline_at);

ALTER TABLE shop_orders
    ADD COLUMN paid_at DATETIME(3) NULL AFTER item_count;

UPDATE order_groups
SET payment_deadline_at = created_at
WHERE payment_deadline_at = '1970-01-01 00:00:00.000';
