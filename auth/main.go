package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
)

var (
	tables                 = map[string]string{"admin": "admins", "login": "logins"}
	adminTable, loginTable *dbInfo
	secretName             = "salt"
	cookieAge              = 86400
	awsCfg                 aws.Config
	saltStr                string
	genericError           = `{"STATUS":"ERROR"}`
)

type AuthFunc func(events.APIGatewayProxyRequest) (string, int, error)

func main() {
	lambda.Start(handler)

}

func init() {

	var err error

	awsCfg, err = NewAWSCfg()

	if err != nil {
		log.Fatal(err)
	}

	adminTable, err = NewDB(tables["admin"], awsCfg)

	if err != nil {
		log.Fatal(err)
	}

	loginTable, err = NewDB(tables["login"], awsCfg)

	if err != nil {
		log.Fatal(err)
	}

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
		token := dbAdminItem{}

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

	var err error

	saltStr, err = GetSalt(awsCfg)

	if err != nil {
		return genericError, http.StatusInternalServerError, err
	}

	login, err := NewLogin(request, saltStr)

	if err != nil {
		return login.responseHTTPMsg, login.responseHTTPCode, err
	}

	token, err := adminTable.getToken(login)

	if err != nil {
		return token, http.StatusInternalServerError, err
	}

	if loginTable.recordLoginToken(login); err != nil {
		login.responseHTTPCode = http.StatusInternalServerError
		log.Println("error adding login record")
	}

	return token, login.responseHTTPCode, err

}

func pong(request events.APIGatewayProxyRequest) (string, int, error) {
	return "pong", http.StatusOK, nil
}
