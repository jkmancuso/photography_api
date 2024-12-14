package shared

import (
	"encoding/json"
	"time"
)

func NewJobItem() *DBJobItem {
	return &DBJobItem{
		Id: GenerateUUID(),
	}
}

func NewOrderItem() *DBOrderItem {
	return &DBOrderItem{
		Id:        GenerateUUID(),
		CreatedAt: time.Now(),
	}
}

func NewQAJobItem() []byte {
	jobItem := &DBJobItem{}

	jobItem.JobName = "IntegrationTest_Job"
	jobItem.JobYear = time.Now().Year()
	jobItem.ExpireAt = ExpireIn + time.Now().Unix()

	result, _ := json.Marshal(jobItem)

	return result
}

func NewQAOrderItem() []byte {
	orderItem := &DBOrderItem{}

	orderItem.JobId = GenerateUUID()
	orderItem.RecordNum = 1
	orderItem.Fname = "Integration"
	orderItem.Lname = "Test"
	orderItem.ExpireAt = ExpireIn + time.Now().Unix()

	result, _ := json.Marshal(orderItem)

	return result
}

func NewDBItem(table string) []byte {
	var b []byte

	switch table {
	case "jobs":
		b = NewQAJobItem()
	case "orders":
		b = NewQAOrderItem()

	}

	return b
}

func NewLoginItem(email string) *DBLoginItem {
	return &DBLoginItem{
		Email:     email,
		LoginDate: int(time.Now().Unix()),
		Success:   false, //start it as false and update to true when complete
	}
}
