package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jkmancuso/photography_api/shared"
)

func jobs(request events.APIGatewayProxyRequest, db *shared.DBInfo) (string, int, error) {
	var returnStr string
	var err error

	switch request.HTTPMethod {
	case "GET":
		returnStr, err = getJobs(context.Background(), db)
	}

	return returnStr, http.StatusOK, err
}

func getJobs(ctx context.Context, db *shared.DBInfo) (string, error) {

	var jobPage []*shared.DBJobItem
	var jobItems []*shared.DBJobItem
	var err error

	//lek := map[string]string{}

	scanInput := dynamodb.ScanInput{
		TableName: &db.Tablename,
		//Limit:     aws.Int32(1),
	}

	scanPaginator := dynamodb.NewScanPaginator(db.Client, &scanInput)

	for scanPaginator.HasMorePages() {
		resp, err := scanPaginator.NextPage(ctx)

		if err != nil {
			return genericError, err
		}

		//_ = attributevalue.UnmarshalMap(resp.LastEvaluatedKey, &lek)

		//log.Printf("KEY : %+v", lek)

		err = attributevalue.UnmarshalListOfMaps(resp.Items, &jobPage)

		if err != nil {
			return genericError, err
		}

		jobItems = append(jobItems, jobPage...)
	}

	jobsStr, err := json.Marshal(jobItems)

	return string(jobsStr), err
}
