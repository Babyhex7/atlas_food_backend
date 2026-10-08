package ai

import "math"

// Status per zat gizi dan status keseluruhan. Nilainya menjadi kontrak dengan frontend.
const (
	StatusLow  = "low"
	StatusGood = "good"
	StatusHigh = "high"

	OverallGood   = "good"
	OverallLess   = "less"
	OverallExcess = "excess"
)

// Ambang klasifikasi tingkat kecukupan terhadap rujukan (persen).
// Di bawah lowerBound = kurang, di atas upperBound = lebih.
const (
	lowerBound = 80.0
	upperBound = 110.0
)

// fullDayMealCount - jumlah waktu makan minimal agar laporan dianggap mewakili satu hari
const fullDayMealCount = 3

// Intake - total asupan harian dari satu submission
type Intake struct {
	Energy  float64
	Protein float64
	Carbs   float64
	Fat     float64
}

// AssessmentItem - hasil perbandingan satu zat gizi terhadap rujukan
type AssessmentItem struct {
	Key       string  `json:"key"`
	Label     string  `json:"label"`
	Unit      string  `json:"unit"`
	Intake    float64 `json:"intake"`
	Reference float64 `json:"reference"`
	Percent   float64 `json:"percent"`
	Status    string  `json:"status"`
}

// Assess - bandingkan asupan dengan rujukan dan tentukan status tiap zat gizi.
//
// Perhitungan ini sengaja dilakukan di sini, bukan oleh LLM: model bahasa lemah
// dalam aritmetika dan hasilnya tidak dapat direproduksi. LLM hanya menerima
// hasil perhitungan ini sebagai fakta untuk dinarasikan.
func Assess(in Intake, ref Reference) []AssessmentItem {
	rows := []struct {
		key, label, unit  string
		intake, reference float64
	}{
		{"energy", "Energi", "kkal", in.Energy, ref.Energy},
		{"protein", "Protein", "g", in.Protein, ref.Protein},
		{"carbs", "Karbohidrat", "g", in.Carbs, ref.Carbs},
		{"fat", "Lemak", "g", in.Fat, ref.Fat},
	}

	items := make([]AssessmentItem, 0, len(rows))
	for _, r := range rows {
		// Dibulatkan lebih dulu supaya status selalu konsisten dengan angka yang
		// ditampilkan (110,0% tidak boleh tergolong "lebih" karena galat pecahan).
		percent := 0.0
		if r.reference > 0 {
			percent = round1(r.intake / r.reference * 100)
		}
		items = append(items, AssessmentItem{
			Key:       r.key,
			Label:     r.label,
			Unit:      r.unit,
			Intake:    round1(r.intake),
			Reference: r.reference,
			Percent:   percent,
			Status:    classify(percent),
		})
	}
	return items
}

// classify - petakan persen kecukupan ke status
func classify(percent float64) string {
	switch {
	case percent < lowerBound:
		return StatusLow
	case percent > upperBound:
		return StatusHigh
	default:
		return StatusGood
	}
}

// OverallStatus - status keseluruhan dari hasil penilaian.
//
// Energi menjadi penentu utama. Bila energi sudah sesuai, status baru bergeser
// kalau minimal dua dari tiga makronutrien menyimpang ke arah yang sama —
// satu zat gizi yang meleset tidak cukup untuk menyebut asupan sehari "kurang".
func OverallStatus(items []AssessmentItem) string {
	var low, high int
	for _, it := range items {
		if it.Key == "energy" {
			switch it.Status {
			case StatusLow:
				return OverallLess
			case StatusHigh:
				return OverallExcess
			}
			continue
		}
		switch it.Status {
		case StatusLow:
			low++
		case StatusHigh:
			high++
		}
	}
	switch {
	case low >= 2:
		return OverallLess
	case high >= 2:
		return OverallExcess
	default:
		return OverallGood
	}
}

// round1 - bulatkan ke satu desimal
func round1(v float64) float64 {
	return math.Round(v*10) / 10
}
