package gemini

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Client - client untuk berinteraksi dengan Google Gemini dan Imagen API
type Client struct {
	apiKey      string
	hfAPIKey    string
	model       string
	imagenModel string
	httpClient  *http.Client
}

// FoodAIAnalysis - hasil estimasi gizi dan porsi dari Gemini
type FoodAIAnalysis struct {
	FoodName        string           `json:"food_name"`
	LocalName       string           `json:"local_name"`
	Description     string           `json:"description"`
	CategoryName    string           `json:"category_name"`
	PhotoType       string           `json:"photo_type"` // series, range, guide
	CaloriesPer100g float64          `json:"calories_per_100g"`
	ProteinPer100g  float64          `json:"protein_per_100g"`
	FatPer100g      float64          `json:"fat_per_100g"`
	CarbsPer100g    float64          `json:"carbs_per_100g"`
	Portions        []PortionConcept `json:"portions"`
}

// PortionConcept - konsep porsi A/B/C/D
type PortionConcept struct {
	Code         string  `json:"code"` // A, B, C, D
	Label        string  `json:"label"` // Kecil, Sedang, Besar, Jumbo
	WeightGram   float64 `json:"weight_gram"`
	Description  string  `json:"description"`
	VisualPrompt string  `json:"visual_prompt"`
	ImageURL     string  `json:"image_url,omitempty"`
}

// NewClient - inisialisasi Gemini client
func NewClient(apiKey, model, imagenModel string) *Client {
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}
	if model == "" {
		model = os.Getenv("GEMINI_MODEL")
		if model == "" {
			model = "gemini-3.1-flash-lite"
		}
	}
	if imagenModel == "" {
		imagenModel = os.Getenv("IMAGEN_MODEL")
		if imagenModel == "" {
			imagenModel = "gemini-3.1-flash-image"
		}
	}

	return &Client{
		apiKey:      apiKey,
		hfAPIKey:    strings.TrimSpace(os.Getenv("HF_API_KEY")),
		model:       model,
		imagenModel: imagenModel,
		httpClient:  &http.Client{Timeout: 45 * time.Second},
	}
}

