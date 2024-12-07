package shared

import (
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/mail"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

type SimpleLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ExtractCredsFromEvent(e events.APIGatewayProxyRequest) (string, string) {
	if len(e.Body) == 0 {
		return "", ""
	}

	login := SimpleLogin{}

	err := json.Unmarshal([]byte(e.Body), &login)

	if err != nil {
		log.Println(err)
		return "", ""
	}

	_, err = mail.ParseAddress(login.Email)

	if err != nil {
		log.Println(err)
		return "", ""
	}

	return login.Email, login.Password
}

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
