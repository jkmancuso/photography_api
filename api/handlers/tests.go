package handlers

import "github.com/jkmancuso/photography_api/shared"

type GenericTest struct {
	Name           string
	Id             string
	Body           string
	WantStatusCode int
	WantErrorMsg   shared.GenericMsg
}
