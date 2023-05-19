package httpclients

type CustomerInfo struct {
	Id       int64  `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Address  struct {
		Street string `json:"street"`
		Suite  string `json:"suite"`
		City   string `json:"city"`
	} `json:"address"`

	Phone   string `json:"phone"`
	Website string `json:"website"`
	Company struct {
		Name        string `json:"name" db:"name"`
		CatchPhrase string `json:"catchPhrase"`
	} `json:"company"`
}
