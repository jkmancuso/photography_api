package shared

import(
	"testing"
	"fmt"
	"io"
	"net/http"
	"encoding/json"
)

const URL = "aws here"
const email = "test@test.com"
const contentType = "application/json"

type integrationTest struct{
	Email string
	Password string
	BaseUrl string
}

func (i integrationTest) setup() {
	var err error

	awsCfg, err = NewAWSCfg()

	if err != nil {
		t.Fatal(err)
	}

	password, err := GetSecret(awsCfg, "testpassword") 

	if err != nil {
		t.Fatal(err)
	}

	i.Password = password

}
func (i integrationTest) teardown() {
	
}

func newIntegrationTest() *integrationTest{
	return &integrationTest{
		Email: email,
		BaseUrl: URL,
	}
}

func TestAuth(t *testing.T) {
	
	test := newIntegrationTest()
	test.setup()

	auth := Auth{
		Email: test.Email,
		Password: test.Password
	}

	body, err := json.Marshal(auth)

	if err !=nil {
		t.Fatal(err)
	}
	reader := bytes.NewReader(body)

	postURL := fmt.Sprintf("%s/%s",integrationTest.URL,"auth", reader)
	resp, err := http.Post(postURL, contentType)

	if err !=nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Got status %d, wanted %d", resp.StatusCode, http.StatusOK)
	}

	if len(resp.Header.Get("Set-Cookie"))==0{
		t.Fatal("Set-Cookie header not returned")
	}

	defer test.teardown() 

}