package entities

type AdministratorLoginSession struct {
	ID            int64  `json:"id" db:"id"`
	UserID        int64  `json:"user_id,omitempty" db:"user_id,omitempty"`
	Agent         string `json:"agent,omitempty" db:"agent,omitempty"`
	AccessKey     string `json:"access_key,omitempty" db:"access_key,omitempty"`
	Cookies       string `json:"cookies,omitempty" db:"cookies,omitempty"`
	LoginType     int64  `json:"login_type,omitempty" db:"login_type,omitempty"`
	LastActivity  int64  `json:"last_activity,omitempty" db:"last_activity,omitempty"`
	GeneratedTime int64  `json:"generated_time,omitempty" db:"generated_time,omitempty"`
}
