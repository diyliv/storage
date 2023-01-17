package hash

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPass(pass []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	if err != nil {
		log.Printf("Error while generating password: %v\n", err)
	}

	return string(hash)
}

func ComparePass(hashedPass string, plainPass []byte) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPass), plainPass); err != nil {
		log.Printf("Error while comparing password: %v\n", err)
		return false
	}

	return true
}
