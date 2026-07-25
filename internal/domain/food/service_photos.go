package food

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"atlas_food/internal/domain/annotation"
	"atlas_food/internal/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Prefiks di as_served_images.description untuk menautkan ke food_images.id
const foodImageLinkPrefix = "food_image:"

func foodImageLink(foodImageID string) string {
	return foodImageLinkPrefix + foodImageID
}

func photoMaxForType(photoType string) int {
	if photoType == "range" {
		return MaxRangePhotos
	}
	return MaxSeriesPhotos
}

// ListFoodPhotos - kartu foto terpadu (anotasi + gram) untuk satu makanan
func (s *foodService) ListFoodPhotos(foodID string) (*FoodPhotoListResponse, error) {
	food, err := s.repo.GetFoodByID(foodID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewAppError(http.StatusNotFound, "NOT_FOUND", "Makanan tidak ditemukan")
		}
		return nil, err
	}

	if s.ann == nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "Layanan anotasi belum terpasang")
	}

	listed, err := s.ann.List(annotation.ListFoodImagesQuery{
		PrimaryFoodID: foodID,
		Page:          1,
		Limit:         MaxSeriesPhotos,
	})
	if err != nil {
		return nil, err
	}

	set, _ := s.ensurePrimaryAsServedSet(food)
	var portionByFoodImage = map[string]AsServedImage{}
	var portionByURL = map[string]AsServedImage{}
	if set != nil && set.ID != "" {
		images, imgErr := s.repo.GetAsServedImagesBySetID(set.ID)
		if imgErr == nil {
			for _, img := range images {
				if strings.HasPrefix(img.Description, foodImageLinkPrefix) {
					key := strings.TrimPrefix(img.Description, foodImageLinkPrefix)
					if key != "" {
						portionByFoodImage[key] = img
					}
				}
				if img.ImageURL != "" {
					portionByURL[img.ImageURL] = img
				}
			}
		}
	}

	items := make([]FoodPhotoResponse, 0, len(listed.Items))
	for _, item := range listed.Items {
		portion, ok := portionByFoodImage[item.ID]
		if !ok {
			portion, ok = portionByURL[item.ImageURL]
		}
		resp := FoodPhotoResponse{
			ID:           item.ID,
			Title:        item.Title,
			ImageURL:     item.ImageURL,
			ThumbnailURL: item.ThumbnailURL,
			Width:        item.Width,
			Height:       item.Height,
			Status:       item.Status,
			AreasCount:   item.AreasCount,
			UpdatedAt:    item.UpdatedAt.Format(time.RFC3339),
		}
		if item.PublishedAt != nil {
			pub := item.PublishedAt.Format(time.RFC3339)
			resp.PublishedAt = &pub
		}
		if ok {
			resp.Label = portion.Label
			resp.WeightGram = portion.WeightGram
			resp.AsServedImageID = portion.ID
			resp.DisplayOrder = portion.DisplayOrder
			// Perbaiki tautan rusak (deskripsi kosong/legacy) agar update berikutnya stabil
			if portion.Description != foodImageLink(item.ID) && portion.ID != "" {
				fixed := foodImageLink(item.ID)
				_, _ = s.UpdateAsServedImage(portion.ID, UpdateAsServedImageRequest{Description: &fixed})
			}
		} else {
			resp.Label = item.Title
		}
		items = append(items, resp)
	}

	photoType := food.PhotoType
	if photoType == "" {
		photoType = "series"
	}

	return &FoodPhotoListResponse{
		Items:     items,
		PhotoType: photoType,
		MaxPhotos: photoMaxForType(photoType),
		Count:     len(items),
	}, nil
}

