package viewmodels

type DeviceInfoVM struct {
	MessageUuid       string          `json:"message_uuid,omitempty" db:"message_uuid,omitempty" bson:"message_uuid,omitempty"`
	CompanyCode       string          `protobuf:"bytes,1,opt,name=company_code,json=companyCode,proto3" json:"company_code"`
	DeviceUuid        string          `protobuf:"bytes,2,opt,name=device_uuid,json=deviceUuid,proto3" json:"device_uuid"`
	DeviceName        string          `protobuf:"bytes,3,opt,name=device_name,json=deviceName,proto3" json:"device_name"`
	DeviceAppVersion  string          `protobuf:"bytes,4,opt,name=device_app_version,json=deviceAppVersion,proto3" json:"device_app_version"`
	DeviceDescription string          `protobuf:"bytes,6,opt,name=device_description,json=deviceDescription,proto3" json:"device_description"`
	LocationCode      string          `protobuf:"bytes,7,opt,name=location_code,json=locationCode,proto3" json:"location_code"`
	DeviceType        int32           `protobuf:"varint,8,opt,name=device_type,json=deviceType,proto3,enum=device_service.DeviceType" json:"device_type"`
	DeviceState       int32           `protobuf:"varint,9,opt,name=device_state,json=deviceState,proto3,enum=device_service.DeviceState" json:"device_state"`
	DeviceConfig      *DeviceConfigVM `protobuf:"bytes,10,opt,name=device_config,json=deviceConfig,proto3" json:"device_config"`
	UserGroupIds      []int64         `protobuf:"varint,11,rep,packed,name=user_group_ids,json=userGroupIds,proto3" json:"user_group_ids"`
	LastModified      int64           `protobuf:"varint,12,opt,name=last_modified,json=lastModified,proto3" json:"last_modified"`
	DeviceId          int64           `protobuf:"varint,13,opt,name=device_id,json=deviceId,proto3" json:"device_id"`
	GroupId           int64           `protobuf:"varint,14,opt,name=group_id,json=groupId,proto3" json:"group_id"` // string group_name = 15;
}

type DeviceConfigVM struct {
	MaskFeature  bool    `protobuf:"varint,1,opt,name=mask_feature,json=maskFeature,proto3" json:"mask_feature"`
	TempFeature  bool    `protobuf:"varint,2,opt,name=temp_feature,json=tempFeature,proto3" json:"temp_feature"`
	TempValue    float32 `protobuf:"fixed32,3,opt,name=temp_value,json=tempValue,proto3" json:"temp_value"`
	AntiSpoofing bool    `protobuf:"varint,4,opt,name=anti_spoofing,json=antiSpoofing,proto3" json:"anti_spoofing"`
	MatchingMode int32   `protobuf:"varint,5,opt,name=matching_mode,json=matchingMode,proto3,enum=device_service.DeviceMatchingMode" json:"matching_mode"`
}
