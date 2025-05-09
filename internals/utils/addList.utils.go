package utils

import "fmt"

func AddList(table, column string, groups []string) string {
	query := "INSERT INTO " + table + " (" + column + ") " + "VALUES "
	for idx := range groups {
		if idx == len(groups) - 1 {
			query += fmt.Sprintf("($%d)", idx + 1)
		}
		if idx < len(groups) - 1 {
			query += fmt.Sprintf("($%d), ", idx + 1)
		}
	}
	query += " ON CONFLICT (name) DO NOTHING"
	return query
}

