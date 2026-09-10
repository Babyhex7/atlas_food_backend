CREATE TABLE IF NOT EXISTS xp_history (
    id          VARCHAR(36)  PRIMARY KEY DEFAULT (UUID()),
    user_id     VARCHAR(36)  NOT NULL,
    activity    VARCHAR(100) NOT NULL,
    xp_earned   INT          NOT NULL,
    description TEXT,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_xp_history_user (user_id),
    INDEX idx_xp_history_created_at (created_at)
);
