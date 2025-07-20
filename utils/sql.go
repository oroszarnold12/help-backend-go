package utils

import "database/sql"

func ConvertNullString(string sql.NullString, target *string) {
	if string.Valid {
		*target = string.String
	}
}
