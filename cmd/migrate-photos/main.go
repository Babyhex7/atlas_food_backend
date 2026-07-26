package main

// migrate-photos — CLI tool untuk menyalin foto Atlas Makananku ke folder uploads backend.
//
// Usage:
//   go run cmd/migrate-photos/main.go [flags]
//
// Flags:
//   -source   Path ke folder Atlas_Makananku_Photos/ (wajib)
//   -dest     Path tujuan uploads/atlas/ (default: ./uploads/atlas)
//   -dry-run  Preview saja tanpa menyalin file (default: false)
//   -sql      Generate SQL UPDATE untuk record as_served_images yang sudah ada di DB
//   -db       Jalankan SQL UPDATE langsung ke database (butuh koneksi DB)
//
// Contoh:
//   # Preview dulu
//   go run cmd/migrate-photos/main.go -source "C:/Users/Raden/Downloads/Atlas_Makananku_Photos_Struktur 2/Atlas_Makananku_Photos_Struktur/Atlas_Makananku_Photos" -dry-run
//
//   # Copy foto ke uploads/atlas/
//   go run cmd/migrate-photos/main.go -source "C:/Users/Raden/Downloads/Atlas_Makananku_Photos_Struktur 2/Atlas_Makananku_Photos_Struktur/Atlas_Makananku_Photos"
//
//   # Copy + generate SQL update untuk database
//   go run cmd/migrate-photos/main.go -source "..." -sql
//
//   # Copy + langsung update database
//   go run cmd/migrate-photos/main.go -source "..." -db
//
// Struktur output (lokal → MinIO-ready):
//   uploads/atlas/
//   ├── MP/
//   │   ├── MP-01_A.jpg        → MinIO key: photos/MP/MP-01_A.jpg
//   │   ├── MP-01_B.jpg
//   │   └── ...
//   ├── LH/
//   │   └── LH-01_A.jpg
//   └── ...