// CreateFoodPhoto - draft anotasi + foto porsi berbobot (satu kartu)
func (s *foodService) CreateFoodPhoto(foodID, userID string, req CreateFoodPhotoRequest) (*FoodPhotoResponse, error) {
	food, err := s.repo.GetFoodByID(foodID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewAppError(http.StatusNotFound, "NOT_FOUND", "Makanan tidak ditemukan")
		}
		return nil, err
	}
	if s.ann == nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "Layanan anotasi belum terpasang")
	}

	photoType := food.PhotoType
	if photoType == "" {
		photoType = "series"
	}
	maxPhotos := photoMaxForType(photoType)

	existing, err := s.ann.List(annotation.ListFoodImagesQuery{
		PrimaryFoodID: foodID,
		Page:          1,
		Limit:         MaxSeriesPhotos,
	})
	if err != nil {
		return nil, err
	}
	if len(existing.Items) >= maxPhotos {
		if photoType == "range" {
			return nil, utils.NewAppError(http.StatusConflict, "PHOTO_LIMIT", "Tipe range hanya boleh 1 foto")
		}
		return nil, utils.NewAppError(http.StatusConflict, "PHOTO_LIMIT", fmt.Sprintf("Tipe series maksimal %d foto", MaxSeriesPhotos))
	}

	title := strings.TrimSpace(req.Title)
	label := strings.TrimSpace(req.Label)
	if label == "" {
		label = title
	}

	fid := foodID
	created, err := s.ann.Create(annotation.CreateFoodImageRequest{
		Title:         title,
		ImageURL:      req.ImageURL,
		ThumbnailURL:  req.ThumbnailURL,
		Width:         req.Width,
		Height:        req.Height,
		PrimaryFoodID: &fid,
	}, userID)
	if err != nil {
		return nil, err
	}

	set, err := s.ensurePrimaryAsServedSet(food)
	if err != nil {
		_ = s.ann.Delete(created.ID)
		return nil, err
	}

	portionImages, err := s.repo.GetAsServedImagesBySetID(set.ID)
	if err != nil {
		_ = s.ann.Delete(created.ID)
		return nil, err
	}

	portionRow := AsServedImage{
		ID:           uuid.New().String(),
		SetID:        set.ID,
		Label:        label,
		ImageURL:     req.ImageURL,
		ThumbnailURL: req.ThumbnailURL,
		WeightGram:   req.WeightGram,
		Description:  foodImageLink(created.ID),
		DisplayOrder: len(portionImages),
	}
	if created.ID == "" || set.ID == "" {
		_ = s.ann.Delete(created.ID)
		return nil, utils.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal membuat tautan foto (ID kosong)")
	}
	if err := s.repo.CreateAsServedImages([]AsServedImage{portionRow}); err != nil {
		_ = s.ann.Delete(created.ID)
		return nil, err
	}

	// Ambil ulang agar ID terisi (Create batch tidak selalu mirror pointer lokal)
	linked, linkErr := s.repo.GetAsServedImageByDescription(set.ID, foodImageLink(created.ID))
	asID := ""
	weight := req.WeightGram
	order := len(portionImages)
	if linkErr == nil && linked != nil {
		asID = linked.ID
		weight = linked.WeightGram
		order = linked.DisplayOrder
	}

	return &FoodPhotoResponse{
		ID:              created.ID,
		Title:           created.Title,
		Label:           label,
		ImageURL:        created.ImageURL,
		ThumbnailURL:    created.ThumbnailURL,
		Width:           created.Width,
		Height:          created.Height,
		Status:          created.Status,
		AreasCount:      0,
		WeightGram:      weight,
		AsServedImageID: asID,
		DisplayOrder:    order,
		UpdatedAt:       created.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// UpdateFoodPhoto - judul anotasi + label/berat porsi
func (s *foodService) UpdateFoodPhoto(foodID, photoID string, req UpdateFoodPhotoRequest) (*FoodPhotoResponse, error) {
	if _, err := s.repo.GetFoodByID(foodID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewAppError(http.StatusNotFound, "NOT_FOUND", "Makanan tidak ditemukan")
		}
		return nil, err
	}
	if s.ann == nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "Layanan anotasi belum terpasang")
	}

	image, err := s.ann.Get(photoID)
	if err != nil {
		return nil, err
	}
	if image.PrimaryFoodID == nil || *image.PrimaryFoodID != foodID {
		return nil, utils.NewAppError(http.StatusNotFound, "NOT_FOUND", "Foto tidak terkait makanan ini")
	}

	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if _, err := s.ann.Update(photoID, annotation.UpdateFoodImageRequest{Title: &title}); err != nil {
			return nil, err
		}
	}

	set, err := s.ensurePrimaryAsServedSetByFoodID(foodID)
	if err != nil {
		return nil, err
	}

	portion, err := s.findLinkedPortion(set.ID, photoID, image.ImageURL)
	if err != nil {
		return nil, err
	}

	if portion == nil {
		// Foto lama tanpa tautan gram — buat baris porsi sekarang
		weight := 1.0
		if req.WeightGram != nil {
			weight = *req.WeightGram
		}
		label := image.Title
		if req.Label != nil && strings.TrimSpace(*req.Label) != "" {
			label = strings.TrimSpace(*req.Label)
		}
		existing, _ := s.repo.GetAsServedImagesBySetID(set.ID)
		row := AsServedImage{
			ID:           uuid.New().String(),
			SetID:        set.ID,
			Label:        label,
			ImageURL:     image.ImageURL,
			ThumbnailURL: image.ThumbnailURL,
			WeightGram:   weight,
			Description:  foodImageLink(photoID),
			DisplayOrder: len(existing),
		}
		if err := s.repo.CreateAsServedImages([]AsServedImage{row}); err != nil {
			return nil, err
		}
	} else {
		upd := UpdateAsServedImageRequest{}
		link := foodImageLink(photoID)
		upd.Description = &link
		if req.Label != nil {
			upd.Label = req.Label
		}
		if req.WeightGram != nil {
			upd.WeightGram = req.WeightGram
		}
		if _, err := s.UpdateAsServedImage(portion.ID, upd); err != nil {
			return nil, err
		}
	}

	return s.getFoodPhotoCard(foodID, photoID)
}

