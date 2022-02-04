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

// ApiResponse - A generic API response
type ApiResponse struct {
	Code int32 `bson:"code,omitempty"`

	Type string `bson:"type,omitempty"`

	Message string `bson:"message,omitempty"`
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
	gigsCollection := quickstartDatabase.Collection("apiresponse")

	cursor, err := gigsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(ctx, &apiresponse); err != nil {
		log.Fatal(err)
	}
	fmt.Println(apiresponse)
}

var apiresponse = []ApiResponse{}
