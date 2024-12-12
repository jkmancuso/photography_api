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
		json.NewEncoder(w).Encode(shared.RECORD_NOT_FOUND)
		return
	}

	json.NewEncoder(w).Encode(items)

}

func (h handlerDBConn) deleteJobHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	if len(id) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.ID_CANNOT_BE_EMPTY)
		return
	}

	if !shared.IsUUID(id) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.ID_NOT_IN_UUID_FORMAT)
		return
	}

	orderURL := fmt.Sprintf("/jobs/%s/orders", id)
	orderItem, err := checkOrderHandler(orderURL)

	if err != nil {
		log.Printf("Got err1: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	if len(orderItem.Id) != 0 {
		log.Println("Got err2")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.RECORD_IN_USE)
		return
	}

	count, err := database.DeleteJob(context.Background(), h.dbInfo, id)

	if err != nil {
		log.Printf("Got err3: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	if count == 0 {
		log.Println("Got err4")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.RECORD_NOT_FOUND)
		return
	}

	json.NewEncoder(w).Encode(shared.OK)

}

func (h handlerDBConn) getJobByIdHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	if len(id) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.ID_CANNOT_BE_EMPTY)
		return
	}

	if !shared.IsUUID(id) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.ID_NOT_IN_UUID_FORMAT)
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
		json.NewEncoder(w).Encode(shared.RECORD_NOT_FOUND)
		return
	}

	json.NewEncoder(w).Encode(item)

}

func (h handlerDBConn) addJobHandler(w http.ResponseWriter, r *http.Request) {
	bytesBody, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if len(bytesBody) == 0 || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.INVALID_BODY)
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

	orderItem := &shared.DBOrderItem{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Printf("ERROR: %v", err)
		return orderItem, err
	}

	respRecorder := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(respRecorder, req)

	resultBytes, err := io.ReadAll(respRecorder.Result().Body)
	log.Println(respRecorder.Body.String())

	if err != nil {
		return orderItem, err
	}

	if err = json.Unmarshal(resultBytes, orderItem); err != nil {
		return orderItem, err
	}
	return orderItem, nil
}
