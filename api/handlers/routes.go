package handlers

import (
	"net/http"

	"github.com/jkmancuso/photography_api/shared"
)

func SetupRoutes(tableMap map[string]*shared.DBInfo, mux *http.ServeMux) {
	/*
		use this struct to provide the DB client since it cannot be passed
		as a parameter to function signature (w http.ResponseWriter, r *http.Request)
		and we want to avoid sending in global variables
	*/
	jobsDBConn := handlerDBConn{dbInfo: tableMap["jobs"]}
	instrumentsDBConn := handlerDBConn{dbInfo: tableMap["instruments"]}
	groupsDBConn := handlerDBConn{dbInfo: tableMap["groups"]}
	ordersDBConn := handlerDBConn{dbInfo: tableMap["orders"]}
	zipcodesDBConn := handlerDBConn{dbInfo: tableMap["zipcodes"]}

	/*JOBS*/
	mux.Handle("GET /jobs",
		shared.CorsMiddleware(http.HandlerFunc(jobsDBConn.getJobsHandler)))
	mux.Handle("GET /jobs/{id}",
		shared.ValidateIdMiddleware(shared.CorsMiddleware(http.HandlerFunc(jobsDBConn.getJobByIdHandler))))
	mux.Handle("DELETE /jobs/{id}",
		shared.ValidateIdMiddleware(shared.CorsMiddleware(http.HandlerFunc(jobsDBConn.deleteJobHandler))))
	mux.Handle("POST /jobs",
		shared.CorsMiddleware(http.HandlerFunc(jobsDBConn.addJobHandler)))

	/*INSTRUMENTS*/
	mux.Handle("GET /instruments",
		shared.CorsMiddleware(http.HandlerFunc(instrumentsDBConn.getInstrumentsHandler)))
	mux.Handle("GET /instruments/{id}",
		shared.ValidateIdMiddleware(shared.CorsMiddleware(http.HandlerFunc(instrumentsDBConn.getInstrumentByIdHandler))))
	mux.Handle("DELETE /instruments/{id}",
		shared.ValidateIdMiddleware(shared.CorsMiddleware(http.HandlerFunc(instrumentsDBConn.deleteInstrumentHandler))))
	mux.Handle("POST /instruments",
		shared.CorsMiddleware(http.HandlerFunc(instrumentsDBConn.addInstrumentHandler)))
	mux.Handle("PATCH /instruments/{id}",
		shared.ValidateIdMiddleware(shared.CorsMiddleware(http.HandlerFunc(instrumentsDBConn.updateInstrumentHandler))))

	/*GROUPS*/
	mux.Handle("GET /groups",
		shared.CorsMiddleware(http.HandlerFunc(groupsDBConn.getGroupsHandler)))
	mux.Handle("GET /groups/{id}",
		shared.ValidateIdMiddleware(shared.CorsMiddleware(http.HandlerFunc(groupsDBConn.getGroupByIdHandler))))
	mux.Handle("DELETE /groups/{id}",
		shared.ValidateIdMiddleware(shared.CorsMiddleware(http.HandlerFunc(groupsDBConn.deleteGroupHandler))))
	mux.Handle("POST /groups",
		shared.CorsMiddleware(http.HandlerFunc(groupsDBConn.addGroupHandler)))
	mux.Handle("PATCH /groups/{id}",
		shared.ValidateIdMiddleware(shared.CorsMiddleware(http.HandlerFunc(groupsDBConn.updateGroupHandler))))

	/*ORDERS*/
	mux.Handle("GET /jobs/{job_id}/orders/{record_num}",
		shared.CorsMiddleware(http.HandlerFunc(ordersDBConn.getOrderByJobIDAndRecordNumHandler)))
	mux.Handle("GET /jobs/{job_id}/orders",
		shared.CorsMiddleware(http.HandlerFunc(ordersDBConn.getOrdersByJobIDHandler)))
	mux.Handle("GET /orders",
		shared.CorsMiddleware(http.HandlerFunc(ordersDBConn.getOrdersHandler)))
	mux.Handle("GET /orders/{id}",
		shared.ValidateIdMiddleware(shared.CorsMiddleware(http.HandlerFunc(ordersDBConn.getOrderByIdHandler))))
	mux.Handle("DELETE /orders/{id}",
		shared.CorsMiddleware(http.HandlerFunc(ordersDBConn.deleteOrderHandler)))
	mux.Handle("POST /orders",
		shared.CorsMiddleware(http.HandlerFunc(ordersDBConn.addOrderHandler)))
	mux.Handle("PATCH /orders/{id}",
		shared.ValidateIdMiddleware(shared.CorsMiddleware(http.HandlerFunc(ordersDBConn.updateOrderHandler))))

	mux.Handle("GET /zipcodes/{code}",
		shared.CorsMiddleware(http.HandlerFunc(zipcodesDBConn.getZipByIdHandler)))
}
