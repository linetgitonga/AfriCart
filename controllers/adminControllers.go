package controllers

// import (
//     "context"
//     "net/http"
//     "time"

//     "github.com/gin-gonic/gin"
//     // "go.mongodb.org/mongo-driver/bson"
// 	// "go.mongodb.org/mongo-driver/mongo"
//     // "go.mongodb.org/mongo-driver/mongo/options"
//     "go-bookstore/database"
//     "go-bookstore/models"
// 	// "go.mongodb.org/mongo-driver/bson/primitive"

// )
// func AddAdminProduct(c *gin.Context) {
//     var product models.Product
//     if err := c.BindJSON(&product); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
//         return
//     }

//     // Assign a custom ID if needed
//     if product.CustomID == "" {
//         product.CustomID = generateCustomID() // Your custom ID generation logic
//     }

//     collection := database.Client.Database("AfricArt").Collection("items")
//     ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//     defer cancel()

//     _, err := collection.InsertOne(ctx, product)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product"})
//         return
//     }

//     c.JSON(http.StatusCreated, gin.H{"message": "Product added successfully"})
// }

// // Route setup

