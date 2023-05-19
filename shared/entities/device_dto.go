package entities

import (
	"aapi/shared/constants"
)

type DeviceDto struct {
	Id                int64                 `json:"id,omitempty" db:"devices.id,omitempty"`
	CompanyCode       string                `json:"company_code,omitempty" db:"devices.company_code,omitempty"`
	GroupId           int64                 `json:"group_id,omitempty" db:"devices.group_id,omitempty"`
	DeviceUUID        string                `json:"device_uuid,omitempty" db:"devices.device_uuid,omitempty"`
	DeviceName        string                `json:"device_name,omitempty" db:"devices.device_name,omitempty"`
	DeviceAppVersion  string                `json:"device_app_version,omitempty" db:"devices.device_app_version,omitempty"`
	DeviceDescription string                `json:"device_description,omitempty" db:"devices.device_description,omitempty"`
	LocationCode      string                `json:"location_code,omitempty" db:"devices.location_code,omitempty"`
	DeviceType        int64                 `json:"device_type,omitempty" db:"devices.device_type,omitempty"`
	DeviceState       constants.DeviceState `json:"device_state,omitempty" db:"devices.device_state,omitempty"`
	ApprovedUserId    string                `json:"approved_user_id,omitempty" db:"devices.approved_user_id,omitempty"`
	LastModified      int64                 `json:"last_modified,omitempty" db:"devices.last_modified,omitempty"`
	DeviceConfig      DeviceConfig          `json:"device_config,omitempty" db:"device_configuration"`
}
