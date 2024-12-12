package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
	"time"
)

type JobIntegrationTest struct {
	Job   *DBJobItem
	Url   string
	Tests []GenericTest
}

func (i *JobIntegrationTest) jobSetup(t *testing.T) {
	t.Helper()

	job := &DBJobItem{
		Id:       GenerateUUID(),
		JobName:  "integrationtest_job",
		JobYear:  time.Now().Year(),
		ExpireAt: time.Now().Unix() + ExpireIn,
	}

	body, err := json.Marshal(job)

	if err != nil {
		t.Fatal(err)
	}

	i.Tests = []GenericTest{
		{
			Name:           "valid job",
			BodyBytes:      body,
			WantStatusCode: 200,
		},
	}

	i.Url = fmt.Sprintf("%s/%s", API_URL, "jobs")
}

func TestE2EJob(t *testing.T) {

	idsToCheck := []string{}
	returnedJob := &DBJobItem{}

	test := &JobIntegrationTest{}
	test.jobSetup(t)

	for _, tt := range test.Tests {
		t.Run(tt.Name, func(t *testing.T) {
			reader := bytes.NewReader(tt.BodyBytes)

			resp, err := http.Post(test.Url, ContentType, reader)

			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != tt.WantStatusCode {
				t.Errorf("Got status %d, wanted %d", resp.StatusCode, tt.WantStatusCode)
			}

			if resp.StatusCode == http.StatusOK {
				responseBody, err := io.ReadAll(resp.Body)

				resp.Body.Close()

				if err != nil {
					t.Fatal(err)
				}

				if err = json.Unmarshal(responseBody, returnedJob); err != nil {
					t.Fatal(err)
				}

				if len(returnedJob.Id) == 0 {
					t.Fatal("Returned Job is empty")
				}

				idsToCheck = append(idsToCheck, returnedJob.Id)

			}

		})
	}

	for _, jobId := range idsToCheck {
		testName := fmt.Sprintf("%s/%s", "GetJobById", jobId)

		t.Run(testName, func(t *testing.T) {
			URL := fmt.Sprintf("%s/%s", test.Url, jobId)

			resp, err := http.Get(URL)

			if err != nil {
				t.Fatal(err)
			}

			responseBody, err := io.ReadAll(resp.Body)

			resp.Body.Close()

			if err != nil || resp.StatusCode != http.StatusOK {
				t.Fatalf("Error getting job by id")
			}

			log.Println(string(responseBody))

		})
	}

}
