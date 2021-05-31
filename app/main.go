package main

import (
	"app/app/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowWebSockets:  true,
		AllowHeaders:     []string{"Origin", "Content-Type", "X-API-KEY"},
	}))

	r.GET("/v1/order", controllers.GetAllOrder)
	r.GET("/v1/order/:orderId", controllers.GetOrderById)
	r.POST("/v1/order", controllers.CreateOrder)

	r.GET("/v1/sandwich", controllers.GetAllSandwich)
	r.GET("/v1/sandwich/:sandwichId", controllers.GetSandwichById)
	r.POST("/v1/sandwich", controllers.CreateSandwich)

	r.POST("/v1/user", controllers.CreateUser)
	r.POST("/v1/user/login", controllers.Login)

	r.Run(":3001")
}
