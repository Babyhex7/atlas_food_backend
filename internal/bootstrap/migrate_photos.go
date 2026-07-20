package bootstrap

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// PhotoMigrationResult - hasil dari proses migrasi foto
type PhotoMigrationResult struct {
	Copied   []string
	Skipped  []string
	Fixed    []string // file yang namanya diperbaiki (anomali double dot, dll.)
	NotFound []string // file dari JSON yang tidak ada di sumber
	Errors   []string
}

// categoryFolderMap - mapping kode kategori ke nama folder sumber
var categoryFolderMap = map[string]string{
	"MP":  "01_MP_Makanan_Pokok",
	"LH":  "02_LH_Lauk_Hewani",
	"LN":  "03_LN_Lauk_Nabati",
	"AS":  "04_AS_Aneka_Sayur",
	"AB":  "05_AB_Aneka_Buah",
	"AP":  "06_AP_Jajanan_Kue_Roti",
	"AMK": "07_AMK_Makanan_Kemasan",
	"KK":  "08_KK_Keripik_Kerupuk",
	"ABK": "09_ABK_Bumbu_Kondimen",
	"AK":  "10_AK_Makanan_Siap_Saji",
	"MDL": "11_MDL_Minyak_Lemak",
	"GK":  "12_GK_Gula_Konfeksioneri",
	"AH":  "13_AH_Alat_Ukur",
}

// anomalyPattern - regex untuk mendeteksi nama file dengan double dot atau karakter aneh
var anomalyPattern = regexp.MustCompile(`\.\.+jpg$`)

// MigratePhotos - menyalin foto dari folder sumber ke folder uploads backend
//
// sourcePath : path ke Atlas_Makananku_Photos/ (folder yang berisi 01_MP_, 02_LH_, dst.)
// destPath   : path ke uploads/atlas/ di backend (misal ./uploads/atlas)
//
// Struktur tujuan yang dihasilkan (siap dimount ke MinIO bucket):
//
//	uploads/atlas/
//	├── MP/
//	│   ├── MP-01_A.jpg
//	│   ├── MP-01_B.jpg
//	│   └── ...
//	├── LH/
//	│   └── ...
//	└── ...
func MigratePhotos(sourcePath, destPath string, dryRun bool) (*PhotoMigrationResult, error) {
	result := &PhotoMigrationResult{}

	if !dryRun {
		if err := os.MkdirAll(destPath, 0755); err != nil {
			return nil, fmt.Errorf("gagal membuat direktori tujuan: %v", err)
		}
	}

	for catCode, folderName := range categoryFolderMap {
		srcFolder := filepath.Join(sourcePath, folderName)

		// Lewati folder yang tidak ada
		if _, err := os.Stat(srcFolder); os.IsNotExist(err) {
			continue
		}

		// Baca semua file JPG di folder kategori
		entries, err := os.ReadDir(srcFolder)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("gagal baca folder %s: %v", folderName, err))
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			name := entry.Name()
			lower := strings.ToLower(name)

			// Skip non-jpg dan README
			if !strings.HasSuffix(lower, ".jpg") && !strings.HasSuffix(lower, ".jpeg") {
				continue
			}

			srcFile := filepath.Join(srcFolder, name)

			// Perbaiki anomali double dot: "AS-01_E..jpg" → "AS-01_E.jpg"
			cleanName := name
			if anomalyPattern.MatchString(name) {
				cleanName = strings.Replace(name, "..jpg", ".jpg", 1)
				result.Fixed = append(result.Fixed, fmt.Sprintf("%s → %s", name, cleanName))
			}

			// Buat subdirektori kategori di tujuan (misal: uploads/atlas/AS/)
			destFolder := filepath.Join(destPath, catCode)
			destFile := filepath.Join(destFolder, cleanName)

			if dryRun {
				result.Copied = append(result.Copied, fmt.Sprintf("[DRY-RUN] %s/%s → %s", catCode, cleanName, destFile))
				continue
			}

			// Buat folder kategori kalau belum ada
			if err := os.MkdirAll(destFolder, 0755); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("gagal buat folder %s: %v", destFolder, err))
				continue
			}

			// Skip kalau file tujuan sudah ada
			if _, err := os.Stat(destFile); err == nil {
				result.Skipped = append(result.Skipped, destFile)
				continue
			}

			// Copy file
			if err := copyFile(srcFile, destFile); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("gagal copy %s: %v", name, err))
				continue
			}

			result.Copied = append(result.Copied, fmt.Sprintf("%s/%s", catCode, cleanName))
		}
	}

	return result, nil
}

// copyFile - salin file dari src ke dst
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return out.Sync()
}

// BuildImageURL - membangun URL gambar sesuai pola nama file asli
//
// Untuk tipe series/range: /uploads/atlas/[CATCODE]/[KODE]_[LABEL].jpg
// Untuk tipe guide:        /uploads/atlas/[CATCODE]/[KODE]_guide.jpg
//
// URL ini juga merupakan object key yang akan dipakai di MinIO:
//   atlas-food/photos/[CATCODE]/[KODE]_[LABEL].jpg
func BuildImageURL(catCode, foodCode, label string, isGuide bool) string {
	if isGuide {
		return fmt.Sprintf("/uploads/atlas/%s/%s_guide.jpg", catCode, foodCode)
	}
	return fmt.Sprintf("/uploads/atlas/%s/%s_%s.jpg", catCode, foodCode, label)
}

// BuildMinIOKey - membangun object key untuk MinIO
// Bucket: atlas-food, prefix: photos/
func BuildMinIOKey(catCode, foodCode, label string, isGuide bool) string {
	if isGuide {
		return fmt.Sprintf("photos/%s/%s_guide.jpg", catCode, foodCode)
	}
	return fmt.Sprintf("photos/%s/%s_%s.jpg", catCode, foodCode, label)
}

// ScanAvailablePhotos - scan folder uploads/atlas dan return map foto yang tersedia
// Key: "CATCODE/KODE_LABEL.jpg" → true
func ScanAvailablePhotos(uploadsAtlasPath string) (map[string]bool, error) {
	available := make(map[string]bool)

	entries, err := os.ReadDir(uploadsAtlasPath)
	if err != nil {
		return available, err
	}

	for _, catDir := range entries {
		if !catDir.IsDir() {
			continue
		}
		catPath := filepath.Join(uploadsAtlasPath, catDir.Name())
		files, err := os.ReadDir(catPath)
		if err != nil {
			continue
		}
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			key := catDir.Name() + "/" + f.Name()
			available[key] = true
		}
	}

	return available, nil
}
