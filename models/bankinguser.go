package models

import (
	"time"
)

type BankingUser struct {
	BankingUserId		int 		`json:"banking_user_id"`
	Username			string		`json:"username"`
	Email				string 		`json:"email"`
	PasswordHash		string		`json:"password_hash"`
	DateCreated			time.Time	`json:"date_created"`
	DateUpdated			time.Time	`json:"date_updated"`
}
