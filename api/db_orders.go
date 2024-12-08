package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jkmancuso/photography_api/shared"
)

func addOrder(ctx context.Context, db *shared.DBInfo, order *shared.DBOrderItem) error {

	item, err := attributevalue.MarshalMap(order)

	if err != nil {
		return err
	}

	err = db.AddItem(ctx, item)
	return err
}

// Global secondary index supports Query, not GetItem
// id is your GSI
func getOrderByGSI(ctx context.Context, db *shared.DBInfo, k string, v string, gsi string) (*shared.DBOrderItem, int, error) {

	orderItem := &shared.DBOrderItem{}
	orderItems := []shared.DBOrderItem{}

	resp, err := db.QueryItem(ctx, k, v, gsi)

	if err != nil {
		return orderItem, 0, err
	}

	if err = attributevalue.UnmarshalListOfMaps(resp.Items, &orderItems); err != nil {
		return orderItem, 0, err
	}

	return &orderItems[0], 1, nil
}

// GetItem
// query for record_num and job_id
func getOrderByPKey(ctx context.Context, db *shared.DBInfo, pKey map[string]types.AttributeValue) (*shared.DBOrderItem, int, error) {

	orderItem := &shared.DBOrderItem{}

	resp, err := db.GetItem(ctx, pKey)

	if err != nil {
		return orderItem, 0, err
	}

	if err = attributevalue.UnmarshalMap(resp.Item, &orderItem); err != nil {
		return orderItem, 0, err
	}

	return orderItem, 1, nil
}
