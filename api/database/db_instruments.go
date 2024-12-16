package database

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jkmancuso/photography_api/shared"
)

func AddInstrument(ctx context.Context, db *shared.DBInfo, instrument *shared.DBInstrumentItem) error {

	item, err := attributevalue.MarshalMap(instrument)

	if err != nil {
		return err
	}

	err = db.AddItem(ctx, item)
	return err
}

func DeleteInstrument(ctx context.Context, db *shared.DBInfo, id string) (int, error) {

	idAttr, err := attributevalue.Marshal(id)

	if err != nil {
		return 0, err
	}

	key := map[string]types.AttributeValue{"id": idAttr}

	count, err := db.DeleteItem(ctx, key)
	return count, err
}

func GetInstruments(ctx context.Context, db *shared.DBInfo) ([]*shared.DBInstrumentItem, int, error) {

	var lek map[string]types.AttributeValue
	var items []*shared.DBInstrumentItem

	const MAX_DB_ITEMS = 200
	const MAX_LOOP = 200

	//add max just in case of inifinte loop, "should break" before then
	for i := 0; i < MAX_LOOP; i++ {

		instrumentPage := []*shared.DBInstrumentItem{}

		resp, err := db.DoFullScan(ctx, MAX_DB_ITEMS, lek)

		if err != nil {
			return items, 0, err
		}

		err = attributevalue.UnmarshalListOfMaps(resp.Items, &instrumentPage)

		if err != nil {
			return items, 0, err
		}

		items = append(items, instrumentPage...)

		lek = resp.LastEvaluatedKey

		if len(lek) == 0 {
			break
		}
	}

	return items, len(items), nil
}

func GetInstrumentById(ctx context.Context, db *shared.DBInfo, pKey map[string]types.AttributeValue) (*shared.DBInstrumentItem, int, error) {

	instrumentItem := &shared.DBInstrumentItem{}

	resp, err := db.GetItem(ctx, pKey)

	if err != nil {
		return instrumentItem, 0, err
	}

	if err = attributevalue.UnmarshalMap(resp.Item, instrumentItem); err != nil {
		return instrumentItem, 0, err
	}

	return instrumentItem, len(resp.Item), nil
}