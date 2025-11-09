package events

import (
	"encoding/json"
	"fmt"
	"log"
	"rooms-api/domain"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewEventPublisher(conn *amqp.Connection) (*EventPublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declarar exchange
	err = ch.ExchangeDeclare(
		"room_events",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &EventPublisher{
		conn:    conn,
		channel: ch,
	}, nil
}

func (p *EventPublisher) PublishRoomCreated(room *domain.Room) error {
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

	event := map[string]interface{}{
		"event_type":      "room.created",
		"room_id":         fmt.Sprintf("%d", room.ID),
		"room_number":     room.Number,
		"room_type":       string(room.Type),
		"capacity":        room.Capacity,
		"price_per_night": room.Price,
		"status":          string(room.Status),
		"description":     room.Description,
		"amenities":       amenities,
		"floor":           room.Floor,
		"is_available":    room.Status == domain.RoomStatusAvailable,
	}

	return p.publish("room.created", event)
}

func (p *EventPublisher) PublishRoomUpdated(room *domain.Room) error {
	event := map[string]interface{}{
		"event_type":      "room.updated",
		"room_id":         fmt.Sprintf("%d", room.ID),
		"room_number":     room.Number,
		"room_type":       room.Type,
		"capacity":        room.Capacity,
		"price_per_night": room.Price,
		"status":          room.Status,
		"description":     room.Description,
		"floor":           room.Floor,
		"is_available":    room.Status == "available",
	}

	return p.publish("room.updated", event)
}

func (p *EventPublisher) PublishRoomDeleted(roomID uint) error {
	event := map[string]interface{}{
		"event_type": "room.deleted",
		"room_id":    fmt.Sprintf("%d", roomID),
	}

	return p.publish("room.deleted", event)
}

func (p *EventPublisher) PublishRoomStatusChanged(roomID uint, oldStatus, newStatus string) error {
	event := map[string]interface{}{
		"event_type": "room.status.changed",
		"room_id":    fmt.Sprintf("%d", roomID),
		"old_status": oldStatus,
		"new_status": newStatus,
	}

	return p.publish("room.status.changed", event)
}

func (p *EventPublisher) publish(routingKey string, event map[string]interface{}) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error marshaling event: %w", err)
	}

	err = p.channel.Publish(
		"room_events",
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		return fmt.Errorf("error publishing event: %w", err)
	}

	log.Printf("✅ Evento publicado: %s", routingKey)
	return nil
}

func (p *EventPublisher) Close() error {
	if p.channel != nil {
		return p.channel.Close()
	}
	return nil
}
