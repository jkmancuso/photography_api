package database

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jkmancuso/photography_api/shared"
)

func DeleteOrder(ctx context.Context, db *shared.DBInfo, id string) (int, error) {

	idAttr, err := attributevalue.Marshal(id)

	if err != nil {
		return 0, err
	}

	key := map[string]types.AttributeValue{"id": idAttr}

	count, err := db.DeleteItem(ctx, key)
	return count, err
}

func AddOrder(ctx context.Context, db *shared.DBInfo, order *shared.DBOrderItem) error {

	item, err := attributevalue.MarshalMap(order)

	if err != nil {
		return err
	}

	err = db.AddItem(ctx, item)
	return err
}

// Global secondary index supports Query, not GetItem
func GetOrderByGSI(ctx context.Context, db *shared.DBInfo, keys map[string]expression.ValueBuilder, gsi string) (*shared.DBOrderItem, int, error) {

	orderItem := &shared.DBOrderItem{}
	orderItems := []shared.DBOrderItem{}

	resp, err := db.QueryItem(ctx, keys, gsi)

	if err != nil {
		return orderItem, 0, err
	}

	if err = attributevalue.UnmarshalListOfMaps(resp.Items, &orderItems); err != nil {
		return orderItem, 0, err
	}

	return &orderItems[0], len(resp.Items), nil
}

// Global secondary index supports Query, not GetItem
func GetOrdersByGSI(ctx context.Context, db *shared.DBInfo, keys map[string]expression.ValueBuilder, gsi string) ([]*shared.DBOrderItem, int, error) {

	orderItems := []*shared.DBOrderItem{}

	resp, err := db.QueryItem(ctx, keys, gsi)

	if err != nil {
		return orderItems, 0, err
	}

	if err = attributevalue.UnmarshalListOfMaps(resp.Items, &orderItems); err != nil {
		return orderItems, 0, err
	}

	return orderItems, len(resp.Items), nil
}

// GetItem
func GetOrderByPKey(ctx context.Context, db *shared.DBInfo, pKey map[string]types.AttributeValue) (*shared.DBOrderItem, int, error) {

	orderItem := &shared.DBOrderItem{}

	resp, err := db.GetItem(ctx, pKey)

	if err != nil {
		return orderItem, 0, err
	}

	if err = attributevalue.UnmarshalMap(resp.Item, &orderItem); err != nil {
		return orderItem, 0, err
	}

	return orderItem, len(resp.Item), nil
}

func GetOrders(ctx context.Context, db *shared.DBInfo) ([]*shared.DBOrderItem, int, error) {

	var lek map[string]types.AttributeValue
	var items []*shared.DBOrderItem

	/*
		/orders endpoint shouldnt be used in production, only E2E test
		orders table can grow to tens of thousands over time
	*/
	const MAX_DB_ITEMS = 10
	const MAX_LOOP = 10

	//add max just in case of inifinte loop, "should break" before then
	for i := 0; i < MAX_LOOP; i++ {

		orderPage := []*shared.DBOrderItem{}

		resp, err := db.DoFullScan(ctx, MAX_DB_ITEMS, lek)

		if err != nil {
			return items, 0, err
		}

		err = attributevalue.UnmarshalListOfMaps(resp.Items, &orderPage)

		if err != nil {
			return items, 0, err
		}

		items = append(items, orderPage...)

		lek = resp.LastEvaluatedKey

		if len(lek) == 0 {
			break
		}
	}

	return items, len(items), nil
}
