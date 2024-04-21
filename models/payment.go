package models

type Payment struct {
	Digital        bool `json:"digital"`
	CashOnDelivery bool `json:"cashOnDelivery"`
}
