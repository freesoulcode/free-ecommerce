CREATE TABLE IF NOT EXISTS refresh_sessions (
    id BIGINT NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token_hash CHAR(64) NOT NULL,
    device_id VARCHAR(128) NULL,
    user_agent VARCHAR(512) NULL,
    client_ip VARCHAR(64) NULL,
    expires_at DATETIME(3) NOT NULL,
    revoked_at DATETIME(3) NULL,
    replaced_by_session_id BIGINT NULL,
    created_at DATETIME(3) NOT NULL,
    updated_at DATETIME(3) NOT NULL,
    UNIQUE KEY uk_refresh_sessions_token_hash (token_hash),
    KEY idx_refresh_sessions_user_id (user_id),
    KEY idx_refresh_sessions_expires_at (expires_at)
);
