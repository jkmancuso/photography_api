package database

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jkmancuso/photography_api/shared"
)

func GetZipById(ctx context.Context, db *shared.DBInfo, pKey map[string]types.AttributeValue) (*shared.DBZipItem, int, error) {

	zipItem := &shared.DBZipItem{}

	resp, err := db.GetItem(ctx, pKey)

	if err != nil {
		return zipItem, 0, err
	}

	if err = attributevalue.UnmarshalMap(resp.Item, zipItem); err != nil {
		return zipItem, 0, err
	}

	return zipItem, len(resp.Item), nil
}
