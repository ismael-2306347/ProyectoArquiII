package events

import (
	"encoding/json"
	"log"
	"search-api/domain"
	"search-api/services"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventConsumer struct {
	conn         *amqp.Connection
	indexService *services.IndexService
}

func NewEventConsumer(conn *amqp.Connection, indexService *services.IndexService) (*EventConsumer, error) {
	return &EventConsumer{
		conn:         conn,
		indexService: indexService,
	}, nil
}

func (ec *EventConsumer) Start() error {
	ch, err := ec.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

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
		return err
	}

	// Declarar cola para search-api
	queue, err := ch.QueueDeclare(
		"search_room_events",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Bind a diferentes routing keys
	routingKeys := []string{
		"room.created",
		"room.updated",
		"room.deleted",
		"room.status.changed",
		"reservation.created",
		"reservation.cancelled",
	}

	for _, key := range routingKeys {
		err = ch.QueueBind(
			queue.Name,
			key,
			"room_events",
			false,
			nil,
		)
		if err != nil {
			log.Printf("Error binding queue to %s: %v", key, err)
		}
	}

	// Consumir mensajes
	msgs, err := ch.Consume(
		queue.Name,
		"search-api-consumer",
		false, // auto-ack false
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	log.Println("✅ Event consumer iniciado, esperando eventos...")

	// Procesar mensajes
	for msg := range msgs {
		ec.handleMessage(msg)
	}

	return nil
}

func (ec *EventConsumer) handleMessage(msg amqp.Delivery) {
	log.Printf("📨 Evento recibido: %s", msg.RoutingKey)

	switch msg.RoutingKey {
	case "room.created", "room.updated":
		ec.handleRoomEvent(msg)
	case "room.deleted":
		ec.handleRoomDeleted(msg)
	case "room.status.changed":
		ec.handleRoomStatusChanged(msg)
	case "reservation.created":
		ec.handleReservationCreated(msg)
	case "reservation.cancelled":
		ec.handleReservationCancelled(msg)
	default:
		log.Printf("⚠️ Evento desconocido: %s", msg.RoutingKey)
	}

	msg.Ack(false)
}

func (ec *EventConsumer) handleRoomEvent(msg amqp.Delivery) {
	var event struct {
		RoomID      string                 `json:"room_id"`
		RoomNumber  string                 `json:"room_number"`
		RoomType    string                 `json:"room_type"`
		Capacity    int                    `json:"capacity"`
		Price       float64                `json:"price_per_night"`
		Status      string                 `json:"status"`
		Description string                 `json:"description"`
		Amenities   []string               `json:"amenities"`
		Floor       int                    `json:"floor"`
		SizeSqm     float64                `json:"size_sqm"`
		ViewType    string                 `json:"view_type"`
		IsAvailable bool                   `json:"is_available"`
		Data        map[string]interface{} `json:"data,omitempty"`
	}

	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("❌ Error deserializando evento: %v", err)
		return
	}

	// Si el evento tiene un campo "data" anidado, extraer de ahí
	if event.Data != nil {
		if roomID, ok := event.Data["id"].(string); ok {
			event.RoomID = roomID
		}
		if roomNumber, ok := event.Data["room_number"].(string); ok {
			event.RoomNumber = roomNumber
		}
		if roomType, ok := event.Data["room_type"].(string); ok {
			event.RoomType = roomType
		}
		if capacity, ok := event.Data["capacity"].(float64); ok {
			event.Capacity = int(capacity)
		}
		if price, ok := event.Data["price_per_night"].(float64); ok {
			event.Price = price
		}
		if status, ok := event.Data["status"].(string); ok {
			event.Status = status
		}
	}

	// Crear documento para Solr
	roomDoc := &domain.RoomDocument{
		ID:            event.RoomID,
		RoomNumber:    event.RoomNumber,
		RoomType:      event.RoomType,
		Capacity:      event.Capacity,
		PricePerNight: event.Price,
		Status:        event.Status,
		Description:   event.Description,
		Amenities:     event.Amenities,
		Floor:         event.Floor,
		SizeSqm:       event.SizeSqm,
		ViewType:      event.ViewType,
		IsAvailable:   event.Status == "available",
	}

	if err := ec.indexService.IndexRoom(roomDoc); err != nil {
		log.Printf("❌ Error indexando habitación %s: %v", event.RoomID, err)
	} else {
		log.Printf("✅ Habitación %s indexada correctamente", event.RoomNumber)
	}
}

func (ec *EventConsumer) handleRoomDeleted(msg amqp.Delivery) {
	var event struct {
		RoomID string `json:"room_id"`
	}

	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("❌ Error deserializando evento: %v", err)
		return
	}

	if err := ec.indexService.DeleteRoom(event.RoomID); err != nil {
		log.Printf("❌ Error eliminando habitación %s del índice: %v", event.RoomID, err)
	} else {
		log.Printf("✅ Habitación %s eliminada del índice", event.RoomID)
	}
}

func (ec *EventConsumer) handleRoomStatusChanged(msg amqp.Delivery) {
	var event struct {
		RoomID    string `json:"room_id"`
		NewStatus string `json:"new_status"`
	}

	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("❌ Error deserializando evento: %v", err)
		return
	}

	// Actualizar disponibilidad en el índice
	update := &domain.RoomDocument{
		ID:          event.RoomID,
		Status:      event.NewStatus,
		IsAvailable: event.NewStatus == "available",
	}

	if err := ec.indexService.UpdateRoomAvailability(update); err != nil {
		log.Printf("❌ Error actualizando disponibilidad: %v", err)
	} else {
		log.Printf("✅ Disponibilidad actualizada para habitación %s", event.RoomID)
	}
}

func (ec *EventConsumer) handleReservationCreated(msg amqp.Delivery) {
	var event struct {
		RoomID string `json:"room_id"`
	}

	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("❌ Error deserializando evento: %v", err)
		return
	}

	// Marcar habitación como no disponible
	update := &domain.RoomDocument{
		ID:          event.RoomID,
		IsAvailable: false,
		Status:      "reserved",
	}

	if err := ec.indexService.UpdateRoomAvailability(update); err != nil {
		log.Printf("❌ Error actualizando disponibilidad: %v", err)
	} else {
		log.Printf("✅ Habitación %s marcada como reservada", event.RoomID)
	}
}

func (ec *EventConsumer) handleReservationCancelled(msg amqp.Delivery) {
	var event struct {
		RoomID string `json:"room_id"`
	}

	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("❌ Error deserializando evento: %v", err)
		return
	}

	// Marcar habitación como disponible
	update := &domain.RoomDocument{
		ID:          event.RoomID,
		IsAvailable: true,
		Status:      "available",
	}

	if err := ec.indexService.UpdateRoomAvailability(update); err != nil {
		log.Printf("❌ Error actualizando disponibilidad: %v", err)
	} else {
		log.Printf("✅ Habitación %s marcada como disponible", event.RoomID)
	}
}
