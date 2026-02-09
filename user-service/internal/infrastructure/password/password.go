package password

import "golang.org/x/crypto/bcrypt"

const HASH_COST = 14

func GenerateHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), HASH_COST)
	return string(bytes), err
}

func CompareHashAndPasword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
