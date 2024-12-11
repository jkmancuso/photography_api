package shared

import(
	"testing"
	"fmt"
	"io"
	"net/http"
	"encoding/json"
)



func (i *integrationTest) setPassword() error{

	var err error

	awsCfg, err = NewAWSCfg()

	if err != nil {
		return(err)
	}

	password, err := GetSecret(awsCfg, "testlogin") 

	if err != nil {
		return(err)
	}

	i.Password = password
}

func (i *integrationTest) authSetup() {

	//1. Get the testlogin password from aws secrets manager
	err := i.setPassword()

	if err != nil {
		t.Fatal(err)
	}

	//2. create a valid auth to test success
	auth := &Auth{
		Email: i.Email,
		Password: i.Password
	}

	validBody, err := json.Marshal(auth)

	if err !=nil {
		t.Fatal(err)
	}

	//3. create an invalid auth to test failure 
	auth := &Auth{
		Email: "invalid_email@something.com",
		Password: "invalid_pass@something.com",
	}

	invalidBody, err := json.Marshal(auth)

	if err !=nil {
		t.Fatal(err)
	}

	i.tests = []GenericTest{
		{
			Name:           "check valid auth",
			BodyBytes:           validBody,
			WantStatusCode: 200,
		},
		{
			Name:           "check invalid auth rejects",
			BodyBytes:           invalidBody,
			WantStatusCode: 400,
		},
	}


}

func (i *integrationTest) teardown() {
	
}


func newIntegrationTest() *integrationTest{
	return &integrationTest{
		Email: Email,
		BaseUrl: URL,
	}
}

func TestAuth(t *testing.T) {
	
	test := newIntegrationTest()
	test.authSetup()

	for _, tt := range test.tests {
		t.Run(tt.Name, func(t *testing.T) {
			reader := bytes.NewReader(body)

			postURL := fmt.Sprintf("%s/%s",test.BaseUrl,"auth")
			resp, err := http.Post(postURL, ContentType, reader)

			if err !=nil {
				t.Fatal(err)
			}

			if resp.StatusCode != tt.WantStatusCode {
				t.Fatalf("Got status %d, wanted %d", resp.StatusCode, WantStatusCode)
			}

			if resp.StatusCode ==http.StatusOK && len(resp.Header.Get("Set-Cookie"))==0 {
				t.Fatal("Set-Cookie header not returned")
			}
		})
	}

	defer test.teardown() 

}