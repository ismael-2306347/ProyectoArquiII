package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomType string

const (
	RoomTypeSingle   RoomType = "single"
	RoomTypeDouble   RoomType = "double"
	RoomTypeSuite    RoomType = "suite"
	RoomTypeDeluxe   RoomType = "deluxe"
	RoomTypeStandard RoomType = "standard"
)

type RoomStatus string

const (
	RoomStatusAvailable   RoomStatus = "available"
	RoomStatusOccupied    RoomStatus = "occupied"
	RoomStatusMaintenance RoomStatus = "maintenance"
	RoomStatusReserved    RoomStatus = "reserved"
)

type Room struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Number      string             `bson:"number" json:"number"`
	Type        RoomType           `bson:"type" json:"type"`
	Status      RoomStatus         `bson:"status" json:"status"`
	Price       float64            `bson:"price" json:"price"`
	Description string             `bson:"description" json:"description"`
	Capacity    int                `bson:"capacity" json:"capacity"`
	Floor       int                `bson:"floor" json:"floor"`
	HasWifi     bool               `bson:"has_wifi" json:"has_wifi"`
	HasAC       bool               `bson:"has_ac" json:"has_ac"`
	HasTV       bool               `bson:"has_tv" json:"has_tv"`
	HasMinibar  bool               `bson:"has_minibar" json:"has_minibar"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
