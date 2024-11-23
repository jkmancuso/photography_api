package main

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"log"
	"net/mail"

	"github.com/aws/aws-lambda-go/events"
)

type SimpleLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func parseEventBody(e events.APIGatewayProxyRequest) (string, string) {
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

func generateHash(s string) string {
	h := sha512.New()
	h.Write([]byte(s))

	return fmt.Sprintf("%x", h.Sum(nil))
}
