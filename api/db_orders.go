package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jkmancuso/photography_api/shared"
)

func deleteOrder(ctx context.Context, db *shared.DBInfo, id string) (int, error) {

	count, err := db.DeleteItem(ctx, id)
	return count, err
}

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
func getOrderByGSI(ctx context.Context, db *shared.DBInfo, keys map[string]expression.ValueBuilder, gsi string) (*shared.DBOrderItem, int, error) {

	orderItem := &shared.DBOrderItem{}
	orderItems := []shared.DBOrderItem{}

	resp, err := db.QueryItem(ctx, keys, gsi)

	if len(resp.Items) == 0 {
		return orderItem, 0, nil
	}

	if err != nil {
		return orderItem, 0, err
	}

	if err = attributevalue.UnmarshalListOfMaps(resp.Items, &orderItems); err != nil {
		return orderItem, 0, err
	}

	return &orderItems[0], 1, nil
}

// GetItem
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
