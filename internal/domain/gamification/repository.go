package gamification

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Repository interface untuk modul gamification
type Repository interface {
	// Profile & User Gamification
	GetUserGamification(userID string) (*UserGamification, error)
	CreateUserGamification(ug *UserGamification) error
	UpdateUserGamification(ug *UserGamification) error
	
	// History
	AddXPHistory(history *XPHistory) error
	GetXPHistory(userID string, limit int) ([]XPHistory, error)
	GetXPEarnedSince(userID string, since time.Time) (int, error)

	// Badges
	GetAllBadges() ([]Badge, error)
	GetUserBadges(userID string) ([]UserBadge, error)
	AwardBadge(userBadge *UserBadge) error
	HasBadge(userID, badgeID string) (bool, error)

	// Leaderboard
	GetLeaderboardAllTime(limit int) ([]LeaderboardEntry, error)
	GetLeaderboardWeekly(limit int) ([]LeaderboardEntry, error)
	GetLeaderboardMonthly(limit int) ([]LeaderboardEntry, error)

	// Transactions
	WithTransaction(fn func(txRepo Repository) error) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new gamification repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) WithTransaction(fn func(txRepo Repository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &repository{db: tx}
		return fn(txRepo)
	})
}

func (r *repository) GetUserGamification(userID string) (*UserGamification, error) {
	var ug UserGamification
	err := r.db.Where("user_id = ?", userID).First(&ug).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil if not found
		}
		return nil, err
	}
	return &ug, nil
}

func (r *repository) CreateUserGamification(ug *UserGamification) error {
	return r.db.Create(ug).Error
}

func (r *repository) UpdateUserGamification(ug *UserGamification) error {
	return r.db.Save(ug).Error
}

func (r *repository) AddXPHistory(history *XPHistory) error {
	return r.db.Create(history).Error
}

func (r *repository) GetXPHistory(userID string, limit int) ([]XPHistory, error) {
	var history []XPHistory
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Limit(limit).Find(&history).Error
	return history, err
}

func (r *repository) GetXPEarnedSince(userID string, since time.Time) (int, error) {
	var totalXP struct {
		Total int
	}
	err := r.db.Model(&XPHistory{}).
		Select("SUM(xp_earned) as total").
		Where("user_id = ? AND created_at >= ?", userID, since).
		Scan(&totalXP).Error
	return totalXP.Total, err
}

func (r *repository) GetAllBadges() ([]Badge, error) {
	var badges []Badge
	err := r.db.Find(&badges).Error
	return badges, err
}

func (r *repository) GetUserBadges(userID string) ([]UserBadge, error) {
	var userBadges []UserBadge
	err := r.db.Preload("Badge").Where("user_id = ?", userID).Order("earned_at desc").Find(&userBadges).Error
	return userBadges, err
}

func (r *repository) AwardBadge(userBadge *UserBadge) error {
	return r.db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(userBadge).Error
}

func (r *repository) HasBadge(userID, badgeID string) (bool, error) {
	var count int64
	err := r.db.Model(&UserBadge{}).Where("user_id = ? AND badge_id = ?", userID, badgeID).Count(&count).Error
	return count > 0, err
}

func (r *repository) GetLeaderboardAllTime(limit int) ([]LeaderboardEntry, error) {
	var entries []LeaderboardEntry
	err := r.db.Table("user_gamification").
		Select("users.id as user_id, users.name, users.photo_url as avatar_url, user_gamification.level, user_gamification.rank_name, user_gamification.total_points as total_xp, user_gamification.current_streak").
		Joins("JOIN users ON users.id = user_gamification.user_id").
		Order("user_gamification.total_points DESC").
		Limit(limit).
		Scan(&entries).Error
	return entries, err
}

// GetLeaderboardWeekly dan GetLeaderboardMonthly dihitung dari total xp_history dalam periode tertentu
func (r *repository) GetLeaderboardWeekly(limit int) ([]LeaderboardEntry, error) {
	var entries []LeaderboardEntry
	since := time.Now().AddDate(0, 0, -7)
	err := r.db.Table("xp_history").
		Select("users.id as user_id, users.name, users.photo_url as avatar_url, user_gamification.level, user_gamification.rank_name, SUM(xp_history.xp_earned) as total_xp, user_gamification.current_streak").
		Joins("JOIN users ON users.id = xp_history.user_id").
		Joins("JOIN user_gamification ON user_gamification.user_id = xp_history.user_id").
		Where("xp_history.created_at >= ?", since).
		Group("users.id, users.name, users.photo_url, user_gamification.level, user_gamification.rank_name, user_gamification.current_streak").
		Order("total_xp DESC").
		Limit(limit).
		Scan(&entries).Error
	return entries, err
}

func (r *repository) GetLeaderboardMonthly(limit int) ([]LeaderboardEntry, error) {
	var entries []LeaderboardEntry
	since := time.Now().AddDate(0, -1, 0)
	err := r.db.Table("xp_history").
		Select("users.id as user_id, users.name, users.photo_url as avatar_url, user_gamification.level, user_gamification.rank_name, SUM(xp_history.xp_earned) as total_xp, user_gamification.current_streak").
		Joins("JOIN users ON users.id = xp_history.user_id").
		Joins("JOIN user_gamification ON user_gamification.user_id = xp_history.user_id").
		Where("xp_history.created_at >= ?", since).
		Group("users.id, users.name, users.photo_url, user_gamification.level, user_gamification.rank_name, user_gamification.current_streak").
		Order("total_xp DESC").
		Limit(limit).
		Scan(&entries).Error
	return entries, err
}
