package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/jkmancuso/photography_api/shared"
)

type handlerMetadata struct {
	Salt  string
	DBMap map[string]*shared.DBInfo
}

func newHandlerMetadata(salt string, DB map[string]*shared.DBInfo) *handlerMetadata {
	return &handlerMetadata{
		Salt:  salt,
		DBMap: DB,
	}
}

func ping(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(shared.GenericMsg{Message: "pong"})

}

func (h handlerMetadata) postAuth(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

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

	loginItem := shared.NewLoginItem(auth.Email)

	//add a login entry for audit and set success to false
	err = addLogin(ctx, h.DBMap["logins"], loginItem)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	hashpass, err := shared.GenerateHash(auth.Password, h.Salt)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	sess, err := returnSessionForValidAuth(context.Background(),
		auth.Email,
		hashpass,
		h.DBMap["admins"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	// new session to the DB
	err = addSession(ctx, h.DBMap["sessions"], sess)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	//update audit login success to true
	_, err = updateLogin(ctx, h.DBMap["logins"], loginItem)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(sess)

}
