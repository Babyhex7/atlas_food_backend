package ai

import (
	"encoding/json"
	"fmt"
	"strings"
)

// buildAssistantSystemPrompt - menyusun system prompt berbasis step recall
func buildAssistantSystemPrompt(step AssistantStep) string {
	basePrompt := `Kamu adalah "Asisten AI Atlas Makananku", ahli gizi ramah yang bertindak sebagai enumerator digital dalam survei Food Recall 24 Jam di Indonesia.

BAHASA & GAYA KOMUNIKASI:
- Gunakan Bahasa Indonesia yang ramah, santun, lugas, dan mudah dimengerti orang awam (maksimal 2-3 kalimat penjelasan).
- Fokus membantu responden mengingat dan mencocokkan apa yang mereka makan secara akurat.
- DILARANG MENGARANG KALORI ATAU ZAT GIZI MAKRO (perhitungan gizi dilakukan terpisah oleh database TKPI). Fokus pada nama makanan, porsi URT, dan gramatur.

OUTPUT FORMAT:
Kembalikan HANYA format JSON valid tanpa tanda markdown (tanpa ` + "```json" + `).
Skema JSON:
{
  "explanation": "penjelasan ramah dan ringkas",
  "recommended_label": "A|B|C|D atau null",
  "recommended_gram": number atau null,
  "suggested_foods": ["nama makanan 1", "nama makanan 2"],
  "quick_replies": ["pertanyaan cepat 1", "pertanyaan cepat 2"]
}
`

	switch step {
	case AssistantStepPortion:
		return basePrompt + `
ATURAN KHUSUS LANGKAH 3 (ESTIMASI PORSI):
Tugasmu adalah mengonversi Ukuran Rumah Tangga (URT) yang disebut pengguna ke porsi foto yang tersedia (A/B/C/D) atau perkiraan gram terdekat.

STANDAR URT INDONESIA (Kemenkes RI):
- Nasi: 1 centong beras/rice cooker sedang ≈ 100g. 1 piring sedang rata ≈ 150-200g. 1 piring munjung ≈ 250-300g.
- Tempe: 1 potong kecil ≈ 15-20g. 1 potong sedang (warteg/gorengan) ≈ 25-30g. 2 potong sedang ≈ 50-60g.
- Tahu: 1 potong kecil ≈ 25-30g. 1 potong sedang ≈ 50-60g. 1 potong besar ≈ 100g.
- Ayam: 1 potong sedang (dada/paha atas tanpa tulang) ≈ 40-50g. Dengan tulang ≈ 80-100g.
- Telur: 1 butir telur ayam ras utuh ≈ 50-60g. Telur ceplok/dadar dengan minyak ≈ 55-65g.
- Ikan: 1 ekor sedang (kembung/lele) ≈ 60-80g daging bersih.
- Daging Sapi (rendang/semur): 1 potong sedang ≈ 35-50g.
- Sayur matang: 1 sendok sayur/kuah tiris ≈ 30-40g. 1 mangkok sedang berserat ≈ 100g.
- Minyak/Mentega: 1 sendok makan (sdm) ≈ 10g. 1 sendok teh (sdt) ≈ 5g.
- Gula pasir: 1 sdm ≈ 13g. 1 sdt ≈ 5g.
- Kerupuk kaleng putih: 1 buah ≈ 10-15g.

Jika daftar available_portions diberikan, COCOKKAN estimasi URT ke label (A, B, C, D) yang paling mendekati nilainya, dan isi "recommended_label" serta "recommended_gram".`

	case AssistantStepFood:
		return basePrompt + `
ATURAN KHUSUS LANGKAH 2 (DAFTAR MAKANAN / QUICK LIST):
Tugasmu adalah mengekstrak makanan-makanan yang dikonsumsi responden dari kalimat bahasa sehari-hari ke dalam daftar nama makanan Indonesia baku di field "suggested_foods".
Contoh: "tadi siang makan nasi padang pake rendang sama sayur daun singkong trus minum es teh manis"
Ekstrak menjadi: ["Nasi Putih", "Rendang Daging Sapi", "Sayur Daun Singkong Rebus", "Es Teh Manis"].
Pastikan tiap item spesifik agar mudah dicari di katalog.`

	case AssistantStepAdditional:
		return basePrompt + `
ATURAN KHUSUS LANGKAH 4 (BAHAN PELENGKAP & BUMBU):
Tugasmu adalah mengingatkan bahan pelengkap, bumbu, atau minyak tersembunyi yang lazim menyertai makanan yang sedang aktif.
Contoh:
- Soto / Rawon: sambal, kecap manis, kerupuk, perasan jeruk nipis, koya, emping.
- Bakso: saus tomat, sambal, kecap, pangsit goreng.
- Roti: mentega/margarin, selai, susu kental manis, keju.
- Salad / Gado-gado: bumbu kacang, kerupuk.
Sajikan daftar bumbu/pelengkap yang mungkin dikonsumsi di "suggested_foods" dan "quick_replies".`

	case AssistantStepReview:
		return basePrompt + `
ATURAN KHUSUS LANGKAH 5 (TINJAUAN & MAKANAN TERLEWAT):
Tugasmu adalah mengidentifikasi potensi makanan/minuman yang sering terlupakan oleh responden berdasarkan daftar makanan yang sudah tercatat.
Cek:
1. Apakah sudah ada air putih / minuman di tiap waktu makan?
2. Apakah ada camilan atau buah di antara jam makan utama?
3. Apakah ada suplemen atau vitamin yang dikonsumsi?
Tanyakan dengan ramah dan cantumkan contoh di "suggested_foods".`

	default:
		return basePrompt
	}
}

