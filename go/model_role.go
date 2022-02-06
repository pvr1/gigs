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

// Role - A role struct to cater for different roles, e.g. gigworker, employer etc
type Role struct {
	Id string `bson:"id,omitempty"`

	Name string `bson:"name,omitempty"`
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
	ctx, errCtxTime := context.WithTimeout(context.Background(), 10*time.Second)
	if errCtxTime != nil {
		fmt.Println(errCtxTime)
	}
	defer client.Disconnect(ctx)

	quickstartDatabase := client.Database("gigs")
	gigsCollection := quickstartDatabase.Collection("roles")

	cursor, err := gigsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(ctx, &roles); err != nil {
		log.Fatal(err)
	}
	//fmt.Println(roles)
}

var roles = []Role{}
