package controllers

import (
	"context"
	"time"

	"app/app/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           int64  `json:"id" validate:"required,unique"`
	Username     string `json:"username" validate:"required"`
	Email        string `json:"email" validate:"required"`
	PasswordHash string `json:"passwordHash" validate:"required"`
	ApiKey       string `json:"apiKey" validate:"required"`
}

type UserRequest struct {
	Id       int64
	Username string
	Password string
	Email    string
}

type UserResponse struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

var userCollection *mongo.Collection = database.GetCollection("users")

func generateApiKey() string {
	return uuid.NewString()
}

func verifyApiKey(apiKey string) bool {
	var rctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if apiKey == "" {
		return false
	}

	user := UserResponse{}
	if err := userCollection.FindOne(rctx, bson.M{"apiKey": apiKey}).Decode(&user); err != nil {
		return false
	}

	return true
}

func CreateUser(c *gin.Context) {
	var rctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if !verifyApiKey(c.GetHeader("x-api-key")) {
		c.JSON(401, gin.H{"error": "Invalid API key"})
		return
	}

	userReq := UserRequest{}
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(userReq.Password), 12)
	passwordHash := string(hashed)

	saveUser := User{
		Id:           userReq.Id,
		Username:     userReq.Username,
		Email:        userReq.Email,
		PasswordHash: passwordHash,
	}

	userCollection.InsertOne(rctx, saveUser)

	c.JSON(200, gin.H{
		"id":       saveUser.Id,
		"username": saveUser.Username,
		"email":    saveUser.Email,
	})
}

func Login(c *gin.Context) {
	var rctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	userReq := UserRequest{}
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	if userReq.Username == "" || userReq.Password == "" {
		c.JSON(400, gin.H{"error": "Invalid username or password"})
		return
	}

	user := User{}
	if err := userCollection.FindOne(rctx, bson.M{"username": userReq.Username}).Decode(&user); err != nil {
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(userReq.Password)); err != nil {
		c.JSON(400, gin.H{"error": "Wrong password"})
		return
	}

	if user.ApiKey != "" {
		c.JSON(200, user.ApiKey)
		return
	}

	apiKey := generateApiKey()
	userCollection.FindOneAndUpdate(rctx, bson.M{"username": userReq.Username}, bson.M{"apiKey": apiKey})

	c.JSON(200, apiKey)
}

func Logout(c *gin.Context) {
	var rctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	userReq := UserRequest{}
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	user := User{}
	userCollection.FindOneAndUpdate(rctx, bson.M{"username": userReq.Username}, bson.M{"apiKey": ""}).Decode(&user)

	c.Status(200)
}
