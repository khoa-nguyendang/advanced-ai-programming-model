package entities

type Office struct {
	EntityBase
	Name           string `json:"name,omitempty" db:"name,omitempty"`
	Description    string `json:"description,omitempty" db:"description,omitempty"`
	OfficeInfoUUID string `json:"office_info_uuid,omitempty" db:"office_info_uuid,omitempty"`
	OfficeState    int32  `json:"office_state,omitempty" db:"office_state,omitempty"`
}
