package interceptors

import "github.com/dgrijalva/jwt-go"

type UserClaims struct {
	jwt.StandardClaims
	Id          int64   `json:"id" db:"id"`
	CompanyId   int64   `json:"company_id" db:"company_id"`
	OfficeIds   []int64 `json:"office_ids"`
	RoleId      int64   `json:"role_ids"`
	GroupIds    []int64 `json:"group_id"`
	EmployeeId  string  `json:"employee_id" db:"employee_id"`
	CardId      string  `json:"card_id" db:"card_id"`
	Username    string  `json:"username" db:"username"`
	CompanyCode string  `json:"company_code" db:"company_code"`
}

type DeviceClaims struct {
	jwt.StandardClaims
	Id             int64  `json:"id" db:"id"`
	DeviceUuid     string `json:"device_uuid" db:"device_uuid"`
	DeviceType     int64  `json:"device_type" db:"device_type"`
	ApprovedUserId string `json:"approved_user_id" db:"approved_user_id"`
	CompanyCode    string `json:"company_code" db:"company_code"`
	CompanyId      int64  `json:"company_id" db:"company_id"`
	GroupId        int64  `json:"group_id"`
}

type DeviceIdentity struct {
	Id             int64  `json:"id" db:"id"`
	LocationCode   string `json:"location_code" db:"location_code"`
	ApprovedUserId string `json:"approved_user_id" db:"approved_user_id"`
	DeviceType     int64  `json:"device_type" db:"device_type"`
	GroupId        int64  `json:"group_id"`
	DeviceUUID     string `json:"device_uuid" db:"device_uuid"`
	CompanyCode    string `json:"company_code" db:"company_code"`
	CompanyId      int64  `json:"company_id" db:"company_id"`
}

type AdministratorIdentity struct {
	Id          int64   `json:"id" db:"id"`
	CompanyId   int64   `json:"company_id" db:"company_id"`
	OfficeIds   []int64 `json:"office_ids"`
	RoleId      int64   `json:"role_id"`
	GroupIds    []int64 `json:"group_ids"`
	EmployeeId  string  `json:"employee_id" db:"employee_id"`
	CardId      string  `json:"card_id" db:"card_id"`
	Username    string  `json:"username" db:"username"`
	CompanyCode string  `json:"company_code" db:"company_code"`
}
