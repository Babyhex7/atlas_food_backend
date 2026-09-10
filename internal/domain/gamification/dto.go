package gamification

import "time"

// ProfileResponse - respon untuk data profile gamifikasi pengguna
type ProfileResponse struct {
	UserID              string            `json:"user_id"`
	Name                string            `json:"name"`
	AvatarURL           *string           `json:"avatar_url"`
	TotalPoints         int               `json:"total_points"`
	Level               int               `json:"level"`
	RankName            string            `json:"rank_name"`
	NextLevelPoints     int               `json:"next_level_points"`
	ProgressToNextLevel int               `json:"progress_to_next_level"`
	CurrentStreak       int               `json:"current_streak"`
	MaxStreak           int               `json:"max_streak"`
	BadgesEarned        int               `json:"badges_earned"`
	Badges              []UserBadgeDTO    `json:"badges"`
}

// UserBadgeDTO - DTO untuk badge yang diraih
type UserBadgeDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	IconEmoji string    `json:"icon_emoji"`
	EarnedAt  time.Time `json:"earned_at"`
}

// LeaderboardEntry - DTO untuk satu baris leaderboard
type LeaderboardEntry struct {
	Rank          int     `json:"rank"`
	UserID        string  `json:"user_id"`
	Name          string  `json:"name"`
	AvatarURL     *string `json:"avatar_url"`
	Level         int     `json:"level"`
	RankName      string  `json:"rank_name"`
	TotalXP       int     `json:"total_xp"`
	CurrentStreak int     `json:"current_streak"`
}

// MyPosition - DTO untuk posisi pengguna saat ini
type MyPosition struct {
	Rank    int `json:"rank"`
	TotalXP int `json:"total_xp"`
}

// LeaderboardResponse - respon endpoint leaderboard
type LeaderboardResponse struct {
	Period      string             `json:"period"`
	Leaderboard []LeaderboardEntry `json:"leaderboard"`
	MyPosition  *MyPosition        `json:"my_position,omitempty"`
}

// AddXPRequest - internal usage
type AddXPRequest struct {
	UserID   string
	Activity string
	XP       int
}
