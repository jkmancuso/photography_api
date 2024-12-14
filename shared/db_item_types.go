package shared

import "time"

type DBAdminItem struct {
	Email    string `dynamodbav:"email" json:"email"`
	Hashpass string `dynamodbav:"hashpass" json:"hashpass"`
	Token    string `dynamodbav:"token" json:"Token"`
}

type DBLoginItem struct {
	Email     string `dynamodbav:"email"`
	LoginDate int    `dynamodbav:"login_date"`
	Success   bool   `dynamodbav:"success"`
}

type DBJobItem struct {
	Id       string `dynamodbav:"id" json:"id,omitempty"`
	JobName  string `dynamodbav:"job_name" json:"job_name"`
	JobYear  int    `dynamodbav:"job_year" json:"job_year"`
	ExpireAt int64  `dynamodbav:"expire_at,omitempty" json:"expire_at,omitempty"`
	/*
		DynamoDB expire TTL
		Using this for purging integration tests automatically

		omitempty is SUPER important!
		Without it, it will default to 0 which means every entry gets purged
	*/
}

type DBOrderItem struct {
	Id                   string    `dynamodbav:"id" json:"id,omitempty"`
	JobId                string    `dynamodbav:"job_id" json:"job_id"`
	RecordNum            int       `dynamodbav:"record_num" json:"record_num"`
	Fname                string    `dynamodbav:"fname" json:"fname"`
	Lname                string    `dynamodbav:"lname" json:"lname"`
	Address              string    `dynamodbav:"address" json:"address"`
	City                 string    `dynamodbav:"city" json:"city"`
	State                string    `dynamodbav:"state" json:"state"`
	Zip                  string    `dynamodbav:"zip" json:"zip"`
	Phone                string    `dynamodbav:"phone" json:"phone"`
	GroupQuantity        int       `dynamodbav:"group_quantity" json:"group_quantity"`
	Group                string    `dynamodbav:"group" json:"group"`
	GroupPictureNum      string    `dynamodbav:"group_picture_num" json:"group_picture_num"`
	Instrument           string    `dynamodbav:"instrument" json:"instrument"`
	InstrumentQuantity   int       `dynamodbav:"instrument_quantity" json:"instrument_quantity"`
	InstrumentPosition   int       `dynamodbav:"instrument_position" json:"instrument_position"`
	InstrumentPictureNum string    `dynamodbav:"instrument_picture_num" json:"instrument_picture_num"`
	CheckNum             int       `dynamodbav:"check_num" json:"check_num"`
	Amount               int       `dynamodbav:"amount" json:"amount"`
	CreatedAt            time.Time `dynamodbav:"created_at"`
	ExpireAt             int64     `dynamodbav:"expire_at,omitempty" json:"expire_at,omitempty"`
}