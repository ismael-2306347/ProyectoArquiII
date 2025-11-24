package main

import (
	"log"
	"os"
	"rooms-api/config"
	"rooms-api/controllers"
	"rooms-api/domain"
	"rooms-api/events"
	"rooms-api/repositories"
	"rooms-api/services"
	"rooms-api/utils"

	"github.com/gin-gonic/gin"
)

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func main() {
	// Initialize MySQL connection
	db := config.InitMySQL()

	// Auto migrate the schema
	if err := db.AutoMigrate(&domain.Room{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	defer sqlDB.Close()

	rabbitURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/")

	var publisher *events.EventPublisher

	rabbitConn, err := config.InitRabbitMQ(rabbitURL)
	if err != nil {
		log.Printf("⚠️  Warning: No se pudo conectar a RabbitMQ: %v", err)
		log.Println("⚠️  Los eventos no serán publicados")
	} else {
		defer rabbitConn.Close()

		publisher, err = events.NewEventPublisher(rabbitConn)
		if err != nil {
			log.Printf("⚠️  Warning: Error inicializando event publisher: %v", err)
		} else {
			defer publisher.Close()
			log.Println("✅ Event publisher inicializado correctamente")
		}
	}

	// Initialize repository
	roomRepo := repositories.NewRoomRepository(db)

	// Initialize service (ahora con publisher)
	roomService := services.NewRoomService(roomRepo, publisher)

	// Initialize controller
	roomController := controllers.NewRoomController(roomService)

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"estado":   "saludable",
			"servicio": "rooms-api",
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Public routes (read-only)
		rooms := api.Group("/rooms")
		{
			rooms.GET("", roomController.GetRooms)
			rooms.GET("/available", roomController.GetAvailableRooms)
			rooms.GET("/number/:number", roomController.GetRoomByNumber)
			rooms.GET("/:id", roomController.GetRoomByID)
		}

		// Protected admin routes
		admin := api.Group("/admin")
		admin.Use(utils.AuthMiddleware())
		admin.Use(utils.AdminMiddleware())
		{
			admin.POST("/rooms", roomController.CreateRoom)
			admin.PUT("/rooms/:id", roomController.UpdateRoom)
			admin.PATCH("/rooms/:id/status", roomController.UpdateRoomStatus)
			admin.DELETE("/rooms/:id", roomController.DeleteRoom)
			admin.GET("/rooms", roomController.GetRooms)
			admin.GET("/rooms/:id", roomController.GetRoomByID)
		}
	}

	// Start server
	port := ":8080"
	log.Printf("🚀 Rooms API server starting on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
