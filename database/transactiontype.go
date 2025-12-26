package database

import (
	"database/sql"
	"moneyd/api/models"
)

type TransactionTypeLookup = models.TransactionTypeLookup

func GetTransactionTypes(db *sql.DB) ([]TransactionTypeLookup, error) {
	var transactiontypes []TransactionTypeLookup
	query := `
		SELECT transaction_type_lookup_code, description
		FROM transaction_type_lookup;
		`
	rows, err := db.Query(
		query,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transactiontype TransactionTypeLookup
		if err := rows.Scan(
			&transactiontype.TransactionTypeLookupCode,
			&transactiontype.Description,
		); err != nil {
			return transactiontypes, err
		}
		transactiontypes = append(transactiontypes, transactiontype)
	}

	if err = rows.Err(); err != nil {
		return transactiontypes, nil
	}

	return transactiontypes, nil
}
