CREATE TABLE IF NOT EXISTS ai_chat_sessions (
    id          VARCHAR(36)  PRIMARY KEY DEFAULT (UUID()),
    user_id     VARCHAR(36)  NOT NULL,
    title       VARCHAR(200),
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_chat_sessions_user (user_id)
);
