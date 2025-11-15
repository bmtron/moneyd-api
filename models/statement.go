package models

import (
	"time"
)

type Statement struct {
	StatementId		int			`json:"statement_id"`
	BankingUserId	int			`json:"banking_user_id"`
	InstitutionId   int			`json:"institution_id"`
	PeriodStart		time.Time	`json:"period_start"`
	PeriodEnd		time.Time	`json:"period_end"`
	DateAdded		time.Time	`json:"date_added"`
}
