package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var table = "admins"
var secretName = "salt"
var genericErrorJSON = `{"STATUS":"ERROR"}`
var cookieAge = 86400

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

	returnBody, statusCode, err := routes[request.Path](request)

	headers := make(map[string]string)

	if err != nil {
		log.Println(err)
	} else {
		token := dbItem{}

		if err = json.Unmarshal([]byte(returnBody), &token); err == nil {
			headers["Set-Cookie"] = fmt.Sprintf("token=%q; max-age=%d", token.Token, cookieAge)
		}
	}

	return events.APIGatewayProxyResponse{
		Body:       returnBody,
		StatusCode: statusCode,
		Headers:    headers,
	}, nil
}

func auth(request events.APIGatewayProxyRequest) (string, int, error) {

	log.Println("Entering auth")

	cfg, err := NewAWSCfg()

	if err != nil {
		return genericErrorJSON, http.StatusInternalServerError, err
	}

	saltStr, err := GetSalt(cfg)

	if err != nil || len(saltStr) == 0 {
		log.Println(err)
		return genericErrorJSON, http.StatusInternalServerError, err
	}

	login, err := NewLogin(request, saltStr)

	if err != nil {
		return login.responseHTTPMsg, login.responseHTTPCode, err
	}

	db, err := NewDB(table, cfg)

	if err != nil {
		log.Println(err)
		return genericErrorJSON, http.StatusInternalServerError, err
	}

	token, err := db.getToken(login)

	if err != nil {
		return genericErrorJSON, http.StatusInternalServerError, err
	}

	return token, login.responseHTTPCode, nil

}

func pong(request events.APIGatewayProxyRequest) (string, int, error) {
	return "pong", http.StatusOK, nil
}
