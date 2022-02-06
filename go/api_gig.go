package openapi

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	ctx, errCtxTime := context.WithTimeout(context.Background(), 10*time.Second)
	if errCtxTime != nil {
		println("Context error")
	}
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
		c.JSON(http.StatusBadRequest, gin.H{"message": "JSON error"})
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

	//Sanitize the Tags
	for _, aTag := range gigs {
		for i, bTag := range aTag.Tags {
			mygig.Tags[i].Name = string(bluemonday.UGCPolicy().SanitizeBytes([]byte(blackfriday.SanitizedAnchorName(bTag.Name))))
		}
	}

	//Sanitize the Category
	for _, aCategory := range gigs {
		mygig.Category.Name = string(bluemonday.UGCPolicy().SanitizeBytes([]byte(blackfriday.SanitizedAnchorName(aCategory.Name))))
	}

	//Sanitize the Description
	for _, aDesc := range gigs {
		for i, bDesc := range aDesc.Description {
			mygig.Description[i] = string(bluemonday.UGCPolicy().SanitizeBytes([]byte(blackfriday.SanitizedAnchorName(bDesc))))
		}
	}

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

// DeepCopy deepcopies a to b using json marshaling
func DeepCopy(a, b interface{}) {
	byt, _ := json.Marshal(a)
	json.Unmarshal(byt, b)
}

