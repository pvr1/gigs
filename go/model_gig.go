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

// Gig - A gig struct
type Gig struct {
	Id string `bson:"id"`

	Category Category `bson:"category,omitempty"`

	Name string `bson:"name"`

	Description []string `bson:"description"`

	Measurableoutcome []string `bson:"measurableoutcome"`

	Tags []Tag `bson:"tags,omitempty"`

	// gig status in the store
	Status string `bson:"status,omitempty"`

	UserId string `bson:"userid"`
}

type Gigfile struct {
	Id       string `bson:"id"`
	Filename string `bson:"filename"`
}

func init() {
	resyncGigs()
}

func resyncGigs() {
	credential := options.Credential{
		Username: "gigbe",
		Password: "gigbe",
	}
	clientOpts := options.Client().ApplyURI("mongodb://mongodb.mongodb.svc.cluster.local:27017").
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
	gigsCollection := quickstartDatabase.Collection("gigs")

	cursor, err := gigsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(ctx, &gigs); err != nil {
		log.Fatal(err)
	}

	gigsCollection = quickstartDatabase.Collection("gigsfiles")

	cursor, err = gigsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(ctx, &gigsfiles); err != nil {
		log.Fatal(err)
	}
}

var gigs = []Gig{}
var gigsfiles = []Gigfile{}
