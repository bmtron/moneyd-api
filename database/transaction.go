package database

import (
	"database/sql"
	"fmt"
	"log"
	"moneyd/api/models"
	"strings"
)

type Transaction = models.Transaction

func CreateTransaction(txn Transaction, db *sql.DB) (Transaction, error) {
	query := `
		INSERT INTO transaction (statement_id, transaction_type_lookup_code, description, amount, transaction_date, date_added, date_updated)
		VALUES ($1, $2, $3, ($4)::NUMERIC(14,2) / 100, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING transaction_id, statement_id, description, transaction_date, date_added, date_updated
	`
	err := db.QueryRow(
		query,
		txn.StatementId,
		txn.TransactionTypeLookupCode,
		txn.Description,
		txn.Amount,
		txn.TransactionDate,
	).Scan(
		&txn.TransactionId,
		&txn.StatementId,
		&txn.Description,
		&txn.TransactionDate,
		&txn.DateAdded,
		&txn.DateUpdated,
	)

	if err != nil {
		log.Print(err)
		return txn, err
	}
	return txn, nil
}

// CreateTransactionAuthorized creates a transaction only if the statement belongs to the authenticated user
func CreateTransactionAuthorized(txn Transaction, authenticatedUserID int, db *sql.DB) (Transaction, error) {
	// First verify the statement belongs to the user
	var count int
	verifyQuery := `SELECT COUNT(*) FROM statement WHERE statement_id = $1 AND banking_user_id = $2`
	err := db.QueryRow(verifyQuery, txn.StatementId, authenticatedUserID).Scan(&count)
	if err != nil {
		log.Print(err)
		return txn, err
	}
	if count == 0 {
		return txn, fmt.Errorf("statement not found or access denied")
	}

	return CreateTransaction(txn, db)
}

func CreateTransactionsBatch(txns []Transaction, db *sql.DB) ([]Transaction, error) {
	txnCols := 4
	var sb strings.Builder
	args := make([]interface{}, 0, len(txns)*txnCols)

	placeholder := 1
	sb.WriteString("INSERT INTO transaction (statement_id, transaction_type_lookup_code, description, amount, transaction_date, date_added, date_updated)")
	sb.WriteString("VALUES ")

	for index, txn := range txns {
		sb.WriteString("(")
		sb.WriteString(fmt.Sprintf("$%d, $%d, $%d, ($%d)::NUMERIC(14,2) / 100, $%d, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP", placeholder, placeholder+1, placeholder+2, placeholder+3, placeholder+4))
		sb.WriteString(")")

		if index < len(txns)-1 {
			sb.WriteString(",")
		}
		args = append(args, txn.StatementId, txn.TransactionTypeLookupCode, txn.Description, txn.Amount, txn.TransactionDate)
		placeholder += 5

	}
	sb.WriteString(" RETURNING transaction_id, statement_id, description, (amount * 100)::INTEGER, transaction_date, date_added, date_updated")

	query := sb.String()

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inserted []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(
			&t.TransactionId,
			&t.StatementId,
			&t.Description,
			&t.Amount,
			&t.TransactionDate,
			&t.DateAdded,
			&t.DateUpdated,
		); err != nil {
			return nil, err
		}
		inserted = append(inserted, t)
	}

	return inserted, rows.Err()

}

// CreateTransactionsBatchAuthorized creates multiple transactions only if all statements belong to the authenticated user
func CreateTransactionsBatchAuthorized(txns []Transaction, authenticatedUserID int, db *sql.DB) ([]Transaction, error) {
	// Collect unique statement IDs
	statementIds := make(map[int]bool)
	for _, txn := range txns {
		statementIds[txn.StatementId] = true
	}

	// Verify all statements belong to the user
	for stmtId := range statementIds {
		var count int
		verifyQuery := `SELECT COUNT(*) FROM statement WHERE statement_id = $1 AND banking_user_id = $2`
		err := db.QueryRow(verifyQuery, stmtId, authenticatedUserID).Scan(&count)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		if count == 0 {
			return nil, fmt.Errorf("statement %d not found or access denied", stmtId)
		}
	}

	return CreateTransactionsBatch(txns, db)
}

func GetTransaction(transactionId int, db *sql.DB) (Transaction, error) {
	var txn Transaction
	query := `
	SELECT transaction_id, transaction_type_lookup_code, statement_id, description, (amount * 100)::INTEGER, transaction_date, date_added, date_updated
		FROM transaction
		WHERE transaction_id = $1
	`
	err := db.QueryRow(query, transactionId).Scan(
		&txn.TransactionId,
		&txn.StatementId,
		&txn.TransactionTypeLookupCode,
		&txn.Description,
		&txn.Amount,
		&txn.TransactionDate,
		&txn.DateAdded,
		&txn.DateUpdated,
	)
	if err != nil {
		log.Print(err)
		return txn, err
	}
	return txn, nil
}

