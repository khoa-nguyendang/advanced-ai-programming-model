package utilities

import (
	"aapi/shared/constants"
	"errors"

	"github.com/go-sql-driver/mysql"
)

func GetMySqlErrorEnum(err error) constants.MySqlErrorEnum {
	if err != nil {
		var mysqlErr *mysql.MySQLError = &mysql.MySQLError{}
		if errors.As(err, &mysqlErr) && mysqlErr != nil {
			switch mysqlErr.Number {
			case 1037:
				return constants.ER_OUTOFMEMORY
			case 1040:
				return constants.ER_CON_COUNT_ERROR
			case 1046:
				return constants.ER_NO_DB_ERROR
			case 1047:
				return constants.ER_UNKNOWN_COM_ERROR
			case 1048:
				return constants.ER_BAD_NULL_ERROR
			case 1053:
				return constants.ER_SERVER_SHUTDOWN
			case 1062:
				return constants.ER_DUP_ENTRY
			case 1072:
				return constants.ER_KEY_COLUMN_DOES_NOT_EXITS
			case 1146:
				return constants.ER_NO_SUCH_TABLE

			case 1149:
				return constants.ER_SYNTAX_ERROR
			case 1758:
				return constants.ER_DA_INVALID_CONDITION_NUMBER
			default:
				return constants.UNKNOWN_ERR
			}

		}

	}
	return constants.UNKNOWN_ERR
}
