package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type dbInfo struct {
	tablename string
	cfg       aws.Config
	client    *dynamodb.Client
}

type dbAdminItem struct {
	Token string `dynamodbav:"token" json:"Token"`
}

type dbLoginItem struct {
	Email     string `dynamodbav:"email"`
	LoginDate int    `dynamodbav:"login_date"`
	Success   bool   `dynamodbav:"success"`
}

func NewDB(table string, cfg aws.Config) (*dbInfo, error) {
	log.Println("entering NewDB")
	db := &dbInfo{tablename: table}

	client := dynamodb.NewFromConfig(cfg)

	log.Printf("Client created success")

	db.cfg = cfg
	db.client = client

	return db, nil
}

func (db *dbInfo) getToken(login *UserLogin) (string, error) {
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
		TableName: &db.tablename,
		Key:       key,
	}

	response, err := db.client.GetItem(context.Background(), &query)

	if err != nil {
		return "", err
	}

	if response.Item == nil {
		return `{"STATUS":"INVALID_USER_PASS"}`, errors.New("invalid user/pass")
	}

	adminItem := dbAdminItem{}

	if err = attributevalue.UnmarshalMap(response.Item, &adminItem); err != nil {
		return "", err
	}

	rowJSON, err := json.Marshal(adminItem)

	if err != nil {
		return "", err
	}

	return string(rowJSON), err

}

func (db *dbInfo) recordLoginToken(login *UserLogin) error {

	success := true

	if login.responseHTTPCode != 200 {
		success = false
	}

	loginItem := dbLoginItem{
		Email:     login.email,
		LoginDate: int(login.creationTime.UnixMilli()),
		Success:   success,
	}

	item, err := attributevalue.MarshalMap(loginItem)

	if err != nil {
		return err
	}

	_, err = db.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(db.tablename),
		Item:      item,
	})

	if err != nil {
		return err
	}

	return nil
}
