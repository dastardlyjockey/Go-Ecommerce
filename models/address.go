package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Address struct {
	AddressID primitive.ObjectID `bson:"address_id"`
	House     *string            `json:"house"`
	Street    *string            `json:"street"`
	City      *string            `json:"city"`
	PinCode   *string            `json:"pin_code"`
}
