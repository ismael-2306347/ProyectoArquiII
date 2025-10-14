package domain

type ReservationStatus string

const (
	ReservationStatusActive   ReservationStatus = "active"
	ReservationStatusCanceled ReservationStatus = "canceled"
)

type Reservation struct {
	ID        uint              `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint              `gorm:"not null"                json:"user_id"`
	StartDate string            `gorm:"not null"                json:"start_date"`
	EndDate   string            `gorm:"not null"                json:"end_date"`
	RoomID    uint              `gorm:"not null"                json:"room_id"`
	Status    ReservationStatus `gorm:"type:varchar(20);default:'active'" json:"status"`
}
