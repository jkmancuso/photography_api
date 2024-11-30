package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func jobs(request events.APIGatewayProxyRequest, db *dynamodb.Client) (string, int, error) {
	var returnStr string
	var err error

	switch request.HTTPMethod {
	case "GET":
		returnStr, err = getJobs(context.Background(), db)
	}

	return returnStr, http.StatusOK, err
}

func getJobs(ctx context.Context, db *dynamodb.Client) (string, error) {
	return "", nil
}
