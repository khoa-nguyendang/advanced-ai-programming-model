package entities

import (
	"aapi/shared/constants"
)

type User struct {
	EntityBase
	Id                int64               `json:"id,omitempty" db:"id,omitempty"`
	CompanyCode       string              `json:"company_code,omitempty" db:"company_code,omitempty"`
	UserId            string              `json:"user_id,omitempty" db:"user_id,omitempty"`
	UserName          string              `json:"user_name,omitempty" db:"user_name,omitempty"`
	UserRoleId        int64               `json:"user_role_id,omitempty" db:"user_role_id,omitempty"`
	UserGroupIds      []int64             `json:"user_group_ids,omitempty" db:"user_group_ids,omitempty"`
	UserInfo          string              `json:"user_info,omitempty" db:"user_info,omitempty"`
	UserState         constants.UserState `json:"user_state,omitempty" db:"user_state,omitempty"`
	ThumbnailImageUrl string              `json:"thumbnail_image_url,omitempty" db:"thumbnail_image_url,omitempty"`
	LastModified      int64               `json:"last_modified,omitempty" db:"last_modified,omitempty"`
	IssuedDate        int64               `json:"issued_date,omitempty" db:"issued_date,omitempty"`
	ExpiryDate        int64               `json:"expiry_date,omitempty" db:"expiry_date,omitempty"`
	ActivationDate    int64               `json:"activation_date,omitempty" db:"activation_date,omitempty"`
	ReferenceId       string              `json:"reference_id,omitempty" db:"reference_id,omitempty"`
}

type MongoUserBase struct {
	EmployeeId string `bson:"employeeId"`
	Company    string `bson:"company"`
}

type MongoUser struct {
	Partition string `bson:"_partition"`

	Name string `bson:"name"`

	// required=True
	EmployeeId string `bson:"employeeId"`
	Company    string `bson:"company"`

	Password string `bson:"password"`
	Role     string `bson:"role"`
	CardId   string `bson:"cardId"`
	UserInfo string `bson:"userInfo"`

	// required=True
	Created  float64 `bson:"created"`
	IsActive bool    `bson:"isActive"`

	Group []string `bson:"group"`

	FeatureVector0 []float64 `bson:"featureVector0,omitempty"`
	FeatureVector1 []float64 `bson:"featureVector1,omitempty"`
	FeatureVector2 []float64 `bson:"featureVector2,omitempty"`
	FeatureVector3 []float64 `bson:"featureVector3,omitempty"`
	FeatureVector4 []float64 `bson:"featureVector4,omitempty"`

	FeatureVector5 []float64 `bson:"featureVector5,omitempty"`
	FeatureVector6 []float64 `bson:"featureVector6,omitempty"`
	FeatureVector7 []float64 `bson:"featureVector7,omitempty"`
	FeatureVector8 []float64 `bson:"featureVector8,omitempty"`
	FeatureVector9 []float64 `bson:"featureVector9,omitempty"`

	FeatureVector10 []float64 `bson:"featureVector10,omitempty"`
	FeatureVector11 []float64 `bson:"featureVector11,omitempty"`
	FeatureVector12 []float64 `bson:"featureVector12,omitempty"`
	FeatureVector13 []float64 `bson:"featureVector13,omitempty"`
	FeatureVector14 []float64 `bson:"featureVector14,omitempty"`

	FeatureVector15 []float64 `bson:"featureVector15,omitempty"`
	FeatureVector16 []float64 `bson:"featureVector16,omitempty"`
	FeatureVector17 []float64 `bson:"featureVector17,omitempty"`
	FeatureVector18 []float64 `bson:"featureVector18,omitempty"`
	FeatureVector19 []float64 `bson:"featureVector19,omitempty"`
}
