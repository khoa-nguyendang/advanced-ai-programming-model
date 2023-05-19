package entities

type Company struct {
	EntityBase
	Name            string   `json:"name,omitempty" db:"name,omitempty"`
	CompanyCode     string   `json:"company_code,omitempty" db:"company_code,omitempty"`
	Description     string   `json:"description,omitempty" db:"description,omitempty"`
	License         string   `json:"license,omitempty" db:"license,omitempty"`
	Version         string   `json:"version,omitempty" db:"version,omitempty"`
	CompanyInfoUUID string   `json:"company_info_uuid,omitempty" db:"company_info_uuid,omitempty"`
	CompanyState    int32    `json:"company_state,omitempty" db:"company_state,omitempty"`
	Groups          []*Group `json:"groups,omitempty"`
}
