package bootstrap

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"gorm.io/gorm"
)

// ─── Format file data/tkpi/tkpi_2020.json ─────────────────────────────────────

// TKPIFile - isi file TKPI 2020 (nilai per 100 g BDD).
type TKPIFile struct {
	Meta TKPIMeta    `json:"meta"`
	Data []TKPIEntry `json:"data"`
}

type TKPIMeta struct {
	Judul  string `json:"judul"`
	Basis  string `json:"basis"`
	Jumlah int    `json:"jumlah"`
}

// TKPIEntry - satu baris TKPI. Gizi memakai pointer supaya null (tidak
// dianalisis) bisa dibedakan dari 0 (tidak terdeteksi).
type TKPIEntry struct {
	Kode     string              `json:"kode"`
	Nama     string              `json:"nama"`
	Sumber   string              `json:"sumber"`
	Kelompok string              `json:"kelompok"`
	Jenis    string              `json:"jenis"`
	Gizi     map[string]*float64 `json:"gizi"`
}

// ─── Format file data/tkpi/atlas_tkpi_mapping.json ────────────────────────────

const (
	MappingStatusCocok    = "cocok"
	MappingStatusPerluCek = "perlu_cek"
)

type TKPIMappingFile struct {
	Mapping []TKPIMapping `json:"mapping"`
}

type TKPIMapping struct {
	AtlasKode string `json:"atlas_kode"`
	AtlasNama string `json:"atlas_nama"`
	TKPIKode  string `json:"tkpi_kode"`
	Status    string `json:"status"`
	Alasan    string `json:"alasan,omitempty"`
}

// ─── Jenis nutrien ────────────────────────────────────────────────────────────

// tkpiNutrientDef - satu kolom TKPI beserta nutrient_types yang menampungnya.
type tkpiNutrientDef struct {
	TKPIKey      string // kunci di objek "gizi"
	Code         string // nutrient_types.code
	Name         string
	UnitCode     string // nutrient_units.code
	DisplayOrder int
}

// tkpiNutrients - kode 1–10 (kecuali sugar) sudah ada sejak migrasi 003;
// sisanya ditambahkan agar seluruh kolom TKPI tersimpan. "bdd" bukan zat gizi
// sehingga tidak disimpan. "sugar" tidak ada di TKPI dan tidak disentuh.
var tkpiNutrients = []tkpiNutrientDef{
	{"energi", "energy", "Energi", "kcal", 1},
	{"protein", "protein", "Protein", "g", 2},
	{"kh", "carbs", "Karbohidrat", "g", 3},
	{"lemak", "fat", "Lemak Total", "g", 4},
	{"serat", "fiber", "Serat", "g", 5},
	{"natrium", "sodium", "Natrium", "mg", 7},
	{"kalsium", "calcium", "Kalsium", "mg", 8},
	{"besi", "iron", "Zat Besi", "mg", 9},
	{"vit_c", "vit_c", "Vitamin C", "mg", 10},
	{"air", "water", "Air", "g", 11},
	{"abu", "ash", "Abu", "g", 12},
	{"fosfor", "phosphorus", "Fosfor", "mg", 13},
	{"kalium", "potassium", "Kalium", "mg", 14},
	{"tembaga", "copper", "Tembaga", "mg", 15},
	{"seng", "zinc", "Seng", "mg", 16},
	{"retinol", "retinol", "Retinol (Vit. A)", "mcg", 17},
	{"b_kar", "beta_carotene", "Beta-Karoten", "mcg", 18},
	{"kar_total", "total_carotene", "Karoten Total", "mcg", 19},
	{"thiamin", "thiamin", "Thiamin (Vit. B1)", "mg", 20},
	{"riboflavin", "riboflavin", "Riboflavin (Vit. B2)", "mg", 21},
	{"niasin", "niacin", "Niasin", "mg", 22},
}

var tkpiUnits = map[string][2]string{ // code -> {name, symbol}
	"kcal": {"Kilokalori", "kcal"},
	"g":    {"Gram", "g"},
	"mg":   {"Miligram", "mg"},
	"mcg":  {"Mikrogram", "μg"},
}

// ─── Fungsi murni (diuji tanpa database) ──────────────────────────────────────

// LoadTKPIFile - baca & validasi file TKPI.
func LoadTKPIFile(path string) (*TKPIFile, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca file TKPI: %w", err)
	}
	var f TKPIFile
	if err := json.Unmarshal(raw, &f); err != nil {
		return nil, fmt.Errorf("format JSON TKPI tidak valid: %w", err)
	}
	if err := ValidateTKPI(&f); err != nil {
		return nil, err
	}
	return &f, nil
}

