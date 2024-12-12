package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func (i *IntegrationTest) setPassword() error {

	var err error

	awsCfg, err := NewAWSCfg()

	if err != nil {
		return (err)
	}

	password, err := GetSecret(awsCfg, "testlogin")

	if err != nil {
		return (err)
	}

	i.Password = password
	return nil
}

func (i *IntegrationTest) authSetup(t *testing.T) {
	t.Helper()

	//1. Get the testlogin password from aws secrets manager
	err := i.setPassword()

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

func (i *IntegrationTest) teardown() {

}

func NewIntegrationTest() *IntegrationTest {
	return &IntegrationTest{
		Email:   Email,
		BaseUrl: URL,
	}
}

func TestAuth(t *testing.T) {

	test := NewIntegrationTest()
	test.authSetup(t)

	for _, tt := range test.Tests {
		t.Run(tt.Name, func(t *testing.T) {
			reader := bytes.NewReader(tt.BodyBytes)

			postURL := fmt.Sprintf("%s/%s", test.BaseUrl, "auth")
			resp, err := http.Post(postURL, ContentType, reader)

			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != tt.WantStatusCode {
				t.Fatalf("Got status %d, wanted %d", resp.StatusCode, tt.WantStatusCode)
			}

			if resp.StatusCode == http.StatusOK && len(resp.Header.Get("Set-Cookie")) == 0 {
				t.Fatal("Set-Cookie header not returned")
			}
		})
	}

	defer test.teardown()

}
