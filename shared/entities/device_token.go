package entities

import "aapi/shared/constants"

type DeviceToken struct {
	Id                 int64                  `json:"id" db:"id"`                                                         // tracking identity
	DeviceId           int64                  `json:"device_id,omitempty" db:"device_id,omitempty"`                       // device identity
	Token              string                 `json:"token,omitempty" db:"token,omitempty"`                               // token
	TokenExpiry        int64                  `json:"token_expiry,omitempty" db:"token_expiry,omitempty"`                 // expiry of token
	RefreshToken       string                 `json:"refresh_token,omitempty" db:"refresh_token,omitempty"`               // refresh token
	RefreshTokenExpiry int64                  `json:"refresh_token_expiry,omitempty" db:"refresh_token_expiry,omitempty"` // expiry of refresh token
	Action             constants.DeviceAction `json:"action,omitempty" db:"action,omitempty"`
	ActionState        constants.ActionState  `json:"action_state,omitempty" db:"action_state,omitempty"`
	ExecutedIp         string                 `json:"executed_ip,omitempty" db:"executed_ip,omitempty"`
	Timestamp          int64                  `json:"timestamp,omitempty" db:"timestamp,omitempty"`
}
