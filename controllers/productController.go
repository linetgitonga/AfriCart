package controllers

import (
    "context"
    "net/http"
    "time"
    "strconv" 
    

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
    // "go.mongodb.org/mongo-driver/mongo/options"
    "go-bookstore/database"
    "go-bookstore/models"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

// GetProducts retrieves all products from the database
func GetProducts(c *gin.Context) {
    // Extract query parameters
    customID := c.Query("custom_id")
    id := c.Query("id")           // MongoDB ObjectID as a string
     // Custom ID as a string

    filter := bson.M{} // Default filter to fetch all products

    // Apply filters if `custom_id` or `id` is provided
    if customID != "" {
        filter = bson.M{"custom_id": customID}
    } else if id != "" {
        objectId, err := primitive.ObjectIDFromHex(id)
        if err == nil {
            filter = bson.M{"_id": objectId}
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ObjectID"})
            return
        }
    }

    // Connect to the products collection
    collection := database.Client.Database("AfricArt").Collection("items")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Fetch data from MongoDB
    var products []models.Product
    cursor, err := collection.Find(ctx, filter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
        return
    }
    defer cursor.Close(ctx)

    // Decode products
    for cursor.Next(ctx) {
        var product models.Product
        if err := cursor.Decode(&product); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode product"})
            return
        }
        products = append(products, product)
    }

    if err := cursor.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error"})
        return
    }

    // Handle case where no products are found
    if len(products) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "No products found"})
        return
    }

    // Return products as JSON
    c.JSON(http.StatusOK, products)
}



// Add a new product
func AddProduct(c *gin.Context) {
    var product models.Product
    if err := c.BindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    // Generate custom ID
    customID, err := database.GetNextSequenceValue("products")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate custom ID"})
        return
    }

    product.CustomID = customID

    // Insert product into the database
    collection := database.Client.Database("AfricArt").Collection("items")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err = collection.InsertOne(ctx, product)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product"})
        return
    }

    c.JSON(http.StatusCreated, product)
}

// GetProductById fetches a single product by ID
// func GetProductById(c *gin.Context) {
//     id := c.Param("id")
//     productId, err := primitive.ObjectIDFromHex(id)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
//         return
//     }

//     collection := database.Client.Database("AfricArt").Collection("items")
//     ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//     defer cancel()

//     var product models.Product
//     err = collection.FindOne(ctx, bson.M{"_id": productId}).Decode(&product)
//     if err == mongo.ErrNoDocuments {
//         c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
//         return
//     } else if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
//         return
//     }

//     c.JSON(http.StatusOK, product)
// }

// Get product by custom ID
func GetProductByCustomID(c *gin.Context) {
    // Retrieve the custom_id parameter from the URL
    customIDParam := c.Param("custom_id")

    // Convert the custom_id to an integer
    customID, err := strconv.Atoi(customIDParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid custom_id format. Must be a number."})
        return
    }

    // Create the MongoDB query
    filter := bson.M{"custom_id": customID}

    // Connect to the database
    collection := database.Client.Database("AfricArt").Collection("products")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Query the database
    var product models.Product
    err = collection.FindOne(ctx, filter).Decode(&product)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
        return
    }

    // Return the product as JSON
    c.JSON(http.StatusOK, product)
}




func GetStats(c *gin.Context) {
    client := database.Client
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    db := client.Database("AfricArt")

    // Count collections
    userCount, _ := db.Collection("users").CountDocuments(ctx, bson.M{})
    itemCount, _ := db.Collection("items").CountDocuments(ctx, bson.M{})
    categoryCount, _ := db.Collection("categories").CountDocuments(ctx, bson.M{})

    c.JSON(http.StatusOK, gin.H{
        "totalUsers":    userCount,
        "totalProducts": itemCount,
        "totalCategories": categoryCount,
    })
}



func GetItems(c *gin.Context) {
    collection := database.Client.Database("AfricArt").Collection("items")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var items []bson.M
    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
        return
    }
    defer cursor.Close(ctx)

    if err = cursor.All(ctx, &items); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode items"})
        return
    }

    c.JSON(http.StatusOK, items)
}

func GetProductByID(c *gin.Context) {
    customID := c.Param("id")
    collection := database.Client.Database("AfricArt").Collection("products")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var product models.Product
    err := collection.FindOne(ctx, bson.M{"custom_id": customID}).Decode(&product)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }

    c.JSON(http.StatusOK, product)
}



// func AddToWishlist(c *gin.Context) {
//     userID := c.Param("userID") // Assume user is authenticated
//     var request struct {
//         ProductID string `json:"product_id"`
//     }

//     if err := c.ShouldBindJSON(&request); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
//         return
//     }

//     collection := database.Client.Database("AfricArt").Collection("wishlists")
//     _, err := collection.UpdateOne(
//         context.Background(),
//         bson.M{"user_id": userID},
//         bson.M{"$addToSet": bson.M{"product_ids": request.ProductID}},
//         options.Update().SetUpsert(true),
//     )
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to wishlist"})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{"message": "Item added to wishlist"})
// }



// func RemoveFromWishlist(c *gin.Context) {
//     var request struct {
//         UserID    string `json:"user_id"`
//         ProductID string `json:"product_id"`
//     }

//     if err := c.ShouldBindJSON(&request); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
//         return
//     }

//     collection := database.Client.Database("AfricArt").Collection("wishlists")
//     _, err := collection.UpdateOne(
//         context.Background(),
//         bson.M{"user_id": request.UserID},
//         bson.M{"$pull": bson.M{"product_ids": request.ProductID}},
//     )
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove from wishlist"})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{"message": "Item removed from wishlist"})
// }

// func GetWishlist(c *gin.Context) {
//     userID := c.Query("user_id")

//     collection := database.Client.Database("AfricArt").Collection("wishlists")
//     var wishlist struct {
//         UserID    string   `bson:"user_id"`
//         ProductIDs []string `bson:"product_ids"`
//     }

//     err := collection.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&wishlist)
//     if err == mongo.ErrNoDocuments {
//         c.JSON(http.StatusOK, []bson.M{})
//         return
//     } else if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch wishlist"})
//         return
//     }

//     // Fetch product details for the wishlist
//     productCollection := database.Client.Database("AfricArt").Collection("products")
//     var products []bson.M
//     cursor, err := productCollection.Find(context.Background(), bson.M{"custom_id": bson.M{"$in": wishlist.ProductIDs}})
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
//         return
//     }
//     defer cursor.Close(context.Background())

//     if err = cursor.All(context.Background(), &products); err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode products"})
//         return
//     }

//     c.JSON(http.StatusOK, products)
// }

