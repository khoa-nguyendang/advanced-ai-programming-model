package entities

type UserImage struct {
	EntityBase
	UserID      int64  `json:"user_id,omitempty" db:"user_id,omitempty"`
	Path        string `json:"path,omitempty" db:"path,omitempty"`
	ImageTypeId int64  `json:"image_type_id,omitempty" db:"image_type_id,omitempty"`
}
