package router

import (
	"atlas_food/internal/domain/auth"
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
		// Auth routes (public)
		authHandler := auth.NewHandler(db)
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.POST("/refresh", authHandler.RefreshToken)
			authRoutes.GET("/me", middleware.JWTAuth(), authHandler.GetProfile)
		}

		// Admin routes (protected)
		surveyHandler := survey.NewHandler(db)
		adminRoutes := v1.Group("/admin", middleware.JWTAuth(), middleware.AdminOnly())
		{
			surveyRoutes := adminRoutes.Group("/surveys")
			{
				surveyRoutes.GET("", surveyHandler.List)
				surveyRoutes.POST("", surveyHandler.Create)
				surveyRoutes.GET("/:id", surveyHandler.GetByID)
				surveyRoutes.PUT("/:id", surveyHandler.Update)
				surveyRoutes.DELETE("/:id", surveyHandler.Delete)
				surveyRoutes.POST("/:id/clone", surveyHandler.Clone)
			}
		}

		// Public survey routes (respondent - menggunakan accessToken)
		publicSurveyRoutes := v1.Group("/surveys/:accessToken", middleware.SurveyAccessToken())
		{
			publicSurveyRoutes.GET("", surveyHandler.GetPublic)
			publicSurveyRoutes.POST("/join", surveyHandler.Join)
		}
		// Survey routes
		surveyHandler := survey.NewHandler(db)
		surveyHandler.SetupRoutes(v1, middleware.JWTAuth())

		// TODO: Food routes
		// TODO: Submission routes
	}

	return r
}
