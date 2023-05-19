package dbqueries

const (
	AddAdministratorQuery = `INSERT INTO administrators( 
		company_code,
		administrator_id,
		username,
		password,
		full_name,
		phone_number,
		email,
		administrator_info,
		state,
		created_by,
		created_at,
		last_modified_by,
		last_modified_at,
		role_id,
		reference_id,
		salt)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	UpdateAdministratorQuery = `UPDATE administrators
	SET full_name = ?,
		phone_number = ?,
		email = ?,
		administrator_info = ?,
		state = ?,
		last_modified_by = ?,
		last_modified_at = ?,
		role_id = ?,
		reference_id = ?
	WHERE state <> ? AND username = ?;`

	DeleteAdministratorQuery = `UPDATE administrators SET state = ?, last_modified_by = ?, last_modified_at = ? WHERE username = ?;`

	IsAdministratorExistsQuery = `SELECT id FROM administrators WHERE state <> ? AND (administrator_id = ? OR username = ?);`

	FindAdministratorPasswordQuery = `SELECT id, password, salt FROM administrators WHERE state <> ? AND username = ?;`

	UpdateAdministratorPasswordQuery = `UPDATE administrators set password = ? WHERE username = ?;`

	FindAdministratorByUsernameQuery = `SELECT
		company_code,
		administrator_id,
		username,
		full_name,
		phone_number,
		email,
		administrator_info,
		state,
		created_by,
		created_at,
		last_modified_by,
		last_modified_at,
		role_id,
		reference_id
	FROM administrators
	WHERE state <> ? AND username = ?;`

	FindAdministratorByUsernamesQuery = `SELECT
		company_code,
		administrator_id,
		username,
		full_name,
		phone_number,
		email,
		administrator_info,
		state,
		created_by,
		created_at,
		last_modified_by,
		last_modified_at,
		role_id,
		reference_id
	FROM administrators
	WHERE state <> ? AND username IN (?) ORDER BY id desc LIMIT ?, ?;`

	FindAdministratorByUsernamesWithCompanyCodeQuery = `SELECT
	company_code,
	administrator_id,
	username,
	full_name,
	phone_number,
	email,
	administrator_info,
	state,
	created_by,
	created_at,
	last_modified_by,
	last_modified_at,
	role_id,
	reference_id
	FROM administrators
	WHERE state <> ? AND company_code = ? AND username IN (?) ORDER BY id desc LIMIT ?, ?;`

	FindAllAdministratorQuery = `SELECT
		company_code,
		administrator_id,
		username,
		full_name,
		phone_number,
		email,
		administrator_info,
		state,
		created_by,
		created_at,
		last_modified_by,
		last_modified_at,
		role_id,
		reference_id
	FROM administrators
	WHERE state <> ? ORDER BY id desc LIMIT ?, ?;`
	CountAllAdministratorQuery = `SELECT COUNT(*) FROM administrators WHERE state <> ?`

	FindAllAdministratorWithCompanyCodeQuery = `SELECT
	company_code,
	administrator_id,
	username,
	full_name,
	phone_number,
	email,
	administrator_info,
	state,
	created_by,
	created_at,
	last_modified_by,
	last_modified_at,
	role_id,
	reference_id
	FROM administrators
	WHERE state <> ? AND company_code = ? ORDER BY id desc LIMIT ?, ?;`
	CountAllAdministratorWithCompanyCodeQuery = `SELECT COUNT(*) FROM administrators WHERE state <> ? AND company_code = ?`

	FindPermissionNamesByRoleId = `SELECT permission.name
	FROM role_permission
	LEFT JOIN permission 
	ON role_permission.permission_id = permission.id
	WHERE role_id = ?;`
)
