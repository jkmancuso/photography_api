package shared

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoClientInterface interface {
	Scan(context.Context, *dynamodb.ScanInput, ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	PutItem(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	DeleteItem(context.Context, *dynamodb.DeleteItemInput, ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
}

type DynamoClientMock struct {
	mockedRow map[string]types.AttributeValue
}

func (client *DynamoClientMock) Scan(ctx context.Context, input *dynamodb.ScanInput, opts ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	return &dynamodb.ScanOutput{}, nil
}

func (client *DynamoClientMock) PutItem(ctx context.Context, input *dynamodb.PutItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return &dynamodb.PutItemOutput{}, nil
}

func (client *DynamoClientMock) GetItem(ctx context.Context, input *dynamodb.GetItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	
	//loop through the mock row 
	for mockKey, mockVal := range client.mockedRow {
		// if the key exists (ie input["id"]) in both the mock and the input 
		if  _, ok := input[mockKey]; ok {
			// AND the value is found
			if mockVal == input[mockKey]{
				//then return the mock row
				return &dynamodb.GetItemOutput{
					Item: client.mockedRow
				}, nil
			}
		}
	} 
	//else return empty row
	return &dynamodb.GetItemOutput{
		Item: map[string]types.AttributeValue{}
	}, nil
}

func (client *DynamoClientMock) DeleteItem(ctx context.Context, input *dynamodb.DeleteItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	return &dynamodb.DeleteItemOutput{}, nil
}
