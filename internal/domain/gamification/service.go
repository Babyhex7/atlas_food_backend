package gamification

import (
	"log"
	"time"
)

// Service interface untuk business logic gamifikasi
type Service interface {
	AddXP(userID string, activity string, xp int, desc string) (*UserGamification, error)
	GetProfile(userID string) (*ProfileResponse, error)
	GetLeaderboard(period string, limit int, userID string) (*LeaderboardResponse, error)
	GetBadges(userID string) ([]UserBadgeDTO, error)
	GetXPHistory(userID string) ([]XPHistory, error)
}

type service struct {
	repo Repository
}

// NewService creates a new gamification service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// getLevelAndRank - Hitung level dan rank_name berdasarkan poin (sesuai PRD)
func getLevelAndRank(points int) (int, string, int) {
	switch {
	case points >= 801:
		return 4, "Pahlawan Gizi Atlas", 0 // Level max
	case points >= 401:
		return 3, "Nutri Champion", 801
	case points >= 151:
		return 2, "Gizi Explorer", 401
	default:
		return 1, "Nutri Novice", 151
	}
}

// AddXP - Adds XP to a user, updates streak, level, and checks for badges
func (s *service) AddXP(userID string, activity string, xp int, desc string) (*UserGamification, error) {
	var finalUG *UserGamification

	err := s.repo.WithTransaction(func(txRepo Repository) error {
		// 1. Ambil/buat user_gamification record
		ug, err := txRepo.GetUserGamification(userID)
		if err != nil {
			return err
		}

		now := time.Now()
		isNew := false
		if ug == nil {
			ug = &UserGamification{
				UserID:           userID,
				TotalPoints:      0,
				CurrentStreak:    0,
				MaxStreak:        0,
				LastActivityDate: nil,
				Level:            1,
				RankName:         "Nutri Novice",
			}
			isNew = true
		}

		// 2. Tambah poin
		ug.TotalPoints += xp

		// 3. Recalculate level
		newLevel, newRank, _ := getLevelAndRank(ug.TotalPoints)
		ug.Level = newLevel
		ug.RankName = newRank

		// 4. Update Streak (berdasarkan hari kalender WIB, kita asumsikan UTC+7 atau simple local time)
		if ug.LastActivityDate == nil {
			ug.CurrentStreak = 1
			ug.MaxStreak = 1
		} else {
			lastDate := *ug.LastActivityDate
			y1, m1, d1 := lastDate.Date()
			y2, m2, d2 := now.Date()

			if y1 == y2 && m1 == m2 && d1 == d2 {
				// Sama hari -> skip streak
			} else {
				// Cek apakah kemarin
				yesterday := now.AddDate(0, 0, -1)
				y3, m3, d3 := yesterday.Date()
				if y1 == y3 && m1 == m3 && d1 == d3 {
					ug.CurrentStreak++
					if ug.CurrentStreak > ug.MaxStreak {
						ug.MaxStreak = ug.CurrentStreak
					}
				} else {
					ug.CurrentStreak = 1 // Reset streak karena bolong
				}
			}
		}
		ug.LastActivityDate = &now

		if isNew {
			err = txRepo.CreateUserGamification(ug)
		} else {
			err = txRepo.UpdateUserGamification(ug)
		}
		if err != nil {
			return err
		}

		// 5. Simpan history
		history := &XPHistory{
			UserID:      userID,
			Activity:    activity,
			XPEarned:    xp,
			Description: desc,
		}
		if err := txRepo.AddXPHistory(history); err != nil {
			return err
		}

		// 6. Cek & Award Badges (async atau sync, di sini sync)
		s.checkAndAwardBadges(txRepo, ug)

		finalUG = ug
		return nil
	})

	if err != nil {
		log.Printf("[Gamification] Gagal menambah XP untuk user %s: %v", userID, err)
		return nil, err
	}

	return finalUG, nil
}

func (s *service) checkAndAwardBadges(repo Repository, ug *UserGamification) {
	// Dapatkan semua badge
	badges, err := repo.GetAllBadges()
	if err != nil {
		log.Printf("[Gamification] Error get badges: %v", err)
		return
	}

	for _, badge := range badges {
		hasBadge, _ := repo.HasBadge(ug.UserID, badge.ID)
		if hasBadge {
			continue
		}

		award := false
		switch badge.ID {
		case "FIRST_STEP":
			// Kita anggap submit_survey adalah event-nya. XP history check.
			history, _ := repo.GetXPHistory(ug.UserID, 100)
			count := 0
			for _, h := range history {
				if h.Activity == "submit_survey" {
					count++
				}
			}
			if count >= 1 {
				award = true
			}
		case "AI_ENTHUSIAST":
			history, _ := repo.GetXPHistory(ug.UserID, 100)
			count := 0
			for _, h := range history {
				if h.Activity == "ai_analysis" {
					count++
				}
			}
			if count >= 5 {
				award = true
			}
		case "STREAK_7":
			if ug.MaxStreak >= 7 {
				award = true
			}
		case "STREAK_30":
			if ug.MaxStreak >= 30 {
				award = true
			}
		case "COLLAB_HERO":
			history, _ := repo.GetXPHistory(ug.UserID, 100)
			count := 0
			for _, h := range history {
				if h.Activity == "collab_join" {
					count++
				}
			}
			if count >= 5 {
				award = true
			}
		case "NUTRI_CHAMPION":
			if ug.Level >= 3 {
				award = true
			}
		case "NUTRIBOT_USER":
			history, _ := repo.GetXPHistory(ug.UserID, 10)
			for _, h := range history {
				if h.Activity == "ai_chat" {
					award = true
					break
				}
			}
		}

		if award {
			ub := &UserBadge{
				UserID:  ug.UserID,
				BadgeID: badge.ID,
			}
			if err := repo.AwardBadge(ub); err != nil {
				log.Printf("[Gamification] Error award badge %s: %v", badge.ID, err)
			}
		}
	}
}

