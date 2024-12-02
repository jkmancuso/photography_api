package shared

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoClientInterface interface {
	Scan(context.Context, *dynamodb.ScanInput, ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	PutItem(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	DeleteItem(context.Context, *dynamodb.DeleteItemInput, ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
}

type DBInfo struct {
	Tablename string
	Client    DynamoClientInterface
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

func (db DBInfo) GetItem(ctx context.Context, idStr string) (*dynamodb.GetItemOutput, error) {

	id, err := attributevalue.Marshal(idStr)

	if err != nil {
		return &dynamodb.GetItemOutput{}, err
	}

	key := map[string]types.AttributeValue{"id": id}

	resp, err := db.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &db.Tablename,
		Key:       key,
	})

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