// EstimateFoodNutritionAndPortions - menganalisis makanan Indonesia berdasarkan standar TKPI Kemenkes
func (c *Client) EstimateFoodNutritionAndPortions(ctx context.Context, foodName, description string) (*FoodAIAnalysis, error) {
	// Jika tidak ada API key, gunakan deterministic fallback berbasis knowledge base pangan Indonesia
	if c.apiKey == "" {
		return c.fallbackEstimation(foodName, description), nil
	}

	systemPrompt := `Kamu adalah Ahli Gizi Senior spesialis database gizi Indonesia (TKPI / FCT Kemenkes) dan Food Photography Director.
Tugasmu: Diberikan nama makanan Indonesia (dan deskripsi opsional), hasilkan analisis gizi per 100g dan 4 tingkatan porsi piring standar (Series A, B, C, D) dalam format JSON valid.

Aturan Porsi:
- Porsi A: Kecil (mengisi ~25% piring 23cm)
- Porsi B: Sedang / Standar (mengisi ~50% piring 23cm)
- Porsi C: Besar (mengisi ~75% piring 23cm)
- Porsi D: Sangat Besar / Jumbo (mengisi ~95% piring 23cm membukit)

Aturan Tipe Foto (photo_type):
- "series": makanan olahan campur, hidangan berkuah, butiran, atau makanan yang disajikan dengan sendok/centong (contoh: nasi goreng, capcay, seblak, kolak, bubur, sup).
- "range": bahan tunggal alami atau makanan satuan yang disajikan utuh dan memiliki variasi ukuran fisik kecil/sedang/besar (contoh: ikan kembung/lele ekoran, pisang, paha ayam, tempe/tahu satuan).
- "guide": makanan yang memiliki beragam variasi bentuk potongan geometris (dadu, kipas, balok) atau aneka varian berbeda dalam 1 kategori (contoh: pepaya potong aneka bentuk, aneka kerupuk, aneka biskuit).

Format Output WAJIB JSON persis seperti ini tanpa markdown tambahan:
{
  "food_name": "Nama Baku",
  "local_name": "Nama Daerah / Populer",
  "description": "Deskripsi singkat komposisi bahan utama",
  "category_name": "Makanan Pokok / Lauk Hewani / Lauk Nabati / Sayuran / Buah / Snack / Minuman",
  "photo_type": "series / range / guide",
  "calories_per_100g": 180.5,
  "protein_per_100g": 15.2,
  "fat_per_100g": 8.4,
  "carbs_per_100g": 12.1,
  "portions": [
    {
      "code": "A",
      "label": "Kecil",
      "weight_gram": 75.0,
      "description": "Porsi kecil 1 potong kecil",
      "visual_prompt": "Scientific nutritional food portion study photo, top-down 45-degree angle, placed on a plain circular 23cm white matte ceramic plate, centered on neutral light-slate background, soft diffused overhead lighting, standard stainless steel spoon and fork on sides for scale reference. A small portion of [food] filling 25% of plate area, crisp culinary texture, 8k resolution, realistic."
    },
    {
      "code": "B",
      "label": "Sedang",
      "weight_gram": 140.0,
      "description": "Porsi sedang standar 1 porsi saji",
      "visual_prompt": "Scientific nutritional food portion study photo, top-down 45-degree angle, placed on a plain circular 23cm white matte ceramic plate, centered on neutral light-slate background, soft diffused overhead lighting, standard stainless steel spoon and fork on sides for scale reference. A medium portion of [food] filling 50% of plate area, crisp culinary texture, 8k resolution, realistic."
    },
    {
      "code": "C",
      "label": "Besar",
      "weight_gram": 220.0,
      "description": "Porsi besar 1.5 porsi saji",
      "visual_prompt": "Scientific nutritional food portion study photo, top-down 45-degree angle, placed on a plain circular 23cm white matte ceramic plate, centered on neutral light-slate background, soft diffused overhead lighting, standard stainless steel spoon and fork on sides for scale reference. A generous large portion of [food] filling 75% of plate area, crisp culinary texture, 8k resolution, realistic."
    },
    {
      "code": "D",
      "label": "Sangat Besar",
      "weight_gram": 310.0,
      "description": "Porsi jumbo ekstra mengenyangkan",
      "visual_prompt": "Scientific nutritional food portion study photo, top-down 45-degree angle, placed on a plain circular 23cm white matte ceramic plate, centered on neutral light-slate background, soft diffused overhead lighting, standard stainless steel spoon and fork on sides for scale reference. An extra-large heavy heaping portion of [food] filling 95% of plate area, tall conical mound, crisp culinary texture, 8k resolution, realistic."
    }
  ]
}`

	userPrompt := fmt.Sprintf("Analisis makanan ini:\nNama: %s\nDeskripsi: %s", foodName, description)

	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role": "user",
				"parts": []map[string]interface{}{
					{"text": systemPrompt + "\n\n" + userPrompt},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"responseMimeType": "application/json",
			"temperature":      0.2,
		},
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", c.model, c.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return c.fallbackEstimation(foodName, description), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.fallbackEstimation(foodName, description), nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.fallbackEstimation(foodName, description), nil
	}

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(bodyBytes, &geminiResp); err != nil || len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return c.fallbackEstimation(foodName, description), nil
	}

	textOutput := geminiResp.Candidates[0].Content.Parts[0].Text
	textOutput = strings.TrimPrefix(textOutput, "```json")
	textOutput = strings.TrimPrefix(textOutput, "```")
	textOutput = strings.TrimSuffix(textOutput, "```")
	textOutput = strings.TrimSpace(textOutput)

	var analysis FoodAIAnalysis
	if err := json.Unmarshal([]byte(textOutput), &analysis); err != nil {
		return c.fallbackEstimation(foodName, description), nil
	}

	return &analysis, nil
}

