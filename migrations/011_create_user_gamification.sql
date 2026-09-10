CREATE TABLE IF NOT EXISTS user_gamification (
    user_id           VARCHAR(36)  PRIMARY KEY,
    total_points      INT          NOT NULL DEFAULT 0,
    current_streak    INT          NOT NULL DEFAULT 0,
    max_streak        INT          NOT NULL DEFAULT 0,
    last_activity_date DATE,
    level             INT          NOT NULL DEFAULT 1,
    rank_name         VARCHAR(50)  NOT NULL DEFAULT 'Nutri Novice',
    created_at        TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
