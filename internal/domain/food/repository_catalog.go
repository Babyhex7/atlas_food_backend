package food

import "gorm.io/gorm"

// Implementasi repository untuk katalog admin: kategori, as-served, portion method.

// ============ CATEGORY ============

// CreateCategory - simpan kategori baru
func (r *foodRepository) CreateCategory(category *Category) error {
	return r.db.Create(category).Error
}

// UpdateCategory - simpan perubahan kategori
func (r *foodRepository) UpdateCategory(category *Category) error {
	return r.db.Save(category).Error
}

// DeleteCategory - hapus kategori berdasarkan ID
func (r *foodRepository) DeleteCategory(id string) error {
	return r.db.Delete(&Category{}, "id = ?", id).Error
}

// CountFoodsByCategoryID - dipakai service untuk menolak hapus kategori
// yang masih dipakai food; menghapusnya akan memutus Find Food.
func (r *foodRepository) CountFoodsByCategoryID(id string) (int64, error) {
	var count int64
	err := r.db.Model(&Food{}).Where("category_id = ?", id).Count(&count).Error
	return count, err
}

// ============ AS-SERVED SET ============

// GetAsServedSetByID - ambil satu set as-served berdasarkan ID
func (r *foodRepository) GetAsServedSetByID(id string) (*AsServedSet, error) {
	var set AsServedSet
	if err := r.db.Where("id = ?", id).First(&set).Error; err != nil {
		return nil, err
	}
	return &set, nil
}

// UpdateAsServedSet - simpan perubahan set as-served
func (r *foodRepository) UpdateAsServedSet(set *AsServedSet) error {
	return r.db.Save(set).Error
}

// DeleteAsServedSet - hapus set beserta seluruh fotonya dalam satu transaksi,
// supaya tidak meninggalkan as_served_images yatim.
func (r *foodRepository) DeleteAsServedSet(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("set_id = ?", id).Delete(&AsServedImage{}).Error; err != nil {
			return err
		}
		return tx.Delete(&AsServedSet{}, "id = ?", id).Error
	})
}

// ============ AS-SERVED IMAGE ============

// GetAsServedImageByID - ambil satu foto porsi berdasarkan ID
func (r *foodRepository) GetAsServedImageByID(id string) (*AsServedImage, error) {
	var image AsServedImage
	if err := r.db.Where("id = ?", id).First(&image).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

// UpdateAsServedImage - simpan perubahan foto porsi
func (r *foodRepository) UpdateAsServedImage(image *AsServedImage) error {
	return r.db.Save(image).Error
}

// DeleteAsServedImage - hapus foto porsi berdasarkan ID
func (r *foodRepository) DeleteAsServedImage(id string) error {
	return r.db.Delete(&AsServedImage{}, "id = ?", id).Error
}

// ============ PORTION METHOD ============

// GetPortionMethodByID - ambil satu metode porsi berdasarkan ID
func (r *foodRepository) GetPortionMethodByID(id int) (*PortionSizeMethod, error) {
	var method PortionSizeMethod
	if err := r.db.Where("id = ?", id).First(&method).Error; err != nil {
		return nil, err
	}
	return &method, nil
}
