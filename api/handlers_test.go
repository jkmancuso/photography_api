package main

import (
	"encoding/json"
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
const MISSING_ID = "c0fd7513-39b4-404a-9b6b-26b00e8369ab"
const INVALID_UUID = "123456789"

var (
	tests = []struct {
		name           string
		id             string
		wantStatusCode int
		wantErrorMsg   shared.GenericMsg
	}{
		{"check valid id", VALID_ID, 200, shared.NO_ERR},
		{"check missing id", MISSING_ID, 400, shared.ID_NOT_FOUND},
		{"check invalid uuid", INVALID_UUID, 400, shared.ID_NOT_IN_UUID_FORMAT},
	}

	mux                    http.ServeMux
	mock                   shared.DynamoClientMock
	id, job_name, job_year types.AttributeValue
	jobsDBConn             handlerDBConn
)

func setupTest() {
	mux = *http.NewServeMux()

	id, _ = attributevalue.Marshal(VALID_ID)
	job_name, _ = attributevalue.Marshal("mockedup_row")
	job_year, _ = attributevalue.Marshal(2025)

	mock = shared.DynamoClientMock{
		MockedRow: map[string]types.AttributeValue{
			"id":       id,
			"job_name": job_name,
			"job_year": job_year,
		},
	}

	jobsDBConn = handlerDBConn{
		dbInfo: &shared.DBInfo{
			Tablename: "jobs",
			Client:    mock,
		},
	}

}
func TestGetJobById(t *testing.T) {

	setupTest()

	mux.HandleFunc("/jobs/{id}", jobsDBConn.getJobsByIdHandler)

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			respRecorder := httptest.NewRecorder()

			req, err := http.NewRequest("GET", fmt.Sprintf("/jobs/%s", tt.id), nil)

			if err != nil {
				t.Fatal(err)
			}

			mux.ServeHTTP(respRecorder, req)

			if tt.wantStatusCode != respRecorder.Code {
				t.Errorf("%s: Got wrong response code %d, wanted %d",
					tt.id, respRecorder.Code, tt.wantStatusCode)
			}

			if respRecorder.Result().StatusCode != http.StatusOK {
				gotMsg := shared.GenericMsg{}
				_ = json.Unmarshal(respRecorder.Body.Bytes(), &gotMsg)

				if gotMsg.Message != tt.wantErrorMsg.Message {
					t.Errorf("%s: Got wrong err msg %s, wanted %s",
						tt.id, respRecorder.Body.String(), tt.wantErrorMsg.Message)
				}
			}
		})
	}

}
