package main

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jkmancuso/photography_api/shared"
)

func returnSessionForValidAuth(ctx context.Context, email string, hashedPassword string, db *shared.DBInfo) (*shared.Session, error) {
	adminItem := &shared.DBAdminItem{}

	sess := &shared.Session{}

	emailAttribute, err := attributevalue.Marshal(email)

	if err != nil {
		return sess, err
	}

	pKey := map[string]types.AttributeValue{
		"email": emailAttribute,
	}
	resp, err := db.GetItem(ctx, pKey)

	if err != nil {
		return sess, err
	}

	if err = attributevalue.UnmarshalMap(resp.Item, adminItem); err != nil {
		return sess, err
	}

	if adminItem.Hashpass != hashedPassword {
		return sess, nil
	}

	sess.ExpireAt = time.Now().Add(720 * time.Hour)
	sess.SessionId = adminItem.Token

	return sess, nil
}

func addLogin(ctx context.Context, db *shared.DBInfo, login *shared.DBLoginItem) error {

	item, err := attributevalue.MarshalMap(login)

	if err != nil {
		return err
	}

	err = db.AddItem(ctx, item)
	return err
}

func updateLogin(ctx context.Context, db *shared.DBInfo, login *shared.DBLoginItem) (int, error) {

	emailAttr, err := attributevalue.Marshal(login.Email)

	if err != nil {
		return 0, err
	}

	logindateAttr, err := attributevalue.Marshal(login.LoginDate)

	if err != nil {
		return 0, err
	}

	pKey := map[string]types.AttributeValue{
		"email":      emailAttr,
		"login_date": logindateAttr,
	}

	update := expression.Set(expression.Name("success"), expression.Value(true))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()

	if err != nil {
		return 0, err
	}

	count, err := db.UpdateItem(ctx, pKey, expr)

	return count, err
}
