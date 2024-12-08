package shared

var OK = GenericMsg{Message: "OK"}
var ID_NOT_FOUND = GenericMsg{Message: "id not found"}
var ID_CANNOT_BE_EMPTY = GenericMsg{Message: "id cannot be empty"}
var ID_NOT_IN_UUID_FORMAT = GenericMsg{Message: "id needs to be in uuid format"}
var INVALID_BODY = GenericMsg{Message: "body is missing or invalid"}
var RECORD_NOT_FOUND = GenericMsg{Message: "no records found"}
var RECORD_IN_USE = GenericMsg{Message: "the record you are attempting to delete is in use"}

type GenericMsg struct {
	Message string `json:"message"`
}
