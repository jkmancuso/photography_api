package shared

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
	"time"
)

type IntegrationTest struct {
	Url   string
	Tests []GenericTest
}

func (i *IntegrationTest) setup(t *testing.T) {
	t.Helper()

	testName := flag.String("name", "jobs", "[jobs, orders, etc]")
	flag.Parse()

	//populated item
	validPayload := NewDBItem(*testName)
	log.Println(string(validPayload))

	//empty item
	invalidPayload := []byte(``)

	i.Tests = []GenericTest{
		{
			Name:           "POST valid request body",
			BodyBytes:      validPayload,
			WantStatusCode: 200,
		},
		{
			Name:           "POST invalid request body",
			BodyBytes:      invalidPayload,
			WantStatusCode: 400,
		},
	}

	i.Url = fmt.Sprintf("%s/%s", API_URL, *testName)
}

func TestE2E(t *testing.T) {

	idsToCheck := []string{}

	//dont need a full struct, just check Id is there is fine
	returnedItem := &IdOnly{}

	test := &IntegrationTest{}
	test.setup(t)

	//step 1- add the job
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

				if err = json.Unmarshal(responseBody, returnedItem); err != nil {
					t.Fatal(err)
				}

				if len(returnedItem.Id) == 0 {
					t.Fatal("Returned Item is empty")
				}

				idsToCheck = append(idsToCheck, returnedItem.Id)

			}

		})
	}
	time.Sleep(2 * time.Second)
	//step 2- check its there "GetJobById"

	for _, Id := range idsToCheck {
		testName := fmt.Sprintf("%s/%s", "GetJobById", Id)

		returnedItem := &IdOnly{}

		t.Run(testName, func(t *testing.T) {
			URL := fmt.Sprintf("%s/%s", test.Url, Id)

			resp, err := http.Get(URL)

			if err != nil {
				t.Fatal(err)
			}

			responseBody, err := io.ReadAll(resp.Body)

			resp.Body.Close()

			if err != nil || resp.StatusCode != http.StatusOK {
				t.Fatalf("Error getting job by id")
			}

			if err = json.Unmarshal(responseBody, returnedItem); err != nil {
				t.Fatal(err)
			}

			if len(returnedItem.Id) == 0 {
				t.Fatal("Empty result set")
			}

		})
	}

	//step 3- check GetJobs output
	t.Run("GetJobs", func(t *testing.T) {
		returnedItems := []*IdOnly{}

		resp, err := http.Get(test.Url)

		if err != nil {
			t.Fatal(err)
		}

		responseBody, err := io.ReadAll(resp.Body)

		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Error getting jobs")
		}

		resp.Body.Close()

		if err = json.Unmarshal(responseBody, &returnedItems); err != nil {
			t.Fatal(err)
		}

		if len(returnedItems) == 0 {
			t.Fatal("Empty result set")
		}

	})

	//setp 4- delete the job

	for _, Id := range idsToCheck {
		testName := fmt.Sprintf("%s/%s", "DeleteJob", Id)

		t.Run(testName, func(t *testing.T) {
			URL := fmt.Sprintf("%s/%s", test.Url, Id)

			req, err := http.NewRequest("DELETE", URL, nil)

			if err != nil {
				t.Fatal(err)
			}

			resp, err := http.DefaultClient.Do(req)

			if err != nil || resp.StatusCode != http.StatusOK {
				t.Fatalf("Error deleting job")
			}

		})
	}

}
