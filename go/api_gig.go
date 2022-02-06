package openapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/twinj/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// connectMongoDB - connect to MongoDB
func connectMongoDB() (*mongo.Client, context.Context) {
	credential := options.Credential{
		Username: "gigbe",
		Password: "gigbe",
	}
	clientOpts := options.Client().ApplyURI("mongodb://mymongodb.mongodb.svc.cluster.local:27017").
		SetAuth(credential)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		println("Connect error")
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return client, ctx
}

func InitiateMongoClient() *mongo.Client {
	var err error
	var client *mongo.Client
	credential := options.Credential{
		Username: "gigbe",
		Password: "gigbe",
	}
	uri := "mongodb://mymongodb.mongodb.svc.cluster.local:27017"
	opts := options.Client()
	opts.ApplyURI(uri).SetAuth(credential)
	//TODO: adjust maxpoolsize
	opts.SetMaxPoolSize(5)
	if client, err = mongo.Connect(context.Background(), opts); err != nil {
		fmt.Println(err.Error())
	}
	return client
}

// getCollectionMongoDB - Get a collection from MongoDB
func getCollectionMongoDB(client *mongo.Client, collection string) *mongo.Collection {
	quickstartDatabase := client.Database("gigs")
	gigsCollection := quickstartDatabase.Collection(collection)
	return gigsCollection
}

// AddGig - Add a new gig to the store
func AddGig(c *gin.Context) {
	var mygig Gig

	// Call BindJSON to bind the received JSON to
	// newgig.
	mygig.Id = uuid.NewV4().String()
	if err := c.BindJSON(&mygig); err != nil {
		return
	}

	//Validate the gig
	if mygig.Name == "" || mygig.Description == nil || mygig.Status == "" || mygig.Measurableoutcome == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required gig fields"})
		return
	}

	//TODO: Validate the user so it is a logged in user
	if mygig.UserId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required user id - please login"})
		return
	}
	mygig.UserId = myUser.Id

	// Add the new gig to the slice.
	gigs = append(gigs, mygig)

	//Insert record into mongodb
	client, ctx := connectMongoDB()
	defer client.Disconnect(ctx)
	gigsCollection := getCollectionMongoDB(client, "gigs")
	gigsE, err := gigsCollection.InsertMany(ctx, []interface{}{
		&mygig,
	})
	if err != nil {
		println("add record error")
		log.Fatal(gigsE, err)
	}

	c.JSON(http.StatusCreated, mygig)
}

// RemoveGig - Helper function to remove a gig from the slice
func RemoveGig(s []Gig, index int) []Gig {
	return append(s[:index], s[index+1:]...)
}

