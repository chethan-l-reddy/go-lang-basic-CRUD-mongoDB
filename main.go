package main

import (
	"go-gin-mongodb-crud/database"
	"go-gin-mongodb-crud/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	database.Connect()

	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "hello"})
	})

	router.POST("/newUser", handlers.CreateUser)
	router.GET("/getUser/:id", handlers.GetUser)
	router.GET("/getUsers", handlers.GetUsers)
	router.DELETE("/deleteUser/:id", handlers.DeleteUser)
	router.PUT("/updateUser/:id", handlers.UpdateUser)

	router.Run("localhost:9000")
}
