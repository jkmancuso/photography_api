package shared

import (
	"crypto/sha512"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func GenerateHash(s string, salt string) (string, error) {

	if len(s) == 0 || len(salt) == 0 {
		return "", errors.New("missing string to generate hash")
	}

	h := sha512.New()
	h.Write([]byte(s + salt))
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func GenerateUUID() string {
	id, _ := uuid.NewV7()
	return id.String()
}

func IsUUID(s string) bool {

	if _, err := uuid.Parse(s); err != nil {
		return false
	}

	return true
}
