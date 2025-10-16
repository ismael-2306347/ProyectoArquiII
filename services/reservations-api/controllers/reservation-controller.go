package controllers

import (
	"errors"
	"net/http"
	"reservations-api/domain"
	"reservations-api/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
func (c *ReservationController) DeleteReservation(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id64, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil || id64 == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
		return
	}

	var body domain.CancelReservationDTO
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.service.DeleteReservation(ctx, uint(id64), body.Reason)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "reserva no encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
