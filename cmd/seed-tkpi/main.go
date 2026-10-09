package main

import (
	"fmt"
	"log"
	"os"

	"atlas_food/internal/bootstrap"
	"atlas_food/internal/config"

	"github.com/joho/godotenv"
)

// main - CLI pengisian nilai gizi (food_nutrients) dari TKPI 2020.
//
// Jalankan SETELAH makanan Atlas di-seed (go run cmd/seed/main.go).
//
// Env opsional:
//
//	TKPI_JSON_PATH           default data/tkpi/tkpi_2020.json
//	TKPI_MAPPING_PATH        default data/tkpi/atlas_tkpi_mapping.json
//	TKPI_INCLUDE_UNVERIFIED  "true" untuk ikut menerapkan mapping berstatus perlu_cek
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using environment variables")
	}

	cfg := config.Load()
	db := config.InitDB(cfg)

	opts := bootstrap.DefaultTKPISeedOptions()
	fmt.Printf("🧪 Seeding gizi TKPI dari %s (mapping: %s)\n", opts.TKPIPath, opts.MappingPath)
	if opts.IncludeUnverified {
		fmt.Println("⚠️  TKPI_INCLUDE_UNVERIFIED=true — mapping 'perlu_cek' ikut diterapkan")
	}

	res, err := bootstrap.SeedTKPINutrients(db, opts)
	if err != nil {
		log.Fatalf("❌ Gagal seeding TKPI: %v", err)
	}
	bootstrap.PrintTKPISeedResult(res)

	if len(res.Applied) == 0 {
		os.Exit(1)
	}
}
