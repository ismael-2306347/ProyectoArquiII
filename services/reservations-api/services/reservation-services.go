package services

import (
	"context"
	"fmt"
	"reservations-api/domain"
	"reservations-api/repositories"
)

type ReservationService interface {
	CreateReservation(ctx context.Context, dto domain.CreateReservationDTO) (domain.ReservationResponseDTO, error)
}
type reservationService struct {
	repository repositories.ReservationRepository
}

func NewReservationService(repository repositories.ReservationRepository) ReservationService {
	return &reservationService{repository: repository}
}
func (s *reservationService) CreateReservation(ctx context.Context, dto domain.CreateReservationDTO) (domain.ReservationResponseDTO, error) {
	// Validación lógica básica
	if !dto.EndDate.After(dto.StartDate) {
		return domain.ReservationResponseDTO{}, fmt.Errorf("end_date debe ser posterior a start_date")
	}

	// Mapear DTO → Entidad de dominio
	newReservation := domain.Reservation{
		UserID:    dto.UserID,
		RoomID:    dto.RoomID,
		StartDate: dto.StartDate,
		EndDate:   dto.EndDate,
		Status:    domain.ReservationStatusActive, // valor por defecto
	}

	// Guardar en base de datos
	saved, err := s.repository.Create(ctx, newReservation)
	if err != nil {
		return domain.ReservationResponseDTO{}, fmt.Errorf("failed to create reservation: %w", err)
	}

	// Mapear entidad → DTO de respuesta
	response := domain.ReservationResponseDTO{
		ID:        saved.ID,
		UserID:    saved.UserID,
		RoomID:    saved.RoomID,
		StartDate: saved.StartDate,
		EndDate:   saved.EndDate,
		Status:    string(saved.Status),
	}

	return response, nil
}
