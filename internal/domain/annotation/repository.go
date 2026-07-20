package annotation

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Repository - akses data untuk food_images dan food_areas
type Repository struct {
	db *gorm.DB
}

// NewRepository - factory Repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ErrNotFound - gambar tidak ditemukan (atau tidak published untuk akses publik)
var ErrNotFound = errors.New("food image tidak ditemukan")

// Create - simpan FoodImage baru berstatus draft
func (r *Repository) Create(image *FoodImage) error {
	return r.db.Create(image).Error
}

// FindByID - ambil satu FoodImage beserta areas-nya, urut z_index
func (r *Repository) FindByID(id string) (*FoodImage, error) {
	var image FoodImage

	err := r.db.
		Preload("Areas", func(db *gorm.DB) *gorm.DB {
			return db.Order("z_index ASC, created_at ASC")
		}).
		First(&image, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &image, nil
}

// FindPublishedByID - sama seperti FindByID tapi hanya status published.
// Draft yang diakses lewat endpoint publik harus terlihat seperti tidak ada.
func (r *Repository) FindPublishedByID(id string) (*FoodImage, error) {
	var image FoodImage

	err := r.db.
		Preload("Areas", func(db *gorm.DB) *gorm.DB {
			return db.Order("z_index ASC, created_at ASC")
		}).
		First(&image, "id = ? AND status = ?", id, StatusPublished).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &image, nil
}

// List - daftar FoodImage dengan filter status/pencarian dan pagination.
// areasCount dihitung lewat subquery agar tidak perlu Preload seluruh polygon.
func (r *Repository) List(status, search string, page, limit int) ([]FoodImageSummary, int64, error) {
	// Filter dibangun ulang untuk tiap query. Count() adalah finisher — memakai
	// ulang *gorm.DB yang sama setelahnya membuat kondisi statement bocor ke
	// query berikutnya.
	applyFilter := func(db *gorm.DB) *gorm.DB {
		if status != "" {
			db = db.Where("status = ?", status)
		}
		if s := strings.TrimSpace(search); s != "" {
			db = db.Where("title LIKE ?", "%"+s+"%")
		}
		return db
	}

	var total int64
	if err := applyFilter(r.db.Model(&FoodImage{})).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []FoodImageSummary
	err := applyFilter(r.db.Model(&FoodImage{})).
		Select(`food_images.id, food_images.title, food_images.image_url,
			food_images.thumbnail_url, food_images.width, food_images.height,
			food_images.status, food_images.published_at, food_images.updated_at,
			(SELECT COUNT(*) FROM food_areas WHERE food_areas.food_image_id = food_images.id) AS areas_count`).
		Order("food_images.updated_at DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Scan(&items).Error

	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// ListPublished - daftar published saja untuk konsumsi publik
func (r *Repository) ListPublished(page, limit int) ([]FoodImageSummary, int64, error) {
	return r.List(StatusPublished, "", page, limit)
}

// ListPublishedByFoodID - gambar published yang terkait sebuah food master,
// baik lewat primary_food_id maupun lewat salah satu area-nya.
func (r *Repository) ListPublishedByFoodID(foodID string) ([]FoodImageSummary, error) {
	var items []FoodImageSummary

	err := r.db.Model(&FoodImage{}).
		Select(`DISTINCT food_images.id, food_images.title, food_images.image_url,
			food_images.thumbnail_url, food_images.width, food_images.height,
			food_images.status, food_images.published_at, food_images.updated_at,
			(SELECT COUNT(*) FROM food_areas WHERE food_areas.food_image_id = food_images.id) AS areas_count`).
		Joins("LEFT JOIN food_areas ON food_areas.food_image_id = food_images.id").
		Where("food_images.status = ?", StatusPublished).
		Where("food_images.primary_food_id = ? OR food_areas.food_id = ?", foodID, foodID).
		Order("food_images.updated_at DESC").
		Scan(&items).Error

	if err != nil {
		return nil, err
	}

	return items, nil
}

// UpdateFields - update sebagian kolom FoodImage (PATCH / publish / unpublish)
//
// RowsAffected sengaja TIDAK dipakai untuk mendeteksi "tidak ditemukan":
// MySQL mengembalikan 0 baris terpengaruh ketika baris ada tapi nilainya
// sudah sama. Memakainya akan membuat PATCH judul yang tidak berubah, atau
// unpublish gambar yang memang sudah draft, salah dilaporkan sebagai 404.
// Keberadaan baris diperiksa pemanggil lewat FindByID.
func (r *Repository) UpdateFields(id string, fields map[string]interface{}) error {
	return r.db.Model(&FoodImage{}).Where("id = ?", id).Updates(fields).Error
}

// Delete - hapus FoodImage; areas ikut terhapus lewat cascade di aplikasi
func (r *Repository) Delete(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("food_image_id = ?", id).Delete(&FoodArea{}).Error; err != nil {
			return err
		}

		result := tx.Delete(&FoodImage{}, "id = ?", id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrNotFound
		}
		return nil
	})
}

// ReplaceAreas - ganti seluruh area sebuah gambar dalam satu transaksi.
//
// Hapus-lalu-insert dipilih ketimbang diff per-area: payload editor selalu
// mengirim state lengkap, dan diffing menambah jalur bug tanpa manfaat.
// ID yang dikirim FE dipertahankan supaya referensi di editor tetap stabil.
func (r *Repository) ReplaceAreas(imageID string, areas []AreaInput) (time.Time, error) {
	var updatedAt time.Time

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Kunci baris induk agar dua autosave bersamaan tidak saling menimpa
		var image FoodImage
		if err := tx.First(&image, "id = ?", imageID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}

		if err := tx.Where("food_image_id = ?", imageID).Delete(&FoodArea{}).Error; err != nil {
			return err
		}

		if len(areas) > 0 {
			rows := make([]FoodArea, 0, len(areas))
			for _, in := range areas {
				row := FoodArea{
					FoodImageID: imageID,
					Name:        in.Name,
					FoodID:      in.FoodID,
					Polygon:     in.Polygon,
					ZIndex:      in.ZIndex,
				}
				if in.ID != nil && strings.TrimSpace(*in.ID) != "" {
					row.ID = *in.ID
				}
				rows = append(rows, row)
			}

			if err := tx.Create(&rows).Error; err != nil {
				return err
			}
		}

		// Sentuh updated_at induk supaya indikator "Tersimpan" punya sumber waktu
		now := time.Now()
		if err := tx.Model(&FoodImage{}).
			Where("id = ?", imageID).
			Update("updated_at", now).Error; err != nil {
			return err
		}
		updatedAt = now

		return nil
	})

	return updatedAt, err
}

// CountAreas - jumlah area sebuah gambar
func (r *Repository) CountAreas(imageID string) (int64, error) {
	var count int64
	err := r.db.Model(&FoodArea{}).Where("food_image_id = ?", imageID).Count(&count).Error
	return count, err
}
