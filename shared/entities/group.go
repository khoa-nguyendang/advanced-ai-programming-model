package entities

type Group struct {
	EntityBase
	Id          int64  `json:"id,omitempty" db:"id,omitempty"`
	Name        string `json:"name,omitempty" db:"name,omitempty"`
	CompanyId   int64  `json:"company_id,omitempty" db:"company_id,omitempty"`
	Description string `json:"description,omitempty" db:"description,omitempty"`
	GroupState  int32  `json:"group_state,omitempty" db:"group_state,omitempty"`
	IsStatic    int    `json:"is_static,omitempty" db:"is_static,omitempty"`
}
