package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jkmancuso/photography_api/api/database"
	"github.com/jkmancuso/photography_api/shared"
	log "github.com/sirupsen/logrus"
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
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	if len(orderItem.Id) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.RECORD_IN_USE)
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

	//1. validate payload
	if len(bytesBody) == 0 || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.INVALID_BODY)
		return
	}

	//2. Unmarshall into an job
	jobItem := shared.NewJobItem()
	if err := json.Unmarshal(bytesBody, jobItem); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	//3. validate job is not empty
	if len(jobItem.Id) == 0 || len(jobItem.JobName) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.INVALID_BODY)
		return
	}

	//4. finally, add to DB
	if err = database.AddJob(context.Background(), h.dbInfo, jobItem); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(jobItem)

}

// check if an order exists using this job in case you
// attempt to delete the job
func checkOrderHandler(url string) (*shared.DBOrderItem, error) {
	log.Debugf("checkOrderHandler: GET %s", url)

	orderItem := &shared.DBOrderItem{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Printf("ERROR: %v", err)
		return orderItem, err
	}

	respRecorder := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(respRecorder, req)

	resultBytes, err := io.ReadAll(respRecorder.Result().Body)
	log.Debugf("Response: %s", respRecorder.Body.String())

	if err != nil {
		log.Println(err)
		return orderItem, err
	}

	if err = json.Unmarshal(resultBytes, orderItem); err != nil {
		log.Println(err)
		return orderItem, err
	}

	return orderItem, nil
}
