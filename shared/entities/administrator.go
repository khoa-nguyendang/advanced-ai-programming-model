package entities

type Administrator struct {
	EntityBase
	Id                int64  `json:"id,omitempty" db:"id,omitempty"`
	CompanyCode       string `json:"company_code,omitempty" db:"company_code,omitempty"`
	AdministratorId   string `json:"administrator_id,omitempty" db:"administrator_id,omitempty"`
	Username          string `json:"username,omitempty" db:"username,omitempty"`
	Password          string `json:"password,omitempty" db:"password,omitempty"`
	FullName          string `json:"full_name,omitempty" db:"full_name,omitempty"`
	PhoneNumber       string `json:"phone_number,omitempty" db:"phone_number,omitempty"`
	Email             string `json:"email,omitempty" db:"email,omitempty"`
	AdministratorInfo string `json:"administrator_info,omitempty" db:"administrator_info,omitempty"`
	State             int64  `json:"state,omitempty" db:"state,omitempty"`

	CreatedBy      string `json:"created_by,omitempty" db:"created_by,omitempty"`
	CreatedAt      int64  `json:"created_at,omitempty" db:"created_at,omitempty"`
	LastModifiedBy string `json:"last_modified_by,omitempty" db:"last_modified_by,omitempty"`
	LastModifiedAt int64  `json:"last_modified_at,omitempty" db:"last_modified_at,omitempty"`

	RoleId      int64  `json:"role_id,omitempty" db:"role_id,omitempty"`
	ReferenceId string `json:"reference_id,omitempty" db:"reference_id,omitempty"`

	Salt string `json:"salt,omitempty" db:"salt,omitempty"`
}
