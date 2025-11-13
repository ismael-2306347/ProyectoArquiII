package services

import (
	"context"
	"rooms-api/domain"
	"rooms-api/events"
	"rooms-api/repositories"
)

type RoomService struct {
	roomRepo  *repositories.RoomRepository
	publisher *events.EventPublisher
	cache     repositories.RoomCacheRepository
}

func NewRoomService(roomRepo *repositories.RoomRepository, publisher *events.EventPublisher, cache repositories.RoomCacheRepository) *RoomService {
	return &RoomService{
		roomRepo:  roomRepo,
		publisher: publisher,
		cache:     cache,
	}
}

func (s *RoomService) CreateRoom(ctx context.Context, req domain.CreateRoomRequest) (*domain.RoomResponse, error) {
	room := &domain.Room{
		Number:      req.Number,
		Type:        req.Type,
		Status:      domain.RoomStatusAvailable,
		Price:       req.Price,
		Description: req.Description,
		Capacity:    req.Capacity,
		Floor:       req.Floor,
		HasWifi:     req.HasWifi,
		HasAC:       req.HasAC,
		HasTV:       req.HasTV,
		HasMinibar:  req.HasMinibar,
	}

	err := s.roomRepo.Create(ctx, room)
	if err != nil {
		return nil, err
	}

	err = s.cache.Set(ctx, room.ID, *s.roomToResponse(room))
	if err != nil {
		return nil, err
	}

	// Publicar evento de creación
	if s.publisher != nil {
		go s.publisher.PublishRoomCreated(room)
	}

	return s.roomToResponse(room), nil
}

func (s *RoomService) GetRoomByID(ctx context.Context, id uint) (*domain.RoomResponse, error) {

	cached, err := s.cache.Get(ctx, id)
	if err == nil {
		return &cached, nil
	}

	room, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Try to update the cache but don't fail the request if caching fails
	_ = s.cache.Set(ctx, room.ID, *s.roomToResponse(room))

	return s.roomToResponse(room), nil
}

func (s *RoomService) GetRoomByNumber(ctx context.Context, number string) (*domain.RoomResponse, error) {
	room, err := s.roomRepo.GetByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	return s.roomToResponse(room), nil
}

func (s *RoomService) GetRooms(ctx context.Context, filter domain.RoomFilter, page, limit int) (*domain.RoomListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	rooms, total, err := s.roomRepo.GetAll(ctx, filter, page, limit)
	if err != nil {
		return nil, err
	}

	roomResponses := make([]domain.RoomResponse, len(rooms))
	for i, room := range rooms {
		roomResponses[i] = *s.roomToResponse(&room)
	}

	return &domain.RoomListResponse{
		Rooms: roomResponses,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *RoomService) UpdateRoom(ctx context.Context, id uint, req domain.UpdateRoomRequest) (*domain.RoomResponse, error) {
	// Obtener estado anterior para detectar cambios de status
	oldRoom, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	oldStatus := oldRoom.Status

	// Actualizar
	room, err := s.roomRepo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}

	// Publicar evento de actualización
	if s.publisher != nil {
		go s.publisher.PublishRoomUpdated(room)

		// Si cambió el status, publicar evento específico
		if oldStatus != room.Status {
			go s.publisher.PublishRoomStatusChanged(room.ID, string(oldStatus), string(room.Status))
		}
	}

	return s.roomToResponse(room), nil
}

func (s *RoomService) DeleteRoom(ctx context.Context, id uint) error {
	// Obtener ID antes de eliminar
	room, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.roomRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Publicar evento de eliminación
	if s.publisher != nil {
		go s.publisher.PublishRoomDeleted(room.ID)
	}

	return nil
}

func (s *RoomService) UpdateRoomStatus(ctx context.Context, id uint, status domain.RoomStatus) (*domain.RoomResponse, error) {
	updateReq := domain.UpdateRoomRequest{
		Status: &status,
	}
	return s.UpdateRoom(ctx, id, updateReq)
}

func (s *RoomService) GetAvailableRooms(ctx context.Context, filter domain.RoomFilter, page, limit int) (*domain.RoomListResponse, error) {
	availableStatus := domain.RoomStatusAvailable
	filter.Status = &availableStatus
	return s.GetRooms(ctx, filter, page, limit)
}

func (s *RoomService) GetRoomsByType(ctx context.Context, roomType domain.RoomType, page, limit int) (*domain.RoomListResponse, error) {
	filter := domain.RoomFilter{
		Type: &roomType,
	}
	return s.GetRooms(ctx, filter, page, limit)
}

func (s *RoomService) GetRoomsByFloor(ctx context.Context, floor int, page, limit int) (*domain.RoomListResponse, error) {
	filter := domain.RoomFilter{
		Floor: &floor,
	}
	return s.GetRooms(ctx, filter, page, limit)
}

func (s *RoomService) SearchRooms(ctx context.Context, query string, page, limit int) (*domain.RoomListResponse, error) {
	filter := domain.RoomFilter{}
	return s.GetRooms(ctx, filter, page, limit)
}

func (s *RoomService) roomToResponse(room *domain.Room) *domain.RoomResponse {
	return &domain.RoomResponse{
		ID:          room.ID,
		Number:      room.Number,
		Type:        room.Type,
		Status:      room.Status,
		Price:       room.Price,
		Description: room.Description,
		Capacity:    room.Capacity,
		Floor:       room.Floor,
		HasWifi:     room.HasWifi,
		HasAC:       room.HasAC,
		HasTV:       room.HasTV,
		HasMinibar:  room.HasMinibar,
		CreatedAt:   room.CreatedAt,
		UpdatedAt:   room.UpdatedAt,
	}
}
