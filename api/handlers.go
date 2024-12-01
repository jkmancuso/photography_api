package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/jkmancuso/photography_api/shared"
)

func getJobsHandler(w http.ResponseWriter, r *http.Request) {
	items, count, err := getJobs(context.Background(), tableMap["jobs"])

	if err != nil {
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: "Error"})
		return
	}

	if count == 0 {
		json.NewEncoder(w).Encode(struct{}{})
		return
	}

	json.NewEncoder(w).Encode(items)

}

func deleteJobHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	if len(id) == 0 {
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: "id cannot be empty"})
		return
	}

	if _, err := strconv.Atoi(id); err != nil {
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: "id needs to be an int"})
		return
	}

	err := deleteJob(context.Background(), tableMap["jobs"], id)

	if err != nil {
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(shared.GenericMsg{Message: "OK"})

}

func getJobsByIdHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	if len(id) == 0 {
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: "id cannot be empty"})
		return
	}

	if _, err := strconv.Atoi(id); err != nil {
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: "id needs to be an int"})
		return
	}

	item, count, err := getJobById(context.Background(), tableMap["jobs"], id)

	if err != nil {
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
		return
	}

	if count == 0 {
		json.NewEncoder(w).Encode(struct{}{})
		return
	}

	json.NewEncoder(w).Encode(item)

}

func addJobsHandler(w http.ResponseWriter, r *http.Request) {
	bytesBody, err := io.ReadAll(r.Body)

	if err != nil {
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
	}

	jobItem, err := shared.ParseBodyIntoNewJob(bytesBody)

	if err != nil {
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: err.Error()})
	}

	err = addJob(context.Background(), tableMap["jobs"], jobItem)

	if err != nil {
		json.NewEncoder(w).Encode(shared.GenericMsg{Message: "Error"})
	}

	json.NewEncoder(w).Encode(jobItem)

}
