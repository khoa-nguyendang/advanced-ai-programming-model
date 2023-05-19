package entities

import "aapi/shared/constants"

type AdministratorTracking struct {
	Id              int64                         `json:"id" db:"id"`                                                           //tracking identity
	AdministratorId int64                         `json:"Administrator_user_id,omitempty" db:"administrator_user_id,omitempty"` // administrator identity
	ExecutorUserId  int64                         `json:"executor_user_id,omitempty" db:"executor_user_id,omitempty"`           //executor user id in system
	Action          constants.AdministratorAction `json:"action,omitempty" db:"action,omitempty"`
	ActionState     constants.ActionState         `json:"action_state,omitempty" db:"action_state,omitempty"`
	ExecutedIP      string                        `json:"executed_ip,omitempty" db:"executed_ip,omitempty"`
	Timestamp       int64                         `json:"timestamp,omitempty" db:"timestamp,omitempty"`
}
