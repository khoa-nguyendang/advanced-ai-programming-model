package dbqueries

const (
	CheckCompanyCodeConfigExistQuery = `SELECT company_code FROM company_configuration WHERE company_code = ?;`
)
