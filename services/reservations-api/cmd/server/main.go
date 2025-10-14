package main

import (
	"log"
	"reservations-api/config"
	"reservations-api/controllers"
	"reservations-api/repositories"
	"reservations-api/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inicializar base de datos
	db := config.InitDatabase()
	if db == nil {
		log.Fatal("Error inicializando la base de datos")
	}
	// Inicializar capas
	reservationRepo := repositories.NewReservationRepository(db)
	reservationService := services.NewReservationService(reservationRepo)
	reservationController := controllers.NewReservationController(reservationService)
	// Configurar router
	router := gin.Default()
	// Rutas
	api := router.Group("/api")
	{
		api.POST("/reservations", reservationController.CreateReservation)
		//api.GET("/reservations/:id", reservationController.GetReservationByID)
		//api.DELETE("/reservations/:id", reservationController.CancelReservation)
	}
	// Iniciar servidor
	log.Println("Reservations API running on port 8082")
	if err := router.Run(":8082"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
