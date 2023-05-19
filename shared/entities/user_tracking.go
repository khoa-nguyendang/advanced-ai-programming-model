package entities

type UserTracking struct {
	Id         int64 `json:"id,omitempty" db:"id,omitempty"`
	UserId     int64 `json:"user_id,omitempty" db:"user_id,omitempty"`
	DeviceId   int64 `json:"device_id,omitempty" db:"device_id,omitempty"`
	ActivityId int64 `json:"activity_id,omitempty" db:"activity_id,omitempty"`
	Timestamp  int64 `json:"timestamp,omitempty" db:"timestamp,omitempty"`
}
