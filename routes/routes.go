// routes/routes.go
package routes

import (
    "github.com/gin-gonic/gin"
    "go-bookstore/controllers"
)

func SetupRouter() *gin.Engine {
    router := gin.Default()
    
    // API route for products
    // router.GET("/api/products", controllers.GetProducts)
    router.GET("/api/categories", controllers.GetCategoryCounts)

    // Product routes
    router.POST("/api/products", controllers.AddProduct)
    router.GET("/api/products/:custom_id", controllers.GetProductByCustomID)

    return router
}
