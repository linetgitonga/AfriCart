package auth

import (
    "errors"
    // "time"

    "github.com/golang-jwt/jwt/v4"
)

// Secret key for signing the token
var jwtSecret = []byte("your_secret_key")

// Claims structure for JWT
type Claims struct {
    Email string `json:"email"`
    jwt.StandardClaims
}

// ValidateToken validates a JWT and extracts the email
func ValidateToken(tokenString string) (string, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
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
