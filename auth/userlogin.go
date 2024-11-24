package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jkmancuso/photography_api/shared"
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

	hashpass, err := shared.GenerateHash(password, salt)

	if err != nil {
		log.Println("ERROR generating hash")
	}
	login.hashpass = hashpass
}

func (login *UserLogin) setstatusCode(code int) {
	login.responseHTTPCode = code
}

func (login *UserLogin) setHTTPMsg(msg string) {
	login.responseHTTPMsg = msg
}

func NewLogin(req events.APIGatewayProxyRequest, salt string) (*UserLogin, error) {
	log.Println("Entering NewLogin")

	email, password := shared.ExtractCredsFromEvent(req)

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

func (login *UserLogin) getToken(db *shared.DBInfo) (string, error) {
	log.Println("Entering getToken")

	email, err := attributevalue.Marshal(login.email)

	if err != nil {
		return "", err
	}

	hashpass, err := attributevalue.Marshal(login.hashpass)

	if err != nil {
		return "", err
	}

	key := map[string]types.AttributeValue{"email": email, "hashpass": hashpass}

	query := dynamodb.GetItemInput{
		TableName: &db.Tablename,
		Key:       key,
	}

	response, err := db.Client.GetItem(context.Background(), &query)

	if err != nil {
		return "", err
	}

	if response.Item == nil {
		return `{"STATUS":"INVALID_USER_PASS"}`, errors.New("invalid user/pass")
	}

	adminItem := shared.DBAdminItem{}

	if err = attributevalue.UnmarshalMap(response.Item, &adminItem); err != nil {
		return "", err
	}

	rowJSON, err := json.Marshal(adminItem)

	if err != nil {
		return "", err
	}

	return string(rowJSON), err

}

func (login *UserLogin) recordLoginToken(db *shared.DBInfo) error {

	success := true

	if login.responseHTTPCode != 200 {
		success = false
	}

	loginItem := shared.DBLoginItem{
		Email:     login.email,
		LoginDate: int(login.creationTime.UnixMilli()),
		Success:   success,
	}

	item, err := attributevalue.MarshalMap(loginItem)

	if err != nil {
		return err
	}

	_, err = db.Client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(db.Tablename),
		Item:      item,
	})

	if err != nil {
		return err
	}

	return nil
}
