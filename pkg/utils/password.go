package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password),14)
	return string(b), err
}

func ComparePassword(hash, password string ) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}