package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clientInstance *mongo.Client
var clientInstanceError error
var mongoOnce sync.Once

//mongo database ile iligli bilgilerimizi giriyoruz
const (
	DB_URL     = "mongodb://localhost:27017"
	DB_NAME    = "internship"
	COLLECTION = "users"
)

//front-end'den alacagimiz json verisinin model yapısı
type Kayit struct {
	NAME    string `json:"name"`
	EMAIL   string `json:"email"`
	MESSAGE string `json:"message"`
}

//mongo database baglantimizi yapiyoruz.
func getMongoClient() (*mongo.Client, error) {
	mongoOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(DB_URL)

		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceError = err
			fmt.Println("Connected to MongoDB succesfully!")
		}

		err = client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = err
		}
		clientInstance = client

	})

	return clientInstance, clientInstanceError
}

//Kayit adlı model yapımızı arguman olarak gecip database'e kayit yapiyoruz
func sendToDB(record Kayit) {
	client, err := getMongoClient()
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database(DB_NAME).Collection(COLLECTION)

	_, err = collection.InsertOne(context.TODO(), record)
	if err != nil {
		log.Fatal(err)
	}
}

//front-end'den veriyi cektigimiz kısım.
func postHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	var kayit Kayit

	err2 := c.ShouldBindJSON(&kayit)
	if err2 != nil {
		log.Fatal(err2)
	}

	sendToDB(kayit)
}

func mainPage(c *gin.Context) {
	fmt.Fprintf(c.Writer, "hello!")
	fmt.Printf("fsa\n")
}

func main() {

	router := gin.Default()
	/*
			corsOpts := cors.New(cors.Options{
		        //allow http://127.0.0.1:5500 origin
		        AllowedOrigins: []string{"http://127.0.0.1:5500"},
		        AllowedMethods: []string{"GET", "POST","DELETE","PUT"},
		        AllowedHeaders: []string{"Content-Type", "Access-Control-Allow-Origin"},
		    })
	*/
	router.Use(cors.Default())

	router.POST("/postRequest", postHandler)
	router.GET("/", mainPage)

	fmt.Println("Listenning serving on :8080!")
	router.Run(":8080")
}
