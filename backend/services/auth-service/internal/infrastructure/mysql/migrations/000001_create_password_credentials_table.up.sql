CREATE TABLE IF NOT EXISTS password_credentials (
    user_id BIGINT NOT NULL PRIMARY KEY,
    email VARCHAR(255) NULL,
    phone VARCHAR(32) NULL,
    password_hash VARCHAR(512) NOT NULL,
    password_algo VARCHAR(32) NOT NULL,
    created_at DATETIME(3) NOT NULL,
    updated_at DATETIME(3) NOT NULL,
    UNIQUE KEY uk_password_credentials_email (email),
    UNIQUE KEY uk_password_credentials_phone (phone)
);
