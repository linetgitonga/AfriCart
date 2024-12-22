package controllers

import (
    "context"
    "net/http"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/gin-gonic/gin"
    "go-bookstore/database"
	
)

// CategoryCount represents the category and its item count
type CategoryCount struct {
    Category string `json:"category"`
    Count    int    `json:"count"`
}

// GetCategoryCounts fetches item counts for each category
func GetCategoryCounts(c *gin.Context) {
    collection := database.Client.Database("AfricArt").Collection("items") // Adjust database and collection names
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Aggregation pipeline to group items by category and count them
    pipeline := mongo.Pipeline{
        {
            {"$group", bson.D{
                {"_id", "$category"}, // Group by the "category" field
                {"count", bson.D{{"$sum", 1}}},
            }},
        },
        {
            {"$project", bson.D{
                {"category", "$_id"}, // Rename _id to category
                {"count", 1},
                {"_id", 0}, // Exclude the _id field
            }},
        },
    }

    cursor, err := collection.Aggregate(ctx, pipeline)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch category counts"})
        return
    }
    defer cursor.Close(ctx)

    var results []CategoryCount
    if err = cursor.All(ctx, &results); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse results"})
        return
    }

    c.JSON(http.StatusOK, results)
}
