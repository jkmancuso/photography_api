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
	tests                  []shared.GenericTest
	mock                   shared.DynamoClientMock
	id, job_name, job_year types.AttributeValue
	jobsDBConn             handlerDBConn
)

var mux = *http.NewServeMux()

func setupMock() {

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

func TestGetJobs(t *testing.T) {

	tests = setupGetJobsTest()
	setupMock()

	mux.HandleFunc("/jobs", jobsDBConn.getJobsHandler)

	for _, tt := range tests {

		t.Run(tt.Name, func(t *testing.T) {
			respRecorder := httptest.NewRecorder()

			req, err := http.NewRequest("GET", "/jobs", nil)

			if err != nil {
				t.Fatal(err)
			}

			mux.ServeHTTP(respRecorder, req)

			if tt.WantStatusCode != respRecorder.Code {
				t.Errorf("%s: Got wrong response code %d, wanted %d",
					tt.Id, respRecorder.Code, tt.WantStatusCode)
			}

			if respRecorder.Result().StatusCode != http.StatusOK {
				t.Errorf("Got error %v", respRecorder.Result().StatusCode)
			}
		})
	}

}

func TestGetJobById(t *testing.T) {

	tests = setupGetJobByIdTest()
	setupMock()

	mux.HandleFunc("/jobs/{id}", jobsDBConn.getJobsByIdHandler)

	for _, tt := range tests {

		t.Run(tt.Name, func(t *testing.T) {
			respRecorder := httptest.NewRecorder()

			req, err := http.NewRequest("GET", fmt.Sprintf("/jobs/%s", tt.Id), nil)

			if err != nil {
				t.Fatal(err)
			}

			mux.ServeHTTP(respRecorder, req)

			if tt.WantStatusCode != respRecorder.Code {
				t.Errorf("%s: Got wrong response code %d, wanted %d",
					tt.Id, respRecorder.Code, tt.WantStatusCode)
			}

			if respRecorder.Result().StatusCode != http.StatusOK {
				gotMsg := shared.GenericMsg{}
				_ = json.Unmarshal(respRecorder.Body.Bytes(), &gotMsg)

				if gotMsg.Message != tt.WantErrorMsg.Message {
					t.Errorf("%s: Got wrong err msg %s, wanted %s",
						tt.Id, respRecorder.Body.String(), tt.WantErrorMsg.Message)
				}
			}
		})
	}

}

func setupGetJobByIdTest() []shared.GenericTest {

	return []shared.GenericTest{
		{
			Name:           "check valid id",
			Id:             VALID_ID,
			WantStatusCode: 200,
			WantErrorMsg:   shared.NO_ERR,
		},
		{
			Name:           "check missing id",
			Id:             MISSING_ID,
			WantStatusCode: 400,
			WantErrorMsg:   shared.ID_NOT_FOUND,
		},
		{
			Name:           "check invalid uuid",
			Id:             INVALID_UUID,
			WantStatusCode: 400,
			WantErrorMsg:   shared.ID_NOT_IN_UUID_FORMAT,
		},
	}

}

func setupGetJobsTest() []shared.GenericTest {

	return []shared.GenericTest{
		{
			Name:           "check result is returned",
			Id:             VALID_ID,
			WantStatusCode: 200,
			WantErrorMsg:   shared.NO_ERR,
		},
	}
}
