package utils

import (
	"log"
	"golang.org/x/crypto/bcrypt"
)

func PasswordHasher(password string) string {
	passbytes := []byte(password)

	hashedPassword, err := bcrypt.GenerateFromPassword(passbytes, bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	// quick sanity check
	err = bcrypt.CompareHashAndPassword(hashedPassword, passbytes)
	if err != nil {
		log.Fatal("PASSWORD_MISMATCH ABORTING")
	}

	return string(hashedPassword)
}
