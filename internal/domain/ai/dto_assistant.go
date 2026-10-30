package ai

// AssistantStep - langkah recall yang sedang aktif
type AssistantStep string

const (
	AssistantStepMeal       AssistantStep = "meal"       // Step 1: Waktu & sesi makan
	AssistantStepFood       AssistantStep = "food"       // Step 2: Pencarian / ekstraksi makanan
	AssistantStepPortion    AssistantStep = "portion"    // Step 3: Estimasi porsi (URT to photo/gram)
	AssistantStepAdditional AssistantStep = "additional" // Step 4: Bahan pelengkap & bumbu
	AssistantStepReview     AssistantStep = "review"     // Step 5: Tinjauan makanan terlewat
)

// AvailablePortion - daftar porsi foto yang tersedia di katalog makanan
type AvailablePortion struct {
	Label       string  `json:"label"`                  // "A", "B", "C", "D"
	WeightGram  float64 `json:"weight_gram"`            // gramatur
	Description string  `json:"description,omitempty"`  // deskripsi URT misal "1 potong kecil"
}

// AssistantRequest - request dari frontend untuk chatbot per-step
type AssistantRequest struct {
	Step              AssistantStep      `json:"step" binding:"required"`
	FoodName          string             `json:"food_name,omitempty"`
	UserQuery         string             `json:"user_query" binding:"required"`
	AvailablePortions []AvailablePortion `json:"available_portions,omitempty"`
	ContextFoods      []string           `json:"context_foods,omitempty"` // Daftar makanan yang sudah dicatat
	MealType          string             `json:"meal_type,omitempty"`     // "breakfast", "lunch", "dinner", "snack"
}

// AssistantResponse - response AI terstruktur untuk frontend
type AssistantResponse struct {
	Explanation      string   `json:"explanation"`                 // Penjelasan ramah natural
	RecommendedLabel string   `json:"recommended_label,omitempty"` // Label porsi terdekat misal "B"
	RecommendedGram  *float64 `json:"recommended_gram,omitempty"`  // Gramatur terdekat misal 50
	SuggestedFoods   []string `json:"suggested_foods,omitempty"`   // Rekomendasi makanan yang diekstrak/pelengkap
	QuickReplies     []string `json:"quick_replies,omitempty"`     // Pilihan quick prompt selanjutnya
}
