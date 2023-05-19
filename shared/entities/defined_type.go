package entities

type DefinedType struct {
	EntityBase
	Name        string `json:"name,omitempty" db:"name,omitempty"`
	Description string `json:"description,omitempty" db:"description,omitempty"`
	IsStatic    int    `json:"is_static,omitempty" db:"is_static,omitempty"`
	TargetGroup string `json:"target_group,omitempty" db:"target_group,omitempty"`
}
