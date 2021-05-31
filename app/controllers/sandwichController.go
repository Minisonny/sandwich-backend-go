package controllers

import (
	"context"
	"log"
	"strconv"
	"time"

	"app/app/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Topping struct {
	Id   int64  `json:"id" validate:"required"`
	Name string `json:"string" validate:"required"`
}

type Sandwich struct {
	Id        int64     `json:"id" validate:"required,unique"`
	Name      string    `json:"name" validate:"required"`
	Toppings  []Topping `json:"toppings"`
	BreadType string    `json:"breadType"`
	ImageURL  string    `json:"imageURL"`
}

var sandwichCollection *mongo.Collection = database.GetCollection("sandwiches")

func GetAllSandwich(c *gin.Context) {
	var rctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	cursor, err := sandwichCollection.Find(rctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var episodes []Sandwich
	if err = cursor.All(rctx, &episodes); err != nil {
		log.Fatal(err)
	}

	defer cancel()
	c.JSON(200, episodes)
}

func GetSandwichById(c *gin.Context) {
	var rctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	sandwichId, _ := strconv.Atoi(c.Param("sandwichId"))
	sandwich := Sandwich{}
	if err := sandwichCollection.FindOne(rctx, bson.M{"id": sandwichId}).Decode(&sandwich); err != nil {
		log.Fatal(err)
	}

	defer cancel()
	c.JSON(200, sandwich)
}

func CreateSandwich(c *gin.Context) {
	var rctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	apiKey := c.GetHeader("x-api-key")

	if !verifyApiKey(apiKey) {
		c.JSON(401, gin.H{"error": "Invalid API key"})
		return
	}

	sandwich := Sandwich{}
	if err := c.ShouldBindJSON(&sandwich); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	sandwichCollection.InsertOne(rctx, sandwich)

	c.JSON(200, sandwich)
}

func UpdateSandwich(c *gin.Context) {
	var rctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	apiKey := c.GetHeader("x-api-key")

	if !verifyApiKey(apiKey) {
		c.JSON(401, gin.H{"error": "Invalid API key"})
		return
	}

	sandwich := Sandwich{}
	if err := c.ShouldBindJSON(&sandwich); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	sandwichId, _ := strconv.Atoi(c.Param("sandwichId"))
	result := sandwichCollection.FindOneAndUpdate(rctx, bson.M{"id": sandwichId}, sandwich)

	c.JSON(200, result)
}
