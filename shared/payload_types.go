package shared

import "time"

type GenericMsg struct {
	Message string `json:"message"`
}

type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Session struct {
	SessionId string
	ExpireAt  time.Time
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
