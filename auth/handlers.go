package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jkmancuso/photography_api/shared"
)

type handlerMetadata struct {
	Salt string
}

func newHandlerMetadata(salt string) *handlerMetadata {
	return &handlerMetadata{
		Salt: salt,
	}
}

func ping(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(shared.GenericMsg{Message: "pong"})

}

func (h handlerMetadata) postAuth(w http.ResponseWriter, r *http.Request) {
	bytesBody, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if len(bytesBody) == 0 || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.INVALID_BODY)
		return
	}

	auth := shared.Auth{}

	err = json.Unmarshal(bytesBody, &auth)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	if len(auth.Email) == 0 || len(auth.Password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.INVALID_BODY)
		return
	}

	hashpass, err := shared.GenerateHash(auth.Password, h.Salt)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.INVALID_USER_PASS)
		return
	}

	token, err := checkDBForValidAuth(auth.Email, hashpass)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	if len(token) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.INVALID_USER_PASS)
		return
	}

	w.Header().Add("Set-Cookie", fmt.Sprintf("token=%q; max-age=%d", token, 86400))

	json.NewEncoder(w).Encode(shared.OK)

}
