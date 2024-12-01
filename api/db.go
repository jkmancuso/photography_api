package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jkmancuso/photography_api/shared"
)

func addJob(ctx context.Context, db *shared.DBInfo, job *shared.DBJobItem) error {

	item, err := attributevalue.MarshalMap(job)

	if err != nil {
		return err
	}

	err = db.AddItem(ctx, item)
	return err
}

func getJobs(ctx context.Context, db *shared.DBInfo) ([]*shared.DBJobItem, int, error) {

	var lek map[string]types.AttributeValue
	var items []*shared.DBJobItem

	//add max just in case of inifinte loop, "should break" before then
	for i := 0; i < MAX_LOOP; i++ {

		jobPage := []*shared.DBJobItem{}

		resp, err := db.DoFullScan(ctx, MAX_DB_ITEMS, lek)

		if err != nil {
			return items, 0, err
		}

		err = attributevalue.UnmarshalListOfMaps(resp.Items, &jobPage)

		if err != nil {
			return items, 0, err
		}

		items = append(items, jobPage...)

		lek = resp.LastEvaluatedKey

		if len(lek) == 0 {
			break
		}
	}

	return items, len(items), nil
}

func getJobById(ctx context.Context, db *shared.DBInfo, id string) (*shared.DBJobItem, int, error) {

	jobItem := &shared.DBJobItem{}

	resp, err := db.GetItem(ctx, id)

	if err != nil {
		return jobItem, 0, err
	}

	if err = attributevalue.UnmarshalMap(resp.Item, jobItem); err != nil {
		return jobItem, 0, err
	}

	return jobItem, len(resp.Item), nil
}
