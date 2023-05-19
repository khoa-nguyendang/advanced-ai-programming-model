package constants

// https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html
type MySqlErrorEnum int

const (
	NO_ERR MySqlErrorEnum = iota

	UNKNOWN_ERR MySqlErrorEnum = iota

	//1037 Out of memory; restart server and try again (needed %d bytes)
	ER_OUTOFMEMORY

	//1040 Too many connections
	ER_CON_COUNT_ERROR

	//1046 No database selected
	ER_NO_DB_ERROR

	//1047 Unknown command
	ER_UNKNOWN_COM_ERROR

	//1048  Column '%s' cannot be null
	ER_BAD_NULL_ERROR

	//1053
	ER_SERVER_SHUTDOWN

	//1062  Duplicate entry '%s' for key %d
	ER_DUP_ENTRY

	//1072 Key column '%s' doesn't exist in table
	ER_KEY_COLUMN_DOES_NOT_EXITS

	//1103 Incorrect table name '%s'
	ER_WRONG_TABLE_NAME

	//1104
	ER_TOO_BIG_SELECT

	//1105 Unknown error
	ER_UNKNOWN_ERROR

	//1106 Unknown procedure '%s'
	ER_UNKNOWN_PROCEDURE

	//1146 Table '%s.%s' doesn't exist
	ER_NO_SUCH_TABLE

	//1149 wrong syntax
	ER_SYNTAX_ERROR

	//1162 Result string is longer than 'max_allowed_packet' bytes
	ER_TOO_LONG_STRING

	//1172 Result consisted of more than one row
	ER_TOO_MANY_ROWS

	//1312 PROCEDURE %s can't return a result set in the given context
	ER_SP_BADSELECT

	//1569 ALTER TABLE causes auto_increment resequencing, resulting in duplicate entry '%s' for key '%s'
	ER_DUP_ENTRY_AUTOINCREMENT_CASE

	//1758 Invalid condition number
	ER_DA_INVALID_CONDITION_NUMBER
)
