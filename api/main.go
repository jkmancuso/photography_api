package main

import (
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jkmancuso/photography_api/shared"
)

type handlerFunc func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

var (
	/*create all Dynamo connections as Global variables so
	they are re used during Lambda warm start
	https://docs.aws.amazon.com/lambda/latest/dg/static-initialization.html
	*/

	endpointHandlers = map[string]handlerFunc{
		"/jobs":        jobs,
		"/groups":      groups,
		"/pictures":    pictures,
		"/instruments": instruments,
		"/orders":      orders}

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

	return events.APIGatewayProxyResponse{
		Body:       response.Body,
		StatusCode: response.StatusCode,
	}, err
}

func routeRequestToHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{}

	endpoint := shared.GetTargetEndpoint(request.Path)

	if len(endpoint) == 0 {
		return response, errors.New("invalid api path")
	}

	if _, ok := endpointHandlers[endpoint]; !ok {
		return response, errors.New("no handler for this api path")
	}

	if err := shared.ValidateEvent(request); err != nil {
		return response, err
	}

	response, err := endpointHandlers[endpoint](request)

	return response, err

}
