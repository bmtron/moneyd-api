package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
	"database/sql"
	"os"
	"github.com/joho/godotenv"
	"moneyd/api/database"
)

type UserResponse struct {
	Id       int	`json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func LoginHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		user, err := database.GetUserByEmail(loginRequest.Email, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginRequest.Password)); err != nil {
			log.Printf("Password comparison failed with error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.BankingUserId,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		})

		enverr := godotenv.Load()
		if enverr != nil {
			log.Print(enverr)
			log.Fatal("Error loading env vars.")
		}

		// In production, this should be an environment variable
		expectedJwtSecret := os.Getenv("JWT_SECRET")

		tokenString, err := token.SignedString([]byte(expectedJwtSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		userResponse := UserResponse{
			Id:       user.BankingUserId,
			Email:    user.Email,
			Username: user.Username,
		}

		c.JSON(http.StatusOK, gin.H{
			"token": tokenString,
			"user":  userResponse,
		})
	}

}
