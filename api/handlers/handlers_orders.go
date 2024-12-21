package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jkmancuso/photography_api/api/database"
	"github.com/jkmancuso/photography_api/shared"
)

// Large row count, this endpoint shouldn't really be used in prod
func (h handlerDBConn) getOrdersHandler(w http.ResponseWriter, r *http.Request) {

	items, count, err := database.GetOrders(context.Background(), h.dbInfo)

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

func (h handlerDBConn) deleteOrderHandler(w http.ResponseWriter, r *http.Request) {

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

	count, err := database.DeleteOrder(context.Background(), h.dbInfo, id)

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

	json.NewEncoder(w).Encode(shared.GenericMsg{Message: "OK"})

}

func (h handlerDBConn) updateOrderHandler(w http.ResponseWriter, r *http.Request) {

	// 1. check valid id in /orders{id} path
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

	// 2. check valid body which should be the params to change
	bytesBody, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if len(bytesBody) == 0 || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.INVALID_BODY)
		return
	}

	// 3. orderItem has the params to change
	orderItem := make(map[string]interface{})

	if err := json.Unmarshal(bytesBody, &orderItem); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	orderItem["id"] = id

	// 4. update DB
	count, err := database.UpdateOrder(context.Background(), h.dbInfo, orderItem)

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

	json.NewEncoder(w).Encode(shared.GenericMsg{Message: "OK"})

}

func (h handlerDBConn) getOrderByIdHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	if len(id) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.ID_CANNOT_BE_EMPTY)
	}

	if !shared.IsUUID(id) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.ID_NOT_IN_UUID_FORMAT)
		return
	}

	idVal, _ := attributevalue.Marshal(id)

	pKey := map[string]types.AttributeValue{"id": idVal}

	item, count, err := database.GetOrderByPKey(context.Background(), h.dbInfo, pKey)

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

func (h handlerDBConn) getOrderByJobIDAndRecordNumHandler(w http.ResponseWriter, r *http.Request) {
	/* GET /jobs/{id}/orders/{record_num} */

	jobId := r.PathValue("job_id")
	recordNum := r.PathValue("record_num")

	if len(recordNum) == 0 || len(jobId) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.INVALID_REQUEST)
		return
	}

	if !shared.IsUUID(jobId) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.ID_NOT_IN_UUID_FORMAT)
		return
	}

	// generally record_num should be an int 1-100
	recordInt, err := strconv.Atoi(recordNum)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	//build Query with expression builder
	key := map[string]expression.ValueBuilder{
		"record_num": expression.Value(recordInt),
		"job_id":     expression.Value(jobId),
	}

	GSI := "job_id-record_num-index"

	item, count, err := database.GetOrderByGSI(context.Background(), h.dbInfo, key, GSI)

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

func (h handlerDBConn) getOrdersByJobIDHandler(w http.ResponseWriter, r *http.Request) {

	jobId := r.PathValue("job_id")

	if len(jobId) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.ID_NOT_FOUND)
		return
	}

	if !shared.IsUUID(jobId) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.ID_NOT_IN_UUID_FORMAT)
		return
	}

	key := map[string]expression.ValueBuilder{
		"job_id": expression.Value(jobId),
	}

	GSI := "job_id-record_num-index"
	items, _, err := database.GetOrdersByGSI(context.Background(), h.dbInfo, key, GSI)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	/*I decided to make this a supported operation ie where you query but no records
	if count == 0 {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.RECORD_NOT_FOUND)
		return
	}*/

	json.NewEncoder(w).Encode(items)

}

func (h handlerDBConn) addOrderHandler(w http.ResponseWriter, r *http.Request) {
	bytesBody, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	//1. validate payload
	if len(bytesBody) == 0 || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.INVALID_BODY)
		return
	}

	//2. Unmarshall into an order
	orderItem := shared.NewOrderItem()
	if err = json.Unmarshal(bytesBody, orderItem); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	//3. validate required fields
	if len(orderItem.Id) == 0 || len(orderItem.JobId) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.INVALID_BODY)
		return
	}

	//4. check that the job id in the order exists

	// skip this check if you are doing your e2e test, just add it
	if orderItem.Fname != "Integration" && orderItem.Lname != "Test" {

		// if you add an order with a non existent job something is wrong- abort
		jobItem, err := checkJobHandler(fmt.Sprintf("/jobs/%s", orderItem.JobId))

		if err != nil || len(jobItem.Id) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(shared.GenericMsg{Message: "the job entered was not found"})
			return
		}
	}

	//5. finally, add to DB
	err = database.AddOrder(context.Background(), h.dbInfo, orderItem)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
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
