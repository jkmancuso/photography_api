package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

	var lek map[string]types.AttributeValue
	var jobItems []*shared.DBJobItem

	//add max just in case of inifinte loop, "should break" before then
	for i := 0; i < MAX_LOOP; i++ {

		jobPage := []*shared.DBJobItem{}

		resp, err := db.DoFullScan(ctx, MAX_DB_ITEMS, lek)

		if err != nil {
			return genericError, err
		}

		err = attributevalue.UnmarshalListOfMaps(resp.Items, &jobPage)

		if err != nil {
			return genericError, err
		}

		jobItems = append(jobItems, jobPage...)

		lek = resp.LastEvaluatedKey

		if len(lek) == 0 {
			break
		}
	}

	jobsStr, err := json.Marshal(jobItems)

	return string(jobsStr), err
}