// GetTransactionAuthorized retrieves a transaction only if it belongs to the authenticated user
func GetTransactionAuthorized(transactionId int, authenticatedUserID int, db *sql.DB) (Transaction, error) {
	var txn Transaction
	query := `
	SELECT t.transaction_id, t.transaction_type_lookup_code, t.statement_id, t.description, (t.amount * 100)::INTEGER, t.transaction_date, t.date_added, t.date_updated
		FROM transaction t
		JOIN statement s ON s.statement_id = t.statement_id
		WHERE t.transaction_id = $1 AND s.banking_user_id = $2
	`
	err := db.QueryRow(query, transactionId, authenticatedUserID).Scan(
		&txn.TransactionId,
		&txn.TransactionTypeLookupCode,
		&txn.StatementId,
		&txn.Description,
		&txn.Amount,
		&txn.TransactionDate,
		&txn.DateAdded,
		&txn.DateUpdated,
	)
	if err != nil {
		log.Print(err)
		return txn, err
	}
	return txn, nil
}

func GetTransactionsByStatementId(statementId int, db *sql.DB) ([]Transaction, error) {
	var txns []Transaction
	query := `
	SELECT t.transaction_id, t.statement_id, t.transaction_type_lookup_code, t.description, (t.amount * 100)::INTEGER, t.transaction_date, t.date_added, t.date_updated
		FROM transaction t
		JOIN statement s on s.statement_id = t.statement_id
		WHERE s.statement_id = $1;
		`

	rows, err := db.Query(
		query,
		statementId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var txn Transaction
		if err := rows.Scan(
			&txn.TransactionId,
			&txn.StatementId,
			&txn.TransactionTypeLookupCode,
			&txn.Description,
			&txn.Amount,
			&txn.TransactionDate,
			&txn.DateAdded,
			&txn.DateUpdated,
		); err != nil {
			return txns, err
		}
		txns = append(txns, txn)
	}

	if err = rows.Err(); err != nil {
		return txns, nil
	}

	return txns, nil
}

// GetTransactionsByStatementIdAuthorized retrieves transactions only if the statement belongs to the authenticated user
func GetTransactionsByStatementIdAuthorized(statementId int, authenticatedUserID int, db *sql.DB) ([]Transaction, error) {
	var txns []Transaction
	query := `
	SELECT t.transaction_id, t.statement_id, t.transaction_type_lookup_code, t.description, (t.amount * 100)::INTEGER, t.transaction_date, t.date_added, t.date_updated
		FROM transaction t
		JOIN statement s on s.statement_id = t.statement_id
		WHERE s.statement_id = $1 AND s.banking_user_id = $2;
		`

	rows, err := db.Query(
		query,
		statementId,
		authenticatedUserID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var txn Transaction
		if err := rows.Scan(
			&txn.TransactionId,
			&txn.StatementId,
			&txn.TransactionTypeLookupCode,
			&txn.Description,
			&txn.Amount,
			&txn.TransactionDate,
			&txn.DateAdded,
			&txn.DateUpdated,
		); err != nil {
			return txns, err
		}
		txns = append(txns, txn)
	}

	if err = rows.Err(); err != nil {
		return txns, nil
	}

	return txns, nil
}