// GeneratePortionImage - panggil Gemini/Imagen API untuk render gambar porsi, fallback ke Pollinations gratis, lalu SVG placeholder
func (c *Client) GeneratePortionImage(ctx context.Context, prompt, foodName, portionLabel string, weightGram float64) (string, error) {
	if c.apiKey != "" {
		// Jika model adalah Gemini native image (mis. gemini-3.1-flash-image), gunakan :generateContent
		if strings.HasPrefix(strings.ToLower(c.imagenModel), "gemini-") {
			endpoint := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", c.imagenModel, c.apiKey)
			reqPayload := map[string]interface{}{
				"contents": []map[string]interface{}{
					{
						"role":  "user",
						"parts": []map[string]interface{}{{"text": prompt}},
					},
				},
				"generationConfig": map[string]interface{}{
					"responseModalities": []string{"IMAGE"},
				},
			}
			jsonBytes, _ := json.Marshal(reqPayload)
			req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonBytes))
			if err == nil {
				req.Header.Set("Content-Type", "application/json")
				resp, err := c.httpClient.Do(req)
				if err == nil {
					if resp.StatusCode == http.StatusOK {
						bodyBytes, _ := io.ReadAll(resp.Body)
						resp.Body.Close()
						var gemImgResp struct {
							Candidates []struct {
								Content struct {
									Parts []struct {
										InlineData *struct {
											MimeType string `json:"mimeType"`
											Data     string `json:"data"`
										} `json:"inlineData"`
									} `json:"parts"`
								} `json:"content"`
							} `json:"candidates"`
						}
						if err := json.Unmarshal(bodyBytes, &gemImgResp); err == nil {
							for _, cand := range gemImgResp.Candidates {
								for _, part := range cand.Content.Parts {
									if part.InlineData != nil && part.InlineData.Data != "" {
										mime := part.InlineData.MimeType
										if mime == "" {
											mime = "image/jpeg"
										}
										return fmt.Sprintf("data:%s;base64,%s", mime, part.InlineData.Data), nil
									}
								}
							}
						}
					} else {
						resp.Body.Close()
					}
				}
			}
		} else {
			// Endpoint legacy Google Imagen :predict
			endpoint := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:predict?key=%s", c.imagenModel, c.apiKey)
			reqPayload := map[string]interface{}{
				"instances": []map[string]interface{}{
					{"prompt": prompt},
				},
				"parameters": map[string]interface{}{
					"sampleCount":      1,
					"aspectRatio":      "1:1",
					"outputMimeType":   "image/jpeg",
					"personGeneration": "dont_allow",
					"safetySetting":    "block_medium_and_above",
				},
			}
			jsonBytes, _ := json.Marshal(reqPayload)
			req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonBytes))
			if err == nil {
				req.Header.Set("Content-Type", "application/json")
				resp, err := c.httpClient.Do(req)
				if err == nil {
					if resp.StatusCode == http.StatusOK {
						bodyBytes, _ := io.ReadAll(resp.Body)
						resp.Body.Close()
						var imagenResp struct {
							Predictions []struct {
								BytesBase64Encoded string `json:"bytesBase64Encoded"`
							} `json:"predictions"`
						}
						if err := json.Unmarshal(bodyBytes, &imagenResp); err == nil && len(imagenResp.Predictions) > 0 && imagenResp.Predictions[0].BytesBase64Encoded != "" {
							return fmt.Sprintf("data:image/jpeg;base64,%s", imagenResp.Predictions[0].BytesBase64Encoded), nil
						}
					} else {
						resp.Body.Close()
					}
				}
			}
		}
	}

	// Tier 2 Fallback: HuggingFace Inference Router (FLUX.1-schnell via fal-ai) jika HF_API_KEY tersedia
	if c.hfAPIKey != "" {
		if hfURL, err := c.generateHuggingFaceImage(ctx, prompt); err == nil && hfURL != "" {
			return hfURL, nil
		}
	}

	// Tier 3 Fallback: Pollinations.ai
	if polURL, err := c.generatePollinationsImage(ctx, prompt, portionLabel); err == nil && polURL != "" {
		return polURL, nil
	}

	// Tier 4 Last-Resort Fallback: SVG piring standar
	return generateSVGPlatePlaceholder(foodName, portionLabel, weightGram), nil
}

// generateHuggingFaceImage - generate foto JPEG via HuggingFace Inference Router (fal-ai/flux/schnell)
func (c *Client) generateHuggingFaceImage(ctx context.Context, prompt string) (string, error) {
	if strings.TrimSpace(prompt) == "" {
		return "", fmt.Errorf("prompt kosong")
	}
	endpoint := "https://router.huggingface.co/fal-ai/fal-ai/flux/schnell"
	payload := map[string]interface{}{
		"prompt":     prompt,
		"image_size": "square",
	}
	jsonBytes, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.hfAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("hf HTTP %d", resp.StatusCode)
	}

	ctype := strings.ToLower(resp.Header.Get("Content-Type"))
	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 5*1024*1024))
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(ctype, "image/") && len(bodyBytes) >= 1024 {
		return fmt.Sprintf("data:%s;base64,%s", strings.Split(ctype, ";")[0], base64.StdEncoding.EncodeToString(bodyBytes)), nil
	}

	// Beberapa provider mengembalikan JSON { images: [{ url: "..." }] }
	var falResp struct {
		Images []struct {
			URL string `json:"url"`
		} `json:"images"`
	}
	if err := json.Unmarshal(bodyBytes, &falResp); err == nil && len(falResp.Images) > 0 && falResp.Images[0].URL != "" {
		return falResp.Images[0].URL, nil
	}
	return "", fmt.Errorf("hf unexpected response")
}

