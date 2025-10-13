package repositories

import (
	"context"
	"rooms-api/domain"
	"rooms-api/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RoomRepository struct {
	collection *mongo.Collection
}

func NewRoomRepository(db *mongo.Database) *RoomRepository {
	return &RoomRepository{
		collection: db.Collection("rooms"),
	}
}

func (r *RoomRepository) Create(ctx context.Context, room *domain.Room) error {
	room.CreatedAt = time.Now()
	room.UpdatedAt = time.Now()

	filter := bson.M{"number": room.Number}
	var existingRoom domain.Room
	err := r.collection.FindOne(ctx, filter).Decode(&existingRoom)
	if err == nil {
		return utils.ErrRoomAlreadyExists
	}
	if err != mongo.ErrNoDocuments {
		return utils.ErrDatabaseError
	}

	_, err = r.collection.InsertOne(ctx, room)
	if err != nil {
		return utils.ErrDatabaseError
	}
	return nil
}

func (r *RoomRepository) GetByID(ctx context.Context, id string) (*domain.Room, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, utils.ErrInvalidID
	}

	filter := bson.M{"_id": objectID}
	var room domain.Room
	err = r.collection.FindOne(ctx, filter).Decode(&room)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.ErrRoomNotFound
		}
		return nil, utils.ErrDatabaseError
	}
	return &room, nil
}

func (r *RoomRepository) GetByNumber(ctx context.Context, number string) (*domain.Room, error) {
	filter := bson.M{"number": number}
	var room domain.Room
	err := r.collection.FindOne(ctx, filter).Decode(&room)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.ErrRoomNotFound
		}
		return nil, utils.ErrDatabaseError
	}
	return &room, nil
}

func (r *RoomRepository) GetAll(ctx context.Context, filter domain.RoomFilter, page, limit int) ([]domain.Room, int64, error) {

	mongoFilter := bson.M{}

	if filter.Type != nil {
		mongoFilter["type"] = *filter.Type
	}
	if filter.Status != nil {
		mongoFilter["status"] = *filter.Status
	}
	if filter.Floor != nil {
		mongoFilter["floor"] = *filter.Floor
	}
	if filter.MinPrice != nil || filter.MaxPrice != nil {
		priceFilter := bson.M{}
		if filter.MinPrice != nil {
			priceFilter["$gte"] = *filter.MinPrice
		}
		if filter.MaxPrice != nil {
			priceFilter["$lte"] = *filter.MaxPrice
		}
		mongoFilter["price"] = priceFilter
	}
	if filter.HasWifi != nil {
		mongoFilter["has_wifi"] = *filter.HasWifi
	}
	if filter.HasAC != nil {
		mongoFilter["has_ac"] = *filter.HasAC
	}
	if filter.HasTV != nil {
		mongoFilter["has_tv"] = *filter.HasTV
	}
	if filter.HasMinibar != nil {
		mongoFilter["has_minibar"] = *filter.HasMinibar
	}

	total, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, utils.ErrDatabaseError
	}

	opts := options.Find()
	opts.SetSkip(int64((page - 1) * limit))
	opts.SetLimit(int64(limit))
	opts.SetSort(bson.D{{"number", 1}})

	cursor, err := r.collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, 0, utils.ErrDatabaseError
	}
	defer cursor.Close(ctx)

	var rooms []domain.Room
	if err = cursor.All(ctx, &rooms); err != nil {
		return nil, 0, utils.ErrDatabaseError
	}

	return rooms, total, nil
}

func (r *RoomRepository) Update(ctx context.Context, id string, updateData domain.UpdateRoomRequest) (*domain.Room, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, utils.ErrInvalidID
	}

	update := bson.M{"$set": bson.M{"updated_at": time.Now()}}

	if updateData.Number != nil {

		if *updateData.Number != "" {
			filter := bson.M{"number": *updateData.Number, "_id": bson.M{"$ne": objectID}}
			var existingRoom domain.Room
			err := r.collection.FindOne(ctx, filter).Decode(&existingRoom)
			if err == nil {
				return nil, utils.ErrRoomAlreadyExists
			}
			if err != mongo.ErrNoDocuments {
				return nil, utils.ErrDatabaseError
			}
		}
		update["$set"].(bson.M)["number"] = *updateData.Number
	}
	if updateData.Type != nil {
		update["$set"].(bson.M)["type"] = *updateData.Type
	}
	if updateData.Status != nil {
		update["$set"].(bson.M)["status"] = *updateData.Status
	}
	if updateData.Price != nil {
		update["$set"].(bson.M)["price"] = *updateData.Price
	}
	if updateData.Description != nil {
		update["$set"].(bson.M)["description"] = *updateData.Description
	}
	if updateData.Capacity != nil {
		update["$set"].(bson.M)["capacity"] = *updateData.Capacity
	}
	if updateData.Floor != nil {
		update["$set"].(bson.M)["floor"] = *updateData.Floor
	}
	if updateData.HasWifi != nil {
		update["$set"].(bson.M)["has_wifi"] = *updateData.HasWifi
	}
	if updateData.HasAC != nil {
		update["$set"].(bson.M)["has_ac"] = *updateData.HasAC
	}
	if updateData.HasTV != nil {
		update["$set"].(bson.M)["has_tv"] = *updateData.HasTV
	}
	if updateData.HasMinibar != nil {
		update["$set"].(bson.M)["has_minibar"] = *updateData.HasMinibar
	}

	filter := bson.M{"_id": objectID}
	result := r.collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))

	var room domain.Room
	if err := result.Decode(&room); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.ErrRoomNotFound
		}
		return nil, utils.ErrDatabaseError
	}

	return &room, nil
}

func (r *RoomRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.ErrInvalidID
	}

	filter := bson.M{"_id": objectID}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return utils.ErrDatabaseError
	}

	if result.DeletedCount == 0 {
		return utils.ErrRoomNotFound
	}

	return nil
}
