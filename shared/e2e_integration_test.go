package shared

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

type IntegrationTest struct {
	Url          string
	Client       *http.Client
	XSessionId   string
	Tests        []GenericTest
	EndpointName string //this is [orders|jobs|etc]
}

func (i *IntegrationTest) setCreds() error {

	//set as env variables since most likely running locally
	i.XSessionId = os.Getenv("SESSION_ID")

	if len(i.XSessionId) == 0 {
		return errors.New("Cannot find local variable SESSION_ID")
	}

	return nil
}

func (i *IntegrationTest) setup(t *testing.T) {
	t.Helper()

	if err := i.setCreds(); err != nil {
		t.Fatal(err)
	}

	i.Client = &http.Client{Timeout: time.Second * 3}

	//populated item
	validPayload := NewDBItem(i.EndpointName)

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

	i.Url = fmt.Sprintf("%s/%s", API_URL, i.EndpointName)
}

func TestE2E(t *testing.T) {

	idsToCheck := []string{}

	//dont need a full struct, just check Id is there is fine
	returnedItem := &IdOnly{}

	test := &IntegrationTest{
		EndpointName: "orders",
	}

	test.setup(t)

	//step 1- add the item
	for _, tt := range test.Tests {
		t.Run(tt.Name, func(t *testing.T) {
			reader := bytes.NewReader(tt.BodyBytes)

			req, err := http.NewRequest("POST", test.Url, reader)

			if err != nil {
				t.Fatal(err)
			}

			req.Header.Add("x-session-id", test.XSessionId)

			resp, err := test.Client.Do(req)

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
	//step 2- check its there "Get[Job|Order|etc]ById"

	for _, Id := range idsToCheck {
		testName := fmt.Sprintf("Get%sById/%s", test.EndpointName, Id)

		returnedItem := &IdOnly{}

		t.Run(testName, func(t *testing.T) {
			URL := fmt.Sprintf("%s/%s", test.Url, Id)

			req, err := http.NewRequest("GET", URL, nil)

			if err != nil {
				t.Fatal(err)
			}

			req.Header.Add("x-session-id", test.XSessionId)

			resp, err := test.Client.Do(req)

			if err != nil {
				t.Fatal(err)
			}

			responseBody, err := io.ReadAll(resp.Body)

			resp.Body.Close()

			if err != nil || resp.StatusCode != http.StatusOK {
				t.Fatalf("Error getting by id")
			}

			if err = json.Unmarshal(responseBody, returnedItem); err != nil {
				t.Fatal(err)
			}

			if len(returnedItem.Id) == 0 {
				t.Fatal("Empty result set")
			}

		})
	}

	//step 3- check Get[Item]s output
	testName := fmt.Sprintf("Get%s", test.EndpointName)

	t.Run(testName, func(t *testing.T) {
		returnedItems := []*IdOnly{}

		req, err := http.NewRequest("GET", test.Url, nil)

		if err != nil {
			t.Fatal(err)
		}

		req.Header.Add("x-session-id", test.XSessionId)

		resp, err := test.Client.Do(req)

		if err != nil {
			t.Fatal(err)
		}

		responseBody, err := io.ReadAll(resp.Body)

		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Error getting %s", test.EndpointName)
		}

		resp.Body.Close()

		if err = json.Unmarshal(responseBody, &returnedItems); err != nil {
			t.Fatal(err)
		}

		if len(returnedItems) == 0 {
			t.Fatal("Empty result set")
		}

	})

	//setp 4- delete the item

	for _, Id := range idsToCheck {
		testName := fmt.Sprintf("Delete%s/%s", test.EndpointName, Id)

		t.Run(testName, func(t *testing.T) {
			URL := fmt.Sprintf("%s/%s", test.Url, Id)

			req, err := http.NewRequest("DELETE", URL, nil)

			if err != nil {
				t.Fatal(err)
			}

			req.Header.Add("x-session-id", test.XSessionId)

			resp, err := test.Client.Do(req)

			if err != nil || resp.StatusCode != http.StatusOK {
				t.Fatalf("Error deleting %s", test.EndpointName)
			}

		})
	}

}
