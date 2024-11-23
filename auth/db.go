package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

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

type dbItem struct {
	Token string `dynamodbav:"token"`
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
		login.setstatusCode(http.StatusBadRequest)
		return `{"STATUS":"NO_RESULTS"}`, nil
	}

	row := dbItem{}

	if err = attributevalue.UnmarshalMap(response.Item, &row); err != nil {
		return "", err
	}

	rowJSON, err := json.Marshal(row)

	if err != nil {
		return "", err
	}

	return string(rowJSON), err

}
