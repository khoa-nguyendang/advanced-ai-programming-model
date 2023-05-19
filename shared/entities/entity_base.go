package entities

type EntityBase struct {
	Id                 int64 `json:"id" db:"id"`
	CreatedAt          int64 `json:"created_at,omitempty" db:"created_at,omitempty"`
	LastModified       int64 `json:"last_modified,omitempty" db:"last_modified,omitempty"`
	CreatedByUserID    int64 `json:"created_user_id,omitempty" db:"created_user_id,omitempty"`
	LastModifiedUserID int64 `json:"last_modified_user_id,omitempty" db:"last_modified_user_id,omitempty"`
}
