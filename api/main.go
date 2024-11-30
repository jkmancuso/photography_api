package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jkmancuso/photography_api/shared"
)

type handlerFunc func(events.APIGatewayProxyRequest, *shared.DBInfo) (string, int, error)

const MAX_DB_ITEMS = 100
const MAX_LOOP = 100

var (
	endpointHandlers = map[string]handlerFunc{
		"/jobs": jobs,
		/*	"/groups":      groups,
			"/pictures":    pictures,
			"/instruments": instruments,
			"/orders":      orders*/}

	// same for aws config
	awsCfg aws.Config

	genericError = `{"STATUS":"ERROR"}`
)

func init() {
	var err error

	awsCfg, err = shared.NewAWSCfg()

	if err != nil {
		log.Fatal(err)
	}

}
func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Println("Entering handler")

	response, err := routeRequestToHandler(request)

	if err != nil {
		log.Print(err)
	}

	return response, err
}

func routeRequestToHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       genericError,
	}

	// has a leading slash
	endpoint := shared.GetTargetEndpoint(request.Path)

	if len(endpoint) == 0 {
		return response, fmt.Errorf("invalid api path %v", request.Path)
	}

	if _, ok := endpointHandlers[endpoint]; !ok {
		return response, fmt.Errorf("no handler for api path %v and endpoint %v",
			request.Path, endpoint)
	}

	if err := shared.ValidateEvent(request); err != nil {
		return response, err
	}

	// table name is endpoint without the leading slash
	db, err := shared.NewDB(endpoint[1:], awsCfg)

	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		return response, err
	}

	returnStr, statusCode, err := endpointHandlers[endpoint](request, db)

	response.Body = returnStr
	response.StatusCode = statusCode

	return response, err

}
