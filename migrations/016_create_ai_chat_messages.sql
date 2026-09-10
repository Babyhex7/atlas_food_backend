CREATE TABLE IF NOT EXISTS ai_chat_messages (
    id           VARCHAR(36)  PRIMARY KEY DEFAULT (UUID()),
    session_id   VARCHAR(36)  NOT NULL,
    role         ENUM('user', 'assistant', 'system') NOT NULL,
    content      TEXT         NOT NULL,
    token_used   INT          DEFAULT 0,
    model_used   VARCHAR(100),
    latency_ms   INT          DEFAULT 0,
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES ai_chat_sessions(id) ON DELETE CASCADE,
    INDEX idx_chat_messages_session (session_id)
);
