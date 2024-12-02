package main

import (
	"net/http"
)

func setupRoutes() {
	/*
		use this struct to provide the DB client since it cannot be passed
		as a parameter to function signature (w http.ResponseWriter, r *http.Request)
		and we want to avoid sending in global variables
	*/
	jobsDBConn := handlerDBConn{dbInfo: tableMap["jobs"]}

	http.HandleFunc("GET /jobs", jobsDBConn.getJobsHandler)
	http.HandleFunc("GET /jobs/{id}", jobsDBConn.getJobsByIdHandler)
	http.HandleFunc("DELETE /jobs/{id}", jobsDBConn.deleteJobHandler)
	http.HandleFunc("POST /jobs", jobsDBConn.addJobsHandler)
}
