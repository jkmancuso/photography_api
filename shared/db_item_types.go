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
	Id        string   `dynamodbav:"id" json:"id,omitempty"`
	JobName   string   `dynamodbav:"job_name" json:"job_name"`
	JobYear   int      `dynamodbav:"job_year" json:"job_year"`
	JobGroups []string `dynamodbav:"job_groups" json:"job_groups,omitempty"`
	ExpireAt  int64    `dynamodbav:"expire_at,omitempty" json:"expire_at,omitempty"`
	/*
		DynamoDB expire TTL
		Using this for purging integration tests automatically

		omitempty is SUPER important!
		Without it, it will default to 0 which means every entry gets purged
	*/
}

type DBOrderItem struct {
	Id                   string    `dynamodbav:"id" json:"id,omitempty"`
	JobId                string    `dynamodbav:"job_id" json:"job_id,omitempty"`
	RecordNum            int       `dynamodbav:"record_num" json:"record_num,omitempty"`
	Fname                string    `dynamodbav:"fname" json:"fname,omitempty"`
	Lname                string    `dynamodbav:"lname" json:"lname,omitempty"`
	Address              string    `dynamodbav:"address" json:"address,omitempty"`
	City                 string    `dynamodbav:"city" json:"city,omitempty"`
	State                string    `dynamodbav:"state" json:"state,omitempty"`
	Zip                  string    `dynamodbav:"zip" json:"zip,omitempty"`
	Phone                string    `dynamodbav:"phone" json:"phone,omitempty"`
	GroupQuantity        int       `dynamodbav:"group_quantity" json:"group_quantity,omitempty"`
	Group                string    `dynamodbav:"group" json:"group,omitempty"`
	GroupPictureNum      string    `dynamodbav:"group_picture_num" json:"group_picture_num,omitempty"`
	Instrument           string    `dynamodbav:"instrument" json:"instrument,omitempty"`
	InstrumentQuantity   int       `dynamodbav:"instrument_quantity" json:"instrument_quantity,omitempty"`
	InstrumentPosition   int       `dynamodbav:"instrument_position" json:"instrument_position,omitempty"`
	InstrumentPictureNum string    `dynamodbav:"instrument_picture_num" json:"instrument_picture_num,omitempty"`
	CheckNum             int       `dynamodbav:"check_num" json:"check_num,omitempty"`
	Amount               int       `dynamodbav:"amount" json:"amount,omitempty"`
	CreatedAt            time.Time `dynamodbav:"created_at,omitempty"`
	ExpireAt             int64     `dynamodbav:"expire_at,omitempty" json:"expire_at,omitempty"`
}

type DBGroupItem struct {
	Id        string `dynamodbav:"id" json:"id,omitempty"`
	GroupName string `dynamodbav:"group_name" json:"group_name,omitempty"`
	ExpireAt  int64  `dynamodbav:"expire_at,omitempty" json:"expire_at,omitempty"`
}

type DBInstrumentItem struct {
	Id             string `dynamodbav:"id" json:"id,omitempty"`
	InstrumentName string `dynamodbav:"instrument_name" json:"instrument_name,omitempty"`
	Section        string `dynamodbav:"section" json:"section,omitempty"`
	ExpireAt       int64  `dynamodbav:"expire_at,omitempty" json:"expire_at,omitempty"`
}

type DBPictureItem struct {
	Id         string `dynamodbav:"id" json:"id,omitempty"`
	PictureNum string `dynamodbav:"picture_num" json:"picture_num,omitempty"`
	Section    string `dynamodbav:"section" json:"section,omitempty"`
	ExpireAt   int64  `dynamodbav:"expire_at,omitempty" json:"expire_at,omitempty"`
}

type DBZipItem struct {
	Code  string `dynamodbav:"code" json:"code"`
	City  string `dynamodbav:"city" json:"city"`
	State string `dynamodbav:"state" json:"state"`
}
