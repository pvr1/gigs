package openapi

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func init() {
	credential := options.Credential{
		Username: "gigbe",
		Password: "gigbe",
	}
	clientOpts := options.Client().ApplyURI("mongodb://mymongodb.mongodb.svc.cluster.local:27017").
		SetAuth(credential)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer client.Disconnect(ctx)

	quickstartDatabase := client.Database("gigs")
	gigsCollection := quickstartDatabase.Collection("transactions")

	cursor, err := gigsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(ctx, &transactions); err != nil {
		log.Fatal(err)
	}
	fmt.Println(transactions)
}

var ddate, _ = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
var transactions = []transaction{
	{Id: "1", GigId: "1", Price: 100, ShipDate: ddate, Status: "pending", Complete: false},
	{Id: "2", GigId: "2", Price: 200, ShipDate: time.Now(), Status: "pending", Complete: false},
	{Id: "3", GigId: "3", Price: 300, ShipDate: time.Now(), Status: "pending", Complete: false},
}
