CREATE TABLE IF NOT EXISTS user_badges (
    id          VARCHAR(36)  PRIMARY KEY DEFAULT (UUID()),
    user_id     VARCHAR(36)  NOT NULL,
    badge_id    VARCHAR(50)  NOT NULL,
    earned_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id)  REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (badge_id) REFERENCES badges(id) ON DELETE CASCADE,
    UNIQUE KEY uq_user_badge (user_id, badge_id)
);
