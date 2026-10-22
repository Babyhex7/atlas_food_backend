package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"atlas_food/internal/config"
	"atlas_food/internal/infra/storage"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	sourceDir := flag.String("source", "./uploads", "Folder sumber lokal yang akan di-upload ke R2")
	prefix := flag.String("prefix", "", "Prefix opsional di R2 (contoh: uploads atau kosong)")
	dryRun := flag.Bool("dry-run", false, "Hanya preview file yang akan diupload")
	flag.Parse()

	endpoint := cfg.StorageEndpoint
	if strings.HasPrefix(endpoint, "https://") {
		endpoint = strings.TrimPrefix(endpoint, "https://")
	} else if strings.HasPrefix(endpoint, "http://") {
		endpoint = strings.TrimPrefix(endpoint, "http://")
	}
	endpoint = strings.TrimSuffix(endpoint, "/")
	if idx := strings.Index(endpoint, "/"); idx != -1 {
		endpoint = endpoint[:idx]
	}

	if cfg.StorageAccessKeyID == "" || cfg.StorageSecretKey == "" || endpoint == "" {
		log.Fatalf("❌ Kredensial R2 belum lengkap di .env (periksa STORAGE_ACCESS_KEY_ID, STORAGE_SECRET_KEY, STORAGE_ENDPOINT)")
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.StorageAccessKeyID, cfg.StorageSecretKey, ""),
		Secure: cfg.StorageUseSSL,
	})
	if err != nil {
		log.Fatalf("❌ Gagal membuat client R2: %v", err)
	}

	ctx := context.Background()
	bucket := cfg.StorageBucket

	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║       Atlas Food — Cloudflare R2 Storage Sync Tool       ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Printf("📦 Bucket     : %s\n", bucket)
	fmt.Printf("🌐 Endpoint   : %s\n", endpoint)
	fmt.Printf("🔗 Public URL : %s\n", cfg.StoragePublicURL)
	fmt.Printf("📂 Folder     : %s\n", *sourceDir)
	if *dryRun {
		fmt.Println("🔍 Mode       : DRY-RUN (tidak ada file yang diunggah)")
	} else {
		fmt.Println("🚀 Mode       : UPLOAD AKTIF")
	}
	fmt.Println()

	var filesToUpload []string
	err = filepath.WalkDir(*sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			filesToUpload = append(filesToUpload, path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("❌ Gagal membaca folder %s: %v", *sourceDir, err)
	}

	if len(filesToUpload) == 0 {
		fmt.Printf("ℹ️ Tidak ada file ditemukan di %s.\n", *sourceDir)
		return
	}

	fmt.Printf("Ditemukan %d file untuk disinkronkan ke R2...\n\n", len(filesToUpload))

	uploaded := 0
	failed := 0

	for _, localPath := range filesToUpload {
		relPath, err := filepath.Rel(*sourceDir, localPath)
		if err != nil {
			relPath = localPath
		}
		relPath = filepath.ToSlash(relPath)

		r2Key := relPath
		if *prefix != "" {
			r2Key = strings.Trim(*prefix, "/") + "/" + strings.TrimLeft(relPath, "/")
		}

		ext := filepath.Ext(localPath)
		contentType := mime.TypeByExtension(ext)
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		if *dryRun {
			fmt.Printf("  [DRY-RUN] %s ➔ %s (%s)\n", localPath, r2Key, contentType)
			uploaded++
			continue
		}

		file, err := os.Open(localPath)
		if err != nil {
			fmt.Printf("  ❌ Gagal buka %s: %v\n", localPath, err)
			failed++
			continue
		}

		stat, err := file.Stat()
		if err != nil {
			file.Close()
			failed++
			continue
		}

		_, err = client.PutObject(ctx, bucket, r2Key, file, stat.Size(), minio.PutObjectOptions{
			ContentType: contentType,
		})
		file.Close()

		if err != nil {
			fmt.Printf("  ❌ Gagal upload %s ke R2: %v\n", r2Key, err)
			failed++
		} else {
			fmt.Printf("  ✅ Uploaded: %s ➔ %s (%d bytes)\n", localPath, r2Key, stat.Size())
			uploaded++
		}
	}

	fmt.Println()
	fmt.Printf("Selesai! Berhasil: %d, Gagal: %d\n", uploaded, failed)
	_ = time.Second
	_ = storage.Config{}
}
