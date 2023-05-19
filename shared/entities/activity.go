package entities

type Activity struct {
	EntityBase
	Name        string `json:"name,omitempty" db:"name,omitempty"`
	Description string `json:"description,omitempty" db:"description,omitempty"`
	TargetGroup string `json:"target_group,omitempty" db:"target_group,omitempty"`
}
