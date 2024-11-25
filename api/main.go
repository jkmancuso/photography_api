package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jkmancuso/photography_api/shared"
)

var (
	/*create all Dynamo connections as Global variables so
	they are re used during Lambda warm start
	https://docs.aws.amazon.com/lambda/latest/dg/static-initialization.html
	*/

	tablesMap  = make(map[string]*shared.DBInfo)
	tableNames = []string{"jobs", "groups", "pictures", "instruments", "orders"}

	// same for aws config
	awsCfg aws.Config

	//genericError = `{"STATUS":"ERROR"}`
)

func init() {
	var err error

	awsCfg, err = shared.NewAWSCfg()

	if err != nil {
		log.Fatal(err)
	}

	for _, name := range tableNames {
		db, err := shared.NewDB(name, awsCfg)
		tablesMap[name] = db

		if err != nil {
			log.Fatal(err)
		}
	}

}
func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Println("Entering handler")

	/*	routes := map[string]AuthFunc{
			"/auth":      auth,
			"/auth/ping": pong,
		}

		returnBody, statusCode, _ := routes[request.Path](request)
	*/

	return events.APIGatewayProxyResponse{
		Body:       `{"hi":"there"}`,
		StatusCode: 200,
	}, nil
}
