package viewmodels

import "encoding/json"

func UnmarshalGetUserResponse(data []byte) (GetUserResponse, error) {
	var r GetUserResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GetUserResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GetUserResponse struct {
	Users       []User `json:"users"`
	CurrentPage int    `json:"current_page"`
	PageSize    int    `json:"page_size"`
	TotalCount  int    `json:"total_count"`
}

type User struct {
	UserGroups          []string `json:"user_groups"`
	RegisteredImageUrls []string `json:"registered_image_urls"`
	UserGroupIDS        []string `json:"user_group_ids"`
	Images              []string `json:"images"`
	UserID              string   `json:"user_id"`
	UserName            string   `json:"user_name"`
	UserRole            string   `json:"user_role"`
	UserInfo            string   `json:"user_info"`
	LastModified        string   `json:"last_modified"`
	State               string   `json:"state"`
	ThumbnailImageURL   string   `json:"thumbnail_image_url"`
	UserRoleID          string   `json:"user_role_id"`
	ReferenceID         string   `json:"reference_id"`
	IssuedDate          string   `json:"issued_date"`
	ExpiryDate          string   `json:"expiry_date"`
	ActivationDate      string   `json:"activation_date"`
	IsActive            bool     `json:"is_active"`
}
