package openapi

import (
	"time"
)

type Order struct {
	Id int64 `json:"id,omitempty"`

	GigId int64 `json:"gigId,omitempty"`

	Quantity int32 `json:"quantity,omitempty"`

	ShipDate time.Time `json:"shipDate,omitempty"`

	// Order Status
	Status string `json:"status,omitempty"`

	Complete bool `json:"complete,omitempty"`
}
