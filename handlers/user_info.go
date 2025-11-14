package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserInfoHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("user_id")
		user, err := GetUser(userID, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
			return
		}

		userResponse := UserResponse{
			Id:    user.BankingUserId,
			Email: user.Email,
			Username:  user.Username,
		}

		c.JSON(http.StatusOK, userResponse)
	}
}
