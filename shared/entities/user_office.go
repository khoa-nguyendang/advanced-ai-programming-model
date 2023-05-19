package entities

type UserOffice struct {
	EntityBase
	UserID   int64 `json:"user_id,omitempty" db:"user_id,omitempty"`
	OfficeID int64 `json:"office_id,omitempty" db:"office_id,omitempty"`
	Enabled  bool  `json:"enabled,omitempty" db:"enabled,omitempty"`
}