func (s *service) GetProfile(userID string) (*ProfileResponse, error) {
	ug, err := s.repo.GetUserGamification(userID)
	if err != nil {
		return nil, err
	}
	if ug == nil {
		ug = &UserGamification{
			UserID:      userID,
			TotalPoints: 0,
			Level:       1,
			RankName:    "Nutri Novice",
		}
	}

	userBadges, _ := s.repo.GetUserBadges(userID)
	var badgesDTO []UserBadgeDTO
	for _, ub := range userBadges {
		emoji := ""
		if ub.Badge.IconEmoji != nil {
			emoji = *ub.Badge.IconEmoji
		}
		badgesDTO = append(badgesDTO, UserBadgeDTO{
			ID:        ub.Badge.ID,
			Name:      ub.Badge.Name,
			IconEmoji: emoji,
			EarnedAt:  ub.EarnedAt,
		})
	}

	_, _, nextLevelPoints := getLevelAndRank(ug.TotalPoints)
	progress := 0
	if nextLevelPoints > 0 {
		progress = (ug.TotalPoints * 100) / nextLevelPoints
	} else {
		progress = 100 // Level max
	}

	// Butuh nama dan avatar, idealnya di-join dengan User table, tapi
	// untuk kesederhanaan respon, kita bisa populate sebagian atau 
	// harus inject AuthRepo. Disini dibiarkan kosong, atau Frontend mengisi dari auth/me
	
	resp := &ProfileResponse{
		UserID:              userID,
		Name:                "", // Akan diisi di handler / via Auth
		AvatarURL:           nil,
		TotalPoints:         ug.TotalPoints,
		Level:               ug.Level,
		RankName:            ug.RankName,
		NextLevelPoints:     nextLevelPoints,
		ProgressToNextLevel: progress,
		CurrentStreak:       ug.CurrentStreak,
		MaxStreak:           ug.MaxStreak,
		BadgesEarned:        len(userBadges),
		Badges:              badgesDTO,
	}

	return resp, nil
}

func (s *service) GetLeaderboard(period string, limit int, userID string) (*LeaderboardResponse, error) {
	var entries []LeaderboardEntry
	var err error

	switch period {
	case "weekly":
		entries, err = s.repo.GetLeaderboardWeekly(limit)
	case "monthly":
		entries, err = s.repo.GetLeaderboardMonthly(limit)
	default:
		entries, err = s.repo.GetLeaderboardAllTime(limit)
		period = "all_time"
	}

	if err != nil {
		return nil, err
	}

	// Cari posisi user
	var myPos *MyPosition
	for i, entry := range entries {
		// Override Rank sesuai urutan array
		entries[i].Rank = i + 1
		if entry.UserID == userID {
			myPos = &MyPosition{
				Rank:    entries[i].Rank,
				TotalXP: entry.TotalXP,
			}
		}
	}

	// Jika user tidak ada di top N, kita tidak mem-query secara spesifik untuk rank aslinya demi efisiensi
	// Tapi setidaknya jika dia tidak ada, rank nil atau placeholder.

	return &LeaderboardResponse{
		Period:      period,
		Leaderboard: entries,
		MyPosition:  myPos,
	}, nil
}

func (s *service) GetBadges(userID string) ([]UserBadgeDTO, error) {
	userBadges, err := s.repo.GetUserBadges(userID)
	if err != nil {
		return nil, err
	}

	var badgesDTO []UserBadgeDTO
	for _, ub := range userBadges {
		emoji := ""
		if ub.Badge.IconEmoji != nil {
			emoji = *ub.Badge.IconEmoji
		}
		badgesDTO = append(badgesDTO, UserBadgeDTO{
			ID:        ub.Badge.ID,
			Name:      ub.Badge.Name,
			IconEmoji: emoji,
			EarnedAt:  ub.EarnedAt,
		})
	}
	return badgesDTO, nil
}

func (s *service) GetXPHistory(userID string) ([]XPHistory, error) {
	return s.repo.GetXPHistory(userID, 50)
}
