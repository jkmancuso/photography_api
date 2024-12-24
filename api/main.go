package main

import (
	"context"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/jkmancuso/photography_api/api/handlers"
	"github.com/jkmancuso/photography_api/shared"
)

/* These global variables initialized below should already
be available on lambda warm start saving on startup resource*/

var (
	httpLambda *httpadapter.HandlerAdapter
	tables     = []string{"jobs", "groups", "instruments", "orders", "zipcodes"}
	tableMap   = map[string]*shared.DBInfo{}

	awsCfg aws.Config
)

func init() {

	if loglevel, found := os.LookupEnv("LOGLEVEL"); found {
		level, _ := log.ParseLevel(loglevel)
		log.SetLevel(level)
	}

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

			db.ConsistentRead = true

			tableMap[name] = db
		}
	}

	mux := http.DefaultServeMux
	handlers.SetupRoutes(tableMap, mux)

	httpLambda = httpadapter.New(mux)

}
func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Println("Entering handler")

	return httpLambda.ProxyWithContext(ctx, req)
}