// DeleteGig - Deletes a gig
func DeleteGig(c *gin.Context) {
	id := c.Param("gigId")

	/*
		file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}

		log.SetOutput(file)
	*/
	//First try at blackfriday and bluemonday
	//log.Println("id: ", id)
	unsafe := blackfriday.SanitizedAnchorName(id)
	//log.Println("unsafe: ", unsafe)
	html := string(bluemonday.UGCPolicy().SanitizeBytes([]byte(unsafe)))
	//log.Println("html: ", html)

	// Loop over the list of gigs, looking for
	// an gig whose ID value matches the parameter.
	for i, a := range gigs {
		if a.Id == html {
			gigs = RemoveGig(gigs, i)

			//Delete record in mongodb
			client, ctx := connectMongoDB()
			defer client.Disconnect(ctx)
			gigsCollection := getCollectionMongoDB(client, "gigs")
			gigsE, err := gigsCollection.DeleteOne(ctx,
				bson.M{"id": html},
			)
			if err != nil {
				println("add record error")
				log.Fatal(gigsE, err)
			}

			c.JSON(http.StatusOK, a)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Gig not found"})
}

// FindGigsByStatus - Finds Gigs by status
func FindGigsByStatus(c *gin.Context) {
	status, err := c.GetQuery("status")
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	unsafe := blackfriday.SanitizedAnchorName(status)
	html := string(bluemonday.UGCPolicy().SanitizeBytes([]byte(unsafe)))

	// Loop over the list of gigs, looking for
	// an gig whose ID value matches the parameter.

	tmp := []Gig{}
	tmp2 := Gig{}

	for _, a := range gigs {
		if a.Status == html {
			tmp2 = a
			tmp2.UserId = "obfuscated"
			tmp = append(tmp, tmp2)
		}
	}

	c.JSON(http.StatusOK, tmp)
}

// FindGigsByStatus - Finds Gigs by status
func FindGigsByUser(c *gin.Context) {
	status, err := c.GetQuery("userid")
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	unsafe := blackfriday.SanitizedAnchorName(status)
	html := string(bluemonday.UGCPolicy().SanitizeBytes([]byte(unsafe)))

	// Loop over the list of gigs, looking for
	// an gig whose ID value matches the parameter.

	tmp := []Gig{}
	fmt.Println("userid: ", html)

	for _, a := range gigs {
		if a.UserId == html {
			tmp = append(tmp, a)
		}
	}

	c.JSON(http.StatusOK, tmp)
}

// GetGigById - Find gig by ID
func GetGigById(c *gin.Context) {
	id := c.Param("gigId")

	unsafe := blackfriday.SanitizedAnchorName(id)
	html := string(bluemonday.UGCPolicy().SanitizeBytes([]byte(unsafe)))

	// Loop over the list of gigs, looking for
	// an gig whose ID value matches the parameter.
	tmp2 := Gig{}

	for _, a := range gigs {
		if a.Id == html {
			tmp2 = a
			tmp2.UserId = myUser.Id
			c.JSON(http.StatusOK, tmp2)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Gig not found"})
}

// UpdateGigWithForm - Updates a gig in the store with form data
func UpdateGigWithForm(c *gin.Context) {
	// If the file doesn't exist, create it or append to the file
	/*
		file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}

		log.SetOutput(file)
	*/
	id := c.Param("gigId")

	unsafe := blackfriday.SanitizedAnchorName(id)
	html := string(bluemonday.UGCPolicy().SanitizeBytes([]byte(unsafe)))

	bodyjson, err := c.GetRawData()
	var body Gig
	json.Unmarshal(bodyjson, &body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if body.Id != id {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Id mismatch"})
		return
	}

	// Loop over the list of gigs, looking for
	// an gig whose ID value matches the parameter.
	for i, a := range gigs {
		if a.Id == html {
			// Update the gig
			DeepCopy(body, &gigs[i])
			c.JSON(http.StatusOK, body)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Gig not found"})
}

// UploadFile - uploads an image
func UploadFile(c *gin.Context) {
	// Single file
	id := c.Param("gigId")
	file, err := c.FormFile("file")
	//log.Println("file: ", file)
	if err != nil {
		log.Println("Error uploading file - ", err)
		c.String(http.StatusBadRequest, fmt.Sprintf("Error uploading: %s", err.Error()))
		return
	}

	// Retrieve file information
	extension := filepath.Ext(file.Filename)
	// Generate random file name for the new uploaded file so it doesn't override the old file with same name
	newFileName := uuid.NewV4().String() + extension

	//Insert record into mongodb
	client, ctx := connectMongoDB()
	defer client.Disconnect(ctx)
	gigsCollection := getCollectionMongoDB(client, "gigsfiles")
	gigsE, err := gigsCollection.InsertMany(ctx, []interface{}{
		bson.D{
			{"id", &id},
			{"filename", &newFileName},
		},
	},
	)
	if err != nil {
		println("add record error")
		fmt.Println(gigsE, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Err: " + err.Error()})
		return
	}

	// The file is received, so let's save it
	conn := InitiateMongoClient()
	bucket, errB := gridfs.NewBucket(
		conn.Database("gigs"),
	)
	if errB != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Err: " + err.Error()})
		return
	}

	uploadStream, errS := bucket.OpenUploadStream(
		newFileName, // this is the name of the file which will be saved in the database
	)
	if errS != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Err: " + err.Error()})
		return
	}
	defer uploadStream.Close()

	fileC, errC := file.Open()
	if errC != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Err: " + err.Error()})
		return
	}
	bytes, errR := ioutil.ReadAll(fileC)
	if errR != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Err: " + err.Error()})
		return
	}

	fileSize, err := uploadStream.Write(bytes)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Printf("Write file to DB was successful. File size: %d \n", fileSize)
	if err != nil {
		println("add record error")
		fmt.Println(gigsE, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Err: " + err.Error()})
		return
	}
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

func DownloadFile(c *gin.Context) {
	id := c.Param("gigId")

	unsafe := blackfriday.SanitizedAnchorName(id)
	html := string(bluemonday.UGCPolicy().SanitizeBytes([]byte(unsafe)))

	// Loop over the list of gigs, looking for
	// an gig whose ID value matches the parameter.
	var resultFile = ""
	for _, a := range gigsfiles {
		fmt.Println("File: ", a.Filename)
		if a.Id == html {
			fmt.Println("Found: ", a.Id, " ", a.Filename)
			resultFile = a.Filename
			break
		}
	}

	if resultFile == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": "GigFile not found"})
		return
	}

	conn := InitiateMongoClient()

	// For CRUD operations, here is an example
	db := conn.Database("gigs")
	fsFiles := db.Collection("fs.files")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var results bson.M
	err := fsFiles.FindOne(ctx, bson.M{}).Decode(&results)
	if err != nil {
		log.Fatal(err)
	}
	// you can print out the result
	fmt.Println(results)

	bucket, _ := gridfs.NewBucket(
		db,
	)
	var buf bytes.Buffer
	dStream, err := bucket.DownloadToStreamByName(resultFile, &buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("File size to download: %v \n", dStream)
	ioutil.WriteFile(resultFile, buf.Bytes(), 0600)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+resultFile)
	c.Header("Content-Type", "application/octet-stream")
	c.File(resultFile)
}
