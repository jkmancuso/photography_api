package handlers

import (
	"net/http"

	"github.com/jkmancuso/photography_api/shared"
)

func SetupRoutes(tableMap map[string]*shared.DBInfo) {
	/*
		use this struct to provide the DB client since it cannot be passed
		as a parameter to function signature (w http.ResponseWriter, r *http.Request)
		and we want to avoid sending in global variables
	*/
	jobsDBConn := handlerDBConn{dbInfo: tableMap["jobs"]}
	ordersDBConn := handlerDBConn{dbInfo: tableMap["orders"]}

	http.HandleFunc("GET /jobs", jobsDBConn.getJobsHandler)
	http.HandleFunc("GET /jobs/{id}", jobsDBConn.getJobsByIdHandler)
	http.HandleFunc("DELETE /jobs/{id}", jobsDBConn.deleteJobHandler)
	http.HandleFunc("POST /jobs", jobsDBConn.addJobsHandler)

	http.HandleFunc("GET /orders/{id}", ordersDBConn.getOrdersByIdHandler)
	http.HandleFunc("GET /orders", ordersDBConn.getOrdersByGSIHandler)
	http.HandleFunc("POST /orders", ordersDBConn.addOrdersHandler)
	http.HandleFunc("DELETE /orders/{id}", ordersDBConn.deleteOrderHandler)
}
