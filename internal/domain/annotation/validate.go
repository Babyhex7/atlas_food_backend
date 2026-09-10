package annotation

import (
	"fmt"
	"net/http"
	"strings"

	"atlas_food/internal/pkg/utils"
)

// MinPolygonPoints - polygon minimal segitiga agar punya luas dan bisa di-hit-test
const MinPolygonPoints = 3

// MaxPolygonPoints - batas atas wajar; menahan payload autosave meledak
const MaxPolygonPoints = 500

// clamp - jepit nilai ke rentang [min, max]
func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// normalizeAreasForDraft - validasi longgar untuk autosave draft.
//
// Sengaja permisif: editor menyimpan setiap 1.5 detik, termasuk saat polygon
// masih digambar. Menolak payload di sini berarti kehilangan pekerjaan admin,
// jadi koordinat di luar kanvas di-clamp, bukan ditolak. Aturan ketat
// (minimal 3 titik, minimal 1 area) baru ditegakkan saat publish.
func normalizeAreasForDraft(areas []AreaInput, width, height int) ([]AreaInput, error) {
	w, h := float64(width), float64(height)

	for i := range areas {
		name := strings.TrimSpace(areas[i].Name)
		if name == "" {
			return nil, utils.NewAppError(http.StatusBadRequest, "VALIDATION_ERROR",
				fmt.Sprintf("Area ke-%d belum punya nama", i+1))
		}
		areas[i].Name = name

		if len(areas[i].Polygon) > MaxPolygonPoints {
			return nil, utils.NewAppError(http.StatusBadRequest, "VALIDATION_ERROR",
				fmt.Sprintf("Area %q punya %d titik, maksimal %d",
					name, len(areas[i].Polygon), MaxPolygonPoints))
		}

		// Kosongkan food_id string kosong agar tersimpan NULL, bukan '' —
		// string kosong akan melanggar foreign key ke foods(id).
		if areas[i].FoodID != nil && strings.TrimSpace(*areas[i].FoodID) == "" {
			areas[i].FoodID = nil
		}

		for j := range areas[i].Polygon {
			areas[i].Polygon[j][0] = clamp(areas[i].Polygon[j][0], 0, w)
			areas[i].Polygon[j][1] = clamp(areas[i].Polygon[j][1], 0, h)
		}
	}

	return areas, nil
}

// validateForPublish - aturan ketat sebelum anotasi boleh dibaca publik (§7.4)
func validateForPublish(image *FoodImage) error {
	incoming := make([]AreaInput, 0, len(image.Areas))
	for _, area := range image.Areas {
		incoming = append(incoming, AreaInput{Name: area.Name, Polygon: area.Polygon})
	}
	return validateAreasForPublish(image, incoming)
}

// validateAreasForPublish - aturan publish yang diterapkan ke kumpulan area
// tertentu, bukan hanya ke area yang sudah tersimpan.
//
// Dipakai dua tempat: saat publish dari draft, dan saat autosave menimpa
// gambar yang SUDAH published (perubahannya langsung tayang, jadi tidak boleh
// lolos dengan standar draft yang longgar).
func validateAreasForPublish(image *FoodImage, areas []AreaInput) error {
	if strings.TrimSpace(image.ImageURL) == "" {
		return utils.NewAppError(http.StatusUnprocessableEntity, "PUBLISH_INVALID",
			"Gambar belum punya image_url")
	}

	if image.Width <= 0 || image.Height <= 0 {
		return utils.NewAppError(http.StatusUnprocessableEntity, "PUBLISH_INVALID",
			"Dimensi gambar tidak valid")
	}

	if len(areas) == 0 {
		return utils.NewAppError(http.StatusUnprocessableEntity, "PUBLISH_INVALID",
			"Minimal 1 area harus dianotasi sebelum publish")
	}

	w, h := float64(image.Width), float64(image.Height)

	for _, area := range areas {
		if strings.TrimSpace(area.Name) == "" {
			return utils.NewAppError(http.StatusUnprocessableEntity, "PUBLISH_INVALID",
				"Ada area tanpa nama")
		}

		if len(area.Polygon) < MinPolygonPoints {
			return utils.NewAppError(http.StatusUnprocessableEntity, "PUBLISH_INVALID",
				fmt.Sprintf("Area %q hanya punya %d titik, minimal %d",
					area.Name, len(area.Polygon), MinPolygonPoints))
		}

		for _, pt := range area.Polygon {
			if pt.X() < 0 || pt.X() > w || pt.Y() < 0 || pt.Y() > h {
				return utils.NewAppError(http.StatusUnprocessableEntity, "PUBLISH_INVALID",
					fmt.Sprintf("Area %q punya titik di luar batas gambar (%.0f×%.0f)",
						area.Name, w, h))
			}
		}
	}

	return nil
}
