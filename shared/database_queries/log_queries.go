package dbqueries

const (
	SaveLogQuery = `INSERT INTO logs(app_user_id, activity, device_uuid, full_message, date) VALUES (?, ?, ?, ?, ?);`

	FindLogQuery = `SELECT logs.id as id, app_user_id, activity, device_uuid, full_message, date, company_code, user_id, user_name, user_state
	FROM logs
	LEFT JOIN users ON app_user_id = users.id
	WHERE date >= ? AND date <= ? ORDER BY id desc;`

	FindLogQueryByDeviceUuid = `SELECT logs.id as id, app_user_id, activity, device_uuid, full_message, date, company_code, user_id, user_name, user_state
	FROM logs
	LEFT JOIN users ON app_user_id = users.id
	WHERE date >= ? AND date <= ? AND device_uuid IN (?) ORDER BY id desc;`

	FindLogQueryByUserId = `SELECT logs.id as id, app_user_id, activity, device_uuid, full_message, date, company_code, user_id, user_name, user_state
	FROM logs
	LEFT JOIN users ON app_user_id = users.id
	WHERE date >= ? AND date <= ? AND user_id IN (?) ORDER BY id desc;`

	FindLogQueryByUserIdAndDeviceUuid = `SELECT logs.id as id, app_user_id, activity, device_uuid, full_message, date, company_code, user_id, user_name, user_state
	FROM logs
	LEFT JOIN users ON app_user_id = users.id
	WHERE date >= ? AND date <= ? AND user_id IN (?) AND device_uuid IN (?) ORDER BY id desc;`

	//With company code
	FindLogWithCompanyCodeQuery = `SELECT logs.id as id, app_user_id, activity, device_uuid, full_message, date, company_code, user_id, user_name, user_state
	FROM logs
	LEFT JOIN users ON app_user_id = users.id
	WHERE company_code = ? AND date >= ? AND date <= ? ORDER BY id desc;`

	FindLogWithCompanyCodeQueryByDeviceUuid = `SELECT logs.id as id, app_user_id, activity, device_uuid, full_message, date, company_code, user_id, user_name, user_state
	FROM logs
	LEFT JOIN users ON app_user_id = users.id
	WHERE company_code = ? AND date >= ? AND date <= ? AND device_uuid IN (?) ORDER BY id desc;`

	FindLogWithCompanyCodeQueryByUserId = `SELECT logs.id as id, app_user_id, activity, device_uuid, full_message, date, company_code, user_id, user_name, user_state
	FROM logs
	LEFT JOIN users ON app_user_id = users.id
	WHERE company_code = ? AND date >= ? AND date <= ? AND user_id IN (?) ORDER BY id desc;`

	FindLogWithCompanyCodeQueryByUserIdAndDeviceUuid = `SELECT logs.id as id, app_user_id, activity, device_uuid, full_message, date, company_code, user_id, user_name, user_state
	FROM logs
	LEFT JOIN users ON app_user_id = users.id
	WHERE company_code = ? AND date >= ? AND date <= ? AND user_id IN (?) AND device_uuid IN (?) ORDER BY id desc;`

	LOG_GetAllBy_CompanyCode_DeviceUUID_UserId = `SELECT logs.id as id, activity, device_uuid, date, company_code, user_id, user_name, user_state
	FROM logs
	LEFT JOIN users ON app_user_id = users.id
	WHERE company_code = ? AND user_id = ? AND device_uuid  = ?  AND date >= ? AND date <= ? 
	ORDER BY id desc
	LIMIT ?,?;`

	LOG_CountGetAllBy_CompanyCode_DeviceUUID_UserId = `SELECT COUNT(*)
	FROM logs
	LEFT JOIN users ON app_user_id = users.id
	WHERE users.company_code = ?  AND user_id = ? AND device_uuid  = ? AND date >= ? AND date <= ? ;`

	LOG_GetAllBy_CompanyCode_DeviceUUID = `SELECT logs.id as id, activity, device_uuid, date, company_code, user_id, user_name, user_state
	FROM logs
	LEFT JOIN users ON app_user_id = users.id
	WHERE company_code = ? AND device_uuid = ?  AND date >= ? AND date <= ? 
	ORDER BY id desc
	LIMIT ?,?;`

	LOG_CountGetAllBy_CompanyCode_DeviceUUID = `SELECT COUNT(*)
	FROM logs
	LEFT JOIN users ON app_user_id = users.id
	WHERE users.company_code = ? AND device_uuid = ? AND date >= ? AND date <= ? ;`

	LOG_GetAllBy_CompanyCode_UserId = `SELECT logs.id as id, activity, device_uuid, date, company_code, user_id, user_name, user_state
	FROM logs
	LEFT JOIN users ON app_user_id = users.id
	WHERE users.company_code = ? AND user_id = ?  AND date >= ? AND date <= ? 
	ORDER BY id desc
	LIMIT ?,?;`

	LOG_CountGetAllBy_CompanyCode_UserId = `SELECT COUNT(*)
	FROM logs
	LEFT JOIN users ON app_user_id = users.id
	WHERE users.company_code = ? AND user_id = ? AND date >= ? AND date <= ? ;`

	LOG_GetAllBy_CompanyCode = `SELECT logs.id as id, activity, device_uuid, date, company_code, user_id, user_name, user_state
	FROM logs
	LEFT JOIN users ON app_user_id = users.id
	WHERE users.company_code = ? AND date >= ? AND date <= ? 
	ORDER BY id desc
	LIMIT ?,?;`

	LOG_CountGetAllBy_CompanyCode = `SELECT COUNT(*)
	FROM logs
	LEFT JOIN users ON app_user_id = users.id
	WHERE users.company_code = ? AND date >= ? AND date <= ? ;`
)
