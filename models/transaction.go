package models

import (
	"time"
)


type Transaction struct {
	TransactionId				int			`json:"transaction_id"`
	StatementId					int			`json:"statement_id"`
	TransactionTypeLookupCode 	int   		`json:"transaction_type_lookup_code"`
	Description					string		`json:"description"`
	Amount						int64		`json:"amount"`
	TransactionDate				time.Time	`json:"transaction_date"`
	DateAdded					time.Time	`json:"date_added"`
	DateUpdated					time.Time	`json:"date_updated"`
}
