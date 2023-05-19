package entities

import (
	"aapi/shared/constants"
)

type Role struct {
	EntityBase
	Name        string             `json:"name,omitempty" db:"name,omitempty"`
	Description string             `json:"description,omitempty" db:"description,omitempty"`
	RoleState   int32              `json:"role_state,omitempty" db:"role_state,omitempty"`
	RoleType    constants.RoleType `json:"role_type,omitempty" db:"role_type,omitempty"`
	IsStatic    int                `json:"is_static,omitempty" db:"is_static,omitempty"`
}
