package controllers

import (
	"net/http"
	"reservations-api/domain"
	"reservations-api/services"

	"github.com/gin-gonic/gin"
)

type ReservationController struct {
	service services.ReservationService
}

func NewReservationController(service services.ReservationService) *ReservationController {
	return &ReservationController{service: service}
}
func (c *ReservationController) CreateReservation(ctx *gin.Context) {
	var req domain.CreateReservationDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created, err := c.service.CreateReservation(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"reservation": created})
}
