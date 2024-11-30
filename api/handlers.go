package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jkmancuso/photography_api/shared"
)

func jobs(request events.APIGatewayProxyRequest, db *shared.DBInfo) (string, int, error) {
	var err error
	var jobItem *shared.DBJobItem
	var jobItems []*shared.DBJobItem

	returnBytes := []byte(genericOK)

	ctx := context.Background()

	if request.HTTPMethod == "GET" {
		jobItems, err = getJobs(ctx, db)

		if err != nil {
			return genericError, http.StatusInternalServerError, err
		}

		returnBytes, err = json.Marshal(jobItems)
	} else if request.HTTPMethod == "POST" {
		jobItem, err = shared.ParseBodyIntoNewJob(request.Body)

		if err != nil {
			return genericError, http.StatusInternalServerError, err
		}

		err = addJob(ctx, db, jobItem)
	} else {
		err = fmt.Errorf("HTTP method %v not handled", request.HTTPMethod)
	}

	if err != nil {
		return genericError, http.StatusInternalServerError, err
	}

	return string(returnBytes), http.StatusOK, nil
}

func addJob(ctx context.Context, db *shared.DBInfo, job *shared.DBJobItem) error {

	item, err := attributevalue.MarshalMap(job)

	if err != nil {
		return err
	}

	err = db.AddItem(ctx, item)
	return err
}

func getJobs(ctx context.Context, db *shared.DBInfo) ([]*shared.DBJobItem, error) {

	var lek map[string]types.AttributeValue
	var items []*shared.DBJobItem

	//add max just in case of inifinte loop, "should break" before then
	for i := 0; i < MAX_LOOP; i++ {

		jobPage := []*shared.DBJobItem{}

		resp, err := db.DoFullScan(ctx, MAX_DB_ITEMS, lek)

		if err != nil {
			return items, err
		}

		err = attributevalue.UnmarshalListOfMaps(resp.Items, &jobPage)

		if err != nil {
			return items, err
		}

		items = append(items, jobPage...)

		lek = resp.LastEvaluatedKey

		if len(lek) == 0 {
			break
		}
	}

	return items, nil
}
