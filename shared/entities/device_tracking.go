package entities

import "aapi/shared/constants"

type DeviceTracking struct {
	Id             int64                  `json:"id,omitempty" db:"id,omitempty"`                             //tracking identity
	DeviceId       int64                  `json:"device_id,omitempty" db:"device_id,omitempty"`               // device identity
	ExecutorUserId int64                  `json:"executor_user_id,omitempty" db:"executor_user_id,omitempty"` //executor user id in system
	Action         constants.DeviceAction `json:"action,omitempty" db:"action,omitempty"`
	ActionState    constants.ActionState  `json:"action_state,omitempty" db:"action_state,omitempty"`
	ExecutedIP     string                 `json:"executed_ip,omitempty" db:"executed_ip,omitempty"`
	Timestamp      int64                  `json:"timestamp,omitempty" db:"timestamp,omitempty"`
}
