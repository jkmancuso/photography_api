package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type UserLogin struct {
	email            string
	hashpass         string
	responseHTTPCode int
	responseHTTPMsg  string
}

func (login *UserLogin) setUserPass(email string, password string) {
	login.email = email
	login.hashpass = generateHash(password)
}

func (login *UserLogin) setstatusCode(code int) {
	login.responseHTTPCode = code
}

func (login *UserLogin) setHTTPMsg(msg string) {
	login.responseHTTPMsg = msg
}

func NewLogin(req events.APIGatewayProxyRequest) (*UserLogin, error) {
	log.Println("Entering NewLogin")

	email, password := parseEventBody(req)

	login := UserLogin{}
	var err error

	if len(email) == 0 || len(password) == 0 {
		login.setstatusCode(http.StatusBadRequest)
		login.setHTTPMsg(`{"STATUS":"INVALID_REQUEST"}`)
		err = errors.New("missing email or password in body")
	} else {
		login.setUserPass(email, password)
		login.setstatusCode(http.StatusOK)
	}

	return &login, err

}
