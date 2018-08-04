package dbutil

import (
	"fmt"
	"reflect"
)

// generateInsertQuery query generator
// make sure that entity that got into these functions is a struct type
func generateInsertQuery(entity reflect.Type) string {
	tableName := ToSnakeCase(entity.Name())
	columnsSelect := ""
	fills := ""

	for i := 1; i < entity.NumField(); i++ {
		name := ToSnakeCase(entity.Field(i).Name)
		if i == 1 {
			columnsSelect += name
			fills += fmt.Sprintf("$%d", i)
		} else {
			columnsSelect += ", " + name
			fills += fmt.Sprintf(", $%d", i)
		}
	}

	return fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s) RETURNING id", tableName, columnsSelect, fills)
}

// generateSelectQuery query generator
// make sure that entity that got into these functions is a struct type
func generateSelectQuery(entity reflect.Type) string {
	tableName := ToSnakeCase(entity.Name())
	columnsSelect := ""
	for i := 0; i < entity.NumField(); i++ {
		name := ToSnakeCase(entity.Field(i).Name)
		if i == 0 {
			columnsSelect += name
		} else {
			columnsSelect += ", " + name
		}
	}

	return fmt.Sprintf("SELECT %s FROM %s LIMIT $1 OFFSET $2", columnsSelect, tableName)
}

// generateUpdateQuery query generator
// make sure that entity that got into these functions is a struct type
func generateUpdateQuery(entity reflect.Type) string {
	tableName := ToSnakeCase(entity.Name())
	sets := ""
	idFieldNumber := fmt.Sprintf("$%d", entity.NumField())

	for i := 1; i < entity.NumField(); i++ {
		name := ToSnakeCase(entity.Field(i).Name)
		if i == 1 {
			sets += fmt.Sprintf("%s = $%d", name, i)
		} else {
			sets += fmt.Sprintf(", %s = $%d", name, i)
		}
	}

	return fmt.Sprintf("UPDATE %s SET %s WHERE id = %s", tableName, sets, idFieldNumber)
}

// generateDeleteQuery query generator
func generateDeleteQuery(entity reflect.Type) string {
	tableName := ToSnakeCase(entity.Name())
	return fmt.Sprintf("DELETE FROM %s WHERE id = $1", tableName)
}
