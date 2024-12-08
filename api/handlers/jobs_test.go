package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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
	tests      []GenericTest
	mock       shared.DynamoClientMock
	jobsDBConn handlerDBConn
)

var mux = *http.NewServeMux()

func setupJobMock() {
	var id, job_name, job_year types.AttributeValue
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

func TestAddJob(t *testing.T) {

	tests = setupAddJobTest()
	setupJobMock()

	mux.HandleFunc("POST /jobs", jobsDBConn.addJobsHandler)

	for _, tt := range tests {

		t.Run(tt.Name, func(t *testing.T) {
			respRecorder := httptest.NewRecorder()

			req, err := http.NewRequest("POST", "/jobs", strings.NewReader(tt.Body))

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

func TestGetJobs(t *testing.T) {

	tests = setupGetJobsTest()
	setupJobMock()

	mux.HandleFunc("GET /jobs", jobsDBConn.getJobsHandler)

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

func TestGetJobById(t *testing.T) {

	tests = setupGetJobByIdTest()
	setupJobMock()

	mux.HandleFunc("GET /jobs/{id}", jobsDBConn.getJobsByIdHandler)

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

func setupGetJobByIdTest() []GenericTest {

	return []GenericTest{
		{
			Name:           "check valid id",
			Id:             VALID_ID,
			WantStatusCode: 200,
		},
		{
			Name:           "check missing id",
			Id:             MISSING_ID,
			WantStatusCode: 400,
			WantErrorMsg:   shared.RECORD_NOT_FOUND,
		},
		{
			Name:           "check invalid uuid",
			Id:             INVALID_UUID,
			WantStatusCode: 400,
			WantErrorMsg:   shared.ID_NOT_IN_UUID_FORMAT,
		},
	}

}

func setupGetJobsTest() []GenericTest {

	return []GenericTest{
		{
			Name:           "check result is returned",
			Id:             VALID_ID,
			WantStatusCode: 200,
		},
	}
}

func setupAddJobTest() []GenericTest {

	return []GenericTest{
		{
			Name:           "check valid body",
			Body:           `{"job_name":"test job", "job_year":2025}`,
			WantStatusCode: 200,
		},
		{
			Name:           "check empty body",
			Body:           ``,
			WantStatusCode: 400,
			WantErrorMsg:   shared.INVALID_BODY,
		},
	}
}
