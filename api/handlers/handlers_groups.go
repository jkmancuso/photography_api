package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jkmancuso/photography_api/api/database"
	"github.com/jkmancuso/photography_api/shared"
)

func (h handlerDBConn) updateGroupHandler(w http.ResponseWriter, r *http.Request) {

	// 1. check valid id in /groups/{id} path
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
	groupItem := make(map[string]interface{})

	if err := json.Unmarshal(bytesBody, &groupItem); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	groupItem["id"] = id

	// 4. update DB
	count, err := database.UpdateOrder(context.Background(), h.dbInfo, groupItem)

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

func (h handlerDBConn) getGroupsHandler(w http.ResponseWriter, r *http.Request) {

	items, count, err := database.GetGroups(context.Background(), h.dbInfo)

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

func (h handlerDBConn) deleteGroupHandler(w http.ResponseWriter, r *http.Request) {

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

	count, err := database.DeleteGroup(context.Background(), h.dbInfo, id)

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

func (h handlerDBConn) getGroupByIdHandler(w http.ResponseWriter, r *http.Request) {

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

	item, count, err := database.GetGroupById(context.Background(), h.dbInfo, pKey)

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

func (h handlerDBConn) addGroupHandler(w http.ResponseWriter, r *http.Request) {
	bytesBody, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	//1. validate payload
	if len(bytesBody) == 0 || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.INVALID_BODY)
		return
	}

	//2. Unmarshall into an group
	groupItem := shared.NewGroupItem()
	if err := json.Unmarshal(bytesBody, groupItem); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	//3. validate Group is not empty
	if len(groupItem.Id) == 0 || len(groupItem.GroupName) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(shared.INVALID_BODY)
		return
	}

	//4. finally, add to DB
	if err = database.AddGroup(context.Background(), h.dbInfo, groupItem); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(groupItem)

}
