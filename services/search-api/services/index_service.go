package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"search-api/domain"
	"search-api/repositories"
	"time"
)

type IndexService struct {
	searchRepo *repositories.SearchRepository
}

func NewIndexService(searchRepo *repositories.SearchRepository) *IndexService {
	return &IndexService{
		searchRepo: searchRepo,
	}
}

// IndexRoom indexa una habitación individual
func (s *IndexService) IndexRoom(room *domain.RoomDocument) error {
	room.LastUpdated = time.Now()

	if err := s.searchRepo.IndexRoom(room); err != nil {
		return fmt.Errorf("error indexando habitación: %w", err)
	}

	log.Printf("✅ Habitación %s indexada correctamente", room.RoomNumber)
	return nil
}

// UpdateRoomAvailability actualiza solo la disponibilidad
func (s *IndexService) UpdateRoomAvailability(room *domain.RoomDocument) error {
	// En Solr, una actualización parcial se hace con el mismo método
	// pero enviando solo los campos a actualizar
	room.LastUpdated = time.Now()
	return s.searchRepo.IndexRoom(room)
}

// DeleteRoom elimina una habitación del índice
func (s *IndexService) DeleteRoom(roomID string) error {
	if err := s.searchRepo.DeleteRoom(roomID); err != nil {
		return fmt.Errorf("error eliminando habitación del índice: %w", err)
	}

	log.Printf("✅ Habitación %s eliminada del índice", roomID)
	return nil
}

// FullReindex reindexar todas las habitaciones desde rooms-api
func (s *IndexService) FullReindex() error {
	log.Println("🔄 Iniciando reindexación completa...")

	roomsAPIURL := os.Getenv("ROOMS_API_URL")
	if roomsAPIURL == "" {
		roomsAPIURL = "http://rooms-api:8080"
	}

	// Obtener todas las habitaciones
	url := fmt.Sprintf("%s/api/v1/rooms?limit=100", roomsAPIURL)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error obteniendo habitaciones: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("rooms-api returned status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Rooms []struct {
			ID          string  `json:"id"`
			Number      string  `json:"number"`
			Type        string  `json:"type"`
			Status      string  `json:"status"`
			Price       float64 `json:"price"`
			Description string  `json:"description"`
			Capacity    int     `json:"capacity"`
			Floor       int     `json:"floor"`
			HasWifi     bool    `json:"has_wifi"`
			HasAC       bool    `json:"has_ac"`
			HasTV       bool    `json:"has_tv"`
			HasMinibar  bool    `json:"has_minibar"`
		} `json:"rooms"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("error decodificando respuesta: %w", err)
	}

	// Indexar cada habitación
	indexed := 0
	for _, room := range response.Rooms {
		amenities := []string{}
		if room.HasWifi {
			amenities = append(amenities, "WiFi")
		}
		if room.HasAC {
			amenities = append(amenities, "Aire Acondicionado")
		}
		if room.HasTV {
			amenities = append(amenities, "TV")
		}
		if room.HasMinibar {
			amenities = append(amenities, "Minibar")
		}

		doc := &domain.RoomDocument{
			ID:            room.ID,
			RoomNumber:    room.Number,
			RoomType:      room.Type,
			Capacity:      room.Capacity,
			PricePerNight: room.Price,
			Status:        room.Status,
			Description:   room.Description,
			Amenities:     amenities,
			Floor:         room.Floor,
			IsAvailable:   room.Status == "available",
			LastUpdated:   time.Now(),
		}

		if err := s.IndexRoom(doc); err != nil {
			log.Printf("❌ Error indexando habitación %s: %v", room.Number, err)
			continue
		}

		indexed++
	}

	log.Printf("✅ Reindexación completa: %d/%d habitaciones indexadas", indexed, len(response.Rooms))
	return nil
}

// GetIndexStats obtiene estadísticas del índice
func (s *IndexService) GetIndexStats() (map[string]interface{}, error) {
	params := domain.SearchParams{
		Query:         "*:*",
		Start:         0,
		Rows:          0,
		IncludeFacets: true,
	}

	results, err := s.searchRepo.SearchRooms(params)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_documents": results.TotalResults,
		"last_updated":    time.Now(),
	}

	return stats, nil
}
