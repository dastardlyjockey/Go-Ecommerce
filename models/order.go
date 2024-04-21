package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Order struct {
	OrderID       primitive.ObjectID `bson:"order_id"`
	OrderCart     []ProductUser      `json:"order_cart"`
	OrderedAt     time.Time          `json:"ordered_at"`
	Price         *float64           `json:"price"`
	Discount      *float64           `json:"discount"`
	PaymentMethod Payment            `json:"payment_method"`
}
