ALTER TABLE shop_orders
    DROP COLUMN paid_at;

ALTER TABLE order_groups
    DROP KEY idx_order_groups_status_payment_deadline,
    DROP COLUMN paid_at,
    DROP COLUMN payment_deadline_at;
