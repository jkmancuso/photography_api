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
	instrumentsDBConn := handlerDBConn{dbInfo: tableMap["instruments"]}
	picturesDBConn := handlerDBConn{dbInfo: tableMap["pictures"]}
	groupsDBConn := handlerDBConn{dbInfo: tableMap["groups"]}
	ordersDBConn := handlerDBConn{dbInfo: tableMap["orders"]}

	/*JOBS*/
	http.HandleFunc("GET /jobs", jobsDBConn.getJobsHandler)
	http.HandleFunc("GET /jobs/{id}", jobsDBConn.getJobByIdHandler)
	http.HandleFunc("DELETE /jobs/{id}", jobsDBConn.deleteJobHandler)
	http.HandleFunc("POST /jobs", jobsDBConn.addJobHandler)

	/*INSTRUMENTS*/
	http.HandleFunc("GET /instruments", instrumentsDBConn.getInstrumentsHandler)
	http.HandleFunc("GET /instruments/{id}", instrumentsDBConn.getInstrumentByIdHandler)
	http.HandleFunc("DELETE /instruments/{id}", instrumentsDBConn.deleteInstrumentHandler)
	http.HandleFunc("POST /instruments", instrumentsDBConn.addInstrumentHandler)

	/*PICTURES*/
	http.HandleFunc("GET /pictures", picturesDBConn.getPicturesHandler)
	http.HandleFunc("GET /pictures/{id}", picturesDBConn.getPictureByIdHandler)
	http.HandleFunc("DELETE /pictures/{id}", picturesDBConn.deletePictureHandler)
	http.HandleFunc("POST /pictures", picturesDBConn.addPictureHandler)

	/*GROUPS*/
	http.HandleFunc("GET /groups", groupsDBConn.getGroupsHandler)
	http.HandleFunc("GET /groups/{id}", groupsDBConn.getGroupByIdHandler)
	http.HandleFunc("DELETE /groups/{id}", groupsDBConn.deleteGroupHandler)
	http.HandleFunc("POST /groups", groupsDBConn.addGroupHandler)

	/*ORDERS*/
	http.HandleFunc("GET /jobs/{job_id}/orders/{record_num}", ordersDBConn.getOrderByJobIDAndRecordNumHandler)
	http.HandleFunc("GET /jobs/{job_id}/orders", ordersDBConn.getOrdersByJobIDHandler)
	http.HandleFunc("GET /orders", ordersDBConn.getOrdersHandler)
	http.HandleFunc("GET /orders/{id}", ordersDBConn.getOrderByIdHandler)
	http.HandleFunc("DELETE /orders/{id}", ordersDBConn.deleteOrderHandler)
	http.HandleFunc("POST /orders", ordersDBConn.addOrderHandler)
	http.HandleFunc("PATCH /orders/{id}", ordersDBConn.updateOrderHandler)
}
