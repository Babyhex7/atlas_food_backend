package ai

import (
	"fmt"
	"time"
)

// Reference - angka rujukan kebutuhan harian yang dipakai sebagai pembanding asupan
type Reference struct {
	Label        string  `json:"label"`
	Source       string  `json:"source"`
	Personalized bool    `json:"personalized"`
	Energy       float64 `json:"energy_kcal"`
	Protein      float64 `json:"protein_g"`
	Carbs        float64 `json:"carbs_g"`
	Fat          float64 `json:"fat_g"`
}

// Profile - data responden yang relevan untuk memilih rujukan.
// Sengaja tidak memuat nama, email, atau pengenal lain: hanya dua atribut ini
// yang dikirim ke penyedia LLM.
type Profile struct {
	Sex string `json:"sex,omitempty"` // male|female
	Age int    `json:"age,omitempty"` // tahun; 0 = tidak diketahui
}

const referenceSourceAKG = "Angka Kecukupan Gizi (Permenkes No. 28 Tahun 2019)"

// defaultReference - rujukan umum bila jenis kelamin atau usia responden tidak diketahui
var defaultReference = Reference{
	Label:   "Rujukan umum dewasa",
	Source:  referenceSourceAKG,
	Energy:  2150,
	Protein: 60,
	Carbs:   320,
	Fat:     65,
}

type akgRow struct {
	minAge, maxAge              int
	energy, protein, carbs, fat float64
}

// akgTable - AKG 2019 untuk energi dan makronutrien per kelompok usia.
// Kelompok di bawah 16 tahun tidak dimuat: sistem ditujukan bagi responden dewasa,
// dan di luar rentang ini rujukan umum yang dipakai.
var akgTable = map[string][]akgRow{
	"male": {
		{16, 18, 2650, 75, 400, 85},
		{19, 29, 2650, 65, 430, 75},
		{30, 49, 2550, 65, 415, 70},
		{50, 64, 2150, 65, 340, 60},
		{65, 80, 1800, 64, 275, 50},
		{81, 200, 1600, 64, 235, 45},
	},
	"female": {
		{16, 18, 2100, 65, 300, 70},
		{19, 29, 2250, 60, 360, 65},
		{30, 49, 2150, 60, 340, 60},
		{50, 64, 1800, 60, 280, 50},
		{65, 80, 1550, 58, 230, 45},
		{81, 200, 1400, 58, 200, 40},
	},
}

var sexLabel = map[string]string{"male": "laki-laki", "female": "perempuan"}

// ResolveReference - pilih rujukan AKG sesuai profil; jatuh ke rujukan umum bila profil tidak lengkap
func ResolveReference(p Profile) Reference {
	rows, ok := akgTable[p.Sex]
	if !ok || p.Age <= 0 {
		return defaultReference
	}
	for _, row := range rows {
		if p.Age >= row.minAge && p.Age <= row.maxAge {
			ageLabel := fmt.Sprintf("%d–%d tahun", row.minAge, row.maxAge)
			if row.maxAge >= 200 {
				ageLabel = fmt.Sprintf("di atas %d tahun", row.minAge-1)
			}
			return Reference{
				Label:        fmt.Sprintf("AKG %s %s", sexLabel[p.Sex], ageLabel),
				Source:       referenceSourceAKG,
				Personalized: true,
				Energy:       row.energy,
				Protein:      row.protein,
				Carbs:        row.carbs,
				Fat:          row.fat,
			}
		}
	}
	return defaultReference
}

// ageAt - usia dalam tahun penuh pada tanggal tertentu; 0 bila tanggal lahir tidak masuk akal
func ageAt(birth, at time.Time) int {
	if birth.IsZero() || birth.After(at) {
		return 0
	}
	age := at.Year() - birth.Year()
	if at.Month() < birth.Month() || (at.Month() == birth.Month() && at.Day() < birth.Day()) {
		age--
	}
	if age < 0 || age > 120 {
		return 0
	}
	return age
}
