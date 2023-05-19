package entities

type MqttSubscriber struct {
	Id             int64  `json:"id" db:"id"`
	LastModified   int64  `json:"last_modified,omitempty" db:"last_modified,omitempty"`
	DeviceUUID     string `json:"device_uuid,omitempty" db:"device_uuid,omitempty"`
	CompanyCode    string `json:"company_code,omitempty" db:"company_code,omitempty"`
	Topic          string `json:"topic,omitempty" db:"topic,omitempty"`
	SubscribeState int32  `json:"subscribe_state,omitempty" db:"subscribe_state,omitempty"` //0 = unsubscribe , 1 = subscribe
}
