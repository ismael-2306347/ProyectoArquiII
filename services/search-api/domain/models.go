package domain

import "time"

// RoomDocument representa un documento de habitación en Solr
type RoomDocument struct {
	ID            string    `json:"id"`
	RoomNumber    string    `json:"room_number"`
	RoomType      string    `json:"room_type"`
	Capacity      int       `json:"capacity"`
	PricePerNight float64   `json:"price_per_night"`
	Status        string    `json:"status"`
	Description   string    `json:"description,omitempty"`
	Amenities     []string  `json:"amenities,omitempty"`
	Floor         int       `json:"floor"`
	SizeSqm       float64   `json:"size_sqm,omitempty"`
	ViewType      string    `json:"view_type,omitempty"`
	IsAvailable   bool      `json:"is_available"`
	LastUpdated   time.Time `json:"last_updated,omitempty"`
}

// SearchParams representa los parámetros de búsqueda
type SearchParams struct {
	Query         string  `json:"query"`
	MinPrice      float64 `json:"min_price"`
	MaxPrice      float64 `json:"max_price"`
	MinCapacity   int     `json:"min_capacity"`
	RoomType      string  `json:"room_type"`
	Status        string  `json:"status"`
	IsAvailable   bool    `json:"is_available"`
	Floor         int     `json:"floor"`
	Start         int     `json:"start"`
	Rows          int     `json:"rows"`
	Sort          string  `json:"sort"`
	IncludeFacets bool    `json:"include_facets"`
	HasWifi       bool    `json:"has_wifi"`
	HasAC         bool    `json:"has_ac"`
	HasTV         bool    `json:"has_tv"`
	HasMinibar    bool    `json:"has_minibar"`
}

// SearchResults representa los resultados de búsqueda
type SearchResults struct {
	TotalResults int                    `json:"total_results"`
	Results      []RoomDocument         `json:"results"`
	Facets       map[string]interface{} `json:"facets,omitempty"`
	Page         int                    `json:"page"`
	PageSize     int                    `json:"page_size"`
	TotalPages   int                    `json:"total_pages"`
}

// SuggestionResponse representa sugerencias de autocompletado
type SuggestionResponse struct {
	Suggestions []string `json:"suggestions"`
}

// FacetResponse representa facetas para filtros
type FacetResponse struct {
	RoomTypes     map[string]int `json:"room_types"`
	StatusCounts  map[string]int `json:"status_counts"`
	PriceRanges   []PriceRange   `json:"price_ranges"`
	FloorCounts   map[int]int    `json:"floor_counts"`
	AmenityCounts map[string]int `json:"amenity_counts"`
}

type PriceRange struct {
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Count int     `json:"count"`
}
