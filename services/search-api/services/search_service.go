package services

import (
	"encoding/json"
	"fmt"
	"log"
	"search-api/domain"
	"search-api/repositories"
)

type SearchService struct {
	searchRepo *repositories.SearchRepository
	cacheRepo  *repositories.CacheRepository
}

func NewSearchService(searchRepo *repositories.SearchRepository, cacheRepo *repositories.CacheRepository) *SearchService {
	return &SearchService{
		searchRepo: searchRepo,
		cacheRepo:  cacheRepo,
	}
}

func (s *SearchService) SearchRooms(params domain.SearchParams) (*domain.SearchResults, error) {
	// Validar parámetros
	if params.Rows <= 0 || params.Rows > 100 {
		params.Rows = 10
	}
	if params.Start < 0 {
		params.Start = 0
	}

	// Intentar obtener del cache
	cacheKey := s.buildCacheKey(params)
	if cachedData, err := s.cacheRepo.GetString(cacheKey); err == nil {
		log.Printf("✅ Cache HIT para: %s", cacheKey)
		var results domain.SearchResults
		if err := json.Unmarshal([]byte(cachedData), &results); err == nil {
			return &results, nil
		}
	}

	log.Printf("❌ Cache MISS para: %s", cacheKey)

	// Buscar en Solr
	results, err := s.searchRepo.SearchRooms(params)
	if err != nil {
		return nil, err
	}

	// Calcular paginación
	page := (params.Start / params.Rows) + 1
	totalPages := (results.TotalResults + params.Rows - 1) / params.Rows

	results.Page = page
	results.PageSize = params.Rows
	results.TotalPages = totalPages

	// Guardar en cache
	if data, err := json.Marshal(results); err == nil {
		s.cacheRepo.SetString(cacheKey, string(data))
	}

	return results, nil
}

func (s *SearchService) GetSuggestions(prefix string, limit int) ([]string, error) {
	if limit <= 0 || limit > 20 {
		limit = 10
	}

	// Intentar obtener del cache
	cacheKey := fmt.Sprintf("suggestions:%s:%d", prefix, limit)
	if cachedData, err := s.cacheRepo.GetString(cacheKey); err == nil {
		var suggestions []string
		if err := json.Unmarshal([]byte(cachedData), &suggestions); err == nil {
			return suggestions, nil
		}
	}

	// Obtener de Solr
	suggestions, err := s.searchRepo.GetSuggestions(prefix, limit)
	if err != nil {
		return nil, err
	}

	// Guardar en cache
	if data, err := json.Marshal(suggestions); err == nil {
		s.cacheRepo.SetString(cacheKey, string(data))
	}

	return suggestions, nil
}

func (s *SearchService) GetFacets() (*domain.FacetResponse, error) {
	// Buscar con facetas activadas
	params := domain.SearchParams{
		Query:         "*:*",
		Start:         0,
		Rows:          0, // No necesitamos resultados
		IncludeFacets: true,
	}

	results, err := s.searchRepo.SearchRooms(params)
	if err != nil {
		return nil, err
	}

	// Parsear facetas
	facetResponse := &domain.FacetResponse{
		RoomTypes:     make(map[string]int),
		StatusCounts:  make(map[string]int),
		FloorCounts:   make(map[int]int),
		AmenityCounts: make(map[string]int),
	}

	if results.Facets != nil {
		// Procesar facet_fields si existen
		if facetFields, ok := results.Facets["facet_fields"].(map[string]interface{}); ok {
			// Room types
			if roomTypes, ok := facetFields["room_type"].([]interface{}); ok {
				for i := 0; i < len(roomTypes)-1; i += 2 {
					if key, ok := roomTypes[i].(string); ok {
						if count, ok := roomTypes[i+1].(float64); ok {
							facetResponse.RoomTypes[key] = int(count)
						}
					}
				}
			}

			// Status counts
			if statusCounts, ok := facetFields["status"].([]interface{}); ok {
				for i := 0; i < len(statusCounts)-1; i += 2 {
					if key, ok := statusCounts[i].(string); ok {
						if count, ok := statusCounts[i+1].(float64); ok {
							facetResponse.StatusCounts[key] = int(count)
						}
					}
				}
			}
		}
	}

	return facetResponse, nil
}

func (s *SearchService) InvalidateCache(pattern string) error {
	// Memcached no soporta invalidación por patrón
	// En producción, considera usar Redis
	log.Printf("⚠️  Invalidación de cache solicitada para patrón: %s", pattern)
	return nil
}

func (s *SearchService) buildCacheKey(params domain.SearchParams) string {
	return fmt.Sprintf("search:q=%s:minp=%.2f:maxp=%.2f:minc=%d:type=%s:avail=%t:start=%d:rows=%d",
		params.Query,
		params.MinPrice,
		params.MaxPrice,
		params.MinCapacity,
		params.RoomType,
		params.IsAvailable,
		params.Start,
		params.Rows,
	)
}
