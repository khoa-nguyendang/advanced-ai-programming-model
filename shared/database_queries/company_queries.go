package dbqueries

const (
	CheckCompanyCodeExistQuery        = "SELECT id, company_code FROM company WHERE company_code = ? limit 1;"
	COMPANY_GET_ONE_BY_CODE           = "SELECT * FROM company WHERE company_state > -1 AND company_code = ? limit 1;"
	COMPANY_GET_ONE_BY_ID             = "SELECT * FROM company WHERE id = ? limit 1;"
	COMPANY_GET_ONE_BY_COMPANY_CODE   = "SELECT * FROM company WHERE company_code = ? limit 1;"
	COMPANY_GET_MANY                  = "SELECT * FROM company WHERE company_state > -1"
	COMPANY_GET_MANY_BY_COMPANY_CODES = "SELECT * FROM company WHERE company_state > -1 AND company_code in (?);"
	COMPANY_DELETE_BY_CODE            = "UPDATE company SET company_state = -1, last_modified = ?, last_modified_user_id = ? WHERE company_code = ?"
	COMPANY_DELETE_BY_ID              = "UPDATE company SET company_state = -1, last_modified = ?, last_modified_user_id = ? WHERE id = ?"
	COMPANY_UPDATE                    = "UPDATE company SET name = ?, description = ?, last_modified = ?, last_modified_user_id = ?  WHERE id = ?"
	COMPANY_ADD                       = `INSERT INTO company(name, company_code, description, created_at, created_user_id, last_modified, last_modified_user_id) VALUES (?, ?, ?, ?, ?, ?, ?)`
	GROUP_GET_ONE_BY_ID               = "SELECT * FROM app.groups WHERE group_state > -1 AND id = ? limit 1;"
	GROUP_GET_MANY_BY_COMPANY_ID      = "SELECT * FROM app.groups  WHERE group_state > -1 AND company_id IN(0,?);"
	GROUP_GET_MANY_BY_COMPANY_CODE    = `SELECT app.groups.* 
		FROM app.groups 
		JOIN  company on company.id = app.groups.company_id
		WHERE group_state <> 5 AND company.company_code = ?;`
	GROUP_DELETE = "UPDATE FROM `groups` SET group_state = -1, last_modified = ?, last_modified_user_id = ?  WHERE id = ?"
	GROUP_UPDATE = "UPDATE `groups` SET name = ?, description = ?, last_modified = ?, last_modified_user_id = ?  WHERE id = ?"
	GROUP_ADD    = "INSERT INTO `groups` (name, company_id, description, created_at, created_user_id, last_modified, last_modified_user_id) VALUES (?, ?, ?, ?, ?, ?, ?)"
)
