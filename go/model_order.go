package openapi

import (
	"time"
)

// Order - A order struct
type Order struct {
	Id string `json:"id,omitempty"`

	GigId string `json:"gigId,omitempty"`

	Price int32 `json:"price,omitempty"`

	ShipDate time.Time `json:"shipDate,omitempty"`

	// Order Status
	Status string `json:"status,omitempty"`

	Complete bool `json:"complete,omitempty"`
}
