package controllers

import (
    "golang.org/x/crypto/bcrypt"
    "github.com/dgrijalva/jwt-go"
    "time"
    "errors"
	"context"
    "net/http"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/gin-gonic/gin"
    "go-bookstore/database"
    "go-bookstore/models"
    "go.mongodb.org/mongo-driver/bson/primitive"
    // "github.com/golang-jwt/jwt/v4"
    
)

var jwtKey = []byte("your_secret_key") // Replace with a secure key

// HashPassword hashes the plain-text password
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// CheckPassword compares a plain-text password with a hashed password
func CheckPassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// GenerateJWT generates a new JWT token
func GenerateJWT(email string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "email": email,
        "exp":   time.Now().Add(time.Hour * 72).Unix(),
    })
    return token.SignedString(jwtKey)
}


// Signup handler
func Signup(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check if user already exists
    collection := database.Client.Database("AfricArt").Collection("users")
    var existingUser models.User
    err := collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
    if err == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
        return
    }

    // Hash password
    hashedPassword, err := HashPassword(user.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }
    user.Password = hashedPassword

    // Save user to MongoDB
    _, err = collection.InsertOne(context.TODO(), user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// Login handler
func Login(c *gin.Context) {
    var credentials struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := c.ShouldBindJSON(&credentials); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Fetch user from MongoDB
    collection := database.Client.Database("AfricArt").Collection("users")
    var user models.User
    err := collection.FindOne(context.TODO(), bson.M{"email": credentials.Email}).Decode(&user)
    if err == mongo.ErrNoDocuments {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    // Check password
    if err := CheckPassword(user.Password, credentials.Password); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    // Generate JWT token
    token, err := GenerateJWT(user.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}

func GetUsers(c *gin.Context) {
    collection := database.Client.Database("AfricArt").Collection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var users []bson.M
    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
        return
    }
    defer cursor.Close(ctx)

    if err = cursor.All(ctx, &users); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode users"})
        return
    }

    c.JSON(http.StatusOK, users)
}


// AddUser handles the creation of a new user
func AddUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    // Hash the password before storing it
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    // Set the hashed password
    user.Password = string(hashedPassword)

    collection := database.Client.Database("AfricArt").Collection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Insert the new user into the database
    _, err = collection.InsertOne(ctx, user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User added successfully"})
}

// DeleteUser deletes a user from the database
func DeleteUser(c *gin.Context) {
    userID := c.Param("id") // Get the ID from the URL parameter

    collection := database.Client.Database("AfricArt").Collection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var filter bson.M
    if primitive.IsValidObjectID(userID) {
        // If the userID is a valid MongoDB ObjectID, use it
        objectID, err := primitive.ObjectIDFromHex(userID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ObjectID format"})
            return
        }
        filter = bson.M{"_id": objectID}
    } else {
        // Otherwise, assume it’s a custom ID
        filter = bson.M{"custom_id": userID}
    }

    result, err := collection.DeleteOne(ctx, filter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
        return
    }

    if result.DeletedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}


// var jwtKey = []byte("your_secret_key") // Replace with a secure key



// Claims structure for JWT
type Claims struct {
    Email string `json:"email"`
    jwt.StandardClaims
}

// ValidateToken validates a JWT and extracts the email
func ValidateToken(tokenString string) (string, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil // Ensure jwtKey is used
    })

    if err != nil {
        if err == jwt.ErrSignatureInvalid {
            return "", errors.New("invalid token signature")
        }
        return "", errors.New("invalid token")
    }

    if !token.Valid {
        return "", errors.New("invalid token")
    }

    return claims.Email, nil
}

// GetCurrentUser retrieves the currently logged-in user's details
func GetCurrentUser(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
        return
    }

    // Assuming the token is stored as "Bearer <token>"
    token := authHeader[len("Bearer "):]

    // Validate the token and retrieve the user's email
    userEmail, err := ValidateToken(token)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    collection := database.Client.Database("AfricArt").Collection("users")
    var user models.User
    err = collection.FindOne(c, bson.M{"email": userEmail}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user details"})
        return
    }

    // Return user details (exclude sensitive data like password)
    c.JSON(http.StatusOK, gin.H{
        "name":   user.Name,
        "email":  user.Email,
        "role":   user.Role,
    })
}
