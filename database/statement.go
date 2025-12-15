package database

import (
	"database/sql"
	"log"
	"moneyd/api/models"
	"strconv"
)

type Statement = models.Statement

func CreateStatement(statement Statement, db *sql.DB) (Statement, error) {
	log.Print("creating statement...")
	query := `
		INSERT INTO statement (banking_user_id, institution_id, period_start, period_end, date_added)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
		RETURNING statement_id, banking_user_id, institution_id, period_start, period_end, date_added
		`
	err := db.QueryRow(query,
		statement.BankingUserId,
		statement.InstitutionId,
		statement.PeriodStart,
		statement.PeriodEnd,
	).Scan(
		&statement.StatementId,
		&statement.BankingUserId,
		&statement.InstitutionId,
		&statement.PeriodStart,
		&statement.PeriodEnd,
		&statement.DateAdded,
	)
	if err != nil {
		log.Print(err)
		return statement, err
	}

	return statement, nil
}

// CreateStatementAuthorized creates a statement using the authenticated user's ID
func CreateStatementAuthorized(statement Statement, authenticatedUserID int, db *sql.DB) (Statement, error) {
	// Override the banking_user_id with the authenticated user's ID
	statement.BankingUserId = authenticatedUserID
	return CreateStatement(statement, db)
}

func GetStatement(statementId int, db *sql.DB) (Statement, error) {
	var statement Statement
	log.Print("statement id is " + strconv.Itoa(statementId))
	query := `
		SELECT statement_id, banking_user_id, institution_id, period_start, period_end, date_added
		FROM statement
		WHERE statement_id = $1
		`
	err := db.QueryRow(
		query,
		statementId,
	).Scan(
		&statement.StatementId,
		&statement.BankingUserId,
		&statement.InstitutionId,
		&statement.PeriodStart,
		&statement.PeriodEnd,
		&statement.DateAdded,
	)
	if err != nil {
		log.Print(err)
		return statement, err
	}

	return statement, nil
}

// GetStatementAuthorized retrieves a statement only if it belongs to the authenticated user
func GetStatementAuthorized(statementId int, authenticatedUserID int, db *sql.DB) (Statement, error) {
	var statement Statement
	query := `
		SELECT statement_id, banking_user_id, institution_id, period_start, period_end, date_added
		FROM statement
		WHERE statement_id = $1 AND banking_user_id = $2
		`
	err := db.QueryRow(
		query,
		statementId,
		authenticatedUserID,
	).Scan(
		&statement.StatementId,
		&statement.BankingUserId,
		&statement.InstitutionId,
		&statement.PeriodStart,
		&statement.PeriodEnd,
		&statement.DateAdded,
	)
	if err != nil {
		log.Print(err)
		return statement, err
	}

	return statement, nil
}

func GetStatementsByUserId(userId int, db *sql.DB) ([]Statement, error) {
	var statements []Statement
	query := `
		SELECT statement_id, banking_user_id, institution_id, period_start, period_end, date_added
		FROM statement
		WHERE banking_user_id = $1
		`
	rows, err := db.Query(
		query,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stmnt Statement
		if err := rows.Scan(
			&stmnt.StatementId,
			&stmnt.BankingUserId,
			&stmnt.InstitutionId,
			&stmnt.PeriodStart,
			&stmnt.PeriodEnd,
			&stmnt.DateAdded,
		); err != nil {
			return statements, err
		}
		statements = append(statements, stmnt)
	}

	if err = rows.Err(); err != nil {
		return statements, nil
	}

	return statements, nil

}

// GetStatementsByUserIdAuthorized retrieves statements only for the authenticated user
func GetStatementsByUserIdAuthorized(userId int, authenticatedUserID int, db *sql.DB) ([]Statement, error) {
	if userId != authenticatedUserID {
		return []Statement{}, nil
	}
	return GetStatementsByUserId(userId, db)
}

func UpdateStatement(stmtId int, stmt Statement, db *sql.DB) (Statement, error) {
	query := `
		UPDATE statement
		SET banking_user_id = $1,
		    institution_id  = $2,
		    period_start    = $3,
		    period_end      = $4
		WHERE statement_id = $5
		RETURNING statement_id, banking_user_id, institution_id, period_start, period_end, date_added
	`
	err := db.QueryRow(
		query,
		stmt.BankingUserId,
		stmt.InstitutionId,
		stmt.PeriodStart,
		stmt.PeriodEnd,
		stmtId,
	).Scan(
		&stmt.StatementId,
		&stmt.BankingUserId,
		&stmt.InstitutionId,
		&stmt.PeriodStart,
		&stmt.PeriodEnd,
		&stmt.DateAdded,
	)
	if err != nil {
		log.Print(err)
		return stmt, err
	}
	return stmt, nil
}

// UpdateStatementAuthorized updates a statement only if it belongs to the authenticated user
func UpdateStatementAuthorized(stmtId int, stmt Statement, authenticatedUserID int, db *sql.DB) (Statement, error) {
	// Override the banking_user_id to prevent reassignment
	stmt.BankingUserId = authenticatedUserID
	query := `
		UPDATE statement
		SET institution_id  = $1,
		    period_start    = $2,
		    period_end      = $3
		WHERE statement_id = $4 AND banking_user_id = $5
		RETURNING statement_id, banking_user_id, institution_id, period_start, period_end, date_added
	`
	err := db.QueryRow(
		query,
		stmt.InstitutionId,
		stmt.PeriodStart,
		stmt.PeriodEnd,
		stmtId,
		authenticatedUserID,
	).Scan(
		&stmt.StatementId,
		&stmt.BankingUserId,
		&stmt.InstitutionId,
		&stmt.PeriodStart,
		&stmt.PeriodEnd,
		&stmt.DateAdded,
	)
	if err != nil {
		log.Print(err)
		return stmt, err
	}
	return stmt, nil
}

func DeleteStatement(statementId int, db *sql.DB) (Statement, error) {
	query := `
		DELETE FROM statement
		WHERE statement_id = $1
		RETURNING statement_id, banking_user_id, institution_id, period_start, period_end, date_added
		`
	var deletedStatement Statement
	err := db.QueryRow(
		query,
		statementId,
	).Scan(
		&deletedStatement.StatementId,
		&deletedStatement.BankingUserId,
		&deletedStatement.InstitutionId,
		&deletedStatement.PeriodStart,
		&deletedStatement.PeriodEnd,
		&deletedStatement.DateAdded,
	)
	if err != nil {
		log.Print(err)
		return deletedStatement, err
	}

	return deletedStatement, nil
}

// DeleteStatementAuthorized deletes a statement only if it belongs to the authenticated user
func DeleteStatementAuthorized(statementId int, authenticatedUserID int, db *sql.DB) (Statement, error) {
	query := `
		DELETE FROM statement
		WHERE statement_id = $1 AND banking_user_id = $2
		RETURNING statement_id, banking_user_id, institution_id, period_start, period_end, date_added
		`
	var deletedStatement Statement
	err := db.QueryRow(
		query,
		statementId,
		authenticatedUserID,
	).Scan(
		&deletedStatement.StatementId,
		&deletedStatement.BankingUserId,
		&deletedStatement.InstitutionId,
		&deletedStatement.PeriodStart,
		&deletedStatement.PeriodEnd,
		&deletedStatement.DateAdded,
	)
	if err != nil {
		log.Print(err)
		return deletedStatement, err
	}

	return deletedStatement, nil
}