// generatePollinationsImage - generate foto JPEG realistis gratis via image.pollinations.ai
func (c *Client) generatePollinationsImage(ctx context.Context, prompt, portionLabel string) (string, error) {
	if strings.TrimSpace(prompt) == "" {
		return "", fmt.Errorf("prompt kosong")
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(prompt + ":" + portionLabel))
	seed := h.Sum32() % 100000

	imgURL := fmt.Sprintf(
		"https://image.pollinations.ai/prompt/%s?width=512&height=512&nologo=true&seed=%d",
		url.PathEscape(prompt),
		seed,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", imgURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "AtlasFood-AI/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("pollinations HTTP %d", resp.StatusCode)
	}

	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(strings.ToLower(ctype), "image/") {
		return "", fmt.Errorf("unexpected content-type: %s", ctype)
	}

	imgBytes, err := io.ReadAll(io.LimitReader(resp.Body, 5*1024*1024))
	if err != nil || len(imgBytes) < 1024 {
		return "", fmt.Errorf("invalid image payload")
	}

	encoded := base64.StdEncoding.EncodeToString(imgBytes)
	return fmt.Sprintf("data:%s;base64,%s", strings.Split(ctype, ";")[0], encoded), nil
}

func containsWord(haystack, word string) bool {
	for _, token := range strings.Fields(haystack) {
		if strings.Trim(token, ".,!?()[]{}\"'") == word {
			return true
		}
	}
	return false
}

// fallbackEstimation - estimasi berbasis heuristik gizi jika API sedang offline / belum ada key
func (c *Client) fallbackEstimation(foodName, description string) *FoodAIAnalysis {
	nameLower := strings.ToLower(foodName)
	category := "Lauk Hewani"
	cal := 210.0
	prot := 18.5
	fat := 12.0
	carb := 5.0
	baseWeight := 120.0

	if strings.Contains(nameLower, "nasi") || strings.Contains(nameLower, "mie") || strings.Contains(nameLower, "bubur") {
		category = "Makanan Pokok"
		cal = 160.0
		prot = 3.5
		fat = 2.0
		carb = 34.0
		baseWeight = 150.0
	} else if strings.Contains(nameLower, "sayur") || strings.Contains(nameLower, "sop") || strings.Contains(nameLower, "tumis") {
		category = "Sayuran"
		cal = 65.0
		prot = 2.0
		fat = 3.5
		carb = 6.0
		baseWeight = 100.0
	} else if strings.Contains(nameLower, "tahu") || strings.Contains(nameLower, "tempe") {
		category = "Lauk Nabati"
		cal = 175.0
		prot = 14.0
		fat = 8.5
		carb = 11.0
		baseWeight = 80.0
	} else if containsWord(nameLower, "es") || containsWord(nameLower, "teh") || containsWord(nameLower, "jus") || containsWord(nameLower, "kopi") {
		category = "Minuman"
		cal = 90.0
		prot = 0.5
		fat = 0.2
		carb = 22.0
		baseWeight = 200.0
	}

	photoType := "series"
	if strings.Contains(nameLower, "ikan") || strings.Contains(nameLower, "pisang") || strings.Contains(nameLower, "paha") || strings.Contains(nameLower, "ayam goreng") || strings.Contains(nameLower, "telur") {
		photoType = "range"
	} else if strings.Contains(nameLower, "potong") || strings.Contains(nameLower, "kerupuk") || strings.Contains(nameLower, "biskuit") || strings.Contains(nameLower, "kue") {
		photoType = "guide"
	}

	wA := baseWeight * 0.6
	wB := baseWeight * 1.0
	wC := baseWeight * 1.6
	wD := baseWeight * 2.3

	masterPrompt := func(label string, fillPct int) string {
		return fmt.Sprintf("Scientific nutritional food portion study photo, top-down 45-degree angle, plain circular 23cm white matte ceramic plate, neutral light-slate background, soft studio lighting, stainless steel spoon and fork on sides for scale reference. Indonesian %s occupying %d%% of plate surface, %s portion, crisp culinary texture, 8k resolution, photorealistic.", foodName, fillPct, label)
	}

	return &FoodAIAnalysis{
		FoodName:        foodName,
		LocalName:       foodName,
		Description:     description,
		CategoryName:    category,
		PhotoType:       photoType,
		CaloriesPer100g: cal,
		ProteinPer100g:  prot,
		FatPer100g:      fat,
		CarbsPer100g:    carb,
		Portions: []PortionConcept{
			{
				Code:         "A",
				Label:        "Kecil",
				WeightGram:   wA,
				Description:  fmt.Sprintf("Porsi kecil (~%.0fg)", wA),
				VisualPrompt: masterPrompt("Kecil", 25),
			},
			{
				Code:         "B",
				Label:        "Sedang",
				WeightGram:   wB,
				Description:  fmt.Sprintf("Porsi sedang standar (~%.0fg)", wB),
				VisualPrompt: masterPrompt("Sedang", 50),
			},
			{
				Code:         "C",
				Label:        "Besar",
				WeightGram:   wC,
				Description:  fmt.Sprintf("Porsi besar (~%.0fg)", wC),
				VisualPrompt: masterPrompt("Besar", 75),
			},
			{
				Code:         "D",
				Label:        "Sangat Besar",
				WeightGram:   wD,
				Description:  fmt.Sprintf("Porsi jumbo (~%.0fg)", wD),
				VisualPrompt: masterPrompt("Sangat Besar", 95),
			},
		},
	}
}

