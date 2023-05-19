package viewmodels

type NewDeviceConfig struct {
	DeviceUUID   string  `json:"device_uuid,omitempty" db:"device_uuid,omitempty"`
	MaskFeature  bool    `json:"mask_feature,omitempty" db:"mask_feature,omitempty"`
	TempFeature  bool    `json:"temp_feature,omitempty" db:"temp_feature,omitempty"`
	TempValue    float32 `json:"temp_value,omitempty" db:"temp_value,omitempty"`
	AntiSpoofing bool    `json:"anti_spoofing,omitempty" db:"anti_spoofing,omitempty"`
	MatchingMode int64   `json:"matching_mode,omitempty" db:"matching_mode,omitempty"`
}
