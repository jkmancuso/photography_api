package database

import (
	"context"

	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jkmancuso/photography_api/shared"
)

func UpdateGroup(ctx context.Context, db *shared.DBInfo, group map[string]interface{}) (int, error) {

	idAttr, err := attributevalue.Marshal(group["id"])

	if err != nil {
		return 0, err
	}

	pKey := map[string]types.AttributeValue{"id": idAttr}

	delete(group, "id") //so you dont loop over it below

	var update expression.UpdateBuilder

	for k, v := range group {
		update = update.Set(expression.Name(k), expression.Value(v))
		fmt.Printf("KEY: %v  VAL: %v", k, v)
	}

	expr, err := expression.NewBuilder().WithUpdate(update).Build()

	if err != nil {
		log.Printf("Couldn't build expression for update. Here's why: %v\n", err)
		return 0, err
	}

	count, err := db.UpdateItem(ctx, pKey, expr)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func AddGroup(ctx context.Context, db *shared.DBInfo, group *shared.DBGroupItem) error {

	item, err := attributevalue.MarshalMap(group)

	if err != nil {
		return err
	}

	err = db.AddItem(ctx, item)
	return err
}

func DeleteGroup(ctx context.Context, db *shared.DBInfo, id string) (int, error) {

	idAttr, err := attributevalue.Marshal(id)

	if err != nil {
		return 0, err
	}

	key := map[string]types.AttributeValue{"id": idAttr}

	count, err := db.DeleteItem(ctx, key)
	return count, err
}

func GetGroups(ctx context.Context, db *shared.DBInfo) ([]*shared.DBGroupItem, int, error) {

	var lek map[string]types.AttributeValue
	var items []*shared.DBGroupItem

	const MAX_DB_ITEMS = 200
	const MAX_LOOP = 200

	//add max just in case of inifinte loop, "should break" before then
	for i := 0; i < MAX_LOOP; i++ {

		groupPage := []*shared.DBGroupItem{}

		resp, err := db.DoFullScan(ctx, MAX_DB_ITEMS, lek)

		if err != nil {
			return items, 0, err
		}

		err = attributevalue.UnmarshalListOfMaps(resp.Items, &groupPage)

		if err != nil {
			return items, 0, err
		}

		items = append(items, groupPage...)

		lek = resp.LastEvaluatedKey

		if len(lek) == 0 {
			break
		}
	}

	return items, len(items), nil
}

func GetGroupById(ctx context.Context, db *shared.DBInfo, pKey map[string]types.AttributeValue) (*shared.DBGroupItem, int, error) {

	groupItem := &shared.DBGroupItem{}

	resp, err := db.GetItem(ctx, pKey)

	if err != nil {
		return groupItem, 0, err
	}

	if err = attributevalue.UnmarshalMap(resp.Item, groupItem); err != nil {
		return groupItem, 0, err
	}

	return groupItem, len(resp.Item), nil
}
