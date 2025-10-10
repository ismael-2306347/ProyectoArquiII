package main

import (
	"log"
	"users-api/config"
	"users-api/controllers"
	"users-api/repositories"
	"users-api/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inicializar base de datos
	db := config.InitDatabase()
	if db == nil {
		log.Fatal("Error inicializando la base de datos")
	}
	// Inicializar capas
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// Configurar router
	router := gin.Default()

	// Rutas
	api := router.Group("/api")
	{
		api.POST("/users", userController.CreateUser)
		api.GET("/users/:id", userController.GetUserByID)
		api.POST("/login", userController.Login)
	}

	// Iniciar servidor
	log.Println("Users API running on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
