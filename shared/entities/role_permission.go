package entities

type RolePermission struct {
	EntityBase
	RoleID       int64 `json:"role_id,omitempty" db:"role_id,omitempty"`
	PermissionID int64 `json:"permission_id,omitempty" db:"permission_id,omitempty"`
}

type RolePermissions struct {
	EntityBase
	RoleID      int64    `json:"role_id,omitempty" db:"role_id,omitempty"`
	Permissions []string `json:"permissions,omitempty" db:"permissions,omitempty"`
}
