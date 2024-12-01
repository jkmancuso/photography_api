package main

import "net/http"

func setupRoutes() {
	http.HandleFunc("GET /jobs", getJobsHandler)
	http.HandleFunc("GET /jobs/{id}", getJobsByIdHandler)
	http.HandleFunc("POST /jobs", addJobsHandler)
}
