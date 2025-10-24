package main

import (
	"log"
	"rooms-api/config"
	"rooms-api/controllers"
	"rooms-api/domain"
	"rooms-api/repositories"
	"rooms-api/services"

	"github.com/gin-gonic/gin"
)

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

	// Initialize repository
	roomRepo := repositories.NewRoomRepository(db)

	// Initialize service
	roomService := services.NewRoomService(roomRepo)

	// Initialize controller
	roomController := controllers.NewRoomController(roomService)

	// Setup Gin router
	r := gin.Default()

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
		rooms := api.Group("/rooms")
		{
			rooms.POST("", roomController.CreateRoom)
			rooms.GET("", roomController.GetRooms)
			rooms.GET("/available", roomController.GetAvailableRooms)
			rooms.GET("/number/:number", roomController.GetRoomByNumber)
			rooms.GET("/:id", roomController.GetRoomByID)
			rooms.PUT("/:id", roomController.UpdateRoom)
			rooms.PATCH("/:id/status", roomController.UpdateRoomStatus)
			rooms.DELETE("/:id", roomController.DeleteRoom)
		}
	}

	// Start server
	port := ":8080"
	log.Printf("🚀 Rooms API server starting on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
