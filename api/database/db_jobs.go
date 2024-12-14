package database

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jkmancuso/photography_api/shared"
)

func AddJob(ctx context.Context, db *shared.DBInfo, job *shared.DBJobItem) error {

	item, err := attributevalue.MarshalMap(job)

	if err != nil {
		return err
	}

	err = db.AddItem(ctx, item)
	return err
}

func DeleteJob(ctx context.Context, db *shared.DBInfo, id string) (int, error) {

	idAttr, err := attributevalue.Marshal(id)

	if err != nil {
		return 0, err
	}

	key := map[string]types.AttributeValue{"id": idAttr}

	count, err := db.DeleteItem(ctx, key)
	return count, err
}

func GetJobs(ctx context.Context, db *shared.DBInfo) ([]*shared.DBJobItem, int, error) {

	var lek map[string]types.AttributeValue
	var items []*shared.DBJobItem

	const MAX_DB_ITEMS = 200
	const MAX_LOOP = 200

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

func GetJobById(ctx context.Context, db *shared.DBInfo, pKey map[string]types.AttributeValue) (*shared.DBJobItem, int, error) {

	jobItem := &shared.DBJobItem{}

	resp, err := db.GetItem(ctx, pKey)

	if err != nil {
		return jobItem, 0, err
	}

	if err = attributevalue.UnmarshalMap(resp.Item, jobItem); err != nil {
		return jobItem, 0, err
	}

	return jobItem, len(resp.Item), nil
}
