package viewmodels

type MatchingInput struct {
	UserId       string `json:"userId,omitempty"`
	UserName     string `json:"userName,omitempty"`
	Temperature  string `json:"temperature,omitempty"`
	DeviceName   string `json:"deviceName,omitempty"`
	ImagePathS3  string `json:"imagePathS3,omitempty"`
	PreSignedUrl string `json:"preSignedUrl,omitempty"`
	VerifyStatus bool   `json:"verifyStatus,omitempty"`
	UserStatus   bool   `json:"userStatus,omitempty"`
}
