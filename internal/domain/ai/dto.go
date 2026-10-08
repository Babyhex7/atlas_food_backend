package ai

import "time"

// NutritionAnalysisRequest - request body endpoint AI
type NutritionAnalysisRequest struct {
	// SubmissionID boleh berisi id server maupun local_id (kunci idempotensi dari klien),
	// sehingga laporan yang dikirim saat luring tetap bisa dianalisis setelah tersinkron.
	SubmissionID string `json:"submission_id" binding:"required"`
	// ForceRefresh meminta analisis baru meskipun hasil lama sudah tersimpan.
	ForceRefresh bool `json:"force_refresh"`
}

// NutritionAnalysisItem - satu baris rincian gizi.
// Angka dan status berasal dari perhitungan server; hanya Description yang ditulis LLM.
type NutritionAnalysisItem struct {
	Key         string  `json:"key"`
	Label       string  `json:"label"`
	Status      string  `json:"status"`
	Description string  `json:"description"`
	Unit        string  `json:"unit"`
	Intake      float64 `json:"intake"`
	Reference   float64 `json:"reference"`
	Percent     float64 `json:"percent"`
}

// HealthInsight - insight kesehatan
type HealthInsight struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Coverage - seberapa lengkap laporan yang dianalisis.
// Dipakai frontend untuk menampilkan peringatan bila hasil berpotensi terlalu rendah.
type Coverage struct {
	MealCount        int  `json:"meal_count"`
	FoodCount        int  `json:"food_count"`
	MissingFoodCount int  `json:"missing_food_count"`
	IsPartialDay     bool `json:"is_partial_day"`
}

// NutritionAnalysisData - payload hasil analisis
type NutritionAnalysisData struct {
	OverallStatus       string                  `json:"overall_status"`
	OverallMessage      string                  `json:"overall_message"`
	NutritionalAnalysis []NutritionAnalysisItem `json:"nutritional_analysis"`
	AIRecommendation    string                  `json:"ai_recommendation"`
	RecommendedFoods    []string                `json:"recommended_foods"`
	HealthInsight       HealthInsight           `json:"health_insight"`
	SuggestedActivities []string                `json:"suggested_activities"`
	Reference           Reference               `json:"reference"`
	Coverage            Coverage                `json:"coverage"`
}

// AnalysisMeta - jejak asal hasil analisis
type AnalysisMeta struct {
	Model         string    `json:"model"`
	PromptVersion string    `json:"prompt_version"`
	GeneratedAt   time.Time `json:"generated_at"`
}

// NutritionAnalysisResult - hasil endpoint AI
type NutritionAnalysisResult struct {
	Source string                `json:"source"` // groq|cache
	Data   NutritionAnalysisData `json:"data"`
	Meta   AnalysisMeta          `json:"meta"`
}
