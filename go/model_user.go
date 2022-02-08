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

// User - A user struct used to store the user information
type User struct {
	Id string `bson:"id,omitempty"`

	Username string `bson:"username,omitempty"`

	FirstName string `bson:"firstName,omitempty"`

	LastName string `bson:"lastName,omitempty"`

	Email string `bson:"email,omitempty"`

	Password string `bson:"password,omitempty"`

	Phone string `bson:"phone,omitempty"`

	// User Status
	UserStatus int32 `bson:"userStatus,omitempty"`

	// User Role - e.g. gigworker, employer etc
	Role []Role `bson:"role,omitempty"`
}

func init() {
	resyncUsers()
}

func resyncUsers() {
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
	gigsCollection := quickstartDatabase.Collection("users")

	cursor, err := gigsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(ctx, &users); err != nil {
		log.Fatal(err)
	}

}

var users = []User{
	{
		FirstName:  "firstName",
		LastName:   "lastName",
		Password:   "password",
		UserStatus: 6,
		Phone:      "888-888-8888",
		Id:         "0",
		Email:      "a.a@a.com",
		Username:   "aaa",
	},
}