// ValidateTKPI - pastikan kode unik, tidak kosong, dan tidak ada nilai gizi negatif.
func ValidateTKPI(f *TKPIFile) error {
	if len(f.Data) == 0 {
		return fmt.Errorf("file TKPI tidak berisi data")
	}
	seen := make(map[string]bool, len(f.Data))
	for _, e := range f.Data {
		if e.Kode == "" {
			return fmt.Errorf("ada entri TKPI tanpa kode (nama: %q)", e.Nama)
		}
		if seen[e.Kode] {
			return fmt.Errorf("kode TKPI duplikat: %s", e.Kode)
		}
		seen[e.Kode] = true
		for key, v := range e.Gizi {
			if v != nil && *v < 0 {
				return fmt.Errorf("nilai %s negatif pada %s", key, e.Kode)
			}
		}
	}
	return nil
}

// LoadTKPIMapping - baca file pemetaan Atlas → TKPI.
func LoadTKPIMapping(path string) (*TKPIMappingFile, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca file mapping: %w", err)
	}
	var m TKPIMappingFile
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("format JSON mapping tidak valid: %w", err)
	}
	return &m, nil
}

// SelectMappings - pilih baris mapping yang akan diterapkan dan pastikan
// kode TKPI-nya benar-benar ada. Baris "perlu_cek" hanya ikut bila diminta.
// Satu makanan Atlas tidak boleh dipetakan dua kali.
func SelectMappings(m *TKPIMappingFile, tkpi map[string]TKPIEntry, includeUnverified bool) ([]TKPIMapping, error) {
	var selected []TKPIMapping
	seenAtlas := make(map[string]bool)

	for _, row := range m.Mapping {
		switch row.Status {
		case MappingStatusCocok:
		case MappingStatusPerluCek:
			if !includeUnverified {
				continue
			}
		default:
			return nil, fmt.Errorf("status mapping %s tidak dikenal: %q", row.AtlasKode, row.Status)
		}

		if row.AtlasKode == "" || row.TKPIKode == "" {
			return nil, fmt.Errorf("baris mapping tidak lengkap: %+v", row)
		}
		if seenAtlas[row.AtlasKode] {
			return nil, fmt.Errorf("makanan Atlas %s dipetakan lebih dari sekali", row.AtlasKode)
		}
		if _, ok := tkpi[row.TKPIKode]; !ok {
			return nil, fmt.Errorf("mapping %s merujuk kode TKPI %s yang tidak ada di file TKPI", row.AtlasKode, row.TKPIKode)
		}
		seenAtlas[row.AtlasKode] = true
		selected = append(selected, row)
	}
	return selected, nil
}

// NutrientValuesFromTKPI - ubah objek gizi TKPI menjadi map nutrient_types.code → nilai.
// Nilai null dilewati (tidak dianalisis ≠ 0); 0 tetap disimpan.
func NutrientValuesFromTKPI(e TKPIEntry) map[string]float64 {
	out := make(map[string]float64)
	for _, def := range tkpiNutrients {
		if v, ok := e.Gizi[def.TKPIKey]; ok && v != nil {
			out[def.Code] = *v
		}
	}
	return out
}

// ─── Seeder ───────────────────────────────────────────────────────────────────

// TKPISeedOptions - lokasi file & opsi seeding TKPI.
type TKPISeedOptions struct {
	TKPIPath          string
	MappingPath       string
	IncludeUnverified bool
}

// DefaultTKPISeedOptions - opsi dari env, dengan path default di data/tkpi/.
func DefaultTKPISeedOptions() TKPISeedOptions {
	opts := TKPISeedOptions{
		TKPIPath:          os.Getenv("TKPI_JSON_PATH"),
		MappingPath:       os.Getenv("TKPI_MAPPING_PATH"),
		IncludeUnverified: os.Getenv("TKPI_INCLUDE_UNVERIFIED") == "true",
	}
	if opts.TKPIPath == "" {
		opts.TKPIPath = "data/tkpi/tkpi_2020.json"
	}
	if opts.MappingPath == "" {
		opts.MappingPath = "data/tkpi/atlas_tkpi_mapping.json"
	}
	return opts
}

// PrintTKPISeedResult - ringkasan hasil seeding untuk CLI.
func PrintTKPISeedResult(res *TKPISeedResult) {
	fmt.Printf("   - Entri TKPI dibaca      : %d\n", res.TKPIEntries)
	fmt.Printf("   - Makanan diisi gizinya  : %d %v\n", len(res.Applied), res.Applied)
	fmt.Printf("   - Nilai gizi tersimpan   : %d\n", res.NutrientsSaved)
	if len(res.MissingFoods) > 0 {
		fmt.Printf("   ⚠️  Belum ada di tabel foods (seed Atlas dulu): %v\n", res.MissingFoods)
	}
}

// TKPISeedResult - ringkasan untuk ditampilkan di CLI.
type TKPISeedResult struct {
	TKPIEntries    int
	Applied        []string // atlas_kode yang gizinya diisi
	MissingFoods   []string // atlas_kode di mapping yang belum ada di tabel foods
	NutrientsSaved int
}

