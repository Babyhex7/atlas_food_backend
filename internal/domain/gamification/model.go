package gamification

import (
	"time"
)

// UserGamification - model untuk tabel user_gamification
type UserGamification struct {
	UserID           string    `gorm:"type:char(36);primaryKey" json:"user_id"`
	TotalPoints      int       `gorm:"type:int;not null;default:0" json:"total_points"`
	CurrentStreak    int       `gorm:"type:int;not null;default:0" json:"current_streak"`
	MaxStreak        int       `gorm:"type:int;not null;default:0" json:"max_streak"`
	LastActivityDate *time.Time `gorm:"type:date" json:"last_activity_date,omitempty"`
	Level            int       `gorm:"type:int;not null;default:1" json:"level"`
	RankName         string    `gorm:"type:varchar(50);not null;default:'Nutri Novice'" json:"rank_name"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (UserGamification) TableName() string {
	return "user_gamification"
}

// Badge - model untuk tabel badges
type Badge struct {
	ID          string    `gorm:"type:varchar(50);primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:text;not null" json:"description"`
	IconEmoji   *string   `gorm:"type:varchar(10)" json:"icon_emoji,omitempty"`
	Criteria    *string   `gorm:"type:text" json:"criteria,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Badge) TableName() string {
	return "badges"
}

// UserBadge - model untuk tabel user_badges
type UserBadge struct {
	ID        string    `gorm:"type:char(36);primaryKey;default:(UUID())" json:"id"`
	UserID    string    `gorm:"type:char(36);not null" json:"user_id"`
	BadgeID   string    `gorm:"type:varchar(50);not null" json:"badge_id"`
	EarnedAt  time.Time `json:"earned_at"`

	Badge     Badge     `gorm:"foreignKey:BadgeID;references:ID" json:"badge"`
}

func (UserBadge) TableName() string {
	return "user_badges"
}

// XPHistory - model untuk tabel xp_history
type XPHistory struct {
	ID          string    `gorm:"type:char(36);primaryKey;default:(UUID())" json:"id"`
	UserID      string    `gorm:"type:char(36);not null" json:"user_id"`
	Activity    string    `gorm:"type:varchar(100);not null" json:"activity"`
	XPEarned    int       `gorm:"type:int;not null" json:"xp_earned"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func (XPHistory) TableName() string {
	return "xp_history"
}
