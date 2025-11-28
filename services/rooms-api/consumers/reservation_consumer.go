package consumers

import (
	"context"
	"encoding/json"
	"log"
	"rooms-api/domain"
	"rooms-api/services"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ReservationConsumer escucha eventos de reserva desde RabbitMQ
type ReservationConsumer struct {
	channel      *amqp.Channel
	roomService  *services.RoomService
	exchangeName string
	queueName    string
}

// ReservationEvent representa el evento de reserva desde reservations-api
type ReservationEvent struct {
	EventType     string  `json:"event_type"`
	ReservationID string  `json:"reservation_id"`
	UserID        uint    `json:"user_id"`
	RoomID        uint    `json:"room_id"`
	StartDate     string  `json:"start_date"`
	EndDate       string  `json:"end_date"`
	Status        string  `json:"status"`
	CancelReason  *string `json:"cancel_reason,omitempty"`
}

// NewReservationConsumer crea un nuevo consumer de reservas
func NewReservationConsumer(channel *amqp.Channel, roomService *services.RoomService) *ReservationConsumer {
	return &ReservationConsumer{
		channel:      channel,
		roomService:  roomService,
		exchangeName: "reservations",
		queueName:    "rooms.reservation.events",
	}
}

// Start comienza a escuchar los eventos de reserva
func (rc *ReservationConsumer) Start(ctx interface{}) error {
	// Declarar el exchange
	err := rc.channel.ExchangeDeclare(
		rc.exchangeName, // nombre del exchange
		"topic",         // tipo
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // argumentos
	)
	if err != nil {
		return err
	}

	// Declarar la cola
	queue, err := rc.channel.QueueDeclare(
		rc.queueName, // nombre de la cola
		true,         // durable
		false,        // auto-delete
		false,        // exclusive
		false,        // no-wait
		nil,          // argumentos
	)
	if err != nil {
		return err
	}

	// Bindear la cola al exchange
	err = rc.channel.QueueBind(
		queue.Name,      // nombre de la cola
		"reservation.*", // routing key (escuchar todos los eventos de reserva)
		rc.exchangeName, // nombre del exchange
		false,           // no-wait
		nil,             // argumentos
	)
	if err != nil {
		return err
	}

	// Consumir mensajes
	messages, err := rc.channel.Consume(
		queue.Name, // cola
		"",         // consumer tag
		false,      // auto-ack (procesaremos manualmente)
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // argumentos
	)
	if err != nil {
		return err
	}

	// Escuchar en goroutine
	go func() {
		for message := range messages {
			rc.handleMessage(message)
		}
	}()

	log.Println("✅ Reservation consumer iniciado")
	return nil
}

// handleMessage procesa un mensaje de evento de reserva
func (rc *ReservationConsumer) handleMessage(msg amqp.Delivery) {
	var event ReservationEvent

	// Decodificar el mensaje
	err := json.Unmarshal(msg.Body, &event)
	if err != nil {
		log.Printf("❌ Error decodificando evento de reserva: %v", err)
		msg.Nack(false, false) // Rechazar el mensaje sin reencolar
		return
	}

	log.Printf("📨 Evento de reserva recibido: tipo=%s, roomID=%d, status=%s", event.EventType, event.RoomID, event.Status)

	// Procesar según el tipo de evento
	switch event.EventType {
	case "reservation.created":
		rc.handleReservationCreated(event)
	case "reservation.canceled":
		rc.handleReservationCanceled(event)
	default:
		log.Printf("⚠️  Tipo de evento desconocido: %s", event.EventType)
	}

	// Confirmar que el mensaje fue procesado
	msg.Ack(false)
}

// handleReservationCreated cambia el estado de la habitación a ocupada
func (rc *ReservationConsumer) handleReservationCreated(event ReservationEvent) {
	log.Printf("🔄 Procesando creación de reserva para room %d", event.RoomID)

	// Cambiar estado a ocupada
	status := domain.RoomStatusOccupied
	updateReq := domain.UpdateRoomRequest{
		Status: &status,
	}

	// Usar contexto vacío ya que esto es un evento asincrónico
	ctx := context.Background()
	_, err := rc.roomService.UpdateRoom(ctx, event.RoomID, updateReq)
	if err != nil {
		log.Printf("❌ Error actualizando estado de habitación %d a ocupada: %v", event.RoomID, err)
		return
	}

	log.Printf("✅ Habitación %d marcada como ocupada", event.RoomID)
}

// handleReservationCanceled cambia el estado de la habitación a disponible
func (rc *ReservationConsumer) handleReservationCanceled(event ReservationEvent) {
	log.Printf("🔄 Procesando cancelación de reserva para room %d", event.RoomID)

	// Cambiar estado a disponible
	status := domain.RoomStatusAvailable
	updateReq := domain.UpdateRoomRequest{
		Status: &status,
	}

	// Usar contexto vacío ya que esto es un evento asincrónico
	ctx := context.Background()
	_, err := rc.roomService.UpdateRoom(ctx, event.RoomID, updateReq)
	if err != nil {
		log.Printf("❌ Error actualizando estado de habitación %d a disponible: %v", event.RoomID, err)
		return
	}

	log.Printf("✅ Habitación %d marcada como disponible", event.RoomID)
}
