package entities

import (
	"aapi/shared/constants"
)

type Location struct {
	LocationCode        string                  `json:"location_code,omitempty" db:"location_code,omitempty"`
	PICID               string                  `json:"pic_id,omitempty" db:"pic_id,omitempty"`
	LocationType        int64                   `json:"type,omitempty" db:"type,omitempty"`
	LocationName        string                  `json:"name,omitempty" db:"name,omitempty"`
	LocationDescription string                  `json:"description,omitempty" db:"description,omitempty"`
	LocationState       constants.LocationState `json:"state,omitempty" db:"state,omitempty"`
	CompanyCode         string                  `json:"company_code,omitempty" db:"company_code,omitempty"`
}
