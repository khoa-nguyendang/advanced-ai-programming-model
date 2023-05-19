package dbqueries

const (
	CheckLocationExistQuery = `SELECT location_code FROM locations WHERE location_code = ?;`

	AddLocationQuery = `INSERT INTO locations(location_code, pic_id, type, name, description, state, company_code)
	VALUES (?, ?, ?, ?, ?, ?, ?);`

	UpdateLocationQuery = `UPDATE locations
	SET pic_id = ?, type = ?, name = ?, description = ?, state = ?, company_code = ?
	WHERE location_code = ?;`

	FindLocationQueryByIds = `SELECT location_code, pic_id, type, name, description, state, company_code
	FROM locations
	WHERE company_code = ? AND location_code IN (?);`

	FindAllLocationQuery = `SELECT location_code, pic_id, type, name, description, state, company_code
	FROM locations
	WHERE company_code = ?;`
)
