package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var statusCode = 200

func main() {
	lambda.Start(handler)

}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Println("Entering handler")

	routes := map[string]func() (string, error){
		"/auth":      auth,
		"/auth/ping": pong,
	}

	returnBody, err := routes[request.Path]()

	if err != nil {
		statusCode = 500
	}

	return events.APIGatewayProxyResponse{
		Body:       returnBody,
		StatusCode: statusCode,
	}, nil
}

func auth() (string, error) {
	ctx := context.Background()

	log.Println("Entering auth")

	//table := "admins"

	/*dynamoOutput, err := client.Scan(ctx, &dynamodb.ScanInput{
		TableName: &table,
	})*/

	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		log.Println(err)
		return "", err
	}

	client := *dynamodb.NewFromConfig(cfg)

	log.Printf("Client created success")

	dynamoOutput, err := client.ListTables(ctx, &dynamodb.ListTablesInput{})

	if err != nil {
		log.Println(err)
		return "", err
	}

	log.Printf("Got query results %+v", dynamoOutput.TableNames)

	/*for item := range dynamoOutput.Items {
		log.Printf("GOT: %+v\n", item)
	}*/
	return "Got to auth endpoint", nil
}

func pong() (string, error) {
	return "pong", nil
}