func (s *foodService) findLinkedPortion(setID, photoID, imageURL string) (*AsServedImage, error) {
	portion, err := s.repo.GetAsServedImageByDescription(setID, foodImageLink(photoID))
	if err == nil && portion != nil {
		return portion, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	images, listErr := s.repo.GetAsServedImagesBySetID(setID)
	if listErr != nil {
		return nil, listErr
	}
	for i := range images {
		if images[i].ImageURL == imageURL {
			return &images[i], nil
		}
		// Legacy bug: description tersimpan "food_image:" tanpa UUID
		if images[i].Description == foodImageLinkPrefix && images[i].ImageURL == imageURL {
			return &images[i], nil
		}
	}
	return nil, nil
}

// DeleteFoodPhoto - hapus anotasi + foto porsi tertaut
func (s *foodService) DeleteFoodPhoto(foodID, photoID string) error {
	if _, err := s.repo.GetFoodByID(foodID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NewAppError(http.StatusNotFound, "NOT_FOUND", "Makanan tidak ditemukan")
		}
		return err
	}
	if s.ann == nil {
		return utils.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "Layanan anotasi belum terpasang")
	}

	image, err := s.ann.Get(photoID)
	if err != nil {
		return err
	}
	if image.PrimaryFoodID == nil || *image.PrimaryFoodID != foodID {
		return utils.NewAppError(http.StatusNotFound, "NOT_FOUND", "Foto tidak terkait makanan ini")
	}

	sets, _ := s.repo.GetAsServedSetsByFoodID(foodID)
	for _, set := range sets {
		portion, err := s.repo.GetAsServedImageByDescription(set.ID, foodImageLink(photoID))
		if err == nil && portion != nil {
			_ = s.repo.DeleteAsServedImage(portion.ID)
		}
	}

	return s.ann.Delete(photoID)
}

// PublishFoodPhoto - finalkan anotasi (wajib punya area valid)
func (s *foodService) PublishFoodPhoto(foodID, photoID string) (*FoodPhotoResponse, error) {
	if err := s.assertPhotoBelongsToFood(foodID, photoID); err != nil {
		return nil, err
	}
	if _, err := s.ann.Publish(photoID); err != nil {
		return nil, err
	}
	return s.getFoodPhotoCard(foodID, photoID)
}

// UnpublishFoodPhoto - kembalikan ke draft
func (s *foodService) UnpublishFoodPhoto(foodID, photoID string) (*FoodPhotoResponse, error) {
	if err := s.assertPhotoBelongsToFood(foodID, photoID); err != nil {
		return nil, err
	}
	if _, err := s.ann.Unpublish(photoID); err != nil {
		return nil, err
	}
	return s.getFoodPhotoCard(foodID, photoID)
}

