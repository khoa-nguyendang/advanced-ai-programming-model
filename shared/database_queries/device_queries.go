package dbqueries

const (
	CheckDevicesExistQuery = `SELECT id FROM devices WHERE device_uuid = ? AND device_state <> ?;`

	AddDeviceQuery = `INSERT INTO devices(company_code, group_id, device_uuid, device_name, device_app_version, device_description, location_code, device_type, device_state, last_modified)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	UpdateDeviceQuery = `UPDATE devices 
	SET device_name = ?, device_app_version = ?, device_description = ?, location_code = ?, device_type = ?, device_state = ?, approved_user_id = ?, last_modified = ?, group_id = ?
	WHERE device_uuid = ? AND device_state <> ?;`

	UpdateDeviceStatusQuery = `UPDATE devices 
	SET device_state = ?, approved_user_id = ?, last_modified = ?
	WHERE device_uuid = ? AND device_state <> ?;`

	DeleteDeviceQuery = `UPDATE devices 
	SET device_state = ?, last_modified = ?
	WHERE device_uuid = ?;`

	FindDeviceQueryByUuids = `SELECT company_code, group_id, device_uuid, device_name, device_app_version, device_description, location_code, device_type, device_state, approved_user_id, last_modified
	FROM devices
	WHERE device_state <> ? AND device_uuid IN (?)
	ORDER BY id desc;`

	FindAllDeviceQuery = `SELECT company_code, group_id, device_uuid, device_name, device_app_version, device_description,  location_code, device_type, device_state, approved_user_id, last_modified
	WHERE device_state <> ?
	FROM devices;`

	AddDeviceConfigQuery = `INSERT INTO device_configuration(device_uuid, mask_feature, temp_feature, temp_value, anti_spoofing, matching_mode)
	VALUES (?, ?, ?, ?, ?, ?)`

	UpdateDeviceConfigQuery = `UPDATE device_configuration 
	SET mask_feature = ?, temp_feature = ?, temp_value = ?, anti_spoofing = ?, matching_mode = ?
	WHERE device_uuid = ?;`

	DeleteDeviceConfig = `DELETE FROM device_configuration WHERE device_uuid = ?;`

	FindDeviceAndConfigQueryByUuids = `SELECT *
	FROM devices
	LEFT JOIN device_configuration
	ON devices.device_uuid = device_configuration.device_uuid
	WHERE device_state <> ? AND devices.device_uuid IN (?)
	ORDER BY id desc;`

	FindDeviceAndConfigWithCompanyCodeQueryByUuids = `SELECT *
	FROM devices
	LEFT JOIN device_configuration
	ON devices.device_uuid = device_configuration.device_uuid
	WHERE company_code = ? AND device_state <> ? AND devices.device_uuid IN (?)
	ORDER BY id desc;`

	FindAllDeviceAndConfigQuery = `SELECT *
	FROM devices
	LEFT JOIN device_configuration
	ON devices.device_uuid = device_configuration.device_uuid
	WHERE device_state <> ?
	ORDER BY id desc;`

	FindAllDeviceAndConfigWithLocationCodeQuery = `SELECT *
	FROM devices
	LEFT JOIN device_configuration
	ON devices.device_uuid = device_configuration.device_uuid
	WHERE location_code = ? AND device_state <> ?
	ORDER BY id desc;`

	FindAllDeviceAndConfigWithCompanyCodeQuery = `SELECT *
	FROM devices
	LEFT JOIN device_configuration
	ON devices.device_uuid = device_configuration.device_uuid
	WHERE company_code = ? AND device_state <> ?
	ORDER BY id desc;`

	FindAllDeviceAndConfigWithCompanyCodeAndLocationCodeQuery = `SELECT *
	FROM devices
	LEFT JOIN device_configuration
	ON devices.device_uuid = device_configuration.device_uuid
	WHERE company_code = ? AND location_code = ? AND device_state <> ?
	ORDER BY id desc;`

	DEVICE_GetAllQuery = `SELECT 
		devices.id as id,
		devices.company_code as company_code,
		devices.group_id as group_id,
		devices.device_uuid as device_uuid,
		devices.device_name as device_name,
		devices.device_app_version as device_app_version,
		devices.device_description as device_description,
		devices.location_code as location_code,
		devices.device_type as device_type,
		devices.device_state as device_state,
		devices.approved_user_id as approved_user_id,
		devices.last_modified as last_modified,
		device_configuration.mask_feature as mask_feature, 
		device_configuration.temp_feature as temp_feature, 
		device_configuration.temp_value as temp_value, 
		device_configuration.anti_spoofing as anti_spoofing, 
		device_configuration.matching_mode as matching_mode
	FROM devices
	LEFT JOIN device_configuration
	ON devices.device_uuid = device_configuration.device_uuid
	WHERE devices.device_state <> 5
	ORDER BY devices.id desc
	LIMIT ?,?;`

	DEVICE_CountGetAllQuery = `SELECT COUNT(*)
	FROM devices
	WHERE devices.device_state <> 5;`

	DEVICE_GetAllByCompanyCodeQuery = `SELECT 
		devices.id as id,
		devices.company_code as company_code,
		devices.group_id as group_id,
		devices.device_uuid as device_uuid,
		devices.device_name as device_name,
		devices.device_app_version as device_app_version,
		devices.device_description as device_description,
		devices.location_code as location_code,
		devices.device_type as device_type,
		devices.device_state as device_state,
		devices.approved_user_id as approved_user_id,
		devices.last_modified as last_modified,
		device_configuration.mask_feature as mask_feature, 
		device_configuration.temp_feature as temp_feature, 
		device_configuration.temp_value as temp_value, 
		device_configuration.anti_spoofing as anti_spoofing, 
		device_configuration.matching_mode as matching_mode
	FROM devices
	LEFT JOIN device_configuration
	ON devices.device_uuid = device_configuration.device_uuid
	WHERE devices.company_code = ? AND devices.device_state <> 5
	ORDER BY devices.id desc
	LIMIT ?,?;`

	DEVICE_CountGetAllByCompanyCodeQuery = `SELECT COUNT(*) FROM devices WHERE company_code = ? AND device_state <> 5;`

	DEVICE_GetAllByCompanyLocationCodeQuery = `SELECT 
		devices.id as id,
		devices.company_code as company_code,
		devices.group_id as group_id,
		devices.device_uuid as device_uuid,
		devices.device_name as device_name,
		devices.device_app_version as device_app_version,
		devices.device_description as device_description,
		devices.location_code as location_code,
		devices.device_type as device_type,
		devices.device_state as device_state,
		devices.approved_user_id as approved_user_id,
		devices.last_modified as last_modified,
		device_configuration.mask_feature as mask_feature, 
		device_configuration.temp_feature as temp_feature, 
		device_configuration.temp_value as temp_value, 
		device_configuration.anti_spoofing as anti_spoofing, 
		device_configuration.matching_mode as matching_mode
	FROM devices
	LEFT JOIN device_configuration
	ON devices.device_uuid = device_configuration.device_uuid
	WHERE devices.company_code = ? AND devices.location_code = ? AND devices.device_state <> 5
	ORDER BY devices.id desc
	LIMIT ?,?;`

	DEVICE_CountGetAllByCompanyLocationCodeQuery = `SELECT Count(*)
	FROM devices
	LEFT JOIN device_configuration
	ON devices.device_uuid = device_configuration.device_uuid
	WHERE devices.company_code = ? AND devices.location_code = ? AND devices.device_state <> 5 ;`

	SubscibeAck = `call sp_device_subscribeack(?,?,?,?,?);`

	GetSubscribers = `call sp_device_GetSubscribers(?, ?, ?)`

	UpdateDeviceVersionQuery = `UPDATE devices  
		SET device_app_version = ?, last_modified = ?
		WHERE device_uuid = ?;`
)
