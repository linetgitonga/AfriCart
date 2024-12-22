package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "go-bookstore/controllers"
    "go-bookstore/database"
    "github.com/gin-contrib/cors"

    // "context"
    // "fmt"
    // "log"
    // "time"


    // "go.mongodb.org/mongo-driver/bson"
    // "go.mongodb.org/mongo-driver/mongo"
    // "go.mongodb.org/mongo-driver/mongo/options"
    // "go-bookstore/database"
)

func main() {
    database.Connect() // Initialize MongoDB connection





    // Initialize the Gin router
    router := gin.Default()
    router.SetTrustedProxies([]string{"127.0.0.1"})


    // Serve static files (CSS, JS, images)
    router.Use(cors.Default())
    router.Static("/static", "./static")
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://127.0.0.1:8000"}, // Allow frontend origin
        // AllowOrigins:     []string{"http://127.0.0.1:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Content-Type", "Authorization"},
        AllowCredentials: true,
    }))
    
    // Load HTML templates from the templates folder
    router.LoadHTMLGlob("templates/*.html")

    router.GET("/signup", func(c *gin.Context) {
        c.HTML(200, "signup.html", nil)
    })
    router.GET("/login", func(c *gin.Context) {
        c.HTML(200, "login.html", nil)
    })

     // Auth routes
    router.POST("/signup", controllers.Signup)
    router.POST("/login", controllers.Login)

    router.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", nil)
    })
    
    router.GET("/shop", func(c *gin.Context) {
        c.HTML(http.StatusOK, "shop.html", nil)
    })
    router.GET("/checkout", func(c *gin.Context) {
        c.HTML(http.StatusOK, "checkout.html", nil)
    })
    router.GET("/product-details", func(c *gin.Context) {
        c.HTML(http.StatusOK, "product-details.html", nil)
    })
    router.GET("/shop-cart", func(c *gin.Context) {
        c.HTML(http.StatusOK, "shop-cart.html", nil)
    })
    router.GET("/dashboard", func(c *gin.Context) {
        c.HTML(http.StatusOK, "dashboard.html", nil)
    })
    router.GET("/contact", func(c *gin.Context) {
        c.HTML(http.StatusOK, "contact.html", nil)
    })

    router.GET("/wishlist", func(c *gin.Context) {
        c.HTML(http.StatusOK, "wishlist.html", nil)
    })
    


    router.GET("/admin_products", func(c *gin.Context) {
        c.HTML(http.StatusOK, "admin_products.html", nil)
    })

    router.GET("/admin_users", func(c *gin.Context) {
        c.HTML(http.StatusOK, "admin_users.html", nil)
    })

    router.GET("/admin_profile", func(c *gin.Context) {
        c.HTML(http.StatusOK, "admin_profile.html", nil)
    })

    router.GET("/profile", func(c *gin.Context) {
        c.HTML(http.StatusOK, "profile.html", nil)
    })
    router.GET("/blog", func(c *gin.Context) {
        c.HTML(http.StatusOK, "blog.html", nil)
    })

    router.GET("/api/categories", controllers.GetCategoryCounts)

     // Routes for products
    router.GET("/api/products", controllers.GetProducts)
    router.POST("/api/products", controllers.AddProduct)
    
    
    
    // router.GET("/api/products/:custom_id", controllers.GetProductByCustomID)

    // router.GET("/api/products/:_id", controllers.GetProductByID)

    router.GET("/api/products/:custom_id", controllers.GetProductByCustomID)

    router.GET("/api/stats", controllers.GetStats)
    router.GET("/api/users", controllers.GetUsers)
    router.GET("/api/items", controllers.GetItems)
    router.POST("/api/items", controllers.AddProduct)


    router.DELETE("/api/users/:id", controllers.DeleteUser)
    router.POST("/api/users/", controllers.AddUser)
    router.GET("/api/users/me", controllers.GetCurrentUser)
    router.POST("/api/messages", controllers.AddMessage)
    router.GET("/api/messages", controllers.GetMessage)


    // router.POST("/api/messages", controllers.AddToWishlist)
    // router.GET("/api/wishlist", controllers.GetWishlist)
    // router.DELTE("/api/wishlist", controllers.RemoveFromWishlist)







    
    

    // Start the server on port 8080
    router.Run(":8000")
}
