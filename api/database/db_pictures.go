package database

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jkmancuso/photography_api/shared"
)

func AddPicture(ctx context.Context, db *shared.DBInfo, picture *shared.DBPictureItem) error {

	item, err := attributevalue.MarshalMap(picture)

	if err != nil {
		return err
	}

	err = db.AddItem(ctx, item)
	return err
}

func DeletePicture(ctx context.Context, db *shared.DBInfo, id string) (int, error) {

	idAttr, err := attributevalue.Marshal(id)

	if err != nil {
		return 0, err
	}

	key := map[string]types.AttributeValue{"id": idAttr}

	count, err := db.DeleteItem(ctx, key)
	return count, err
}

func GetPictures(ctx context.Context, db *shared.DBInfo) ([]*shared.DBPictureItem, int, error) {

	var lek map[string]types.AttributeValue
	var items []*shared.DBPictureItem

	const MAX_DB_ITEMS = 200
	const MAX_LOOP = 200

	//add max just in case of inifinte loop, "should break" before then
	for i := 0; i < MAX_LOOP; i++ {

		picturePage := []*shared.DBPictureItem{}

		resp, err := db.DoFullScan(ctx, MAX_DB_ITEMS, lek)

		if err != nil {
			return items, 0, err
		}

		err = attributevalue.UnmarshalListOfMaps(resp.Items, &picturePage)

		if err != nil {
			return items, 0, err
		}

		items = append(items, picturePage...)

		lek = resp.LastEvaluatedKey

		if len(lek) == 0 {
			break
		}
	}

	return items, len(items), nil
}

func GetPictureById(ctx context.Context, db *shared.DBInfo, pKey map[string]types.AttributeValue) (*shared.DBPictureItem, int, error) {

	pictureItem := &shared.DBPictureItem{}

	resp, err := db.GetItem(ctx, pKey)

	if err != nil {
		return pictureItem, 0, err
	}

	if err = attributevalue.UnmarshalMap(resp.Item, pictureItem); err != nil {
		return pictureItem, 0, err
	}

	return pictureItem, len(resp.Item), nil
}
