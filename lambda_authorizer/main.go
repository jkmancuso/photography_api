package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayCustomAuthorizerResponse, error) {

	log.Printf("Request: %+v", req)
	log.Printf("RequestContext: %+v", req.RequestContext)

	//to fix
	principalID := "testuser"
	effect := "Deny"

	resource := fmt.Sprintf("arn:aws:execute-api:%s:%s:%s/%s/*",
		strings.Split(req.RequestContext.DomainName, ".")[2],
		req.RequestContext.AccountID,
		req.RequestContext.APIID,
		req.RequestContext.Stage)

	inDB, err := sessionExistsInDB(context.Background(), req.Headers["x-session-id"])

	if err != nil {
		log.Printf("Error querying session DB: %v", err)
	}

	if inDB {
		effect = "Allow"
	}

	resp := NewAuthorizerResponse(principalID, resource, effect)

	b, _ := json.Marshal(resp)

	log.Printf("Retuning document: %+v", string(b))
	return resp, nil
}

func main() {

	if loglevel, found := os.LookupEnv("LOGLEVEL"); found {
		level, _ := log.ParseLevel(loglevel)
		log.SetLevel(level)
	}

	lambda.Start(handler)

}

func NewAuthorizerResponse(principalID string, resource string, effect string) events.APIGatewayCustomAuthorizerResponse {

	policyDoc := events.APIGatewayCustomAuthorizerPolicy{
		Version: "2012-10-17",
		Statement: []events.IAMPolicyStatement{
			{
				Action:   []string{"execute-api:Invoke"},
				Effect:   effect,
				Resource: []string{resource},
			},
		},
	}

	response := events.APIGatewayCustomAuthorizerResponse{
		PrincipalID:    principalID,
		PolicyDocument: policyDoc,
	}

	fmt.Printf("RESPONSE: %+v", response)

	return response
}
