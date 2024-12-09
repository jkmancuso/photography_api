package shared

import (
	"encoding/json"
	"errors"
	"log"
	"time"
)

func NewJobItem() *DBJobItem {
	return &DBJobItem{
		Id: GenerateUUID(),
	}
}

func NewOrderItem() *DBOrderItem {
	return &DBOrderItem{
		Id: GenerateUUID(),
	}
}

func NewLoginItem(email string) *DBLoginItem {
	return &DBLoginItem{
		Email:     email,
		LoginDate: int(time.Now().Unix()),
		Success:   false, //start it as false and update to true when complete
	}
}

func ParseBodyIntoNewJob(body []byte) (*DBJobItem, error) {
	jobItem := NewJobItem()
	err := json.Unmarshal(body, jobItem)

	if len(jobItem.JobName) == 0 || jobItem.JobYear == 0 {
		err = errors.New("missing field in body")
	}

	log.Println(jobItem)

	return jobItem, err
}

func ParseBodyIntoNewOrder(body []byte) (*DBOrderItem, error) {
	orderItem := NewOrderItem()
	err := json.Unmarshal(body, orderItem)

	if orderItem.RecordNum == 0 || len(orderItem.JobId) == 0 {
		err = errors.New("missing field in body")
	}

	log.Println(orderItem)

	return orderItem, err
}
