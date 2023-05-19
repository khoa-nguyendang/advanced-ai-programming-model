package entities

type Permission struct {
	EntityBase
	Name            string `json:"name,omitempty" db:"name,omitempty"`
	PermissionValue int32  `json:"permission_value,omitempty" db:"permission_value,omitempty"`
	Description     string `json:"description,omitempty" db:"description,omitempty"`
	PermissionState int32  `json:"permission_state,omitempty" db:"permission_state,omitempty"`
	IsStatic        int    `json:"is_static,omitempty" db:"is_static,omitempty"`
}
