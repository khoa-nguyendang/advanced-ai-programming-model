package entities

import "aapi/shared/constants"

type AdministratorInfo struct {
	Id                int64                        `json:"id" db:"id"`
	CompanyCode       string                       `json:"company_code" db:"company_code"`
	AdministratorId   string                       `json:"administrator_id" db:"administrator_id"`
	Username          string                       `json:"username" db:"username"`
	Password          string                       `json:"password" db:"password"`
	FullName          string                       `json:"full_name" db:"full_name"`
	PhoneNumber       string                       `json:"phone_number" db:"phone_number"`
	Email             string                       `json:"email" db:"email"`
	AdministratorInfo string                       `json:"administrator_info" db:"administrator_info"`
	State             constants.AdministratorState `json:"state" db:"state"`
	CreatedBy         string                       `json:"created_by" db:"created_by"`
	CreatedAt         string                       `json:"created_at" db:"created_at"`
	LastModifiedBy    string                       `json:"last_modified_by" db:"last_modified_by"`
	LastModifiedAt    string                       `json:"last_modified_at" db:"last_modified_at"`
	RoleId            int64                        `json:"role_id" db:"role_id"`
	ReferenceId       string                       `json:"reference_id" db:"reference_id"`
	Salt              string                       `json:"salt" db:"salt"`
}
