package shared

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DBInfo struct {
	Tablename      string
	Client         DynamoClientInterface
	ConsistentRead bool
}

func NewDB(table string, cfg aws.Config) (*DBInfo, error) {

	log.Debugf("Creating client for DynamoDB table: %s", table)

	db := &DBInfo{Tablename: table}
	client := dynamodb.NewFromConfig(cfg)
	db.Client = client

	return db, nil
}

func (db DBInfo) DoFullScan(ctx context.Context, limit int32, lek map[string]types.AttributeValue) (*dynamodb.ScanOutput, error) {

	log.Debugf("Full scan for %s", db.Tablename)
	log.Debugf("Start Key: %+v", lek)

	resp, err := db.Client.Scan(ctx, &dynamodb.ScanInput{
		TableName:         &db.Tablename,
		Limit:             aws.Int32(limit),
		ExclusiveStartKey: lek,
		ConsistentRead:    &db.ConsistentRead,
	})

	if err != nil {
		log.Println(err)
	}

	return resp, err

}

func (db DBInfo) AddItem(ctx context.Context, item map[string]types.AttributeValue) error {

	log.Debugf("PutItem for table: %s", db.Tablename)
	log.Debugf("Item is: %+v", item)

	_, err := db.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &db.Tablename,
		Item:      item,
	})

	if err != nil {
		log.Println(err)
	}

	return err
}

func (db DBInfo) GetItem(ctx context.Context, pKey map[string]types.AttributeValue) (*dynamodb.GetItemOutput, error) {

	log.Debugf("PutItem for table: %s", db.Tablename)
	log.Debugf("Key is: %+v", pKey)

	input := &dynamodb.GetItemInput{
		TableName:      &db.Tablename,
		Key:            pKey,
		ConsistentRead: &db.ConsistentRead,
	}

	resp, err := db.Client.GetItem(ctx, input)

	if err != nil {
		log.Println(err)
	}

	return resp, err
}

func (db DBInfo) QueryItem(ctx context.Context, keys map[string]expression.ValueBuilder, gsi string) (*dynamodb.QueryOutput, error) {

	log.Debugf("QueryItem for table: %s", db.Tablename)
	log.Debugf("Key is: %+v", keys)

	if len(keys) != 1 && len(keys) != 2 {
		return &dynamodb.QueryOutput{}, errors.New("unsupported key condition")
	}

	exprBuilder := expression.NewBuilder()

	for k, v := range keys {
		exprBuilder = exprBuilder.WithKeyCondition(expression.Key(k).Equal(v))
	}

	expr, err := exprBuilder.Build()

	if err != nil {
		log.Printf("Couldn't build expression for query. Here's why: %v\n", err)
		return &dynamodb.QueryOutput{}, nil
	}

	input := &dynamodb.QueryInput{
		TableName:                 &db.Tablename,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	if len(gsi) != 0 {
		input.IndexName = &gsi
		log.Debugf("Using Global Secondary Index: %s", gsi)

	} else {
		//https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.ReadConsistency.html
		// Strongly consistent reads from a global secondary index is not supported.
		input.ConsistentRead = &db.ConsistentRead
	}

	resp, err := db.Client.Query(ctx, input)

	if err != nil {
		log.Println(err)
	}

	return resp, err
}

func (db DBInfo) DeleteItem(ctx context.Context, pKey map[string]types.AttributeValue) (int, error) {

	log.Debugf("DeleteItem for table: %s", db.Tablename)
	log.Debugf("Key is: %+v", pKey)

	resp, err := db.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName:    &db.Tablename,
		Key:          pKey,
		ReturnValues: types.ReturnValueAllOld,
	})

	if err != nil {
		log.Println(err)
	}

	return len(resp.Attributes), err
}

func (db DBInfo) UpdateItem(ctx context.Context, pKey map[string]types.AttributeValue, expr expression.Expression) (int, error) {

	log.Debugf("UpdateItem for table: %s", db.Tablename)
	log.Debugf("Key is: %+v", pKey)

	resp, err := db.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 &db.Tablename,
		Key:                       pKey,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              types.ReturnValueAllNew,
	})

	if err != nil {
		log.Println(err)
	}

	return len(resp.Attributes), err
}
