package router

import (
	"atlas_food/internal/domain/auth"
	"atlas_food/internal/domain/food"
	"atlas_food/internal/domain/survey"
	"atlas_food/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Setup - mengkonfigurasi dan mengembalikan Gin router dengan semua route
// db: koneksi database GORM untuk diinject ke handler
func Setup(db *gorm.DB) *gin.Engine {
	// Set mode Gin (debug/release)
	gin.SetMode(gin.DebugMode)

	// Buat router baru
	r := gin.New()

	// Global middleware
	r.Use(gin.Recovery())            // Recovery dari panic
	r.Use(middleware.Logger())       // Log setiap request
	r.Use(middleware.CORS())         // CORS handling
	r.Use(middleware.ErrorHandler()) // Global error handling

	// Health check endpoint (tanpa auth)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "atlas-food-api"})
	})

	// API v1 group
	v1 := r.Group("/api/v1")
	{
		// Auth routes
		authHandler := auth.NewHandler(db)
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.RefreshToken)
			authGroup.GET("/me", middleware.JWTAuth(), authHandler.GetProfile)
		}

		// Survey domain
		surveyRepo := survey.NewRepository(db)
		surveyService := survey.NewService(surveyRepo)
		surveyHandler := survey.NewHandler(surveyService)
		surveyHandler.SetupRoutes(v1, middleware.JWTAuth())

		// Food domain
		foodRepo := food.NewRepository(db)
		foodService := food.NewService(foodRepo)
		foodHandler := food.NewHandler(foodService)
		foodHandler.SetupRoutes(v1, middleware.JWTAuth())

		// TODO: Submission domain
	}

	return r
}
