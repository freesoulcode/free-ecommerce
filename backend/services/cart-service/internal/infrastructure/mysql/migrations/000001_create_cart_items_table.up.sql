CREATE TABLE IF NOT EXISTS cart_items (
    id BIGINT NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    sku_id BIGINT NOT NULL,
    quantity INT NOT NULL,
    selected BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME(3) NOT NULL,
    updated_at DATETIME(3) NOT NULL,
    UNIQUE KEY uk_cart_items_user_sku (user_id, sku_id),
    KEY idx_cart_items_user_selected (user_id, selected, updated_at),
    KEY idx_cart_items_user_updated (user_id, updated_at)
);
