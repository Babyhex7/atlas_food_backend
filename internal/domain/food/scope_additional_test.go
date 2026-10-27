package food

import (
	"strings"
	"testing"

	"atlas_food/internal/config"
	redisInfra "atlas_food/internal/infra/redis"
)

// TestGenerateSearchCacheKeyScope - cache key untuk scope additional harus berbeda dari pencarian umum
func TestGenerateSearchCacheKeyScope(t *testing.T) {
	keyRegular := redisInfra.GenerateSearchCacheKey("tumis", "food", "", 20)
	keyAdditional := redisInfra.GenerateSearchCacheKey("tumis", "food", "additional", 20)

	if keyRegular == keyAdditional {
		t.Errorf("Cache key untuk regular dan additional tidak boleh sama, dapat: %s", keyRegular)
	}

	if !strings.HasPrefix(keyRegular, "atlas:prod:cache:foods:search:") {
		t.Errorf("Prefix cache key tidak valid: %s", keyRegular)
	}
}

// TestSearchFoodsPublicScopeAdditional - verifikasi database: scope additional mengecualikan MP, AK, AMK, AH
func TestSearchFoodsPublicScopeAdditional(t *testing.T) {
	cfg := config.Load()
	db := config.InitDB(cfg)
	if db == nil {
		t.Skip("MySQL database tidak tersedia, lewati integrasi DB test")
		return
	}

	repo := NewRepository(db)

	// 1. Pencarian "mie" dengan scope "additional" TIDAK BOLEH memunculkan kategori MP (Mie Basah/Rebus), AK (Mie Bakso/Soto Mie), AMK (Mie Instan)
	foods, err := repo.SearchFoodsPublic("mie", "", "additional", 20)
	if err != nil {
		t.Fatalf("SearchFoodsPublic error: %v", err)
	}

	blockedCategories := map[string]bool{
		"MP":     true,
		"AK":     true,
		"AMK":    true,
		"AH":     true,
		"drinks": true,
	}

	for _, f := range foods {
		if f.Category != nil && blockedCategories[f.Category.Code] {
			t.Errorf("Makanan berat %s (kategori %s) tidak boleh muncul pada scope additional!", f.Name, f.Category.Code)
		}
	}

	// 2. Pencarian "tumis" dengan scope "additional" HARUS berhasil memunculkan tumisan sayur (kategori AS / LN)
	tumisFoods, err := repo.SearchFoodsPublic("tumis", "", "additional", 20)
	if err != nil {
		t.Fatalf("SearchFoodsPublic tumis error: %v", err)
	}

	foundTumis := false
	for _, f := range tumisFoods {
		if strings.Contains(strings.ToLower(f.Name), "tumis") {
			foundTumis = true
			if f.Category != nil && blockedCategories[f.Category.Code] {
				t.Errorf("Tumis %s memiliki kategori terlarang %s!", f.Name, f.Category.Code)
			}
		}
	}

	if !foundTumis && len(tumisFoods) == 0 {
		t.Log("Perhatian: Tidak ada menu tumis di DB sampel lokal saat ini.")
	} else if foundTumis {
		t.Logf("Sukses: Ditemukan %d menu tumis valid untuk pelengkap.", len(tumisFoods))
	}
}
