package openapi

import (
	"time"
)

// Order - A order struct
type Order struct {
	Id string `bson:"id,omitempty"`

	GigId string `bson:"gigId,omitempty"`

	Price int32 `bson:"price,omitempty"`

	ShipDate time.Time `bson:"shipDate,omitempty"`

	// Order Status
	Status string `bson:"status,omitempty"`

	Complete bool `bson:"complete,omitempty"`
}
