package entities

import (
	"aapi/shared/constants"
)

type Device struct {
	Id                int64                 `json:"id,omitempty" db:"id,omitempty"`
	CompanyCode       string                `json:"company_code,omitempty" db:"company_code,omitempty"`
	GroupId           int64                 `json:"group_id,omitempty" db:"group_id,omitempty"`
	DeviceUUID        string                `json:"device_uuid,omitempty" db:"device_uuid,omitempty"`
	DeviceName        string                `json:"device_name,omitempty" db:"device_name,omitempty"`
	DeviceAppVersion  string                `json:"device_app_version,omitempty" db:"device_app_version,omitempty"`
	DeviceDescription string                `json:"device_description,omitempty" db:"device_description,omitempty"`
	LocationCode      string                `json:"location_code,omitempty" db:"location_code,omitempty"`
	DeviceType        int64                 `json:"device_type,omitempty" db:"device_type,omitempty"`
	DeviceState       constants.DeviceState `json:"device_state,omitempty" db:"device_state,omitempty"`
	ApprovedUserId    string                `json:"approved_user_id,omitempty" db:"approved_user_id,omitempty"`
	LastModified      int64                 `json:"last_modified,omitempty" db:"last_modified,omitempty"`
}
