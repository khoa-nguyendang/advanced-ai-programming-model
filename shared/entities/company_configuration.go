package entities

type CompanyConfiguration struct {
	EntityBase
	CompanyCode string `json:"company_code,omitempty" db:"company_code,omitempty"`
	ConfigKey   string `json:"config_key,omitempty" db:"config_key,omitempty"`
	ConfigValue string `json:"config_value,omitempty" db:"config_value,omitempty"`
	Description string `json:"description,omitempty" db:"description,omitempty"`
}
