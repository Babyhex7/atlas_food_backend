package router

import (
	"testing"

	"atlas_food/internal/config"
	"atlas_food/internal/domain/collab"

	"github.com/gin-gonic/gin"
)

// TestSetupRegistersRoutesWithoutConflict - gin panic saat ada konflik route
// (misal dua nama parameter berbeda di posisi yang sama). Konflik itu tidak
// terdeteksi `go build`, hanya saat registrasi berjalan — jadi test ini
// memanggil Setup untuk memastikan seluruh domain bisa didaftarkan bersama.
func TestSetupRegistersRoutesWithoutConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("registrasi route panic: %v", r)
		}
	}()

	// db nil aman: Setup hanya menyimpan handle dan mendaftarkan route,
	// tidak ada query yang dijalankan.
	engine := Setup(nil, &config.Config{}, collab.NewHub())

	if engine == nil {
		t.Fatal("Setup mengembalikan engine nil")
	}

	// Pastikan endpoint kunci brief §7.4 dan §7.5 benar-benar terdaftar
	want := map[string]string{
		"PUT":  "/api/v1/admin/food-images/:id/areas",
		"GET":  "/api/v1/public/food-images/:id",
		"POST": "/api/v1/admin/categories",
	}

	registered := map[string]bool{}
	for _, route := range engine.Routes() {
		registered[route.Method+" "+route.Path] = true
	}

	for method, path := range want {
		if !registered[method+" "+path] {
			t.Errorf("route %s %s tidak terdaftar", method, path)
		}
	}
}
