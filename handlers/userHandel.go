package handlers

import (
	"context"
	"go-gin-mongodb-crud/database"
	"go-gin-mongodb-crud/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var collectionInstance *mongo.Collection

func CreateUser(c *gin.Context) {
	var user models.UserModel
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	collectionInstance = database.Client.Database("go-mongo-crud").Collection("user")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collectionInstance.InsertOne(ctx, user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	collectionInstance.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&user)

	c.JSON(http.StatusOK, user)
}

func GetUser(c *gin.Context) {
	var user models.UserModel
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ID"})
		return
	}
	collectionInstance = database.Client.Database("go-mongo-crud").Collection("user")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = collectionInstance.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"mongo-error": err.Error()})
	}

	c.JSON(http.StatusOK, user)
}

func GetUsers(c *gin.Context) {
	var users []models.UserModel

	collectionInstance = database.Client.Database("go-mongo-crud").Collection("user")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collectionInstance.Find(ctx, bson.M{})

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "documrnt not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"mongo-error": err.Error()})
		return
	}

	for cursor.Next(ctx) {
		var user models.UserModel
		if err := cursor.Decode(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"internal-error": err.Error()})
			return
		}
		users = append(users, user)
	}

	cursor.Close(ctx)

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"internal-error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ID"})
		return
	}
	collectionInstance = database.Client.Database("go-mongo-crud").Collection("user")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collectionInstance.DeleteOne(ctx, bson.M{"_id": objID})

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}
	if result.DeletedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "document deletion failed"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"mongo-error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"isDeleted": true})

}

func UpdateUser(c *gin.Context) {
	var user models.UserModel
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collectionInstance = database.Client.Database("go-mongo-crud").Collection("user")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collectionInstance.UpdateOne(ctx, bson.M{"_id": objID}, bson.D{{Key: "$set", Value: user}})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"mongo-error": err.Error()})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "document updation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"isUpdated": true})

}
