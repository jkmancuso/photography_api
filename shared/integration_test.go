package shared

import(
	"net/http"
	"testing"
	"fmt"
	"encoding/json"
)

const URL = "aws here"
const email = "test@test.com"

type integrationTest struct{
	Email string
	Password string
	BaseUrl string
}

func (i integrationTest) setup() error{
	//i.setPassword()

}
func (i integrationTest) teardown() error{
	
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

	postURL := fmt.Sprintf("%s/%s",integrationTest.URL,"auth", body)
	resp, err := http.Post(postURL,"application/json")

	if err !=nil {
		t.Fatal(err)
	}	
	defer test.teardown() 

}