package entities

type UserPermission struct {
	EntityBase
	UserID       int64 `json:"user_id,omitempty" db:"user_id,omitempty"`
	PermissionID int64 `json:"permission_id,omitempty" db:"permission_id,omitempty"`
	Enabled      bool  `json:"enabled,omitempty" db:"enabled,omitempty"`
}