// generateSVGPlatePlaceholder - menghasilkan visual piring keramik 23cm + makanan berbobot dalam SVG base64
func generateSVGPlatePlaceholder(foodName, label string, weight float64) string {
	fillRadius := 45
	switch label {
	case "A":
		fillRadius = 40
	case "B":
		fillRadius = 65
	case "C":
		fillRadius = 90
	case "D":
		fillRadius = 115
	}

	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 400 400" width="100%%" height="100%%">
  <defs>
    <radialGradient id="plateShadow" cx="50%%" cy="50%%" r="50%%">
      <stop offset="85%%" stop-color="#cbd5e1" stop-opacity="0.8"/>
      <stop offset="100%%" stop-color="#94a3b8" stop-opacity="0"/>
    </radialGradient>
    <radialGradient id="plateInner" cx="45%%" cy="45%%" r="55%%">
      <stop offset="0%%" stop-color="#ffffff"/>
      <stop offset="90%%" stop-color="#f1f5f9"/>
      <stop offset="100%%" stop-color="#e2e8f0"/>
    </radialGradient>
    <linearGradient id="foodGrad" x1="0%%" y1="0%%" x2="100%%" y2="100%%">
      <stop offset="0%%" stop-color="#ea580c"/>
      <stop offset="50%%" stop-color="#c2410c"/>
      <stop offset="100%%" stop-color="#9a3412"/>
    </linearGradient>
  </defs>
  <!-- Background Studio Slate -->
  <rect width="400" height="400" fill="#f8fafc"/>
  
  <!-- Scale Cutlery Fork & Spoon -->
  <!-- Fork Left -->
  <rect x="25" y="120" width="8" height="160" rx="4" fill="#94a3b8"/>
  <rect x="20" y="100" width="18" height="30" rx="2" fill="#94a3b8"/>
  <!-- Spoon Right -->
  <rect x="367" y="120" width="8" height="160" rx="4" fill="#94a3b8"/>
  <ellipse cx="371" cy="115" rx="14" ry="24" fill="#94a3b8"/>

  <!-- Plate Shadow -->
  <circle cx="200" cy="205" r="160" fill="url(#plateShadow)"/>
  
  <!-- Outer Plate Rim 23cm standard -->
  <circle cx="200" cy="200" r="150" fill="url(#plateInner)" stroke="#cbd5e1" stroke-width="3"/>
  <!-- Inner Plate Depth -->
  <circle cx="200" cy="200" r="115" fill="#fcfcfc" stroke="#e2e8f0" stroke-width="2"/>
  
  <!-- Food Portion Mound -->
  <circle cx="200" cy="200" r="%d" fill="url(#foodGrad)" opacity="0.95" filter="drop-shadow(0px 4px 6px rgba(0,0,0,0.25))"/>
  
  <!-- Scale Indicator Label -->
  <rect x="130" y="20" width="140" height="30" rx="6" fill="#0f172a" opacity="0.85"/>
  <text x="200" y="40" fill="#ffffff" font-family="Arial, sans-serif" font-size="13" font-weight="bold" text-anchor="middle">
    Porsi %s (%.0fg)
  </text>
  
  <!-- Food Name -->
  <text x="200" y="375" fill="#475569" font-family="Arial, sans-serif" font-size="14" font-weight="600" text-anchor="middle">
    %s
  </text>
  <text x="200" y="392" fill="#94a3b8" font-family="Arial, sans-serif" font-size="10" text-anchor="middle">
    Visual Guide AI · 23cm Ceramic Plate
  </text>
</svg>`, fillRadius, label, weight, foodName)

	encoded := base64.StdEncoding.EncodeToString([]byte(svg))
	return fmt.Sprintf("data:image/svg+xml;base64,%s", encoded)
}
