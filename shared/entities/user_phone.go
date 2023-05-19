package entities

type UserPhone struct {
	EntityBase
	UserID      int64  `json:"user_id,omitempty" db:"user_id,omitempty"`
	CountryCode string `json:"country_code,omitempty" db:"country_code,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty" db:"phone_number,omitempty"`
	PhoneTypeID int64  `json:"phone_type_id,omitempty" db:"phone_type_id,omitempty"`
}
