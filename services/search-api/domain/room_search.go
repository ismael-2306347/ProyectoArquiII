package domain

import (
	"fmt"
	"strconv"
	"time"
)

// SolrRoomDocument representa el documento indexado en Solr
// Nota: Solr devuelve arrays para campos multivaluados, usamos el primer elemento
type SolrRoomDocument struct {
	ID          string    `json:"id"`
	Number      []int     `json:"number"`       // Array en Solr
	Type        []string  `json:"type"`         // Array en Solr
	Status      []string  `json:"status"`       // Array en Solr
	Price       []float64 `json:"price"`        // Array en Solr
	Description []string  `json:"description,omitempty"` // Array en Solr
	Capacity    []int     `json:"capacity"`     // Array en Solr
	Floor       []int     `json:"floor"`        // Array en Solr
	HasWifi     []bool    `json:"has_wifi"`     // Array en Solr
	HasAC       []bool    `json:"has_ac"`       // Array en Solr
	HasTV       []bool    `json:"has_tv"`       // Array en Solr
	HasMinibar  []bool    `json:"has_minibar"`  // Array en Solr
	CreatedAt   []string  `json:"created_at"`   // Array de strings RFC3339
	UpdatedAt   []string  `json:"updated_at"`   // Array de strings RFC3339
}

// Room representa el modelo de Room de rooms-api (para mapeo)
type Room struct {
	ID          uint      `json:"id"`
	Number      string    `json:"number"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	Capacity    int       `json:"capacity"`
	Floor       int       `json:"floor"`
	HasWifi     bool      `json:"has_wifi"`
	HasAC       bool      `json:"has_ac"`
	HasTV       bool      `json:"has_tv"`
	HasMinibar  bool      `json:"has_minibar"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RoomEvent representa un evento de RabbitMQ para rooms
type RoomEvent struct {
	EventType string `json:"event_type"` // "created", "updated", "deleted"
	RoomID    uint   `json:"room_id"`
	Timestamp string `json:"timestamp"` // RFC3339 format
}

// SolrRoomWrite representa el documento para escritura en Solr (valores simples)
type SolrRoomWrite struct {
	ID          string `json:"id"`
	Number      int    `json:"number"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Price       float64 `json:"price"`
	Description string `json:"description,omitempty"`
	Capacity    int    `json:"capacity"`
	Floor       int    `json:"floor"`
	HasWifi     bool   `json:"has_wifi"`
	HasAC       bool   `json:"has_ac"`
	HasTV       bool   `json:"has_tv"`
	HasMinibar  bool   `json:"has_minibar"`
	CreatedAt   string `json:"created_at"` // RFC3339
	UpdatedAt   string `json:"updated_at"` // RFC3339
}

// ToSolrDocument convierte Room a SolrRoomWrite para escritura en Solr
func (r *Room) ToSolrDocument() *SolrRoomWrite {
	// Convertir Number string a int
	numberInt, _ := strconv.Atoi(r.Number)

	return &SolrRoomWrite{
		ID:          fmt.Sprintf("%d", r.ID), // Convertir uint a string
		Number:      numberInt,
		Type:        r.Type,
		Status:      r.Status,
		Price:       r.Price,
		Description: r.Description,
		Capacity:    r.Capacity,
		Floor:       r.Floor,
		HasWifi:     r.HasWifi,
		HasAC:       r.HasAC,
		HasTV:       r.HasTV,
		HasMinibar:  r.HasMinibar,
		CreatedAt:   r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   r.UpdatedAt.Format(time.RFC3339),
	}
}
