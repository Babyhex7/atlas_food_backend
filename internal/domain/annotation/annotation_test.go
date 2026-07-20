package annotation

import (
	"encoding/json"
	"testing"

	"atlas_food/internal/pkg/utils"
)

// TestPolygonRoundTrip - polygon harus selamat pulang-pergi ke kolom JSON MySQL.
// Driver MySQL bisa mengembalikan []byte maupun string, keduanya harus jalan.
func TestPolygonRoundTrip(t *testing.T) {
	original := Polygon{{120, 180.5}, {160, 170}, {210, 210}}

	value, err := original.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}

	raw, ok := value.(string)
	if !ok {
		t.Fatalf("Value() harus string, dapat %T", value)
	}

	for name, scanned := range map[string]interface{}{
		"string": raw,
		"bytes":  []byte(raw),
		"nil":    nil,
		"empty":  []byte{},
	} {
		var target Polygon
		if err := target.Scan(scanned); err != nil {
			t.Errorf("Scan(%s) error: %v", name, err)
			continue
		}

		if name == "nil" || name == "empty" {
			if len(target) != 0 {
				t.Errorf("Scan(%s) harus menghasilkan polygon kosong, dapat %v", name, target)
			}
			continue
		}

		if len(target) != len(original) {
			t.Errorf("Scan(%s): panjang %d, harusnya %d", name, len(target), len(original))
			continue
		}
		for i := range original {
			if target[i] != original[i] {
				t.Errorf("Scan(%s) titik %d: %v, harusnya %v", name, i, target[i], original[i])
			}
		}
	}
}

// TestPolygonNilValue - polygon nil harus tersimpan sebagai array kosong,
// bukan NULL: kolom polygon NOT NULL.
func TestPolygonNilValue(t *testing.T) {
	var p Polygon

	value, err := p.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}
	if value != "[]" {
		t.Errorf("polygon nil harus jadi \"[]\", dapat %v", value)
	}
}

// TestPolygonJSONShape - kontrak API: polygon adalah array [x,y], bukan objek.
func TestPolygonJSONShape(t *testing.T) {
	area := FoodArea{Name: "Dada Ayam", Polygon: Polygon{{1, 2}, {3, 4}, {5, 6}}}

	encoded, err := json.Marshal(area)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var decoded struct {
		Polygon [][]float64 `json:"polygon"`
	}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if len(decoded.Polygon) != 3 || decoded.Polygon[0][0] != 1 || decoded.Polygon[0][1] != 2 {
		t.Errorf("bentuk polygon salah: %s", encoded)
	}
}

// TestNormalizeAreasForDraft - autosave harus permisif: koordinat di luar
// kanvas dijepit, bukan ditolak, agar pekerjaan admin tidak hilang.
func TestNormalizeAreasForDraft(t *testing.T) {
	empty := ""
	areas := []AreaInput{{
		Name:    "  Sayap  ",
		FoodID:  &empty,
		Polygon: Polygon{{-50, 20}, {2000, 900}, {100, 100}},
	}}

	got, err := normalizeAreasForDraft(areas, 1200, 800)
	if err != nil {
		t.Fatalf("normalizeAreasForDraft error: %v", err)
	}

	if got[0].Name != "Sayap" {
		t.Errorf("nama harus di-trim, dapat %q", got[0].Name)
	}

	// food_id string kosong harus jadi nil, kalau tidak melanggar FK ke foods(id)
	if got[0].FoodID != nil {
		t.Errorf("food_id string kosong harus jadi nil, dapat %v", *got[0].FoodID)
	}

	want := Polygon{{0, 20}, {1200, 800}, {100, 100}}
	for i := range want {
		if got[0].Polygon[i] != want[i] {
			t.Errorf("titik %d: %v, harusnya dijepit ke %v", i, got[0].Polygon[i], want[i])
		}
	}
}

// TestNormalizeAreasRejectsBlankName - satu-satunya alasan autosave menolak
func TestNormalizeAreasRejectsBlankName(t *testing.T) {
	_, err := normalizeAreasForDraft([]AreaInput{{Name: "   "}}, 100, 100)
	if err == nil {
		t.Fatal("area tanpa nama harus ditolak")
	}
}

// TestValidateForPublish - aturan ketat brief §7.4 baru berlaku saat publish
func TestValidateForPublish(t *testing.T) {
	base := func(areas []FoodArea) *FoodImage {
		return &FoodImage{ImageURL: "/uploads/annotations/a.jpg", Width: 100, Height: 100, Areas: areas}
	}

	cases := []struct {
		name    string
		image   *FoodImage
		wantErr bool
	}{
		{
			name:    "tanpa area",
			image:   base(nil),
			wantErr: true,
		},
		{
			name:    "polygon kurang dari 3 titik",
			image:   base([]FoodArea{{Name: "A", Polygon: Polygon{{1, 1}, {2, 2}}}}),
			wantErr: true,
		},
		{
			name:    "titik di luar batas",
			image:   base([]FoodArea{{Name: "A", Polygon: Polygon{{1, 1}, {2, 2}, {500, 3}}}}),
			wantErr: true,
		},
		{
			name:    "area tanpa nama",
			image:   base([]FoodArea{{Name: " ", Polygon: Polygon{{1, 1}, {2, 2}, {3, 3}}}}),
			wantErr: true,
		},
		{
			name:    "valid",
			image:   base([]FoodArea{{Name: "Dada", Polygon: Polygon{{1, 1}, {2, 2}, {3, 3}}}}),
			wantErr: false,
		},
		{
			name:    "dimensi nol",
			image:   &FoodImage{ImageURL: "/a.jpg", Width: 0, Height: 100},
			wantErr: true,
		},
		{
			name:    "tanpa image_url",
			image:   &FoodImage{Width: 10, Height: 10, Areas: []FoodArea{{Name: "A", Polygon: Polygon{{1, 1}, {2, 2}, {3, 3}}}}},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateForPublish(tc.image)

			if tc.wantErr && err == nil {
				t.Fatal("harusnya ditolak, tapi lolos")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("harusnya lolos, tapi ditolak: %v", err)
			}

			// Penolakan publish harus jadi 422, bukan 500
			if tc.wantErr {
				appErr, ok := err.(*utils.AppError)
				if !ok {
					t.Fatalf("error harus *utils.AppError, dapat %T", err)
				}
				if appErr.StatusCode != 422 {
					t.Errorf("status harus 422, dapat %d", appErr.StatusCode)
				}
			}
		})
	}
}
