package domain

import "time"

type ReservationStatus string

const (
	ReservationStatusActive   ReservationStatus = "active"
	ReservationStatusCanceled ReservationStatus = "canceled"
)

type Reservation struct {
	ID        uint              `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint              `gorm:"not null"                json:"user_id"`
	StartDate time.Time         `gorm:"not null"                json:"start_date"`
	EndDate   time.Time         `gorm:"not null"                json:"end_date"`
	RoomID    uint              `gorm:"not null"                json:"room_id"`
	Status    ReservationStatus `gorm:"type:varchar(20);default:'active'" json:"status"`
	CreatedAt time.Time         `gorm:"autoCreateTime"          json:"created_at"`
	UpdatedAt time.Time         `gorm:"autoUpdateTime"          json:"updated_at"`
}