// UpdateGig - Update an existing gig
func UpdateGig(c *gin.Context) {
	/*
		// If the file doesn't exist, create it or append to the file
		file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}

		log.SetOutput(file)
	*/
	bodyjson, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}
	var body Gig

	//unsafe := blackfriday.Run(bodyjson)
	//html := bluemonday.UGCPolicy().SanitizeBytes([]byte(unsafe))

	json.Unmarshal(bodyjson, &body)

	// Loop over the list of gigs, looking for
	// an gig whose ID value matches the parameter.
	for i, a := range gigs {
		if a.Id == body.Id {
			// Update the gig
			DeepCopy(body, &gigs[i])
			c.JSON(http.StatusOK, body)
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

func FindGigsByTagsAndStatus(c *gin.Context) {
	status, err := c.GetQuery("status")
	tags, err2 := c.GetQueryArray("tags") //TODO: Check if tags is an array
	if !err || !err2 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if len(tags) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	unsafe := blackfriday.SanitizedAnchorName(status)
	html := string(bluemonday.UGCPolicy().SanitizeBytes([]byte(unsafe)))

	// Loop over the list of tags
	unsafetag := ""
	htmltag := ""
	safetags := []string{}
	for _, tag := range tags {
		unsafetag = blackfriday.SanitizedAnchorName(tag)
		htmltag = string(bluemonday.UGCPolicy().SanitizeBytes([]byte(unsafetag)))
		safetags = append(safetags, htmltag)
	}

	fmt.Println("tags: ", tags, safetags)

	// Loop over the list of gigs, looking for
	// an gig whose ID value matches the parameter.

	tmp := []Gig{}
	tmp2 := Gig{}

	for _, a := range gigs {
		if a.Status == html && containsTag(safetags, a.Tags) {
			tmp2 = a
			tmp2.UserId = "obfuscated"
			tmp = append(tmp, tmp2)
		}
	}

	c.JSON(http.StatusOK, tmp)
}

func containsTag(safetags []string, tag []Tag) bool {
	for _, t := range safetags {
		for _, t2 := range tag {
			if t == t2.Name {
				return true
			}
		}
	}
	return false
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

type gigsfilesstruct struct {
	Id          string `bson:"id"`
	Filename    string `bson:"filename"`
	Description string `bson:"description"`
}

// UploadFile - uploads an image
func UploadFile(c *gin.Context) {
	// Single file
	id := c.Param("gigId")
	desc, errQQ := c.GetPostForm("description")
	//desc, errQQ := c.GetQuery("description")
	if errQQ {
		fmt.Println("Error extracting description in POST payload: ", errQQ)
	}
	fmt.Println("desc: ", desc)

	file, errQ := c.FormFile("file")
	//log.Println("file: ", file)
	if errQ != nil {
		log.Println("Error uploading file - ", errQ)
		c.String(http.StatusBadRequest, fmt.Sprintf("Error uploading: %s", errQ.Error()))
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

	// Insert a single document
	gigsE, errQ := gigsCollection.InsertOne(ctx, gigsfilesstruct{
		Id:          id,
		Filename:    newFileName,
		Description: desc,
	})

	if errQ != nil {
		println("add record error")
		fmt.Println(gigsE, errQ)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Err: " + errQ.Error()})
		return
	}

	// The file is received, so let's save it
	conn := InitiateMongoClient()
	bucket, errB := gridfs.NewBucket(
		conn.Database("gigs"),
	)
	if errB != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Err: " + errQ.Error()})
		return
	}

	uploadStream, errS := bucket.OpenUploadStream(
		newFileName, // this is the name of the file which will be saved in the database
	)
	if errS != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Err: " + errQ.Error()})
		return
	}
	defer uploadStream.Close()

	fileC, errC := file.Open()
	if errC != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Err: " + errQ.Error()})
		return
	}
	bytes, errR := ioutil.ReadAll(fileC)
	if errR != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Err: " + errQ.Error()})
		return
	}

	fileSize, errQ := uploadStream.Write(bytes)
	if errQ != nil {
		log.Fatal(errQ)
		os.Exit(1)
	}
	fmt.Printf("Write file to DB was successful. File size: %d \n", fileSize)
	if errQ != nil {
		println("add record error")
		fmt.Println(gigsE, errQ)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Err: " + errQ.Error()})
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

	conn := InitiateMongoClient()
	db := conn.Database("gigs")
	fsFiles := db.Collection("fs.files")
	defer conn.Disconnect(context.TODO())

	archive, err := os.Create("archive.zip")
	if err != nil {
		fmt.Println("Error creating archive: ", err)
	}
	defer archive.Close()
	zipWriter := zip.NewWriter(archive)

	var resultFile = ""
	client, ctx := connectMongoDB()
	defer client.Disconnect(ctx)

	gigsCollection := getCollectionMongoDB(client, "gigsfiles")
	cursor, err := gigsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(ctx, &gigsfiles); err != nil {
		log.Fatal(err)
	}

	for _, a := range gigsfiles {
		fmt.Println("File: ", a.Filename)
		if a.Id == html {
			resultFile = a.Filename
			addFileToZip(zipWriter, resultFile, fsFiles, db)
		}
	}
	zipWriter.Close()

	if resultFile == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": "GigFile not found"})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename=archive.zip")
	c.Header("Content-Type", "application/octet-stream")
	c.File("archive.zip")
}

func addFileToZip(zipWriter *zip.Writer, resultFile string, fsFiles *mongo.Collection, db *mongo.Database) {
	ctx, errCtxTime := context.WithTimeout(context.Background(), 10*time.Second)
	if errCtxTime != nil {
		fmt.Println("Error creating context: ", errCtxTime)
	}
	var results bson.M
	err := fsFiles.FindOne(ctx, bson.M{}).Decode(&results)
	if err != nil {
		fmt.Println(err)
	}

	bucket, errNB := gridfs.NewBucket(
		db,
	)
	if errNB != nil {
		fmt.Println(errNB)
	}

	var buf bytes.Buffer
	dStream, err := bucket.DownloadToStreamByName(resultFile, &buf)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("File size to download: %v \n", dStream)
	ioutil.WriteFile(resultFile, buf.Bytes(), 0600)

	w1, errI := zipWriter.Create(resultFile)
	if errI != nil {
		fmt.Println(err)
	}
	if _, err := io.Copy(w1, &buf); err != nil {
		fmt.Println(err)
	}
}
