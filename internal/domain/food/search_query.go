package food

import "strings"

// maxSearchTermLen - batas panjang satu kata kunci agar query tetap wajar
const maxSearchTermLen = 64

// booleanModeOperators - karakter yang punya arti khusus di MySQL
// "MATCH ... AGAINST (... IN BOOLEAN MODE)".
//
// Input user masuk mentah ke ekspresi itu, jadi tanda kutip yang tidak
// berpasangan atau operator yang menggantung (mis. mencari `ayam "goreng`
// atau `nasi +`) membuat MySQL melempar error sintaks 1064 — seluruh query
// gagal, termasuk cabang LIKE-nya, dan pencarian makanan balas 500.
const booleanModeOperators = `+-><()~*:"&|@`

// indonesianFoodSynonyms - pemetaan ejaan lokal/kolokial dan typo umum pada makanan Indonesia
var indonesianFoodSynonyms = map[string]string{
	"buncic":   "buncis",
	"telor":    "telur",
	"sambel":   "sambal",
	"cabe":     "cabai",
	"duren":    "durian",
	"krupuk":   "kerupuk",
	"indomie":  "mie",
	"popmie":   "mie",
	"tempeh":   "tempe",
	"semangko": "semangka",
	"semangke": "semangka",
}

// normalizeFoodSearchTerm - periksa kata terhadap kamus variasi dan typo umum
func normalizeFoodSearchTerm(w string) string {
	lower := strings.ToLower(w)
	if repl, exists := indonesianFoodSynonyms[lower]; exists {
		return repl
	}
	// Typo umum akhiran 'c' untuk kata berakhiran 's' seperti buncic -> buncis
	if strings.HasSuffix(lower, "c") && len(lower) >= 4 {
		candidate := lower[:len(lower)-1] + "s"
		if candidate == "buncis" {
			return candidate
		}
	}
	return lower
}

// normalizeSearchQuery - normalisasi kata-kata dalam query pencarian user
func normalizeSearchQuery(raw string) string {
	words := strings.Fields(raw)
	if len(words) == 0 {
		return ""
	}
	changed := false
	result := make([]string, len(words))
	for i, w := range words {
		norm := normalizeFoodSearchTerm(w)
		if norm != strings.ToLower(w) {
			changed = true
		}
		result[i] = norm
	}
	if !changed {
		return raw
	}
	return strings.Join(result, " ")
}

// sanitizeFulltextQuery - bersihkan input user menjadi ekspresi BOOLEAN MODE
// yang aman.
//
// Setiap kata dijadikan prefix-match (`nasi*`) supaya perilaku "ketik sebagian
// nama" tetap sama seperti sebelumnya. Mengembalikan string kosong kalau tidak
// ada kata yang tersisa; pemanggil harus jatuh ke pencarian LIKE saja.
func sanitizeFulltextQuery(raw string) string {
	cleaned := strings.Map(func(r rune) rune {
		if strings.ContainsRune(booleanModeOperators, r) {
			return ' '
		}
		return r
	}, raw)

	words := strings.Fields(cleaned)
	terms := make([]string, 0, len(words))
	for _, w := range words {
		if len(w) > maxSearchTermLen {
			w = w[:maxSearchTermLen]
		}
		terms = append(terms, w+"*")
	}
	return strings.Join(terms, " ")
}

// escapeLikePattern - lolos-kan wildcard LIKE (% dan _) supaya kata kunci
// diperlakukan sebagai teks biasa. Tanpa ini, mengetik "%" mencocokkan semua
// baris dan "_" mencocokkan karakter apa pun.
func escapeLikePattern(raw string) string {
	replacer := strings.NewReplacer(`\`, `\`, `%`, `\%`, `_`, `\_`)
	return "%" + replacer.Replace(raw) + "%"
}
