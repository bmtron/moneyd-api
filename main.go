package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"log"
	"moneyd/api/database"
	"moneyd/api/handlers"
	"net/http"
	"os"
	"strconv"
	"time"
	"strings"
)

var db *sql.DB

func main() {

	log.Print("Starting up server...")
	db = database.SetupDb()
	router := gin.Default()

	router.Use(gin.Logger())

	config := cors.DefaultConfig()
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
		log.Fatal("Error loading env vars.")
	}

	expectedJwtSecret := os.Getenv("JWT_SECRET")
	expectedApiKey := os.Getenv("API_KEY")

	config.AllowOrigins = []string{"http://localhost:8085"}
	config.AllowMethods = []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-API-Key"}
	router.Use(cors.New(config))
	router.Use(func(c *gin.Context) {
		fmt.Printf("Request origin: %s\n", c.Request.Header.Get("Origin"))
		c.Next()
	})

	router.POST("/auth/login", apiKeyMiddleware(expectedApiKey), handlers.LoginHandler(db))
	router.GET("/auth/me", apiKeyMiddleware(expectedApiKey), AuthMiddleware([]byte(expectedJwtSecret)), handlers.GetUserInfoHandler(db))

	api := router.Group("/api")
	api.GET("/test", testHandler)

	api.POST("/users", apiKeyMiddleware(expectedApiKey), handlers.CreateHandler(handlers.CreateUser, db))
	api.Use(apiKeyMiddleware(expectedApiKey), AuthMiddleware([]byte(expectedJwtSecret)))
	{
		api.GET("/users/:id", handlers.GetHandler(handlers.GetUser, db))
		api.PUT("/users/:id", handlers.UpdateHandler(handlers.UpdateUser, db))
		api.DELETE("/users/:id", handlers.DeleteHandler(handlers.DeleteUser, db))
	}

	log.Print("Setup complete...")
	log.Print("Running...")

	router.Run("0.0.0.0:8085")
}

func apiKeyMiddleware(expectedApiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
			c.Abort()
			return
		}

		if apiKey != expectedApiKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		c.Next()
	}
}
func AuthMiddleware(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Enforce HS256
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims := token.Claims.(*jwt.RegisteredClaims)

		// Check expiration
		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		userID, _ := strconv.Atoi(claims.Subject)
		c.Set("user_id", userID)

		c.Next()
	}
}
func authMiddleware_old(expectedJwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
			}
			return []byte(expectedJwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userID := int(claims["user_id"].(float64))
		c.Set("user_id", userID)
		c.Next()
	}
}

func testHandler(c *gin.Context) {
	_, err := database.Test(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch test item"})
		return
	}
	c.JSON(http.StatusOK, "success retrieving test data")
}
