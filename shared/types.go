package shared

var NO_ERR = GenericMsg{Message: ""}
var ID_NOT_FOUND = GenericMsg{Message: "id not found"}
var ID_NOT_IN_UUID_FORMAT = GenericMsg{Message: "id needs to be in uuid format"}

type GenericMsg struct {
	Message string `json:"message"`
}
