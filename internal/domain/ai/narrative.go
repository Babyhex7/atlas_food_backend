package ai

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Batas isi narasi LLM yang diteruskan ke klien
const (
	maxMessageLen        = 700
	maxDescriptionLen    = 400
	maxRecommendationLen = 1200
	maxListItemLen       = 80
	maxRecommendedFoods  = 6
	maxActivities        = 4
)

// narrative - bagian hasil analisis yang ditulis LLM
type narrative struct {
	OverallMessage      string
	Descriptions        map[string]string // key zat gizi → deskripsi
	Recommendation      string
	RecommendedFoods    []string
	HealthInsight       HealthInsight
	SuggestedActivities []string
}

// parseNarrative - urai jawaban LLM secara longgar lalu rapikan.
//
// Keluaran model diperlakukan sebagai masukan tidak tepercaya: medan bisa hilang,
// bertipe salah (string di tempat array atau sebaliknya), atau terlalu panjang.
// Mengembalikan false bila bagian inti (ringkasan dan rekomendasi) tidak ada,
// supaya pemanggil bisa mencoba ulang alih-alih menyimpan hasil kosong.
func parseNarrative(content string) (narrative, bool) {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(content), &raw); err != nil {
		return narrative{}, false
	}

	n := narrative{
		OverallMessage:      cleanText(asText(raw["overall_message"]), maxMessageLen),
		Recommendation:      cleanParagraphs(asText(raw["ai_recommendation"]), maxRecommendationLen),
		RecommendedFoods:    asTextList(raw["recommended_foods"], maxListItemLen),
		SuggestedActivities: asTextList(raw["suggested_activities"], maxListItemLen),
		Descriptions:        map[string]string{},
	}

	if insight, ok := raw["health_insight"].(map[string]interface{}); ok {
		n.HealthInsight = HealthInsight{
			Title:       cleanText(asText(insight["title"]), maxListItemLen),
			Description: cleanText(asText(insight["description"]), maxMessageLen),
		}
	}

	if items, ok := raw["nutritional_analysis"].([]interface{}); ok {
		for _, item := range items {
			obj, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			key := normalizeNutrientKey(asText(obj["key"]))
			if key == "" {
				key = normalizeNutrientKey(asText(obj["label"]))
			}
			if desc := cleanText(asText(obj["description"]), maxDescriptionLen); key != "" && desc != "" {
				n.Descriptions[key] = desc
			}
		}
	}

	return n, n.OverallMessage != "" && n.Recommendation != ""
}

// compose - gabungkan fakta deterministik dengan narasi LLM menjadi hasil akhir.
// Status, angka, dan label selalu diambil dari perhitungan server.
func compose(ctx analysisContext, n narrative) NutritionAnalysisData {
	items := make([]NutritionAnalysisItem, 0, len(ctx.Assessment))
	for _, a := range ctx.Assessment {
		desc := n.Descriptions[a.Key]
		if desc == "" {
			desc = fallbackDescription(a)
		}
		items = append(items, NutritionAnalysisItem{
			Key:         a.Key,
			Label:       a.Label,
			Status:      a.Status,
			Description: desc,
			Unit:        a.Unit,
			Intake:      a.Intake,
			Reference:   a.Reference,
			Percent:     a.Percent,
		})
	}

	return NutritionAnalysisData{
		OverallStatus:       ctx.Overall,
		OverallMessage:      n.OverallMessage,
		NutritionalAnalysis: items,
		AIRecommendation:    n.Recommendation,
		RecommendedFoods:    filterFoods(n.RecommendedFoods, ctx.EatenFoods),
		HealthInsight:       n.HealthInsight,
		SuggestedActivities: limit(dedupe(n.SuggestedActivities), maxActivities),
		Reference:           ctx.Reference,
		Coverage:            ctx.Coverage,
	}
}

// fallbackDescription - deskripsi baku bila LLM tidak memberi penjelasan untuk satu zat gizi
func fallbackDescription(a AssessmentItem) string {
	level := map[string]string{
		StatusLow:  "masih di bawah",
		StatusGood: "sudah sesuai dengan",
		StatusHigh: "melebihi",
	}[a.Status]
	return fmt.Sprintf("Asupan %s Anda sekitar %.0f%% dari rujukan, %s kebutuhan harian.",
		strings.ToLower(a.Label), a.Percent, level)
}

// normalizeNutrientKey - petakan key/label bebas dari LLM ke key baku
func normalizeNutrientKey(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "energy", "energi", "kalori", "calories", "calorie":
		return "energy"
	case "protein":
		return "protein"
	case "carbs", "carbohydrate", "carbohydrates", "karbohidrat":
		return "carbs"
	case "fat", "fats", "lemak":
		return "fat"
	}
	return ""
}

// filterFoods - buang saran yang sudah dimakan responden, hilangkan duplikat, batasi jumlah
func filterFoods(foods []string, eaten map[string]struct{}) []string {
	out := make([]string, 0, len(foods))
	for _, f := range dedupe(foods) {
		if _, already := eaten[strings.ToLower(f)]; already {
			continue
		}
		out = append(out, f)
	}
	return limit(out, maxRecommendedFoods)
}

// dedupe - hilangkan duplikat tanpa membedakan huruf besar/kecil, urutan dipertahankan
func dedupe(items []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(items))
	for _, item := range items {
		key := strings.ToLower(item)
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, item)
	}
	return out
}

// limit - ambil paling banyak n elemen pertama; selalu mengembalikan slice non-nil
func limit(items []string, n int) []string {
	if items == nil {
		return []string{}
	}
	if len(items) > n {
		return items[:n]
	}
	return items
}

// asText - ambil teks dari nilai JSON bebas; array string digabung per baris
func asText(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case []interface{}:
		parts := make([]string, 0, len(val))
		for _, p := range val {
			if s, ok := p.(string); ok && strings.TrimSpace(s) != "" {
				parts = append(parts, strings.TrimSpace(s))
			}
		}
		return strings.Join(parts, "\n")
	}
	return ""
}

// asTextList - ambil daftar teks dari nilai JSON bebas; string tunggal dipecah per baris/koma
func asTextList(v interface{}, maxLen int) []string {
	var rawItems []string
	switch val := v.(type) {
	case []interface{}:
		for _, p := range val {
			if s, ok := p.(string); ok {
				rawItems = append(rawItems, s)
			}
		}
	case string:
		rawItems = strings.FieldsFunc(val, func(r rune) bool { return r == '\n' || r == ',' || r == ';' })
	}

	out := make([]string, 0, len(rawItems))
	for _, item := range rawItems {
		item = strings.TrimLeft(strings.TrimSpace(item), "-•* ")
		if cleaned := cleanText(item, maxLen); cleaned != "" {
			out = append(out, cleaned)
		}
	}
	return out
}

// cleanParagraphs - seperti cleanText tetapi mempertahankan pemisah baris,
// karena rekomendasi boleh berbentuk beberapa butir.
func cleanParagraphs(raw string, maxLen int) string {
	lines := strings.Split(strings.ReplaceAll(raw, "\r\n", "\n"), "\n")
	kept := make([]string, 0, len(lines))
	for _, line := range lines {
		if cleaned := cleanText(line, maxLen); cleaned != "" {
			kept = append(kept, cleaned)
		}
	}
	joined := strings.Join(kept, "\n")
	if runes := []rune(joined); len(runes) > maxLen {
		joined = strings.TrimSpace(string(runes[:maxLen]))
	}
	return joined
}
