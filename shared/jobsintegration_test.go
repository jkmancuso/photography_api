package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestAddJob(t *testing.T) {

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
				t.Fatalf("Got status %d, wanted %d", resp.StatusCode, tt.WantStatusCode)
			}

		})
	}

}
