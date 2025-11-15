package database

import (
	"database/sql"
	"log"
	"moneyd/api/models"

	_ "github.com/lib/pq"
	"moneyd/api/utils"
)

type User struct {
	models.BankingUser
	PlainPassword string `json:"password"`
}


func CreateUser(user User, db *sql.DB) (User, error) {
	user.PasswordHash = utils.PasswordHasher(user.PlainPassword);
	query := `
        INSERT INTO banking_user (username, email, password_hash, date_created, date_updated)
        VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        RETURNING banking_user_id, date_created, date_updated
        `
	err := db.QueryRow(query, user.Username, user.Email, user.PasswordHash).Scan(&user.BankingUserId, &user.DateCreated, &user.DateUpdated)
	if err != nil {
		log.Print(err)
		return user, err
	}

	return user, nil
}

func GetUser(userID int, db *sql.DB) (User, error) {
	var user User
	query := `
        SELECT banking_user_id, username, email, password_hash, date_created, date_updated
        FROM banking_user
        WHERE banking_user_id = $1
        `
	err := db.QueryRow(query, userID).Scan(&user.BankingUserId, &user.Username, &user.Email, &user.PasswordHash, &user.DateCreated, &user.DateUpdated)
	if err != nil {
		log.Print(err)
		return user, err
	}

	return user, nil
}

func GetUserByEmail(email string, db *sql.DB) (User, error) {
	var user User
	query := `
        SELECT banking_user_id, username, email, password_hash, date_created, date_updated
        FROM banking_user
        WHERE email = $1
        `
	err := db.QueryRow(query, email).Scan(&user.BankingUserId, &user.Username, &user.Email, &user.PasswordHash, &user.DateCreated, &user.DateUpdated)
	if err != nil {
		log.Print(err)
		return user, err
	}

	return user, nil
}

func UpdateUser(userID int, updatedUser User, db *sql.DB) (User, error) {
	updatedUser.PasswordHash = utils.PasswordHasher(updatedUser.PlainPassword)
	query := `
        UPDATE banking_user
        SET username = $1, email = $2, password_hash = $3, date_created = CURRENT_TIMESTAMP
        WHERE banking_user_id = $4
        RETURNING banking_user_id, username, email, password_hash, date_created, date_updated
        `
	err := db.QueryRow(query, updatedUser.Username, updatedUser.Email, updatedUser.PasswordHash, userID).Scan(&updatedUser.BankingUserId, &updatedUser.Username, &updatedUser.Email, &updatedUser.PasswordHash, &updatedUser.DateCreated, &updatedUser.DateUpdated)
	if err != nil {
		log.Print(err)
		return updatedUser, err
	}

	return updatedUser, nil
}

func DeleteUser(userID int, db *sql.DB) (User, error) {
	query := `
        DELETE FROM banking_user
        WHERE banking_user_id = $1
        RETURNING banking_user_id, username, email, password_hash, date_created, date_updated
        `
	var deletedUser User
	err := db.QueryRow(query, userID).Scan(&deletedUser.BankingUserId, &deletedUser.Username, &deletedUser.Email, &deletedUser.PasswordHash, &deletedUser.DateCreated, &deletedUser.DateUpdated)
	if err != nil {
		log.Print(err)
		return deletedUser, err
	}

	return deletedUser, nil
}

