package collab

import (
	"os"
	"strings"
)

// Saklar MODE UJI untuk eksperimen artikel (ablation & batching). Hanya aktif bila
// COLLAB_E2E=1; tanpa itu seluruh saklar inert dan perilaku produksi tidak berubah.
// BERKAS INI TIDAK UNTUK DI-COMMIT (lihat atlas_food_artikel/skrip/instrumentasi-*.patch).
var e2eMode = os.Getenv("COLLAB_E2E") == "1"

// e2eNoBatch - room berawalan "e2e-nobatch" mematikan batching + coalescing (kondisi pembanding E1).
func e2eNoBatch(roomID string) bool {
	return e2eMode && strings.HasPrefix(roomID, "e2e-nobatch")
}
