package shared

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DBInfo struct {
	Tablename string
	Client    DynamoClientInterface
}

type DBAdminItem struct {
	Email    string `dynamodbav:"email" json:"email"`
	Hashpass string `dynamodbav:"hashpass" json:"hashpass"`
	Token    string `dynamodbav:"token" json:"Token"`
}

type DBLoginItem struct {
	Email     string `dynamodbav:"email"`
	LoginDate int    `dynamodbav:"login_date"`
	Success   bool   `dynamodbav:"success"`
}

type DBJobItem struct {
	Id       string `dynamodbav:"id" json:"id"`
	JobName  string `dynamodbav:"job_name" json:"job_name"`
	JobYear  int    `dynamodbav:"job_year" json:"job_year"`
	ExpireAt int64  `dynamodbav:"expire_at,omitempty" json:"expire_at,omitempty"` //DynamoDB expire TTL
}

type DBOrderItem struct {
	Id        string `dynamodbav:"id" json:"id"`
	JobId     string `dynamodbav:"job_id" json:"job_id"`
	RecordNum int    `dynamodbav:"record_num" json:"record_num"`
	ExpireAt  int64  `dynamodbav:"expire_at,omitempty" json:"expire_at,omitempty"` //DynamoDB expire TTL
}

func NewDB(table string, cfg aws.Config) (*DBInfo, error) {

	db := &DBInfo{Tablename: table}
	client := dynamodb.NewFromConfig(cfg)
	db.Client = client

	return db, nil
}

func (db DBInfo) DoFullScan(ctx context.Context, limit int32, lek map[string]types.AttributeValue) (*dynamodb.ScanOutput, error) {

	resp, err := db.Client.Scan(ctx, &dynamodb.ScanInput{
		TableName:         &db.Tablename,
		Limit:             aws.Int32(limit),
		ExclusiveStartKey: lek,
	})

	if err != nil {
		log.Println(err)
	}

	return resp, err

}

func (db DBInfo) AddItem(ctx context.Context, item map[string]types.AttributeValue) error {

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

	input := &dynamodb.GetItemInput{
		TableName: &db.Tablename,
		Key:       pKey,
	}

	resp, err := db.Client.GetItem(ctx, input)

	if err != nil {
		log.Println(err)
	}

	return resp, err
}

func (db DBInfo) QueryItem(ctx context.Context, keys map[string]expression.ValueBuilder, gsi string) (*dynamodb.QueryOutput, error) {

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
	}

	input := &dynamodb.QueryInput{
		TableName:                 &db.Tablename,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	if len(gsi) != 0 {
		input.IndexName = &gsi
	}

	resp, err := db.Client.Query(ctx, input)

	if err != nil {
		log.Println(err)
	}

	return resp, err
}

func (db DBInfo) DeleteItem(ctx context.Context, pKey map[string]types.AttributeValue) (int, error) {

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
