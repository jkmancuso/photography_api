package shared

var OK = GenericMsg{Message: "OK"}
var ID_NOT_FOUND = GenericMsg{Message: "id not found"}
var ID_CANNOT_BE_EMPTY = GenericMsg{Message: "id cannot be empty"}
var ID_NOT_IN_UUID_FORMAT = GenericMsg{Message: "id needs to be in uuid format"}
var INVALID_BODY = GenericMsg{Message: "body is missing or invalid"}
var INVALID_REQUEST = GenericMsg{Message: "request is invalid"}
var RECORD_NOT_FOUND = GenericMsg{Message: "no records found"}
var RECORD_IN_USE = GenericMsg{Message: "the record you are attempting to delete is in use"}
var INVALID_USER_PASS = GenericMsg{Message: "invalid user or password"}

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
	BodyStr           string
	BodyBytes []byte
	WantStatusCode int
	WantErrorMsg   GenericMsg
}