package repositories

import (
	"context"
	"fmt"

	"reservations-api/domain"

	"gorm.io/gorm"
)

type ReservationRepository interface {
	Create(ctx context.Context, reservation domain.Reservation) (domain.Reservation, error)
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

// Create persiste una Reservation (entidad) y devuelve la entidad guardada.
// Nota: El formateo a DTO debe hacerse en la capa Service.
func (r *reservationRepository) Create(ctx context.Context, reservation domain.Reservation) (domain.Reservation, error) {
	// Estado por defecto si no viene seteado
	if reservation.Status == "" {
		reservation.Status = domain.ReservationStatusActive
	}

	// Usar el context en GORM
	if err := r.db.WithContext(ctx).Create(&reservation).Error; err != nil {
		return domain.Reservation{}, fmt.Errorf("failed to create reservation: %w", err)
	}

	// GORM completa ID/CreatedAt/UpdatedAt en 'reservation'
	return reservation, nil
}
