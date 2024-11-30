package shared

import (
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/mail"
	"slices"
	"strings"

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

func GetTargetEndpoint(path string, withoutLeadingSlash bool) string {
	resultStr := strings.Split(path, "/")[1]

	if withoutLeadingSlash {
		return resultStr
	}

	return "/" + resultStr
}

func ValidateEvent(e events.APIGatewayProxyRequest) error {

	if len(e.HTTPMethod) == 0 {
		return errors.New("no idea what this is")
	}

	if slices.Contains([]string{"POST", "PATCH"}, e.HTTPMethod) && len(e.Body) == 0 {
		return errors.New("event body is empty")

	}

	if slices.Contains([]string{"GET", "DELETE"}, e.HTTPMethod) && len(e.Body) != 0 {
		return fmt.Errorf("no body should be sent for method %v", e.Path)
	}
	/*add back in later
	if !strings.Contains(e.Headers["Set-Cookie"], "token=") {
		return errors.New("missing auth cookie")
	}*/

	return nil
}

func GenerateUUID() string {
	id, _ := uuid.NewV7()
	return id.String()
}
