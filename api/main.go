package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/jkmancuso/photography_api/shared"
)

const MAX_DB_ITEMS = 100
const MAX_LOOP = 100

var httpLambda *httpadapter.HandlerAdapter

var (
	tables   = []string{"jobs", "groups", "pictures", "instruments", "orders"}
	tableMap = map[string]*shared.DBInfo{}

	awsCfg aws.Config
)

func init() {

	var err error

	if len(awsCfg.Region) == 0 {

		log.Println("loading new config from cold start")
		awsCfg, err = shared.NewAWSCfg()

		if err != nil {
			log.Fatal(err)
		}
	}

	if len(tableMap) == 0 {
		log.Println("loading new DB connections from cold start")

		for _, name := range tables {
			db, err := shared.NewDB(name, awsCfg)

			if err != nil {
				log.Fatal(err)
			}

			tableMap[name] = db
		}
	}

	setupRoutes()

	httpLambda = httpadapter.New(http.DefaultServeMux)

}
func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Println("Entering handler")

	return httpLambda.ProxyWithContext(ctx, req)
}
