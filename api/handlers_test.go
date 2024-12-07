package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jkmancuso/photography_api/shared"
)

// just random generate, nothing special
const VALID_ID = "26e9a998-2e16-44df-96e9-76520b6a22d8"
const INVALID_ID = "c0fd7513-39b4-404a-9b6b-26b00e8369ab"

var tests = []struct {
	id             string
	wantStatusCode int
}{
	{VALID_ID, 200},
	{INVALID_ID, 400},
}

func TestGetJobById(t *testing.T) {

	mux := http.NewServeMux()

	id, _ := attributevalue.Marshal(VALID_ID)
	job_name, _ := attributevalue.Marshal("mockedup_row")
	job_year, _ := attributevalue.Marshal(2025)

	mock := shared.DynamoClientMock{
		MockedRow: map[string]types.AttributeValue{
			"id":       id,
			"job_name": job_name,
			"job_year": job_year,
		},
	}

	jobsDBConn := handlerDBConn{
		dbInfo: &shared.DBInfo{
			Tablename: "jobs",
			Client:    mock,
		},
	}

	mux.HandleFunc("/jobs/{id}", jobsDBConn.getJobsByIdHandler)

	for _, tt := range tests {
		respRecorder := httptest.NewRecorder()

		req, err := http.NewRequest("GET", fmt.Sprintf("/jobs/%s", tt.id), nil)

		if err != nil {
			t.Fatal(err)
		}

		mux.ServeHTTP(respRecorder, req)

		if tt.wantStatusCode != respRecorder.Code {
			t.Errorf("%v: Got wrong response code %d, wanted %d",
				tt.id, respRecorder.Code, tt.wantStatusCode)
		}
	}

}
