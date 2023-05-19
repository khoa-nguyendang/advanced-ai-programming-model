package entities

import (
	"aapi/shared/constants"
)

type Log struct {
	Id          int64               `json:"id,omitempty" db:"id"`
	AppUserId   int64               `json:"app_user_id,omitempty" db:"app_user_id"`
	Activity    constants.Activity  `json:"activity,omitempty" db:"activity,omitempty"`
	DeviceUuid  string              `json:"device_uuid,omitempty" db:"device_uuid,omitempty"`
	FullMessage string              `json:"full_message,omitempty" db:"full_message,omitempty"`
	Date        int64               `json:"date,omitempty" db:"date,omitempty"`
	CompanyCode string              `json:"company_code,omitempty" db:"company_code,omitempty"`
	UserId      string              `json:"user_id,omitempty" db:"user_id,omitempty"`
	UserName    string              `json:"user_name,omitempty" db:"user_name,omitempty"`
	UserState   constants.UserState `json:"user_state,omitempty" db:"user_state,omitempty"`
}

type UserVerificationEvent struct {
	AppUserId        int64  `json:"app_user_id,omitempty" db:"app_user_id,omitempty"`
	VerificationTime int64  `json:"verification_time,omitempty" db:"verification_time,omitempty"`
	DeviceId         string `json:"device_id,omitempty" db:"device_id,omitempty"`
}
