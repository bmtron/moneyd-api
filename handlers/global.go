package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// Authorized handlers that enforce user ownership

func GetHandlerAuthorized[T any](getFunc func(id int, authenticatedUserID int, db *sql.DB) (T, error), db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		itemId := c.Param("id")
		itemIdInt, err := strconv.Atoi(itemId)
		if err != nil {
			log.Print(err)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		item, dbErr := getFunc(itemIdInt, userID.(int), db)
		if dbErr != nil {
			log.Print(dbErr)
			if dbErr == sql.ErrNoRows {
				c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Resource not found or access denied"})
			} else {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}
		c.IndentedJSON(http.StatusOK, item)
	}
}

func UpdateHandlerAuthorized[T any](updateFunc func(id int, model T, authenticatedUserID int, db *sql.DB) (T, error), db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		itemId := c.Param("id")
		itemIdInt, err := strconv.Atoi(itemId)
		if err != nil {
			log.Print(err)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var updatedItem T
		if err := c.BindJSON(&updatedItem); err != nil {
			log.Print(err)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		itemResult, dbErr := updateFunc(itemIdInt, updatedItem, userID.(int), db)
		if dbErr != nil {
			log.Print(dbErr)
			if dbErr == sql.ErrNoRows {
				c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Resource not found or access denied"})
			} else {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		c.IndentedJSON(http.StatusOK, itemResult)
	}
}

func DeleteHandlerAuthorized[T any](deleteFunc func(itemId int, authenticatedUserID int, db *sql.DB) (T, error), db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		itemId := c.Param("id")
		itemIdInt, err := strconv.Atoi(itemId)
		if err != nil {
			log.Print(err)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		result, dbErr := deleteFunc(itemIdInt, userID.(int), db)
		if dbErr != nil {
			log.Print(dbErr)
			if dbErr == sql.ErrNoRows {
				c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Resource not found or access denied"})
			} else {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		c.IndentedJSON(http.StatusOK, result)
	}
}

func CreateHandlerAuthorized[T any](createFunc func(model T, authenticatedUserID int, db *sql.DB) (T, error), db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var model T
		if err := c.BindJSON(&model); err != nil {
			log.Print(err)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		result, dbErr := createFunc(model, userID.(int), db)
		if dbErr != nil {
			log.Print(dbErr)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		c.IndentedJSON(http.StatusCreated, result)
	}
}

func CreateBatchHandlerAuthorized[T any](createBatchFunc func(models []T, authenticatedUserID int, db *sql.DB) ([]T, error), db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var models []T
		if err := c.BindJSON(&models); err != nil {
			log.Print(err)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		createdModels, dbErr := createBatchFunc(models, userID.(int), db)
		if dbErr != nil {
			log.Print(dbErr)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		c.IndentedJSON(http.StatusCreated, createdModels)
	}
}

// GetHandlerByUserIdAuthorized handles GET requests where the route parameter :id represents a user_id
// It validates that the requested user_id matches the authenticated user's ID from the JWT
func GetHandlerByUserIdAuthorized[T any](getFunc func(userId int, authenticatedUserID int, db *sql.DB) (T, error), db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestedUserId := c.Param("id")
		requestedUserIdInt, err := strconv.Atoi(requestedUserId)
		if err != nil {
			log.Print(err)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		authenticatedUserID, exists := c.Get("user_id")
		if !exists {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		log.Print("Requested id: ", requestedUserId);
		log.Print("Authenticated id: ", authenticatedUserID);
		// Validate that the requested user_id matches the authenticated user
		if requestedUserIdInt != authenticatedUserID.(int) {
			c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Access denied: cannot access other users' data"})
			return
		}

		item, dbErr := getFunc(requestedUserIdInt, authenticatedUserID.(int), db)
		if dbErr != nil {
			log.Print(dbErr)
			if dbErr == sql.ErrNoRows {
				c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
			} else {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		c.IndentedJSON(http.StatusOK, item)
	}
}

func GetHandlerIndeterminiteArgsAuthorized[T any](getFunc func(db *sql.DB, args []int, authenticatedUserID int) (T, error), db *sql.DB, argcount int, userIdParamIndex int) gin.HandlerFunc {
	return func(c *gin.Context) {
		finalArgs := []int{}
		for i := range argcount {
			tempId := c.Param("id" + strconv.Itoa(i + 1))
			tempIdInt, err := strconv.Atoi(tempId)
			if err != nil {
				log.Print(err)
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
				return
			}
			finalArgs = append(finalArgs, tempIdInt)
		}

		authenticatedUserID, exists := c.Get("user_id")
		if !exists {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Validate that the user_id parameter matches the authenticated user
		// userIdParamIndex is 0-based (0 for id1, 1 for id2, etc.)
		if userIdParamIndex >= 0 && userIdParamIndex < len(finalArgs) {
			if finalArgs[userIdParamIndex] != authenticatedUserID.(int) {
				c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Access denied: cannot access other users' data"})
				return
			}
		}

		item, dbErr := getFunc(db, finalArgs, authenticatedUserID.(int))
		if dbErr != nil {
			log.Print(dbErr)
			if dbErr == sql.ErrNoRows {
				c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
			} else {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		c.IndentedJSON(http.StatusOK, item)
	}
}

func DeleteHandler[T any](deleteFunc func(itemId int, db *sql.DB) (T, error), db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		itemId := c.Param("id")
		itemIdInt, err := strconv.Atoi(itemId)

		if err != nil {
			log.Print(err)
			c.IndentedJSON(http.StatusBadRequest, err)
			return
		}
		result, dbErr := deleteFunc(itemIdInt, db)
		if dbErr != nil {
			log.Print(dbErr)
			c.IndentedJSON(http.StatusInternalServerError, dbErr)
			return
		}

		c.IndentedJSON(http.StatusOK, result)
	}
}

func CreateHandler[T any](createFunc func(model T, db *sql.DB) (T, error), db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var model T
		log.Print("Beginning JSON parse...")
		if err := c.BindJSON(&model); err != nil {
			log.Print(err)
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		log.Print(&model)

		log.Print("beginning db call...")
		result, dbErr := createFunc(model, db)
		if dbErr != nil {
			log.Print(dbErr)
			c.IndentedJSON(http.StatusInternalServerError, dbErr)
			return
		}

		c.IndentedJSON(http.StatusCreated, result)
	}
}

func GetHandlerIndeterminiteArgs[T any](getFunc func(db *sql.DB, args []int) (T, error), db *sql.DB, argcount int) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		finalArgs := []int{}
		for i := range argcount {
			tempId := c.Param("id" + strconv.Itoa(i + 1))
			tempIdInt, err := strconv.Atoi(tempId)
			if err != nil {
				log.Print(err)
				c.IndentedJSON(http.StatusBadRequest, err)
				return
			}
			finalArgs = append(finalArgs, tempIdInt)
		}
		item, dbErr := getFunc(db, finalArgs)
		if dbErr != nil {
			log.Print(dbErr)
			c.IndentedJSON(http.StatusInternalServerError, dbErr)
			return
		}

		c.IndentedJSON(http.StatusOK, item);
	}
}

func GetHandler[T any](getFunc func(id int, db *sql.DB) (T, error), db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		itemId := c.Param("id")
		itemIdInt, err := strconv.Atoi(itemId)
		if err != nil {
			log.Print(err)
			c.IndentedJSON(http.StatusBadRequest, err)
			return
		}
		item, dbErr := getFunc(itemIdInt, db)
		if dbErr != nil {
			log.Print(dbErr)
			c.IndentedJSON(http.StatusInternalServerError, dbErr)
			return
		}
		c.IndentedJSON(http.StatusOK, item)
	}
}

func UpdateBatchHandler[T any](updateFunc func(models []T, db *sql.DB) ([]T, error), db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var models []T
		if err := c.BindJSON(&models); err != nil {
			log.Print(err)
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		updatedModels, dbErr := updateFunc(models, db)
		if dbErr != nil {
			log.Print(dbErr)
			c.IndentedJSON(http.StatusInternalServerError, dbErr)
			return
		}
		c.IndentedJSON(http.StatusOK, updatedModels)
	}
}

func CreateBatchHandler[T any](createBatchFunc func(models []T, db *sql.DB) ([]T, error), db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var models []T
		if err := c.BindJSON(&models); err != nil {
			log.Print(err)
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		createdModels, dbErr := createBatchFunc(models, db)
		if dbErr != nil {
			log.Print(dbErr)
			c.IndentedJSON(http.StatusInternalServerError, dbErr)
			return
		}

		c.IndentedJSON(http.StatusOK, createdModels)
	}
}

func UpdateHandler[T any](updateFunc func(id int, model T, db *sql.DB) (T, error), db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		itemId := c.Param("id")
		itemIdInt, err := strconv.Atoi(itemId)
		if err != nil {
			log.Print(err)
			c.IndentedJSON(http.StatusBadRequest, err)
			return
		}

		var updatedItem T
		if err := c.BindJSON(&updatedItem); err != nil {
			log.Print(err)
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		itemResult, dbErr := updateFunc(itemIdInt, updatedItem, db)
		if dbErr != nil {
			log.Print(dbErr)
			c.IndentedJSON(http.StatusInternalServerError, dbErr)
			return
		}

		c.IndentedJSON(http.StatusOK, itemResult)
	}
}
