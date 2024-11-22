package shared

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

func NewDB(table string, cfg aws.Config) (*DBInfo, error) {
	log.Println("entering NewDB")
	db := &DBInfo{Tablename: table}

	client := dynamodb.NewFromConfig(cfg)

	log.Printf("Client created success")

	db.Cfg = cfg
	db.Client = client

	return db, nil
}
