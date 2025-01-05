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
)

type AuthIntegrationTest struct {
	Email    string
	Password string
	Url      string
	Tests    []GenericTest
}

func (i *AuthIntegrationTest) setCreds() error {

	//set as env variables since most likely running locally
	i.Email = os.Getenv("TESTLOGIN")
	i.Password = os.Getenv("TESTPASSWORD")

	if len(i.Email) == 0 || len(i.Password) == 0 {
		return errors.New("Cannot find local creds")
	}

	return nil
}

func (i *AuthIntegrationTest) authSetup(t *testing.T) {
	t.Helper()

	//1. Get the testlogin password from aws secrets manager
	err := i.setCreds()

	if err != nil {
		t.Fatal(err)
	}

	//2. create a valid auth to test success
	auth := &Auth{
		Email:    i.Email,
		Password: i.Password,
	}

	validBody, err := json.Marshal(auth)

	if err != nil {
		t.Fatal(err)
	}

	//3. create an invalid auth to test failure
	auth = &Auth{
		Email:    "invalid_email@something.com",
		Password: "invalid_pass@something.com",
	}

	invalidBody, err := json.Marshal(auth)

	if err != nil {
		t.Fatal(err)
	}

	i.Tests = []GenericTest{
		{
			Name:           "check valid auth",
			BodyBytes:      validBody,
			WantStatusCode: 200,
		},
		{
			Name:           "check invalid auth rejects",
			BodyBytes:      invalidBody,
			WantStatusCode: 400,
		},
	}

}

func TestAuth(t *testing.T) {

	test := &AuthIntegrationTest{
		Url: fmt.Sprintf("%s/%s", API_URL, "auth"),
	}

	test.authSetup(t)

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

			if resp.StatusCode == http.StatusOK {
				authResponse := &DBSessionItem{}
				responseBody, err := io.ReadAll(resp.Body)

				if err != nil {
					t.Fatal(err)
				}

				err = json.Unmarshal(responseBody, authResponse)

				if err != nil {
					t.Fatal(err)
				}

				if len(authResponse.Id) == 0 {
					t.Fatal("Session Id not returned!")
				}
			}
		})
	}

}
