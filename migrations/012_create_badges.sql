CREATE TABLE IF NOT EXISTS badges (
    id           VARCHAR(50)  PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    description  TEXT         NOT NULL,
    icon_emoji   VARCHAR(10),
    criteria     TEXT,
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Seed Data Badges
INSERT IGNORE INTO badges (id, name, description, icon_emoji, criteria) VALUES
('FIRST_STEP', 'First Step', 'Pertama kali submit survei pangan', '🔰', 'submit_survey=1'),
('AI_ENTHUSIAST', 'AI Enthusiast', 'Gunakan analisis AI sebanyak 5 kali', '🤖', 'ai_analysis=5'),
('STREAK_7', 'Streak Master', 'Menjaga daily streak 7 hari berturut-turut', '🔥', 'streak=7'),
('STREAK_30', 'Flame Legend', 'Menjaga daily streak 30 hari berturut-turut', '🔥🔥', 'streak=30'),
('COLLAB_HERO', 'Collab Hero', 'Bergabung di 5 sesi kolaborasi live', '👥', 'collab_join=5'),
('NUTRI_CHAMPION', 'Nutri Champion', 'Mencapai XP level 3 (401 XP)', '🥗', 'level=3'),
('TOP_WEEKLY', 'Juara Mingguan', 'Berada di peringkat #1 leaderboard mingguan', '🏅', 'rank=1'),
('NUTRIBOT_USER', 'NutriBot User', 'Pertama kali menggunakan AI Chat NutriBot', '💬', 'ai_chat=1');
