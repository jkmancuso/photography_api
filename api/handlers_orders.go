package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jkmancuso/photography_api/shared"
)

func (h handlerDBConn) deleteOrderHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	if len(id) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: "id cannot be empty"})
		return
	}

	if _, err := strconv.Atoi(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: "id needs to be an int"})
		return
	}

	count, err := deleteJob(context.Background(), h.dbInfo, id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	if count == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: "id not found"})
		return
	}

	json.NewEncoder(w).Encode(shared.GenericMsg{Message: "OK"})

}

// Dynamo Query via GSI
func (h handlerDBConn) getOrdersByIdHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	if len(id) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: "id cannot be empty"})
		return
	}

	if !shared.IsUUID(id) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: "id needs to be in uuid format"})
		return
	}

	k := "id"
	v := id

	item, count, err := getOrderByGSI(context.Background(), h.dbInfo, k, v)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	if count == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: "id not found"})
		return
	}

	json.NewEncoder(w).Encode(item)

}

// Dynamo GetItem via Primary Key
func (h handlerDBConn) getOrdersByPKeyHandler(w http.ResponseWriter, r *http.Request) {

	queryParam1 := "record_num"
	queryParam2 := "job_id" //should be some uuid

	queryVal1 := r.URL.Query().Get(queryParam1)
	queryVal2 := r.URL.Query().Get(queryParam2)

	if len(queryVal1) == 0 || len(queryVal2) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: "missing id"})
		return
	}

	if !shared.IsUUID(queryVal2) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: "id needs to be in uuid format"})
		return
	}

	val1, _ := attributevalue.Marshal(queryVal1)
	val2, _ := attributevalue.Marshal(queryVal2)

	pkey := map[string]types.AttributeValue{
		queryParam1: val1,
		queryParam2: val2,
	}

	item, count, err := getOrderByPKey(context.Background(), h.dbInfo, pkey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	if count == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: "record not found"})
		return
	}

	json.NewEncoder(w).Encode(item)

}

func (h handlerDBConn) addOrdersHandler(w http.ResponseWriter, r *http.Request) {
	bytesBody, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
	}

	jobItem, err := shared.ParseBodyIntoNewJob(bytesBody)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
	}

	err = addJob(context.Background(), h.dbInfo, jobItem)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
	}

	json.NewEncoder(w).Encode(jobItem)

}
