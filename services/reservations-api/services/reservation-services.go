package services

import (
	"context"
	"fmt"
	"reservations-api/domain"
	"reservations-api/repositories"
)

type ReservationService interface {
	CreateReservation(ctx context.Context, dto domain.CreateReservationDTO) (domain.ReservationResponseDTO, error)
	DeleteReservation(ctx context.Context, id uint, reason string) error
	GetReservationByID(ctx context.Context, id uint) (domain.ReservationResponseDTO, error)
}
type reservationService struct {
	repository repositories.ReservationRepository
}

func NewReservationService(repository repositories.ReservationRepository) ReservationService {
	return &reservationService{repository: repository}
}

func (s *reservationService) CreateReservation(
	ctx context.Context,
	dto domain.CreateReservationDTO,
) (domain.ReservationResponseDTO, error) {

	// Validaciones simples con strings
	if dto.StartDate == "" || dto.EndDate == "" {
		return domain.ReservationResponseDTO{}, fmt.Errorf("start_date y end_date son requeridos")
	}

	if dto.EndDate <= dto.StartDate {
		return domain.ReservationResponseDTO{}, fmt.Errorf("end_date debe ser posterior a start_date")
	}

	// Mapear DTO → entidad (sin conversión)
	entity := domain.Reservation{
		UserID:    dto.UserID,
		RoomID:    dto.RoomID,
		StartDate: dto.StartDate, // sigue siendo string
		EndDate:   dto.EndDate,
		Status:    domain.ReservationStatusActive,
	}

	// Guardar en base de datos
	saved, err := s.repository.Create(ctx, entity)
	if err != nil {
		return domain.ReservationResponseDTO{}, fmt.Errorf("failed to create reservation: %w", err)
	}

	// Mapear entidad → DTO de respuesta (sin Format)
	resp := domain.ReservationResponseDTO{
		ID:        saved.ID,
		UserID:    saved.UserID,
		RoomID:    saved.RoomID,
		StartDate: saved.StartDate,
		EndDate:   saved.EndDate,
		Status:    string(saved.Status),
	}

	return resp, nil
}

func (s *reservationService) DeleteReservation(ctx context.Context, id uint, reason string) error {
	if id == 0 {
		return fmt.Errorf("id invalido")
	}
	if reason == "" {
		return fmt.Errorf("reason es requerido")
	}
	return s.repository.Delete(ctx, id, reason)
}

func (s *reservationService) GetReservationByID(ctx context.Context, id uint) (domain.ReservationResponseDTO, error) {
	res, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return domain.ReservationResponseDTO{}, err // puede ser gorm.ErrRecordNotFound
	}
	dto := domain.ReservationResponseDTO{
		ID:        res.ID,
		UserID:    res.UserID,
		RoomID:    res.RoomID,
		StartDate: res.StartDate,
		EndDate:   res.EndDate,
		Status:    string(res.Status),
	}
	return dto, nil
}
