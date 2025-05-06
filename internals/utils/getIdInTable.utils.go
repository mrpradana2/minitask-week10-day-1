package utils

import (
	"fmt"
	"log"
)

func GetIdTable(table, column string, groups []string) (string, []any) {
	// build dinamic query untuk mengambil seat_id dari table seats
	querySelect := "SELECT id FROM "
	querySelect += table
	querySelect += " WHERE "
	querySelect += column
	querySelect += " IN ("
	result := []any{}
	for idx, item := range groups {
		log.Println("INDEX", idx, len(groups))
		if idx < len(groups) {
			querySelect += fmt.Sprintf("$%d", idx+1)
			if idx < len(groups) - 1 {
				querySelect += ", "
			}
		}
		if idx == len(groups)-1 {
			querySelect += `)`
		}
		result = append(result, item)
	}
	return querySelect, result
}