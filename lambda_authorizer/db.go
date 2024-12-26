package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jkmancuso/photography_api/shared"
)

func sessionExistsInDB(ctx context.Context, id string) (bool, error) {

	if len(id) == 0 {
		return false, errors.New("missing session id")
	}

	log.Println(id)

	idAttribute, err := attributevalue.Marshal(id)

	if err != nil {
		return false, err
	}

	pKey := map[string]types.AttributeValue{
		"id": idAttribute,
	}

	awsCfg, err := shared.NewAWSCfg()

	if err != nil {
		return false, err
	}

	db, err := shared.NewDB("sessions", awsCfg)

	if err != nil {
		return false, err
	}

	resp, err := db.GetItem(ctx, pKey)

	if err != nil {
		return false, err
	}

	if len(resp.Item) == 0 {
		log.Printf("No session found for this user")
		return false, nil
	}

	sessionItem := &shared.DBSessionItem{}
	if err = attributevalue.UnmarshalMap(resp.Item, sessionItem); err != nil {
		return false, err
	}

	t := time.Now().Unix()

	if sessionItem.ExpireAt < t {
		log.Printf("Found expired session expired %d. \nTime right now is %d", sessionItem.ExpireAt, t)
		return false, nil
	}

	return true, nil
}
