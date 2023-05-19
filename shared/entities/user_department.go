package entities

type UserGroup struct {
	EntityBase
	UserID  int64 `json:"user_id,omitempty" db:"user_id,omitempty"`
	GroupID int64 `json:"group_id,omitempty" db:"group_id,omitempty"`
	Enabled int   `json:"enabled,omitempty" db:"enabled,omitempty"`
}