// buildAssistantUserPrompt - menyusun prompt input dari data request
func buildAssistantUserPrompt(req AssistantRequest) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Step aktif: %s\n", req.Step))
	if req.FoodName != "" {
		sb.WriteString(fmt.Sprintf("Makanan yang sedang dipilih: %s\n", req.FoodName))
	}
	if req.MealType != "" {
		sb.WriteString(fmt.Sprintf("Waktu makan: %s\n", req.MealType))
	}
	if len(req.ContextFoods) > 0 {
		sb.WriteString(fmt.Sprintf("Makanan lain yang sudah dicatat pada sesi ini: %s\n", strings.Join(req.ContextFoods, ", ")))
	}

	if len(req.AvailablePortions) > 0 {
		sb.WriteString("Foto porsi yang tersedia di sistem:\n")
		for _, p := range req.AvailablePortions {
			desc := p.Description
			if desc == "" {
				desc = "-"
			}
			sb.WriteString(fmt.Sprintf("- Porsi %s: %.0fg (URT/deskripsi: %s)\n", p.Label, p.WeightGram, desc))
		}
	}

	sb.WriteString(fmt.Sprintf("\nPertanyaan / Input Responden:\n\"%s\"", req.UserQuery))
	return sb.String()
}

// fallbackHeuristicAssistant - fallback cerdas jika Groq API sedang offline atau tidak ada token
func fallbackHeuristicAssistant(req AssistantRequest) *AssistantResponse {
	queryLower := strings.ToLower(req.UserQuery)

	switch req.Step {
	case AssistantStepPortion:
		// Jika ada available portions, coba cocokkan porsi terdekat
		if len(req.AvailablePortions) > 0 {
			// Deteksi porsi kecil/sedang/besar dari query
			var targetLabel string
			var targetGram float64

			if strings.Contains(queryLower, "kecil") || strings.Contains(queryLower, "sedikit") || strings.Contains(queryLower, "1/2") || strings.Contains(queryLower, "setengah") {
				targetLabel = req.AvailablePortions[0].Label
				targetGram = req.AvailablePortions[0].WeightGram
			} else if strings.Contains(queryLower, "besar") || strings.Contains(queryLower, "banyak") || strings.Contains(queryLower, "munjung") || strings.Contains(queryLower, "jumbo") {
				lastIdx := len(req.AvailablePortions) - 1
				targetLabel = req.AvailablePortions[lastIdx].Label
				targetGram = req.AvailablePortions[lastIdx].WeightGram
			} else {
				// Default ke porsi tengah (sedang)
				midIdx := len(req.AvailablePortions) / 2
				targetLabel = req.AvailablePortions[midIdx].Label
				targetGram = req.AvailablePortions[midIdx].WeightGram
			}

			foodName := req.FoodName
			if foodName == "" {
				foodName = "makanan ini"
			}

			explanation := fmt.Sprintf("Berdasarkan ukuran umum, porsi sedang untuk %s paling mendekati Porsi %s (%.0fg).", foodName, targetLabel, targetGram)
			return &AssistantResponse{
				Explanation:      explanation,
				RecommendedLabel: targetLabel,
				RecommendedGram:  &targetGram,
				QuickReplies: []string{
					"Bagaimana kalau 1 centong?",
					"Porsi mana yang paling kecil?",
					"Saya makan 2 porsi ini",
				},
			}
		}

		return &AssistantResponse{
			Explanation: "Silakan pilih salah satu foto porsi yang paling mirip dengan jumlah yang Anda santap kemarin.",
			QuickReplies: []string{
				"Porsi normal yang mana?",
				"Saya makan sedikit saja",
			},
		}

	case AssistantStepFood:
		return &AssistantResponse{
			Explanation: "Silakan ketik nama makanan yang Anda konsumsi di kolom pencarian di sebelah kiri untuk menambahkannya ke daftar.",
			QuickReplies: []string{
				"Nasi Putih",
				"Ayam Goreng",
				"Tempe Goreng",
				"Air Putih",
			},
		}

	case AssistantStepAdditional:
		return &AssistantResponse{
			Explanation:    "Jangan lupa mencatat minyak goreng, sambal, kecap, atau saus pelengkap jika Anda menggunakannya kemarin.",
			SuggestedFoods: []string{"Sambal", "Kecap Manis", "Kerupuk", "Minyak Goreng"},
			QuickReplies:   []string{"Saya pakai sambal", "Tidak ada bumbu tambahan", "Goreng pakai minyak"},
		}

	case AssistantStepReview:
		return &AssistantResponse{
			Explanation:    "Pastikan air putih atau minuman selingan sudah tercatat di setiap waktu makan Anda ya.",
			SuggestedFoods: []string{"Air Putih", "Teh Hangat", "Kopi", "Buah"},
			QuickReplies:   []string{"Semua sudah lengkap", "Ada camilan terlewat"},
		}

	default:
		return &AssistantResponse{
			Explanation: "Ada yang bisa saya bantu terkait pencatatan makanan Anda?",
		}
	}
}

// parseAssistantResponse - parse string JSON mentah dari Groq ke struct response
func parseAssistantResponse(rawContent string) (*AssistantResponse, error) {
	var resp AssistantResponse
	if err := json.Unmarshal([]byte(rawContent), &resp); err != nil {
		return nil, fmt.Errorf("gagal unmarshal JSON asisten: %w (konten: %s)", err, rawContent)
	}
	return &resp, nil
}
