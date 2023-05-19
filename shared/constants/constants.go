package constants

const (
	DeviceUUID  string = "device_uuid"
	DeviceId    string = "device_id"
	DeviceMAC   string = "device_uuid"
	AppUserId   string = "app_user_id"
	UserId      string = "user_id"
	Username    string = "username"
	CardId      string = "card_id"
	RoleId      string = "role_id"
	EmployeeId  string = "employee_id"
	OfficeIds   string = "office_ids"
	GroupId     string = "group_id"
	GroupIds    string = "group_ids"
	CompanyId   string = "company_id"
	CompanyCode string = "company_code"

	ENABLED  = "ENABLED"
	DISABLED = "DISABLED"

	// Env variable
	// Upload image to MinIO: ENABLED or DISABLED
	VERIFICATION_IMAGE_LOG = "VERIFICATION_IMAGE_LOG"

	// ACL implementation: ENABLED or DISABLED
	ACL_FEATURE = "ACL_FEATURE"

	// Write verification log to CSV file: ENABLED or DISABLED
	VERIFICATION_CSV_LOG = "VERIFICATION_CSV_LOG"
	// Publish verification events: ENABLED or DISABLED
	VERIFICATION_EVENT = "VERIFICATION_EVENT"

	MAX_MESSAGE_LENGTH = 20000000
)
