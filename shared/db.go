package shared

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DBInfo struct {
	Tablename string
	Client    DynamoClientInterface
	GSI       string //optinal global secondary index
}

type DBAdminItem struct {
	Token string `dynamodbav:"token" json:"Token"`
}

type DBLoginItem struct {
	Email     string `dynamodbav:"email"`
	LoginDate int    `dynamodbav:"login_date"`
	Success   bool   `dynamodbav:"success"`
}

type DBJobItem struct {
	Id      string `dynamodbav:"id" json:"id"`
	JobName string `dynamodbav:"job_name" json:"job_name"`
	JobYear int    `dynamodbav:"job_year" json:"job_year"`
}

type DBOrderItem struct {
	Id        string `dynamodbav:"id" json:"id"`
	JobId     string `dynamodbav:"job_id" json:"job_id"`
	RecordNum int    `dynamodbav:"record_num" json:"record_num"`
}

func NewDB(table string, cfg aws.Config) (*DBInfo, error) {

	db := &DBInfo{Tablename: table}
	client := dynamodb.NewFromConfig(cfg)
	db.Client = client

	return db, nil
}

func NewJobItem() *DBJobItem {
	return &DBJobItem{
		Id: GenerateUUID(),
	}
}

func (db DBInfo) DoFullScan(ctx context.Context, limit int32, lek map[string]types.AttributeValue) (*dynamodb.ScanOutput, error) {

	resp, err := db.Client.Scan(ctx, &dynamodb.ScanInput{
		TableName:         &db.Tablename,
		Limit:             aws.Int32(limit),
		ExclusiveStartKey: lek,
	})

	return resp, err

}

func (db DBInfo) AddItem(ctx context.Context, item map[string]types.AttributeValue) error {

	_, err := db.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &db.Tablename,
		Item:      item,
	})

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
		log.Printf("Table: %v", *input.TableName)
	}

	return resp, err
}

func (db DBInfo) QueryItem(ctx context.Context, k string, v string, gsi string) (*dynamodb.QueryOutput, error) {

	keyEx := expression.Key(k).Equal(expression.Value(v))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()

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
		log.Printf("Table: %v", *input.TableName)
	}

	return resp, err
}

func (db DBInfo) DeleteItem(ctx context.Context, idStr string) (int, error) {

	id, err := attributevalue.Marshal(idStr)

	if err != nil {
		return 0, err
	}

	key := map[string]types.AttributeValue{"id": id}

	resp, err := db.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &db.Tablename,
		Key:       key,
	})

	count := len(resp.Attributes)

	return count, err
}

func ParseBodyIntoNewJob(body []byte) (*DBJobItem, error) {
	jobItem := NewJobItem()
	err := json.Unmarshal(body, jobItem)

	if len(jobItem.JobName) == 0 || jobItem.JobYear == 0 {
		err = errors.New("missing field in body")
	}

	log.Println(jobItem)

	return jobItem, err
}
