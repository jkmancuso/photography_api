package shared

import (
	"encoding/json"
	"errors"

	log "github.com/sirupsen/logrus"
)

func ParseBodyIntoNewJob(body []byte) (*DBJobItem, error) {
	jobItem := NewJobItem()

	if err := json.Unmarshal(body, jobItem); err != nil {
		log.Println(err)
		return jobItem, err
	}

	log.Debugf("Got Job body: %s", string(body))

	if len(jobItem.JobName) == 0 || jobItem.JobYear == 0 {
		log.Println(INVALID_BODY.Message)
		return jobItem, errors.New(INVALID_BODY.Message)
	}

	return jobItem, nil
}

func ParseBodyIntoNewOrder(body []byte) (*DBOrderItem, error) {
	orderItem := NewOrderItem()

	if err := json.Unmarshal(body, orderItem); err != nil {
		log.Println(err)
		return orderItem, err
	}

	log.Debugf("Got Job body: %s", string(body))

	if orderItem.RecordNum == 0 || len(orderItem.JobId) == 0 {
		log.Println(INVALID_BODY.Message)
		return orderItem, errors.New(INVALID_BODY.Message)
	}

	return orderItem, nil
}
