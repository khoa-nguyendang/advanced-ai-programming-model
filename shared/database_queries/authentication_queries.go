package dbqueries

const (
	FindDeviceByUuidQuery = `SELECT * FROM devices WHERE device_uuid = ? AND device_state <> ?;`

	AddDeviceTokenQuery = `INSERT INTO devices_token(device_id, 
		token, 
		token_expiry, 
		refresh_token, 
		refresh_token_expiry, 
		action, 
		action_state, 
		executed_ip, 
		timestamp)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`

	AddAdministratorTokenQuery = `INSERT INTO administrators_token(administrator_user_id,
		token, 
		token_expiry, 
		refresh_token, 
		refresh_token_expiry, 
		action, 
		action_state, 
		executed_ip, 
		timestamp)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`

	FindDeviceTokenByDeviceIdQuery = `SELECT 
		dt.id,
		device_id,
		token,
		refresh_token,
		refresh_token_expiry,
		action,
		action_state,
		executed_ip
	FROM devices_token dt JOIN devices d ON d.Id = dt.device_id
	WHERE d.device_uuid = ? AND dt.refresh_token = ?
	ORDER BY timestamp desc LIMIT 1;`

	FindAdministratorTokenByRefreshTokenQuery = `SELECT
		at.id,
		administrator_user_id,
		token,
		token_expiry,
		refresh_token,
		refresh_token_expiry,
		action,
		action_state,
		executed_ip,
		timestamp 
    FROM administrators_token at JOIN administrators s ON s.Id = at.administrator_user_id
    WHERE s.state <> ? AND at.refresh_token = ?
    ORDER BY timestamp desc LIMIT 1;`

	UpdateAccessTokenQuery = `UPDATE devices_token 
	SET token = ?, token_expiry = ?
	WHERE device_id = ? AND refresh_token = ?; `

	UpdateAdministratorTokenQuery = `UPDATE administrators_token 
	SET token = ?, token_expiry = ?
	WHERE administrator_user_id = ? AND refresh_token = ?; `

	FindAdministratorByEmployeeIdQuery = `SELECT * FROM administrators  WHERE employee_id = ?;`

	FindAdministratorByIdQuery = `SELECT * FROM administrators  WHERE state <> ? AND id = ?;`

	FindAdministratorByUsernameOrEmailQuery = `SELECT *
	FROM administrators
	WHERE state <> ? AND username = ?;`
)
