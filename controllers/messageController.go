package controllers

import (
    "context"
    "net/http"
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
    "github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
    "go-bookstore/database"
    "go-bookstore/models"
)

// AddMessage handles adding a message to the database
func AddMessage(c *gin.Context) {
    var message models.Message
    if err := c.ShouldBindJSON(&message); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    message.ID = primitive.NewObjectID()
    message.Timestamp = time.Now()

    collection := database.Client.Database("AfricArt").Collection("messages")
    _, err := collection.InsertOne(context.Background(), message)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Message saved successfully"})
}

func GetMessage(c *gin.Context) {
    collection := database.Client.Database("AfricArt").Collection("messages")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var users []bson.M
    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
        return
    }
    defer cursor.Close(ctx)

    if err = cursor.All(ctx, &users); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode messages"})
        return
    }

    c.JSON(http.StatusOK, users)
}
