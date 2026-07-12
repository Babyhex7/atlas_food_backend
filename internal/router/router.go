package router

import (
	"atlas_food/internal/config"
	"atlas_food/internal/domain/ai"
	"atlas_food/internal/domain/auth"
	"atlas_food/internal/domain/collab"
	"atlas_food/internal/domain/food"
	"atlas_food/internal/domain/submission"
	"atlas_food/internal/domain/survey"
	"atlas_food/internal/domain/upload"
	"atlas_food/internal/pkg/groq"
	"atlas_food/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Setup - mengkonfigurasi dan mengembalikan Gin router dengan semua route
// db: koneksi database GORM untuk diinject ke handler
// cfg: konfigurasi aplikasi (dipakai antara lain untuk generate link survey)
// hub: WebSocket hub untuk real-time collaboration
func Setup(db *gorm.DB, cfg *config.Config, hub *collab.Hub) *gin.Engine {
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
		// ======== PUBLIC ROUTES (NO AUTH) ========
		
		// Public Food Routes (Find Your Food feature)
		foodRepo := food.NewRepository(db)
		publicFoodHandler := food.NewPublicHandler(foodRepo)
		
		publicGroup := v1.Group("/public")
		{
			// Food search & browse
			publicGroup.GET("/foods/search", publicFoodHandler.SearchFoods)
			publicGroup.GET("/foods/:id", publicFoodHandler.GetFoodDetail)
			
			// Categories
			publicGroup.GET("/categories", publicFoodHandler.GetCategories)
			publicGroup.GET("/categories/:code/foods", publicFoodHandler.GetFoodsByCategory)
		}
		
		// ======== AUTHENTICATED ROUTES ========
		
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
		surveyService := survey.NewService(surveyRepo, cfg.FrontendURL)
		surveyHandler := survey.NewHandler(surveyService)
		surveyHandler.SetupRoutes(v1, middleware.JWTAuth())

		// Food domain (authenticated)
		foodService := food.NewService(foodRepo)
		foodHandler := food.NewHandler(foodService)
		foodHandler.SetupRoutes(v1, middleware.JWTAuth())

		// Submission domain
		subRepo := submission.NewRepository(db)
		subService := submission.NewService(subRepo, surveyRepo, foodRepo)
		subHandler := submission.NewHandler(subService)
		subHandler.SetupRoutes(v1, middleware.JWTAuth())

		// AI domain
		aiRepo := ai.NewRepository(db)
		groqClient := groq.NewClient(cfg.GroqAPIKey, cfg.GroqModel, cfg.GroqBaseURL, cfg.GroqTimeoutSecs, cfg.GroqMaxTokens)
		aiService := ai.NewService(aiRepo, groqClient)
		aiHandler := ai.NewHandler(aiService)
		aiHandler.SetupRoutes(v1, middleware.JWTAuth())

		// Upload domain
		uploadHandler := upload.NewHandler("./uploads")
		uploadHandler.SetupRoutes(v1, middleware.JWTAuth())
		
		// ======== WEBSOCKET COLLABORATION ROUTES ========
		collabHandler := collab.NewHandler(hub)
		collabGroup := v1.Group("/collab")
		collabGroup.Use(middleware.JWTAuth()) // WebSocket requires authentication
		{
			collabGroup.GET("/rooms/:room_id/ws", collabHandler.HandleWebSocket)
			collabGroup.GET("/rooms/:room_id", collabHandler.GetRoomInfo)
			collabGroup.GET("/stats", collabHandler.GetHubStats)
		}
	}

	// Serve static files (uploads)
	r.Static("/uploads", "./uploads")

	return r
}
