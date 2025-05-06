package utils

import "fmt"

func InsertTableAssoc(table, column1, column2 string, groups []any) (string) {
	query:= "INSERT INTO " + table + " (" + column1 + ", " + column2 + ") " + "VALUES"
	for idx := range groups {
		if idx == len(groups) - 2 {
			query += fmt.Sprintf("($1, $%d)", idx + 2)
			break
		}

		if idx < len(groups) - 2 {
			query += fmt.Sprintf("($1, $%d), ", idx + 2)
		}
	}

	return query
}