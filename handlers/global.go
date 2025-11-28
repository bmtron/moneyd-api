package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

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
