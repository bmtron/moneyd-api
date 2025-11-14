package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func SetupDb() *sql.DB {
	connStr := getDbInfo()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func Test(db *sql.DB) (string, error) {
	rows, err := db.Query("SELECT data FROM test where testid = $1", 1)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		var res string
		if err := rows.Scan(&res); err != nil {
			log.Fatal(err)
			return "", err
		}
		fmt.Println(res)
	}
	return "success", nil
}

func getDbInfo() string {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
		log.Fatal("Error loading .env file")
	}
	dbInfo := os.Getenv("DATABASE_URL")
	return dbInfo
}
