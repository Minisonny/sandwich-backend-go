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

type Order struct {
	Id         int64  `json:"id" validate:"required,unique"`
	SandwichId int64  `json:"sandwichId" validate:"required"`
	Status     string `json:"status" validate:"required"`
}

var orderCollection *mongo.Collection = database.GetCollection("orders")

func GetAllOrder(c *gin.Context) {
	var rctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	cursor, err := orderCollection.Find(rctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var episodes []Order
	if err = cursor.All(rctx, &episodes); err != nil {
		log.Fatal(err)
	}
	defer cancel()
	c.JSON(200, episodes)
}

func GetOrderById(c *gin.Context) {
	var rctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	orderId, _ := strconv.Atoi(c.Param("orderId"))
	order := Order{}
	if err := orderCollection.FindOne(rctx, bson.M{"id": orderId}).Decode(&order); err != nil {
		log.Fatal(err)
	}

	defer cancel()
	c.JSON(200, order)
}

func CreateOrder(c *gin.Context) {
	var rctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	order := Order{}
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	orderCollection.InsertOne(rctx, order)

	c.JSON(200, order)
}
