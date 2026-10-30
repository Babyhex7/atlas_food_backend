package ai

import (
	"context"
	"testing"
)

func TestBuildAssistantSystemPrompt(t *testing.T) {
	steps := []AssistantStep{
		AssistantStepPortion,
		AssistantStepFood,
		AssistantStepAdditional,
		AssistantStepReview,
	}

	for _, step := range steps {
		prompt := buildAssistantSystemPrompt(step)
		if len(prompt) == 0 {
			t.Errorf("prompt untuk step %s kosong", step)
		}
	}
}

func TestBuildAssistantUserPrompt(t *testing.T) {
	req := AssistantRequest{
		Step:     AssistantStepPortion,
		FoodName: "Tempe Goreng Tepung",
		MealType: "lunch",
		AvailablePortions: []AvailablePortion{
			{Label: "A", WeightGram: 25, Description: "1 potong kecil"},
			{Label: "B", WeightGram: 50, Description: "2 potong sedang"},
		},
		UserQuery: "Saya makan 2 potong sedang, pilih yang mana?",
	}

	userPrompt := buildAssistantUserPrompt(req)
	if len(userPrompt) == 0 {
		t.Fatal("userPrompt kosong")
	}
}

func TestFallbackHeuristicAssistant(t *testing.T) {
	req := AssistantRequest{
		Step:     AssistantStepPortion,
		FoodName: "Tempe Goreng Tepung",
		AvailablePortions: []AvailablePortion{
			{Label: "A", WeightGram: 25},
			{Label: "B", WeightGram: 50},
			{Label: "C", WeightGram: 100},
		},
		UserQuery: "Saya makan 2 potong sedang",
	}

	resp := fallbackHeuristicAssistant(req)
	if resp == nil {
		t.Fatal("response fallback nil")
	}
	if resp.RecommendedLabel != "B" {
		t.Errorf("diharapkan label B, didapat %s", resp.RecommendedLabel)
	}
	if resp.RecommendedGram == nil || *resp.RecommendedGram != 50 {
		t.Errorf("diharapkan 50g, didapat %v", resp.RecommendedGram)
	}
}

func TestParseAssistantResponse(t *testing.T) {
	rawJSON := `{
		"explanation": "2 potong tempe sedang setara 50 gram (Porsi B).",
		"recommended_label": "B",
		"recommended_gram": 50,
		"quick_replies": ["Berapa kalau 1 potong?"]
	}`

	parsed, err := parseAssistantResponse(rawJSON)
	if err != nil {
		t.Fatalf("gagal parse JSON: %v", err)
	}

	if parsed.RecommendedLabel != "B" {
		t.Errorf("label salah: %s", parsed.RecommendedLabel)
	}
	if parsed.RecommendedGram == nil || *parsed.RecommendedGram != 50 {
		t.Errorf("gram salah: %v", parsed.RecommendedGram)
	}
}

func TestAskAssistantFallback(t *testing.T) {
	svc := &service{
		repo: nil,
		groq: nil, // force fallback
	}

	req := AssistantRequest{
		Step:      AssistantStepPortion,
		FoodName:  "Nasi Putih",
		UserQuery: "1 centong",
		AvailablePortions: []AvailablePortion{
			{Label: "A", WeightGram: 100},
		},
	}

	res, err := svc.AskAssistant(context.Background(), req)
	if err != nil {
		t.Fatalf("AskAssistant error: %v", err)
	}
	if res.RecommendedLabel != "A" {
		t.Errorf("expected A, got %s", res.RecommendedLabel)
	}
}
