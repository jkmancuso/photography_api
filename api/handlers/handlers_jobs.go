package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jkmancuso/photography_api/api/database"
	"github.com/jkmancuso/photography_api/shared"
)

func (h handlerDBConn) getJobsHandler(w http.ResponseWriter, r *http.Request) {

	items, count, err := database.GetJobs(context.Background(), h.dbInfo)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	if count == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: "no jobs found"})
		return
	}

	json.NewEncoder(w).Encode(items)

}

func (h handlerDBConn) deleteJobHandler(w http.ResponseWriter, r *http.Request) {

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

	orderURL := fmt.Sprintf("/orders?record_num=1&job_id=%s", id)
	orderItem, err := checkOrderHandler(orderURL)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	if len(orderItem.Id) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: "attempting to delete a job in use"})
		return
	}

	count, err := database.DeleteJob(context.Background(), h.dbInfo, id)

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

func (h handlerDBConn) getJobsByIdHandler(w http.ResponseWriter, r *http.Request) {

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

	idVal, _ := attributevalue.Marshal(id)
	pKey := map[string]types.AttributeValue{"id": idVal}

	item, count, err := database.GetJobById(context.Background(), h.dbInfo, pKey)

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

func (h handlerDBConn) addJobsHandler(w http.ResponseWriter, r *http.Request) {
	bytesBody, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if len(bytesBody) == 0 || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: "invalid request body"})
		return
	}

	jobItem, err := shared.ParseBodyIntoNewJob(bytesBody)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	err = database.AddJob(context.Background(), h.dbInfo, jobItem)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(jobItem)

}

// check if an order exists using this job in case you
// attempt to delete the job
func checkOrderHandler(url string) (*shared.DBOrderItem, error) {
	log.Printf("URL: %s", url)
	orderItem := &shared.DBOrderItem{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Printf("ERROR: %v", err)
		return orderItem, err
	}

	respRecorder := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(respRecorder, req)

	resultBytes, err := io.ReadAll(respRecorder.Result().Body)

	if err != nil {
		log.Printf("ERROR: %v", err)
		return orderItem, err
	}

	err = json.Unmarshal(resultBytes, &orderItem)
	return orderItem, err
}
