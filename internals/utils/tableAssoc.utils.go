package utils

import (
	"fmt"
	"strings"
)

func InsertTableAssoc(table, column1, column2 string, groups []any) string {
	query := "INSERT INTO " + table + " (" + column1 + ", " + column2 + ") VALUES "

	values := make([]string, len(groups))
	for i := range groups {
		// $1 untuk column1, lalu $2, $3, $4, dst untuk column2 dari tiap group
		if i == len(groups) - 1 {
			break
		}
		values[i] = fmt.Sprintf("($1, $%d)", i+2)
	}

	query += strings.Join(values, ", ")
	query = strings.TrimSuffix(query, ", ")
	return query
}


// func InsertTableAssoc(table, column1, column2 string, groups []any) (string) {
// 	query:= "INSERT INTO " + table + " (" + column1 + ", " + column2 + ") " + "VALUES"
// 	for idx := range groups {
// 		if idx == len(groups) - 2 {
// 			query += fmt.Sprintf("($1, $%d)", idx + 2)
// 			break
// 		}

// 		if idx < len(groups) - 2 {
// 			query += fmt.Sprintf("($1, $%d), ", idx + 2)
// 		}
// 	}

// 	log.Println("{DEBUG QUERY}", query)

// 	return query
// }