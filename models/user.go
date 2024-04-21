package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id"`
	FirstName      *string            `json:"first_name" validate:"required"`
	LastName       *string            `json:"last_name" validate:"required"`
	Password       *string            `json:"password" validate:"required"`
	Email          *string            `json:"email" validate:"required"`
	PhoneNumber    *string            `json:"phone_number" validate:"required"`
	Token          *string            `json:"token"`
	RefreshToken   *string            `json:"refresh_token"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	UserID         string             `json:"user_id"`
	UserCart       []Product          `json:"user_cart"`
	AddressDetails []Address          `json:"address_details"`
	OrderStatus    []Order            `json:"order_status"`
}
