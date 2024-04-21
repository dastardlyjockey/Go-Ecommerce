package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ProductID   primitive.ObjectID `bson:"product_id"`
	ProductName *string            `json:"product_name"`
	Price       *float64           `json:"price"`
	Rating      int16              `json:"rating"`
	Image       *string            `json:"image"`
}
