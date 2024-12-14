package shared

type GenericMsg struct {
	Message string `json:"message"`
}

type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GenericTest struct {
	Name           string
	Id             string
	BodyStr        string
	BodyBytes      []byte
	WantStatusCode int
	WantErrorMsg   GenericMsg
}

type IdOnly struct {
	Id string `json:"id"`
}