// SeedTKPINutrients - isi food_nutrients dari TKPI untuk makanan yang dipetakan.
//
// Idempoten: dijalankan ulang menghasilkan data yang sama. Untuk setiap makanan
// yang dipetakan, nilai lama pada jenis nutrien yang dicakup TKPI dihapus lalu
// ditulis ulang dari TKPI — sehingga nilai yang di TKPI null tidak tertinggal
// dari data contoh lama. Nutrien di luar TKPI (mis. sugar) tidak disentuh.
func SeedTKPINutrients(db *gorm.DB, opts TKPISeedOptions) (*TKPISeedResult, error) {
	tkpiFile, err := LoadTKPIFile(opts.TKPIPath)
	if err != nil {
		return nil, err
	}
	mappingFile, err := LoadTKPIMapping(opts.MappingPath)
	if err != nil {
		return nil, err
	}

	index := make(map[string]TKPIEntry, len(tkpiFile.Data))
	for _, e := range tkpiFile.Data {
		index[e.Kode] = e
	}

	mappings, err := SelectMappings(mappingFile, index, opts.IncludeUnverified)
	if err != nil {
		return nil, err
	}

	result := &TKPISeedResult{TKPIEntries: len(tkpiFile.Data)}

	err = db.Transaction(func(tx *gorm.DB) error {
		typeIDs, err := ensureTKPINutrientTypes(tx)
		if err != nil {
			return err
		}
		trackedIDs := make([]int, 0, len(typeIDs))
		for _, id := range typeIDs {
			trackedIDs = append(trackedIDs, id)
		}

		for _, row := range mappings {
			var food struct{ ID string }
			res := tx.Table("foods").Select("id").Where("code = ?", row.AtlasKode).Limit(1).Find(&food)
			if res.Error != nil {
				return fmt.Errorf("gagal mencari makanan %s: %w", row.AtlasKode, res.Error)
			}
			if res.RowsAffected == 0 {
				result.MissingFoods = append(result.MissingFoods, row.AtlasKode)
				continue
			}

			if err := tx.Table("food_nutrients").
				Where("food_id = ? AND nutrient_type_id IN ?", food.ID, trackedIDs).
				Delete(nil).Error; err != nil {
				return fmt.Errorf("gagal menghapus gizi lama %s: %w", row.AtlasKode, err)
			}

			values := NutrientValuesFromTKPI(index[row.TKPIKode])
			codes := make([]string, 0, len(values))
			for code := range values {
				codes = append(codes, code)
			}
			sort.Strings(codes) // urutan insert stabil

			rows := make([]map[string]interface{}, 0, len(codes))
			for _, code := range codes {
				rows = append(rows, map[string]interface{}{
					"food_id":          food.ID,
					"nutrient_type_id": typeIDs[code],
					"value_per_100g":   values[code],
				})
			}
			if len(rows) > 0 {
				if err := tx.Table("food_nutrients").Create(&rows).Error; err != nil {
					return fmt.Errorf("gagal menyimpan gizi %s: %w", row.AtlasKode, err)
				}
			}

			result.Applied = append(result.Applied, row.AtlasKode)
			result.NutrientsSaved += len(rows)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ensureTKPINutrientTypes - pastikan satuan & jenis nutrien TKPI ada, lalu
// kembalikan map code → id. Baris yang sudah ada tidak diubah (nama/urutan
// yang diatur admin tetap dipertahankan).
func ensureTKPINutrientTypes(tx *gorm.DB) (map[string]int, error) {
	unitIDs := make(map[string]int)
	for code, meta := range tkpiUnits {
		id, err := findOrCreateID(tx, "nutrient_units", code, map[string]interface{}{
			"code": code, "name": meta[0], "symbol": meta[1],
		})
		if err != nil {
			return nil, fmt.Errorf("satuan %s: %w", code, err)
		}
		unitIDs[code] = id
	}

	typeIDs := make(map[string]int, len(tkpiNutrients))
	for _, def := range tkpiNutrients {
		id, err := findOrCreateID(tx, "nutrient_types", def.Code, map[string]interface{}{
			"code":          def.Code,
			"name":          def.Name,
			"unit_id":       unitIDs[def.UnitCode],
			"display_order": def.DisplayOrder,
			"is_active":     true,
		})
		if err != nil {
			return nil, fmt.Errorf("jenis nutrien %s: %w", def.Code, err)
		}
		typeIDs[def.Code] = id
	}
	return typeIDs, nil
}

func findOrCreateID(tx *gorm.DB, table, code string, row map[string]interface{}) (int, error) {
	var found struct{ ID int }
	res := tx.Table(table).Select("id").Where("code = ?", code).Limit(1).Find(&found)
	if res.Error != nil {
		return 0, res.Error
	}
	if res.RowsAffected > 0 {
		return found.ID, nil
	}
	if err := tx.Table(table).Create(row).Error; err != nil {
		return 0, err
	}
	res = tx.Table(table).Select("id").Where("code = ?", code).Limit(1).Find(&found)
	if res.Error != nil {
		return 0, res.Error
	}
	if res.RowsAffected == 0 {
		return 0, fmt.Errorf("baris %s.%s tidak ditemukan setelah dibuat", table, code)
	}
	return found.ID, nil
}
