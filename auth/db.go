package main

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jkmancuso/photography_api/shared"
)

func returnTokenForValidAuth(ctx context.Context, email string, hashedPassword string, db *shared.DBInfo) (string, error) {
	adminItem := &shared.DBAdminItem{}
	token := ""

	emailAttribute, err := attributevalue.Marshal(email)

	if err != nil {
		return token, err
	}

	pKey := map[string]types.AttributeValue{
		email: emailAttribute,
	}
	resp, err := db.GetItem(ctx, pKey)

	if err != nil {
		return token, err
	}

	if err = attributevalue.UnmarshalMap(resp.Item, adminItem); err != nil {
		return token, err
	}

	if adminItem.Hashpass != hashedPassword {
		return token, errors.New("invalid user/pass")
	}

	return adminItem.Token, nil
}