import (
	"atlas_food/internal/bootstrap"
	"atlas_food/internal/config"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// main - baca flag CLI lalu jalankan proses salin foto (opsional: dry-run, generate SQL, atau langsung update DB)
func main() {
	// --- flags ---
	sourceFlag := flag.String("source", "", "Path ke folder Atlas_Makananku_Photos/ (wajib)")
	destFlag := flag.String("dest", "", "Path tujuan uploads/atlas/ (default: ./uploads/atlas)")
	dryRun := flag.Bool("dry-run", false, "Preview saja tanpa menyalin file")
	genSQL := flag.Bool("sql", false, "Generate file SQL UPDATE untuk update image_url di database")
	applyDB := flag.Bool("db", false, "Langsung apply SQL UPDATE ke database (butuh koneksi DB)")
	flag.Parse()

	// Validasi source
	if *sourceFlag == "" {
		fmt.Println("❌ Flag -source wajib diisi.")
		fmt.Println("   Contoh: go run cmd/migrate-photos/main.go -source \"C:/path/ke/Atlas_Makananku_Photos\"")
		os.Exit(1)
	}

	// Resolve dest
	destPath := *destFlag
	if destPath == "" {
		destPath = filepath.Join(".", "uploads", "atlas")
	}

	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║         Atlas Makananku — Photo Migration Tool           ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("📂 Sumber  : %s\n", *sourceFlag)
	fmt.Printf("📁 Tujuan  : %s\n", destPath)
	if *dryRun {
		fmt.Println("🔍 Mode    : DRY-RUN (tidak ada file yang disalin)")
	} else {
		fmt.Println("🚀 Mode    : COPY (file akan disalin)")
	}
	fmt.Println()

	// --- Langkah 1: Migrasi foto ---
	fmt.Println("▶ Langkah 1: Menyalin foto ke uploads/atlas/ ...")
	result, err := bootstrap.MigratePhotos(*sourceFlag, destPath, *dryRun)
	if err != nil {
		log.Fatalf("❌ Migrasi gagal: %v", err)
	}

	// Laporan anomali yang diperbaiki
	if len(result.Fixed) > 0 {
		fmt.Println()
		fmt.Printf("🔧 Nama file yang diperbaiki (%d):\n", len(result.Fixed))
		for _, f := range result.Fixed {
			fmt.Printf("   ✏️  %s\n", f)
		}
	}

	// Laporan error
	if len(result.Errors) > 0 {
		fmt.Println()
		fmt.Printf("⚠️  Error (%d):\n", len(result.Errors))
		for _, e := range result.Errors {
			fmt.Printf("   ❌ %s\n", e)
		}
	}

	// Ringkasan copy
	fmt.Println()
	fmt.Println("📊 Ringkasan:")
	fmt.Printf("   ✅ Disalin  : %d file\n", len(result.Copied))
	fmt.Printf("   ⏭️  Dilewati : %d file (sudah ada)\n", len(result.Skipped))
	fmt.Printf("   🔧 Diperbaiki: %d file (anomali nama)\n", len(result.Fixed))
	fmt.Printf("   ❌ Error    : %d\n", len(result.Errors))

	if *dryRun {
		fmt.Println()
		fmt.Println("────────────────────────────────────────────────────────────")
		fmt.Println("ℹ️  Ini adalah DRY-RUN. Jalankan tanpa -dry-run untuk menyalin file.")
		fmt.Println("────────────────────────────────────────────────────────────")
		return
	}

	// --- Langkah 2: Generate SQL ---
	if *genSQL || *applyDB {
		fmt.Println()
		fmt.Println("▶ Langkah 2: Men-generate SQL UPDATE image_url ...")

		sqlContent, err := generateUpdateSQL(destPath)
		if err != nil {
			log.Printf("⚠️  Gagal generate SQL: %v", err)
		} else {
			sqlFile := "migrate_photos_update.sql"
			if err := os.WriteFile(sqlFile, []byte(sqlContent), 0644); err != nil {
				log.Printf("⚠️  Gagal tulis SQL file: %v", err)
			} else {
				fmt.Printf("   📄 SQL ditulis ke: %s\n", sqlFile)
			}
		}

		// --- Langkah 3: Apply ke database ---
		if *applyDB {
			fmt.Println()
			fmt.Println("▶ Langkah 3: Menerapkan perubahan ke database ...")

			_ = godotenv.Load()
			cfg := config.Load()
			db := config.InitDB(cfg)

			stmts := splitSQL(sqlContent)
			applied := 0
			for _, stmt := range stmts {
				stmt = strings.TrimSpace(stmt)
				if stmt == "" {
					continue
				}
				if err := db.Exec(stmt).Error; err != nil {
					fmt.Printf("   ⚠️  Gagal: %s\n      Error: %v\n", truncate(stmt, 80), err)
				} else {
					applied++
				}
			}
			fmt.Printf("   ✅ %d statement berhasil dijalankan\n", applied)
		}
	}

	fmt.Println()
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Println("✅ Migrasi foto selesai!")
	fmt.Println()
	fmt.Println("📌 Langkah selanjutnya:")
	fmt.Println("   1. Jalankan seeder:  go run cmd/seed/main.go")
	fmt.Println("      → Ini akan mengisi as_served_images dengan URL yang benar")
	fmt.Println()
	fmt.Println("   2. Untuk produksi (MinIO), upload isi uploads/atlas/ ke bucket:")
	fmt.Println("      mc cp --recursive ./uploads/atlas/ minio/atlas-food/photos/")
	fmt.Println("      Lalu ubah STORAGE_BASE_URL di .env ke https://minio.example.com/atlas-food")
	fmt.Println("════════════════════════════════════════════════════════════")
}

// generateUpdateSQL - buat SQL UPDATE berdasarkan file yang tersedia di uploads/atlas/
// Untuk setiap file, update as_served_images.image_url dan thumbnail_url
func generateUpdateSQL(uploadsAtlasPath string) (string, error) {
	available, err := bootstrap.ScanAvailablePhotos(uploadsAtlasPath)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString("-- Auto-generated by migrate-photos\n")
	sb.WriteString("-- Update image_url di as_served_images sesuai nama file asli\n")
	sb.WriteString("-- Dihasilkan dari scan folder: " + uploadsAtlasPath + "\n\n")
	sb.WriteString("SET FOREIGN_KEY_CHECKS=0;\n\n")

	count := 0
	for key := range available {
		// key format: "MP/MP-01_A.jpg"
		parts := strings.SplitN(key, "/", 2)
		if len(parts) != 2 {
			continue
		}
		catCode := parts[0]
		filename := parts[1]

		// Ekstrak food code dan label dari nama file
		// MP-01_A.jpg → code=MP-01, label=A
		// AS-22_guide.jpg → code=AS-22, label=guide
		nameParts := strings.SplitN(strings.TrimSuffix(filename, ".jpg"), "_", 2)
		if len(nameParts) != 2 {
			continue
		}
		foodCode := nameParts[0]
		label := nameParts[1]

		imageURL := fmt.Sprintf("/uploads/atlas/%s/%s", catCode, filename)

		// Update as_served_images berdasarkan set_id (yang terhubung ke food.code)
		// dan label yang cocok
		if label == "guide" {
			// Foto guide → update semua label dalam set yang food-nya = foodCode
			sb.WriteString(fmt.Sprintf(
				"UPDATE as_served_images asi\n"+
					"  JOIN as_served_sets ass ON ass.id = asi.set_id\n"+
					"  JOIN foods f ON f.id = ass.food_id\n"+
					"SET asi.image_url = '%s', asi.thumbnail_url = '%s'\n"+
					"WHERE f.code = '%s';\n\n",
				imageURL, imageURL, foodCode,
			))
		} else {
			sb.WriteString(fmt.Sprintf(
				"UPDATE as_served_images asi\n"+
					"  JOIN as_served_sets ass ON ass.id = asi.set_id\n"+
					"  JOIN foods f ON f.id = ass.food_id\n"+
					"SET asi.image_url = '%s', asi.thumbnail_url = '%s'\n"+
					"WHERE f.code = '%s' AND asi.label = '%s';\n\n",
				imageURL, imageURL, foodCode, strings.ToUpper(label),
			))
		}
		count++
	}

	sb.WriteString("SET FOREIGN_KEY_CHECKS=1;\n")
	sb.WriteString(fmt.Sprintf("\n-- Total: %d statements\n", count))

	return sb.String(), nil
}

// splitSQL - pisah SQL multi-statement berdasarkan ";"
func splitSQL(sql string) []string {
	var stmts []string
	for _, s := range strings.Split(sql, ";") {
		s = strings.TrimSpace(s)
		if s != "" && !strings.HasPrefix(s, "--") {
			stmts = append(stmts, s)
		}
	}
	return stmts
}

// truncate - potong string untuk ditampilkan di log
func truncate(s string, max int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}
