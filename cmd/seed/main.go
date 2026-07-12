package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"atlas_food/internal/bootstrap"
	"atlas_food/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using environment variables")
	}

	// Initialize config
	cfg := config.Load()

	// Initialize database
	db := config.InitDB(cfg)

	fmt.Println("🌱 Starting Find Your Food data seeding...")
	fmt.Println("============================================================")

	// Get JSON file path
	jsonFilePath := os.Getenv("ATLAS_JSON_PATH")
	if jsonFilePath == "" {
		// Default to project root
		jsonFilePath = filepath.Join(".", "Atlas_Makananku_FINAL.json")
	}

	// Check if file exists
	if _, err := os.Stat(jsonFilePath); os.IsNotExist(err) {
		log.Fatalf("❌ JSON file not found: %s", jsonFilePath)
	}

	fmt.Printf("📄 Reading data from: %s\n", jsonFilePath)
	
	// Seed data
	if err := bootstrap.SeedFindYourFoodData(db, jsonFilePath); err != nil {
		log.Fatalf("❌ Failed to seed data: %v", err)
	}

	fmt.Println("============================================================")
	fmt.Println("✅ Seeding completed successfully!")
	fmt.Println("")
	fmt.Println("📊 Summary:")
	
	// Get counts
	var foodCount, categoryCount, photoCount int64
	db.Table("foods").Count(&foodCount)
	db.Table("categories").Count(&categoryCount)
	db.Table("as_served_images").Count(&photoCount)
	
	fmt.Printf("   - Categories: %d\n", categoryCount)
	fmt.Printf("   - Foods: %d\n", foodCount)
	fmt.Printf("   - Portion Photos: %d\n", photoCount)
	
	fmt.Println("")
	fmt.Println("🚀 You can now start the API server:")
	fmt.Println("   go run cmd/api/main.go")
}
