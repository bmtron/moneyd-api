package database

import (
	"database/sql"
	"log"
	"moneyd/api/models"
)

type Institution = models.Institution

func CreateInstitution(institution Institution, db *sql.DB) (Institution, error) {
	log.Print("creating institution...")
	query := `
		INSERT INTO institution (institution_id, name, date_added)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		RETURNING institution_id, name 
		`
	err := db.QueryRow(query,
		institution.InstitutionId,
		institution.Name,
	).Scan(
		&institution.InstitutionId,
		&institution.Name,
	)
	if err != nil {
		log.Print(err)
		return institution, err
	}

	return institution, nil
}

// CreateInstitutionAuthorized creates a institution using the authenticated user's ID
func CreateInstitutionAuthorized(institution Institution, authenticatedUserID int, db *sql.DB) (Institution, error) {
	// Override the banking_user_id with the authenticated user's ID
	return CreateInstitution(institution, db)
}

func GetInstitutions(db *sql.DB) ([]Institution, error) {
	var institutions []Institution
	query := `
		SELECT institution_id, name
		FROM institution;
		`
	rows, err := db.Query(
		query,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var inst Institution
		if err := rows.Scan(
			&inst.InstitutionId,
			&inst.Name,
		); err != nil {
			return institutions, err
		}
		institutions = append(institutions, inst)
	}

	if err = rows.Err(); err != nil {
		return institutions, nil
	}

	return institutions, nil
}
