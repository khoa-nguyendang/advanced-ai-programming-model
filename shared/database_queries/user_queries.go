package dbqueries

const (
	//User
	EnrollUserQuery = `INSERT INTO users(company_code, user_id, user_name, user_role_id, user_info, user_state, thumbnail_image_url, last_modified, issued_date, activation_date, expiry_date, reference_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	InsertEnollmentImage = `INSERT INTO user_image (user_column_id, path, image_type_id, created_at, last_modified, created_user_id, last_modified_user_id)
	VALUES(?, ?, ?, ?, ?, ?, ?);`

	UpdateUserByUserIdQuery = `UPDATE users set user_name = ?, user_role_id = ?, user_info = ?, user_state = ?, last_modified = ?, activation_date =?, expiry_date = ?, reference_id = ? WHERE user_id = ? AND company_code = ? AND user_state <> ?`

	FindUserByIdQuery = `SELECT company_code, user_id, user_name, user_role_id, user_info, user_state, thumbnail_image_url, last_modified, issued_date, activation_date, expiry_date, reference_id FROM users WHERE id = ? AND user_state <> ?`

	FindUserByUserIDQuery = `SELECT id, company_code, user_id, user_name, user_role_id, user_info, user_state, thumbnail_image_url, last_modified, issued_date, activation_date, expiry_date, reference_id FROM users WHERE user_id = ? AND user_state <> ?`

	FindUserWithCompanyCodeByUserIDQuery = `SELECT id, company_code, user_id, user_name, user_role_id, user_info, user_state, thumbnail_image_url, last_modified, issued_date, activation_date, expiry_date, reference_id FROM users WHERE company_code = ? AND user_id = ? AND user_state <> ?`

	FindUsersQueryByUserIds = `SELECT id, company_code, user_id, user_name, user_role_id, user_info, user_state, thumbnail_image_url, last_modified, issued_date, activation_date, expiry_date, reference_id
	FROM users
	WHERE user_id IN (?) AND user_state <> ? ORDER BY id desc LIMIT ?, ?;`
	CountUserByUserIdsQuery = `SELECT COUNT(*) FROM users WHERE user_id IN (?) AND user_state <> ?`

	FindUsersQueryByRoleIdAndUserIds = `SELECT id, company_code, user_id, user_name, user_role_id, user_info, user_state, thumbnail_image_url, last_modified, issued_date, activation_date, expiry_date, reference_id
	FROM users
	WHERE user_role_id = ? AND user_id IN (?) AND user_state <> ? ORDER BY id desc LIMIT ?, ?;`
	CountUserByRoleIdAndUserIdsQuery = `SELECT COUNT(*) FROM users WHERE user_role_id = ? AND user_id IN (?) AND user_state <> ?`

	FindUsersWithCompanyCodeQueryByUserIds = `SELECT id, company_code, user_id, user_name, user_role_id, user_info, user_state, thumbnail_image_url, last_modified, issued_date, activation_date, expiry_date, reference_id
	FROM users
	WHERE company_code = ? AND user_id IN (?)  AND user_state <> ? ORDER BY id desc LIMIT ?, ?;`
	CountUserOfCompanyCodeByUserIdsQuery = `SELECT COUNT(*) FROM users WHERE company_code = ? AND user_id IN (?) AND user_state <> ?`

	FindUsersWithCompanyCodeQueryByRoleIdAndUserIds = `SELECT id, company_code, user_id, user_name, user_role_id, user_info, user_state, thumbnail_image_url, last_modified, issued_date, activation_date, expiry_date, reference_id
	FROM users
	WHERE company_code = ? AND user_role_id = ? AND user_id IN (?) AND user_state <> ? ORDER BY id desc LIMIT ?, ?;`
	CountUserOfCompanyCodeByRoleIdAndUserIdsQuery = `SELECT COUNT(*) FROM users WHERE company_code = ? AND user_role_id = ? AND user_id IN (?) AND user_state <> ?`

	FindAllUserQuery = `SELECT id, company_code, user_id, user_name, user_role_id, user_info, user_state, thumbnail_image_url, last_modified, issued_date, activation_date, expiry_date, reference_id FROM users WHERE user_state <> ? ORDER BY id desc LIMIT ?, ?`
	CountUserQuery   = `SELECT COUNT(*) FROM users WHERE user_state <> ?`

	FindAllUserByRoleIdQuery = `SELECT id, company_code, user_id, user_name, user_role_id, user_info, user_state, thumbnail_image_url, last_modified, issued_date, activation_date, expiry_date, reference_id FROM users WHERE user_role_id = ? AND user_state <> ? ORDER BY id desc LIMIT ?, ?`
	CountUserByRoleIdQuery   = `SELECT COUNT(*) FROM users WHERE user_role_id = ? AND user_state <> ?`

	FindAllUserWithCompanyCodeQuery = `SELECT id, company_code, user_id, user_name, user_role_id, user_info, user_state, thumbnail_image_url, last_modified, issued_date, activation_date, expiry_date, reference_id FROM users WHERE company_code = ? AND user_state <> ? ORDER BY id desc LIMIT ?, ?`
	CountUserOfCompanyCodeQuery     = `SELECT COUNT(*) FROM users WHERE company_code = ? AND user_state <> ?`

	FindAllUserByRoleIdWithCompanyCodeQuery = `SELECT id, company_code, user_id, user_name, user_role_id, user_info, user_state, thumbnail_image_url, last_modified, issued_date, activation_date, expiry_date, reference_id FROM users WHERE company_code = ? AND user_role_id = ? AND user_state <> ? ORDER BY id desc LIMIT ?, ?`

	CountUserOfCompanyCodeAndRoleQuery = `SELECT COUNT(*) FROM users WHERE company_code = ? AND user_role_id = ? AND user_state <> ?`

	DeleteUserQuery = `UPDATE users set user_state = ?, last_modified = ? WHERE user_id = ?`

	DeleteUserWithCompanyCodesQuery = `UPDATE users set user_state = ?, last_modified = ? WHERE company_code = ? AND user_id = ?`

	GetGroups            = "SELECT * FROM `app`.`groups`"
	GetGroupsByCompanyId = "SELECT * FROM `app`.`groups` WHERE company_id = ?"

	InsertUserGroups            = "INSERT INTO user_group (user_id, group_id, enabled, created_at, last_modified, created_user_id, last_modified_user_id) VALUES(?, ?, ?, ?, ?, ?, ?);"
	DeleteUserGroupsByUserColId = "DELETE FROM user_group WHERE user_id = ?;"
	GetUserGroupsByUserColIds   = `SELECT user_id, group_id, app.groups.name
	FROM user_group
	LEFT JOIN  app.groups
	ON group_id = app.groups.id
	WHERE user_id IN (?)`

	GetRoles    = "SELECT * FROM `app`.`roles`"
	GetRoleById = `SELECT id, name, role_type FROM app.roles WHERE id = ?;`

	GetEnrollmentImagePath = `SELECT user_column_id, path
	FROM user_image
	WHERE user_column_id IN (?);`

	FindUserByStateInActivationDuration  = `SELECT id, company_code, user_id, user_state FROM users WHERE user_role_id = ? AND user_state = ? AND activation_date >= ? AND activation_date < ? AND expiry_date > ?;`
	FindUserByStateOutActivationDuration = `SELECT id, company_code, user_id, user_state FROM users WHERE user_role_id = ? AND user_state = ? AND (activation_date >= ? OR expiry_date  < ?);`

	UpdateUserStateByColumnIds = `UPDATE users set user_state = ? WHERE id IN (?);`
)
