package openapi

import (
	"time"
)

// transaction - A transaction struct
type transaction struct {
	Id string `bson:"id,omitempty"`

	GigId string `bson:"gigId,omitempty"`

	Price int32 `bson:"price,omitempty"`

	ShipDate time.Time `bson:"shipDate,omitempty"`

	// transaction Status
	Status string `bson:"status,omitempty"`

	Complete bool `bson:"complete,omitempty"`
}

var ddate, _ = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
var transactions = []transaction{

	{Id: "1", GigId: "1", Price: 100, ShipDate: ddate, Status: "pending", Complete: false},
	{Id: "2", GigId: "2", Price: 200, ShipDate: time.Now(), Status: "pending", Complete: false},
	{Id: "3", GigId: "3", Price: 300, ShipDate: time.Now(), Status: "pending", Complete: false},
}
