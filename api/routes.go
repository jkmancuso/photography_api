package main

import "net/http"

func setupRoutes() {
	http.HandleFunc("GET /jobs", getJobsHandler)
	http.HandleFunc("GET /jobs/{id}", getJobsByIdHandler)
	http.HandleFunc("DELETE /jobs/{id}", deleteJobHandler)
	http.HandleFunc("POST /jobs", addJobsHandler)
}