func (s *foodService) assertPhotoBelongsToFood(foodID, photoID string) error {
	if _, err := s.repo.GetFoodByID(foodID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NewAppError(http.StatusNotFound, "NOT_FOUND", "Makanan tidak ditemukan")
		}
		return err
	}
	if s.ann == nil {
		return utils.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "Layanan anotasi belum terpasang")
	}
	image, err := s.ann.Get(photoID)
	if err != nil {
		return err
	}
	if image.PrimaryFoodID == nil || *image.PrimaryFoodID != foodID {
		return utils.NewAppError(http.StatusNotFound, "NOT_FOUND", "Foto tidak terkait makanan ini")
	}
	return nil
}

func (s *foodService) getFoodPhotoCard(foodID, photoID string) (*FoodPhotoResponse, error) {
	list, err := s.ListFoodPhotos(foodID)
	if err != nil {
		return nil, err
	}
	for i := range list.Items {
		if list.Items[i].ID == photoID {
			return &list.Items[i], nil
		}
	}
	return nil, utils.NewAppError(http.StatusNotFound, "NOT_FOUND", "Foto tidak ditemukan")
}

func (s *foodService) ensurePrimaryAsServedSetByFoodID(foodID string) (*AsServedSet, error) {
	food, err := s.repo.GetFoodByID(foodID)
	if err != nil {
		return nil, err
	}
	return s.ensurePrimaryAsServedSet(food)
}

// ensurePrimaryAsServedSet - satu set + portion method as_served per makanan
func (s *foodService) ensurePrimaryAsServedSet(food *Food) (*AsServedSet, error) {
	sets, err := s.repo.GetAsServedSetsByFoodID(food.ID)
	if err != nil {
		return nil, err
	}
	if len(sets) > 0 {
		if err := s.ensureAsServedPortionMethod(food.ID, sets[0].Code); err != nil {
			return nil, err
		}
		return &sets[0], nil
	}

	base := strings.ToLower(food.Code)
	base = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		return '-'
	}, base)
	base = strings.Trim(base, "-")
	if base == "" {
		base = "food"
	}
	code := fmt.Sprintf("%s-porsi", base)
	if existing, err := s.repo.GetAsServedSetByCode(code); err == nil && existing != nil {
		// Kode terpakai makanan lain — unikkan
		code = fmt.Sprintf("%s-%s", base, food.ID[:8])
	}

	set, err := s.CreateAsServedSet(CreateAsServedSetRequest{
		Code:        code,
		Name:        "Porsi " + food.Name,
		Description: "Set foto porsi otomatis untuk " + food.Name,
		FoodID:      food.ID,
	})
	if err != nil {
		return nil, err
	}
	// Pastikan ID terisi (fallback reload by code)
	if set.ID == "" {
		reloaded, reloadErr := s.repo.GetAsServedSetByCode(code)
		if reloadErr != nil {
			return nil, reloadErr
		}
		set = reloaded
	}

	if err := s.ensureAsServedPortionMethod(food.ID, set.Code); err != nil {
		return nil, err
	}
	return set, nil
}

func (s *foodService) ensureAsServedPortionMethod(foodID, setCode string) error {
	methods, err := s.repo.GetPortionMethodsByFoodID(foodID)
	if err != nil {
		return err
	}
	for _, m := range methods {
		if m.MethodType != "as_served" {
			continue
		}
		var cfg struct {
			SetCode string `json:"setCode"`
		}
		_ = json.Unmarshal([]byte(m.Config), &cfg)
		if cfg.SetCode == setCode {
			return nil
		}
	}

	cfg, _ := json.Marshal(map[string]interface{}{
		"selectionType":   "as_served_quantity",
		"setCode":         setCode,
		"allowFractions":  true,
		"showCalculation": true,
	})

	_, err = s.AddPortionMethod(foodID, CreatePortionMethodRequest{
		MethodType:  "as_served",
		Label:       "Foto porsi",
		Description: "Pilih foto porsi yang paling mirip",
		Config:      cfg,
	})
	return err
}
