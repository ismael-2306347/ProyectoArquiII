package domain

import "time"

// Gin parsea strings JSON a time.Time usando time_format
type CreateReservationDTO struct {
	UserID    uint      `json:"user_id"    binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required" time_format:"2006-01-02"`
	EndDate   time.Time `json:"end_date"   binding:"required" time_format:"2006-01-02"`
	RoomID    uint      `json:"room_id"    binding:"required"`
}

type ReservationResponseDTO struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	RoomID    uint      `json:"room_id"`
	Status    string    `json:"status"`
}

type CancelReservationDTO struct {
	Reason string `json:"reason" binding:"required"`
}
