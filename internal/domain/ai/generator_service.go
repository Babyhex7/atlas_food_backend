package ai

import (
	"atlas_food/internal/domain/food"
	"atlas_food/internal/pkg/gemini"
	"atlas_food/internal/pkg/utils"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GeneratorService - interface untuk fungsionalitas Missing Foods & AI Generator
type GeneratorService interface {
	GetMissingFoodsSummary(ctx context.Context) ([]MissingFoodSummary, error)
	GenerateFoodAI(ctx context.Context, req GenerateFoodAIRequest) (*GenerateFoodAIResponse, error)
	PublishFoodAI(ctx context.Context, req PublishFoodAIRequest) (*PublishFoodAIResponse, error)
}

type generatorService struct {
	db     *gorm.DB
	gemini *gemini.Client
}

// NewGeneratorService - inisialisasi generator service
func NewGeneratorService(db *gorm.DB, geminiClient *gemini.Client) GeneratorService {
	return &generatorService{
		db:     db,
		gemini: geminiClient,
	}
}

// GetMissingFoodsSummary - agregasi laporan missing foods dari seluruh submission
func (s *generatorService) GetMissingFoodsSummary(ctx context.Context) ([]MissingFoodSummary, error) {
	type rawSubmission struct {
		MissingFoods string    `gorm:"column:missing_foods"`
		SubmittedAt  time.Time `gorm:"column:submitted_at"`
	}

	var rawSubs []rawSubmission
	err := s.db.WithContext(ctx).Table("survey_submissions").
		Select("missing_foods, submitted_at").
		Where("missing_foods IS NOT NULL AND missing_foods != '' AND missing_foods != '[]'").
		Find(&rawSubs).Error
	if err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "DB_ERROR", "Gagal mengambil data submission: "+err.Error())
	}

	// Agregasi berdasar nama makanan (case-insensitive)
	type aggItem struct {
		name        string
		description string
		count       int
		lastTime    time.Time
	}
	aggMap := make(map[string]*aggItem)

	for _, sub := range rawSubs {
		var items []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		if err := json.Unmarshal([]byte(sub.MissingFoods), &items); err != nil {
			continue
		}

		for _, it := range items {
			nameClean := strings.TrimSpace(it.Name)
			if nameClean == "" {
				continue
			}
			key := strings.ToLower(nameClean)
			if existing, ok := aggMap[key]; ok {
				existing.count++
				if sub.SubmittedAt.After(existing.lastTime) {
					existing.lastTime = sub.SubmittedAt
				}
				if it.Description != "" && existing.description == "" {
					existing.description = it.Description
				}
			} else {
				aggMap[key] = &aggItem{
					name:        nameClean,
					description: it.Description,
					count:       1,
					lastTime:    sub.SubmittedAt,
				}
			}
		}
	}

	// Ambil semua nama makanan yang sudah ada di katalog foods
	var existingFoods []struct {
		ID   string `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}
	s.db.WithContext(ctx).Table("foods").Select("id, name").Find(&existingFoods)
	catalogMap := make(map[string]string)
	for _, f := range existingFoods {
		catalogMap[strings.ToLower(strings.TrimSpace(f.Name))] = f.ID
	}

	result := make([]MissingFoodSummary, 0, len(aggMap))
	for key, it := range aggMap {
		summary := MissingFoodSummary{
			Name:           it.name,
			Description:    it.description,
			ReportCount:    it.count,
			LastReportedAt: it.lastTime,
		}
		if foodID, exists := catalogMap[key]; exists {
			summary.IsInCatalog = true
			summary.CatalogFoodID = &foodID
		}
		result = append(result, summary)
	}

	// Urutkan dari yang paling sering dilaporkan
	sort.Slice(result, func(i, j int) bool {
		if result[i].ReportCount == result[j].ReportCount {
			return result[i].LastReportedAt.After(result[j].LastReportedAt)
		}
		return result[i].ReportCount > result[j].ReportCount
	})

	return result, nil
}

// GenerateFoodAI - estimasi gizi dan generate 4 foto porsi A/B/C/D
func (s *generatorService) GenerateFoodAI(ctx context.Context, req GenerateFoodAIRequest) (*GenerateFoodAIResponse, error) {
	if strings.TrimSpace(req.FoodName) == "" {
		return nil, utils.NewAppError(http.StatusBadRequest, "VALIDATION_ERROR", "Nama makanan wajib diisi")
	}

	// 1. Minta estimasi gizi & konsep porsi ke Gemini
	analysis, err := s.gemini.EstimateFoodNutritionAndPortions(ctx, req.FoodName, req.Description)
	if err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "AI_ERROR", "Gagal memproses estimasi AI: "+err.Error())
	}

	// 2. Generate 4 gambar porsi (A, B, C, D)
	portionViews := make([]GeneratedPortionView, 0, len(analysis.Portions))
	for _, p := range analysis.Portions {
		imgURL, imgErr := s.gemini.GeneratePortionImage(ctx, p.VisualPrompt, analysis.FoodName, p.Code, p.WeightGram)
		if imgErr != nil || imgURL == "" {
			imgURL = p.ImageURL
		}
		portionViews = append(portionViews, GeneratedPortionView{
			Code:         p.Code,
			Label:        p.Label,
			WeightGram:   p.WeightGram,
			Description:  p.Description,
			VisualPrompt: p.VisualPrompt,
			ImageURL:     imgURL,
		})
	}

	// 3. Match kategori jika belum ada ID
	catID := req.CategoryID
	if catID == nil || *catID == "" {
		var matchedCat struct {
			ID string `gorm:"column:id"`
		}
		catKeywords := map[string]string{
			"sayuran": "AS",
			"sayur":   "AS",
			"buah":    "AB",
			"snack":   "AK",
			"kue":     "AK",
			"minuman": "AMK",
			"pokok":   "MP",
			"hewani":  "LH",
			"nabati":  "LN",
		}
		catLower := strings.ToLower(strings.TrimSpace(analysis.CategoryName))
		searchPattern := "%" + analysis.CategoryName + "%"
		codeHint := ""
		for kw, code := range catKeywords {
			if strings.Contains(catLower, kw) {
				codeHint = code
				break
			}
		}
		if codeHint != "" {
			if err := s.db.WithContext(ctx).Table("categories").
				Where("code = ? OR name LIKE ?", codeHint, searchPattern).
				First(&matchedCat).Error; err == nil && matchedCat.ID != "" {
				catID = &matchedCat.ID
			}
		} else if err := s.db.WithContext(ctx).Table("categories").
			Where("name LIKE ? OR code LIKE ?", searchPattern, searchPattern).
			First(&matchedCat).Error; err == nil && matchedCat.ID != "" {
			catID = &matchedCat.ID
		}
	}

	return &GenerateFoodAIResponse{
		FoodName:        analysis.FoodName,
		LocalName:       analysis.LocalName,
		Description:     analysis.Description,
		CategoryID:      catID,
		CategoryName:    analysis.CategoryName,
		PhotoType:       analysis.PhotoType,
		CaloriesPer100g: analysis.CaloriesPer100g,
		ProteinPer100g:  analysis.ProteinPer100g,
		FatPer100g:      analysis.FatPer100g,
		CarbsPer100g:    analysis.CarbsPer100g,
		Portions:        portionViews,
	}, nil
}

// PublishFoodAI - menyimpan makanan hasil AI yang telah disetujui admin ke database
func (s *generatorService) PublishFoodAI(ctx context.Context, req PublishFoodAIRequest) (*PublishFoodAIResponse, error) {
	if strings.TrimSpace(req.FoodName) == "" {
		return nil, utils.NewAppError(http.StatusBadRequest, "VALIDATION_ERROR", "Nama makanan wajib diisi")
	}
	if len(req.Portions) == 0 {
		return nil, utils.NewAppError(http.StatusBadRequest, "VALIDATION_ERROR", "Minimal harus ada 1 foto porsi")
	}

	var foodID, foodCode, setID string

	// Validasi & fallback photo_type (series, range, guide)
	photoType := strings.ToLower(strings.TrimSpace(req.PhotoType))
	if photoType != "range" && photoType != "guide" {
		photoType = "series"
	}

	// Jalankan transaksi database
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		foodID = uuid.New().String()
		randSuffix := make([]byte, 3)
		rand.Read(randSuffix)
		foodCode = fmt.Sprintf("AI-%s", strings.ToUpper(hex.EncodeToString(randSuffix)))

		// 1. Insert ke tabel foods
		newFood := food.Food{
			ID:            foodID,
			Code:          foodCode,
			Name:          req.FoodName,
			LocalName:     req.LocalName,
			Description:   req.Description,
			PhotoType:     photoType,
			CategoryID:    req.CategoryID,
			IsActive:      true,
			IsAiGenerated: true,
			CreatedAt:     time.Now(),
		}
		if err := tx.Create(&newFood).Error; err != nil {
			return err
		}

		// 2. Insert ke as_served_sets
		setID = uuid.New().String()
		setCode := fmt.Sprintf("set-ai-%s", strings.ToLower(hex.EncodeToString(randSuffix)))
		asSet := food.AsServedSet{
			ID:          setID,
			Code:        setCode,
			Name:        fmt.Sprintf("Porsi %s", req.FoodName),
			Description: fmt.Sprintf("Panduan foto porsi AI untuk %s", req.FoodName),
			FoodID:      &foodID,
			Category:    "dishes",
			CreatedAt:   time.Now(),
		}
		if err := tx.Create(&asSet).Error; err != nil {
			return err
		}

		// 3. Insert ke as_served_images
		var asImages []food.AsServedImage
		for i, p := range req.Portions {
			imgID := uuid.New().String()
			asImages = append(asImages, food.AsServedImage{
				ID:            imgID,
				SetID:         setID,
				Label:         p.Code,
				ImageURL:      p.ImageURL,
				ThumbnailURL:  p.ImageURL,
				WeightGram:    p.WeightGram,
				Description:   p.Description,
				DisplayOrder:  i + 1,
				IsAiGenerated: true,
				CreatedAt:     time.Now(),
			})
		}
		if err := tx.Create(&asImages).Error; err != nil {
			return err
		}

		// 4. Insert ke food_nutrients (Energy, Protein, Carbs, Fat)
		// Ambil ID tipe nutrisi dari database
		var nTypes []food.NutrientType
		if err := tx.Find(&nTypes).Error; err == nil {
			typeMap := make(map[string]int)
			for _, nt := range nTypes {
				typeMap[nt.Code] = nt.ID
			}

			nutrientsToInsert := []food.FoodNutrient{}
			if id, ok := typeMap["energy"]; ok && req.CaloriesPer100g > 0 {
				nutrientsToInsert = append(nutrientsToInsert, food.FoodNutrient{FoodID: foodID, NutrientTypeID: id, ValuePer100g: req.CaloriesPer100g})
			}
			if id, ok := typeMap["protein"]; ok && req.ProteinPer100g > 0 {
				nutrientsToInsert = append(nutrientsToInsert, food.FoodNutrient{FoodID: foodID, NutrientTypeID: id, ValuePer100g: req.ProteinPer100g})
			}
			if id, ok := typeMap["carbs"]; ok && req.CarbsPer100g > 0 {
				nutrientsToInsert = append(nutrientsToInsert, food.FoodNutrient{FoodID: foodID, NutrientTypeID: id, ValuePer100g: req.CarbsPer100g})
			}
			if id, ok := typeMap["fat"]; ok && req.FatPer100g > 0 {
				nutrientsToInsert = append(nutrientsToInsert, food.FoodNutrient{FoodID: foodID, NutrientTypeID: id, ValuePer100g: req.FatPer100g})
			}

			if len(nutrientsToInsert) > 0 {
				if err := tx.Create(&nutrientsToInsert).Error; err != nil {
					return err
				}
			}
		}

		// 5. Insert ke food_portion_size_methods
		configJSON, _ := json.Marshal(map[string]interface{}{
			"as_served_set_id": setID,
			"setCode":          setCode,
			"set_code":         setCode,
		})
		method := food.PortionSizeMethod{
			FoodID:       foodID,
			MethodType:   "as_served",
			Label:        "Porsi Foto (As Served)",
			Description:  "Pilih foto porsi yang paling mendekati porsi makan Anda",
			Config:       string(configJSON),
			DisplayOrder: 1,
			IsActive:     true,
			CreatedAt:    time.Now(),
		}
		if err := tx.Create(&method).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "DB_ERROR", "Gagal menyimpan makanan ke katalog: "+err.Error())
	}

	return &PublishFoodAIResponse{
		FoodID:     foodID,
		FoodCode:   foodCode,
		FoodName:   req.FoodName,
		SetID:      setID,
		TotalImage: len(req.Portions),
		Message:    "Makanan dan porsi AI berhasil dipublikasikan ke katalog resmi!",
	}, nil
}
