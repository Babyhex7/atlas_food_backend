package food

import (
	"testing"
	"time"

	"atlas_food/internal/config"
	"github.com/google/uuid"
)

func TestAssociatedFoodRepositoryAndService(t *testing.T) {
	cfg := config.Load()
	db := config.InitDB(cfg)
	if db == nil {
		t.Skip("MySQL database tidak tersedia, lewati integrasi DB test")
		return
	}

	repo := NewRepository(db)
	svc := NewService(repo, nil)

	// Create test foods
	foodA := Food{
		ID:        uuid.New().String(),
		Code:      "TEST-MAIN-" + uuid.New().String()[:8],
		Name:      "Nasi Goreng Spesial Test",
		PhotoType: "series",
		IsActive:  true,
		CreatedAt: time.Now(),
	}
	foodB := Food{
		ID:        uuid.New().String(),
		Code:      "TEST-ASSOC-1-" + uuid.New().String()[:8],
		Name:      "Kerupuk Udang Test",
		PhotoType: "series",
		IsActive:  true,
		CreatedAt: time.Now(),
	}
	foodC := Food{
		ID:        uuid.New().String(),
		Code:      "TEST-ASSOC-2-" + uuid.New().String()[:8],
		Name:      "Acar Timun Test",
		PhotoType: "series",
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	if err := db.Create(&foodA).Error; err != nil {
		t.Fatalf("Failed to create foodA: %v", err)
	}
	if err := db.Create(&foodB).Error; err != nil {
		t.Fatalf("Failed to create foodB: %v", err)
	}
	if err := db.Create(&foodC).Error; err != nil {
		t.Fatalf("Failed to create foodC: %v", err)
	}

	defer func() {
		db.Where("food_id = ?", foodA.ID).Delete(&AssociatedFood{})
		db.Delete(&foodA)
		db.Delete(&foodB)
		db.Delete(&foodC)
	}()

	t.Run("Get when empty", func(t *testing.T) {
		items, err := svc.GetAssociatedFoods(foodA.ID)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(items) != 0 {
			t.Fatalf("Expected 0 items, got %d", len(items))
		}
	})

	t.Run("Prevent self association", func(t *testing.T) {
		_, err := svc.SetAssociatedFoods(foodA.ID, SetAssociatedFoodsRequest{
			Items: []AssociatedFoodInput{
				{AssociatedFoodID: foodA.ID, Priority: 1},
			},
		})
		if err == nil {
			t.Fatal("Expected error for self association, got nil")
		}
	})

	t.Run("Set and Get associated foods", func(t *testing.T) {
		setResp, err := svc.SetAssociatedFoods(foodA.ID, SetAssociatedFoodsRequest{
			Items: []AssociatedFoodInput{
				{AssociatedFoodID: foodB.ID, Priority: 1, IsDefault: true},
				{AssociatedFoodID: foodC.ID, Priority: 2, IsDefault: false},
			},
		})
		if err != nil {
			t.Fatalf("SetAssociatedFoods error: %v", err)
		}
		if len(setResp) != 2 {
			t.Fatalf("Expected 2 items, got %d", len(setResp))
		}
		if setResp[0].AssociatedFoodID != foodB.ID {
			t.Errorf("Expected first item %s, got %s", foodB.ID, setResp[0].AssociatedFoodID)
		}
		if setResp[0].AssociatedFood == nil || setResp[0].AssociatedFood.Name != "Kerupuk Udang Test" {
			t.Errorf("Expected AssociatedFood.Name 'Kerupuk Udang Test', got %+v", setResp[0].AssociatedFood)
		}
		if !setResp[0].IsDefault {
			t.Error("Expected IsDefault to be true for foodB")
		}

		if setResp[1].AssociatedFoodID != foodC.ID {
			t.Errorf("Expected second item %s, got %s", foodC.ID, setResp[1].AssociatedFoodID)
		}
		if setResp[1].AssociatedFood == nil || setResp[1].AssociatedFood.Name != "Acar Timun Test" {
			t.Errorf("Expected AssociatedFood.Name 'Acar Timun Test', got %+v", setResp[1].AssociatedFood)
		}

		// Verify via repository directly
		rawList, err := repo.GetAssociatedFoods(foodA.ID)
		if err != nil {
			t.Fatalf("repo.GetAssociatedFoods error: %v", err)
		}
		if len(rawList) != 2 {
			t.Fatalf("Expected 2 raw items, got %d", len(rawList))
		}
		if rawList[0].AssociatedFood.Name != "Kerupuk Udang Test" {
			t.Errorf("Expected raw name 'Kerupuk Udang Test', got %s", rawList[0].AssociatedFood.Name)
		}
	})

	t.Run("Clear associated foods", func(t *testing.T) {
		setResp, err := svc.SetAssociatedFoods(foodA.ID, SetAssociatedFoodsRequest{
			Items: []AssociatedFoodInput{},
		})
		if err != nil {
			t.Fatalf("Clear error: %v", err)
		}
		if len(setResp) != 0 {
			t.Fatalf("Expected 0 items after clear, got %d", len(setResp))
		}

		rawList, err := repo.GetAssociatedFoods(foodA.ID)
		if err != nil {
			t.Fatalf("repo query error: %v", err)
		}
		if len(rawList) != 0 {
			t.Fatalf("Expected 0 raw items, got %d", len(rawList))
		}
	})
}
