package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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

	if !shared.IsUUID(id) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: "id needs to be an int"})
		return
	}

	count, err := deleteOrder(context.Background(), h.dbInfo, id)

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

	item, count, err := getOrderByGSI(context.Background(), h.dbInfo, k, v, h.dbInfo.GSI)

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

	// generally record_num should be an int 1-100
	intVal1, err := strconv.Atoi(queryVal1)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: "wrong format"})
		return
	}

	val1, _ := attributevalue.Marshal(intVal1)
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

	orderItem, err := shared.ParseBodyIntoNewOrder(bytesBody)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
	}

	// if you add an order with a non existent job something is wrong- abort
	jobItem, err := checkJobHandler(fmt.Sprintf("/jobs/%s", orderItem.JobId))

	if err != nil || len(jobItem.Id) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: "the job entered was not found"})
		return
	}

	err = addOrder(context.Background(), h.dbInfo, orderItem)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
	}

	json.NewEncoder(w).Encode(orderItem)

}

// check if a job exists
func checkJobHandler(url string) (*shared.DBJobItem, error) {
	jobItem := &shared.DBJobItem{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return jobItem, err
	}

	respRecorder := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(respRecorder, req)

	resultBytes, err := io.ReadAll(respRecorder.Result().Body)

	if err != nil {
		return jobItem, err
	}

	err = json.Unmarshal(resultBytes, &jobItem)
	return jobItem, err
}
