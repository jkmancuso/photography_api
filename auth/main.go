package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var statusCode = 200
var table = "admins"
var genericErrorJSON = `{"STATUS":"ERROR"}`

type AuthFunc func(events.APIGatewayProxyRequest) (string, int, error)

func main() {
	lambda.Start(handler)

}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Println("Entering handler")

	routes := map[string]AuthFunc{
		"/auth":      auth,
		"/auth/ping": pong,
	}

	var returnBody string
	var err error

	returnBody, statusCode, err = routes[request.Path](request)

	returnBody = fmt.Sprintf("%v\n%+v", returnBody, request)

	return events.APIGatewayProxyResponse{
		Body:       returnBody,
		StatusCode: statusCode,
	}, err
}

func auth(request events.APIGatewayProxyRequest) (string, int, error) {

	log.Println("Entering auth")

	login, err := NewLogin(request)

	if err != nil {
		return login.responseHTTPMsg, login.responseHTTPCode, err
	}

	db, err := NewDB(table)

	if err != nil {
		log.Println(err)
		return genericErrorJSON, http.StatusInternalServerError, err
	}

	token, err := db.getToken(login)

	if err != nil {
		return genericErrorJSON, http.StatusInternalServerError, err
	}

	return token, http.StatusOK, nil

}

func pong(request events.APIGatewayProxyRequest) (string, int, error) {
	return "pong", http.StatusOK, nil
}
