package shared

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DBInfo struct {
	Tablename string
	Cfg       aws.Config
	Client    *dynamodb.Client
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
	log.Println("entering NewDB")
	db := &DBInfo{Tablename: table}

	client := dynamodb.NewFromConfig(cfg)

	log.Printf("Client created success")

	db.Cfg = cfg
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

func ParseBodyIntoNewJob(body string) (*DBJobItem, error) {
	jobItem := NewJobItem()
	err := json.Unmarshal([]byte(body), jobItem)

	if len(jobItem.JobName) == 0 || jobItem.JobYear == 0 {
		err = errors.New("missing field in body")
	}

	log.Println(jobItem)

	return jobItem, err
}
