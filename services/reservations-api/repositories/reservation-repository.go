package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"reservations-api/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReservationRepository interface {
	Create(ctx context.Context, reservation domain.Reservation) (domain.Reservation, error)
	Delete(ctx context.Context, id string, reason string) error
	GetByID(ctx context.Context, id string) (domain.Reservation, error)
}

type reservationRepository struct {
	collection *mongo.Collection
}

func NewReservationRepository(db *mongo.Database) ReservationRepository {
	return &reservationRepository{
		collection: db.Collection("reservations"),
	}
}

// Create persiste una Reservation en MongoDB y devuelve la entidad guardada.
func (r *reservationRepository) Create(ctx context.Context, reservation domain.Reservation) (domain.Reservation, error) {
	// Estado por defecto si no viene seteado
	if reservation.Status == "" {
		reservation.Status = domain.ReservationStatusActive
	}

	// Generar nuevo ID si no existe
	if reservation.ID.IsZero() {
		reservation.ID = primitive.NewObjectID()
	}

	// Setear timestamps
	now := time.Now()
	reservation.CreatedAt = now
	reservation.UpdatedAt = now

	// Insertar en MongoDB
	_, err := r.collection.InsertOne(ctx, reservation)
	if err != nil {
		return domain.Reservation{}, fmt.Errorf("failed to create reservation: %w", err)
	}

	return reservation, nil
}

// Delete marca una reserva como cancelada (soft delete)
func (r *reservationRepository) Delete(ctx context.Context, id string, reason string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	now := time.Now()

	// Solo cancela si esta activa
	filter := bson.M{
		"_id":    objID,
		"status": domain.ReservationStatusActive,
	}

	update := bson.M{
		"$set": bson.M{
			"status":        domain.ReservationStatusCanceled,
			"cancel_reason": reason,
			"deleted_at":    now,
			"updated_at":    now,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("cancelar reserva: %w", err)
	}

	if result.MatchedCount == 0 {
		// Verificar si existe el documento
		count, err := r.collection.CountDocuments(ctx, bson.M{"_id": objID})
		if err != nil {
			return err
		}
		if count == 0 {
			return errors.New("reservation not found")
		}
		// Ya estaba cancelada
		return nil
	}

	return nil
}

// GetByID obtiene una reserva por su ID
func (r *reservationRepository) GetByID(ctx context.Context, id string) (domain.Reservation, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Reservation{}, fmt.Errorf("invalid id format: %w", err)
	}

	var reservation domain.Reservation
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&reservation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Reservation{}, errors.New("reservation not found")
		}
		return domain.Reservation{}, err
	}

	return reservation, nil
}
