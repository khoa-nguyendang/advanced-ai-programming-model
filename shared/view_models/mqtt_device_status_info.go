package viewmodels

type MqttDeviceStatusInfo struct {
	MessageUuid        string `json:"message_uuid,omitempty" db:"message_uuid,omitempty" bson:"message_uuid,omitempty"`
	DeviceUUID         string `json:"device_uuid"`
	DeviceName         string `json:"device_name"`
	DeviceStatus       int32  `json:"device_status"`
	Token              string `json:"token"`
	RefreshToken       string `json:"refresh_token"`
	TokenExpiry        int64  `json:"token_expiry"`
	RefreshTokenExpiry int64  `json:"refresh_token_expiry"`
}
