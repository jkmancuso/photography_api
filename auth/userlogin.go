package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type UserLogin struct {
	email            string
	hashpass         string
	responseHTTPCode int
	responseHTTPMsg  string
	creationTime     time.Time
}

func (login *UserLogin) setUserPass(email string, password string, salt string) {
	login.email = email
	login.hashpass = generateHash(password, salt)
}

func (login *UserLogin) setstatusCode(code int) {
	login.responseHTTPCode = code
}

func (login *UserLogin) setHTTPMsg(msg string) {
	login.responseHTTPMsg = msg
}

func NewLogin(req events.APIGatewayProxyRequest, salt string) (*UserLogin, error) {
	log.Println("Entering NewLogin")

	email, password := parseEventBody(req)

	login := UserLogin{creationTime: time.Now()}
	var err error

	if len(email) == 0 || len(password) == 0 {
		login.setstatusCode(http.StatusBadRequest)
		login.setHTTPMsg(`{"STATUS":"INVALID_REQUEST"}`)
		err = errors.New("missing email or password in body")
	} else {
		login.setUserPass(email, password, salt)
		login.setstatusCode(http.StatusOK)
	}

	return &login, err

}
