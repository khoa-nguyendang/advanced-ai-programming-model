package entities

type UserRole struct {
	EntityBase
	UserID  int64 `json:"user_id,omitempty" db:"user_id,omitempty"`
	RoleID  int64 `json:"role_id,omitempty" db:"role_id,omitempty"`
	Enabled int   `json:"enabled,omitempty" db:"enabled,omitempty"`
}
