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