func GetTransactionsByUserId(userId int, db *sql.DB) ([]Transaction, error) {
	var txns []Transaction
	query := `
	SELECT t.transaction_id, t.statement_id, t.transaction_type_lookup_code, t.description, (t.amount * 100)::INTEGER, t.transaction_date, t.date_added, t.date_updated
		FROM transaction t
		JOIN statement s on s.statement_id = t.statement_id
		WHERE s.banking_user_id = $1;
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
		var txn Transaction
		if err := rows.Scan(
			&txn.TransactionId,
			&txn.StatementId,
			&txn.TransactionTypeLookupCode,
			&txn.Description,
			&txn.Amount,
			&txn.TransactionDate,
			&txn.DateAdded,
			&txn.DateUpdated,
		); err != nil {
			return txns, err
		}
		txns = append(txns, txn)
	}

	if err = rows.Err(); err != nil {
		return txns, nil
	}

	return txns, nil
}

// GetTransactionsByUserIdAuthorized retrieves transactions only for the authenticated user
func GetTransactionsByUserIdAuthorized(userId int, authenticatedUserID int, db *sql.DB) ([]Transaction, error) {
	if userId != authenticatedUserID {
		return []Transaction{}, nil
	}
	return GetTransactionsByUserId(userId, db)
}

func UpdateTransaction(txnId int, txn Transaction, db *sql.DB) (Transaction, error) {
	query := `
		UPDATE transaction
		SET statement_id = $1,
		    description  = $2,
			amount       = ($3)::NUMERIC(14,2) / 100,
		    transaction_date = $4,
		    date_updated = CURRENT_TIMESTAMP
		WHERE transaction_id = $5
		RETURNING transaction_id, statement_id, description, (amount * 100)::INTEGER, transaction_date, date_added, date_updated
	`
	err := db.QueryRow(
		query,
		txn.StatementId,
		txn.Description,
		txn.Amount,
		txn.TransactionDate,
		txnId,
	).Scan(
		&txn.TransactionId,
		&txn.StatementId,
		&txn.Description,
		&txn.Amount,
		&txn.TransactionDate,
		&txn.DateAdded,
		&txn.DateUpdated,
	)
	if err != nil {
		log.Print(err)
		return txn, err
	}
	return txn, nil
}

// UpdateTransactionAuthorized updates a transaction only if it belongs to the authenticated user
func UpdateTransactionAuthorized(txnId int, txn Transaction, authenticatedUserID int, db *sql.DB) (Transaction, error) {
	// Verify both the transaction and new statement belong to the user
	var existingStmtId int
	checkQuery := `
		SELECT t.statement_id
		FROM transaction t
		JOIN statement s ON s.statement_id = t.statement_id
		WHERE t.transaction_id = $1 AND s.banking_user_id = $2
	`
	err := db.QueryRow(checkQuery, txnId, authenticatedUserID).Scan(&existingStmtId)
	if err != nil {
		log.Print(err)
		return txn, err
	}

	// Verify new statement belongs to user (if changing statement)
	if txn.StatementId != existingStmtId {
		var count int
		verifyQuery := `SELECT COUNT(*) FROM statement WHERE statement_id = $1 AND banking_user_id = $2`
		err := db.QueryRow(verifyQuery, txn.StatementId, authenticatedUserID).Scan(&count)
		if err != nil {
			log.Print(err)
			return txn, err
		}
		if count == 0 {
			return txn, fmt.Errorf("new statement not found or access denied")
		}
	}

	return UpdateTransaction(txnId, txn, db)
}

func DeleteTransaction(transactionId int, db *sql.DB) (Transaction, error) {
	query := `
		DELETE FROM transaction
		WHERE transaction_id = $1
		RETURNING transaction_id, statement_id, description, (amount * 100)::INTEGER, transaction_date, date_added, date_updated
	`
	var txn Transaction
	err := db.QueryRow(query, transactionId).Scan(
		&txn.TransactionId,
		&txn.StatementId,
		&txn.Description,
		&txn.Amount,
		&txn.TransactionDate,
		&txn.DateAdded,
		&txn.DateUpdated,
	)
	if err != nil {
		log.Print(err)
		return txn, err
	}
	return txn, nil
}

// DeleteTransactionAuthorized deletes a transaction only if it belongs to the authenticated user
func DeleteTransactionAuthorized(transactionId int, authenticatedUserID int, db *sql.DB) (Transaction, error) {
	query := `
		DELETE FROM transaction t
		USING statement s
		WHERE t.transaction_id = $1
		AND t.statement_id = s.statement_id
		AND s.banking_user_id = $2
		RETURNING t.transaction_id, t.statement_id, t.description, (t.amount * 100)::INTEGER, t.transaction_date, t.date_added, t.date_updated
	`
	var txn Transaction
	err := db.QueryRow(query, transactionId, authenticatedUserID).Scan(
		&txn.TransactionId,
		&txn.StatementId,
		&txn.Description,
		&txn.Amount,
		&txn.TransactionDate,
		&txn.DateAdded,
		&txn.DateUpdated,
	)
	if err != nil {
		log.Print(err)
		return txn, err
	}
	return txn, nil
}

func GetTransactionsByInstitutionId(db *sql.DB, args []int) ([]Transaction, error) {
	userId := args[0]
	institutionId := args[1]
	var txns []Transaction
	query := `
	SELECT t.transaction_id, t.statement_id, t.transaction_type_lookup_code, t.description, (t.amount * 100)::INTEGER, t.transaction_date, t.date_added, t.date_updated
		FROM transaction t
		JOIN statement s on s.statement_id = t.statement_id
		WHERE s.banking_user_id = $1
		AND s.institution_id = $2;
		`
	rows, err := db.Query(
		query,
		userId,
		institutionId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var txn Transaction
		if err := rows.Scan(
			&txn.TransactionId,
			&txn.StatementId,
			&txn.TransactionTypeLookupCode,
			&txn.Description,
			&txn.Amount,
			&txn.TransactionDate,
			&txn.DateAdded,
			&txn.DateUpdated,
		); err != nil {
			return txns, err
		}
		txns = append(txns, txn)
	}

	if err = rows.Err(); err != nil {
		return txns, nil
	}

	return txns, nil

}

// GetTransactionsByInstitutionIdAuthorized retrieves transactions only for the authenticated user
func GetTransactionsByInstitutionIdAuthorized(db *sql.DB, args []int, authenticatedUserID int) ([]Transaction, error) {
	userId := args[0]
	if userId != authenticatedUserID {
		return []Transaction{}, nil
	}
	return GetTransactionsByInstitutionId(db, args)
}
