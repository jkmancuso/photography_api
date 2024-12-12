package shared

var (
	OK                    = GenericMsg{Message: "OK"}
	ID_NOT_FOUND          = GenericMsg{Message: "id not found"}
	ID_CANNOT_BE_EMPTY    = GenericMsg{Message: "id cannot be empty"}
	ID_NOT_IN_UUID_FORMAT = GenericMsg{Message: "id needs to be in uuid format"}
	INVALID_BODY          = GenericMsg{Message: "body is missing or invalid"}
	INVALID_REQUEST       = GenericMsg{Message: "request is invalid"}
	RECORD_NOT_FOUND      = GenericMsg{Message: "no records found"}
	RECORD_IN_USE         = GenericMsg{Message: "the record you are attempting to delete is in use"}
	INVALID_USER_PASS     = GenericMsg{Message: "invalid user or password"}
)
