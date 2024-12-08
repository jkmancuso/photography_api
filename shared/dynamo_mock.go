package shared

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoClientInterface interface {
	Query(context.Context, *dynamodb.QueryInput, ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	Scan(context.Context, *dynamodb.ScanInput, ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	PutItem(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	DeleteItem(context.Context, *dynamodb.DeleteItemInput, ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
}

type DynamoClientMock struct {
	MockedRow map[string]types.AttributeValue
}

func (client DynamoClientMock) Query(ctx context.Context, input *dynamodb.QueryInput, opts ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	return &dynamodb.QueryOutput{}, nil
}

func (client DynamoClientMock) Scan(ctx context.Context, input *dynamodb.ScanInput, opts ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	return &dynamodb.ScanOutput{
		Items: []map[string]types.AttributeValue{client.MockedRow},
	}, nil
}

func (client DynamoClientMock) PutItem(ctx context.Context, input *dynamodb.PutItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return &dynamodb.PutItemOutput{
		Attributes: client.MockedRow,
	}, nil
}

func (client DynamoClientMock) GetItem(ctx context.Context, input *dynamodb.GetItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	var actualMockVal, actualInputVal string

	//loop through the mock row
	for mockKey, mockVal := range client.MockedRow {

		if err := attributevalue.Unmarshal(mockVal, &actualMockVal); err != nil {
			log.Fatal(err)
		}

		// if the key exists (ie input["id"]) in both the mock and the input
		if _, ok := input.Key[mockKey]; ok {

			if err := attributevalue.Unmarshal(input.Key[mockKey], &actualInputVal); err != nil {
				log.Fatal(err)
			}

			// AND the value is found
			if actualMockVal == actualInputVal {
				//then return the mock row
				return &dynamodb.GetItemOutput{
					Item: client.MockedRow,
				}, nil
			}
		}
	}
	//else return empty row
	return &dynamodb.GetItemOutput{
		Item: map[string]types.AttributeValue{},
	}, nil
}

func (client DynamoClientMock) DeleteItem(ctx context.Context, input *dynamodb.DeleteItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	return &dynamodb.DeleteItemOutput{}, nil
}
