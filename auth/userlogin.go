package main

import (
	"errors"
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

	email, password := parseEventBody(req)

	login := UserLogin{}
	var err error

	if len(email) == 0 || len(password) == 0 {
		login.setstatusCode(http.StatusBadRequest)
		login.setHTTPMsg(`{"STATUS":"INVALID_EVENT"`)
		err = errors.New("bad request")
	} else {
		login.setUserPass(email, password)
	}

	return &login, err

}
