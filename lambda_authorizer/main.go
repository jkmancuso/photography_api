package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {

	//to fix
	principalID := "testuser"
	awsAccountID := "084375571753"
	APIID := "ygaqa1m2xf"
	stage := "v1"
	region := "us-east-2"

	resource := fmt.Sprintf("arn:aws:execute-api:{%s}:{%s}:{%s}/{%s}/*]",
		region,
		awsAccountID,
		APIID,
		stage)

	fmt.Printf("Resource Arn: %s", resource)

	resp := NewAuthorizerResponse(principalID, resource)

	return resp, nil
}

func main() {

	lambda.Start(handler)

}

func NewAuthorizerResponse(principalID string, resource string) events.APIGatewayCustomAuthorizerResponse {

	policyDoc := events.APIGatewayCustomAuthorizerPolicy{
		Version: "2012-10-17",
		Statement: []events.IAMPolicyStatement{
			{
				Action:   []string{"execute-api:Invoke"},
				Effect:   "Allow",
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
