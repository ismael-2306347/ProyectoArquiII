package repositories

import (
	"context"
	"fmt"
	"time"

	"reservations-api/domain"

	"gorm.io/gorm"
)

type ReservationRepository interface {
	Create(ctx context.Context, reservation domain.Reservation) (domain.Reservation, error)
	Delete(ctx context.Context, id uint, reason string) error
	GetByID(ctx context.Context, id uint) (domain.Reservation, error)
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

func (r *reservationRepository) Delete(ctx context.Context, id uint, reason string) error {
	now := time.Now()
	// Solo cancela si esta activa

	res := r.db.WithContext(ctx).Model(&domain.Reservation{}).
		Where("id = ? AND status = ?", id, domain.ReservationStatusActive).
		Updates(map[string]interface{}{
			"status":        domain.ReservationStatusCanceled,
			"cancel_reason": reason,
			"deleted_at":    &now,
		})
	if res.Error != nil {
		return fmt.Errorf("cancelar reserva: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		var count int64
		if err := r.db.WithContext(ctx).Model(&domain.Reservation{}).
			Where("id = ?", id).
			Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil // Ya estaba cancelada
	}
	return nil
}
func (r *reservationRepository) GetByID(ctx context.Context, id uint) (domain.Reservation, error) {
	var res domain.Reservation
	if err := r.db.WithContext(ctx).First(&res, "id = ?").Error; err != nil {
		return domain.Reservation{}, err // puede ser gorm.ErrRecordNotFound
	}
	return res, nil
}
