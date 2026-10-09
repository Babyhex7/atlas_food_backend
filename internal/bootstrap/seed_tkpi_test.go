package bootstrap

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	tkpiPath    = filepath.Join("..", "..", "data", "tkpi", "tkpi_2020.json")
	mappingPath = filepath.Join("..", "..", "data", "tkpi", "atlas_tkpi_mapping.json")
	atlasPath   = filepath.Join("..", "..", "Atlas_Makananku_FINAL.json")
)

func f64(v float64) *float64 { return &v }

func loadIndex(t *testing.T) map[string]TKPIEntry {
	t.Helper()
	f, err := LoadTKPIFile(tkpiPath)
	if err != nil {
		t.Fatalf("LoadTKPIFile: %v", err)
	}
	idx := make(map[string]TKPIEntry, len(f.Data))
	for _, e := range f.Data {
		idx[e.Kode] = e
	}
	return idx
}

func TestTKPIFile_MetaMatchesData(t *testing.T) {
	f, err := LoadTKPIFile(tkpiPath)
	if err != nil {
		t.Fatalf("LoadTKPIFile: %v", err)
	}
	if f.Meta.Jumlah != len(f.Data) {
		t.Errorf("meta.jumlah = %d, jumlah data = %d", f.Meta.Jumlah, len(f.Data))
	}
	for _, e := range f.Data {
		for _, def := range tkpiNutrients {
			if _, ok := e.Gizi[def.TKPIKey]; !ok {
				t.Errorf("%s tidak memiliki kunci gizi %q", e.Kode, def.TKPIKey)
			}
		}
	}
}

func TestTKPIFile_SpotCheckNasi(t *testing.T) {
	nasi, ok := loadIndex(t)["AP001"]
	if !ok {
		t.Fatal("AP001 (Nasi) tidak ada")
	}
	want := map[string]float64{"energy": 180, "protein": 3.0, "fat": 0.3, "carbs": 39.8, "fiber": 0.2, "calcium": 25}
	got := NutrientValuesFromTKPI(nasi)
	for code, v := range want {
		if got[code] != v {
			t.Errorf("Nasi %s = %v, want %v", code, got[code], v)
		}
	}
}

func TestNutrientValuesFromTKPI_NullSkippedZeroKept(t *testing.T) {
	e := TKPIEntry{Kode: "X", Gizi: map[string]*float64{
		"energi":    f64(100),
		"retinol":   f64(0),
		"kar_total": nil,
		"bdd":       f64(90),
	}}
	got := NutrientValuesFromTKPI(e)

	if v, ok := got["retinol"]; !ok || v != 0 {
		t.Errorf("retinol 0 harus disimpan, got %v (ada=%v)", v, ok)
	}
	if _, ok := got["total_carotene"]; ok {
		t.Error("kar_total null tidak boleh disimpan")
	}
	if _, ok := got["bdd"]; ok {
		t.Error("bdd bukan zat gizi dan tidak boleh disimpan")
	}
	if len(got) != 2 {
		t.Errorf("jumlah nilai = %d, want 2: %v", len(got), got)
	}
}

func TestValidateTKPI_RejectsBadData(t *testing.T) {
	cases := map[string]TKPIFile{
		"kosong":     {},
		"duplikat":   {Data: []TKPIEntry{{Kode: "A"}, {Kode: "A"}}},
		"tanpa kode": {Data: []TKPIEntry{{Nama: "x"}}},
		"negatif":    {Data: []TKPIEntry{{Kode: "A", Gizi: map[string]*float64{"energi": f64(-1)}}}},
	}
	for name, f := range cases {
		if err := ValidateTKPI(&f); err == nil {
			t.Errorf("%s: harus error", name)
		}
	}
}

func TestSelectMappings_RealFile(t *testing.T) {
	idx := loadIndex(t)
	m, err := LoadTKPIMapping(mappingPath)
	if err != nil {
		t.Fatalf("LoadTKPIMapping: %v", err)
	}

	verified, err := SelectMappings(m, idx, false)
	if err != nil {
		t.Fatalf("SelectMappings(verified): %v", err)
	}
	all, err := SelectMappings(m, idx, true)
	if err != nil {
		t.Fatalf("SelectMappings(all): %v", err)
	}
	if len(verified) == 0 || len(all) <= len(verified) {
		t.Errorf("verified=%d all=%d — baris perlu_cek harus disaring secara default", len(verified), len(all))
	}
	for _, row := range verified {
		if row.Status != MappingStatusCocok {
			t.Errorf("%s berstatus %s tapi ikut terpilih tanpa IncludeUnverified", row.AtlasKode, row.Status)
		}
	}
}

// Setiap atlas_kode di mapping harus ada di Atlas_Makananku_FINAL.json dengan
// nama yang sama — mencegah gizi masuk ke makanan yang salah.
func TestMapping_AtlasCodesExist(t *testing.T) {
	raw, err := os.ReadFile(atlasPath)
	if err != nil {
		t.Fatalf("baca Atlas JSON: %v", err)
	}
	var atlas []AtlasFoodItem
	if err := json.Unmarshal(raw, &atlas); err != nil {
		t.Fatalf("parse Atlas JSON: %v", err)
	}
	names := make(map[string]string, len(atlas))
	for _, a := range atlas {
		names[a.Kode] = a.NamaID
	}

	m, err := LoadTKPIMapping(mappingPath)
	if err != nil {
		t.Fatalf("LoadTKPIMapping: %v", err)
	}
	for _, row := range m.Mapping {
		name, ok := names[row.AtlasKode]
		if !ok {
			t.Errorf("atlas_kode %s tidak ada di Atlas_Makananku_FINAL.json", row.AtlasKode)
			continue
		}
		if !strings.EqualFold(name, row.AtlasNama) {
			t.Errorf("atlas_kode %s: nama di mapping %q ≠ nama Atlas %q", row.AtlasKode, row.AtlasNama, name)
		}
	}
}

func TestSelectMappings_RejectsBadRows(t *testing.T) {
	idx := map[string]TKPIEntry{"AP001": {Kode: "AP001"}}
	cases := map[string][]TKPIMapping{
		"kode TKPI tidak ada": {{AtlasKode: "MP-01", TKPIKode: "ZZ999", Status: MappingStatusCocok}},
		"atlas ganda": {
			{AtlasKode: "MP-01", TKPIKode: "AP001", Status: MappingStatusCocok},
			{AtlasKode: "MP-01", TKPIKode: "AP001", Status: MappingStatusCocok},
		},
		"status asing": {{AtlasKode: "MP-01", TKPIKode: "AP001", Status: "mungkin"}},
	}
	for name, rows := range cases {
		if _, err := SelectMappings(&TKPIMappingFile{Mapping: rows}, idx, true); err == nil {
			t.Errorf("%s: harus error", name)
		}
	}
}
