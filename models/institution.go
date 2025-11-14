package models

import (
	"time"
)

type Institution struct {
	InstitutionId 	int 		`json:"institution_id"`
	Name 			string		`json:"name"`
	StatementPeriod time.Time	`json:"statement_period"`
	DateAdded		time.Time	`json:"date_added"`
}
