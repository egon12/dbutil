package dbutil

import (
	"fmt"
	"reflect"
)

/*
generateSelectQuery

generate string "SELECT column1, column2 FROM tableName"
*/
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

/*
generateSelectQueryWithout

generate string "SELECT column1 FROM tableName" but remove the filtered column
*/
func generateSelectQueryWithout(entity reflect.Type, without []string) string {
	tableName := ToSnakeCase(entity.Name())
	fieldNames := readEntityFieldsName(entity)
	filteredFieldNames := filterFieldNames(fieldNames, without)
	columnsSelect := generateColumnNames(filteredFieldNames)
	return fmt.Sprintf("SELECT %s FROM %s LIMIT $1 OFFSET $2", columnsSelect, tableName)
}

// generateInsertQuery query generator
// make sure that entity that got into these functions is a struct type
func generateInsertQuery(entity reflect.Type) string {
	tableName := ToSnakeCase(entity.Name())
	fieldNames := readEntityFieldsName(entity)

	filteredFieldNames := filterFieldNamesWithFunc(
		fieldNames,
		func(n string, i int) bool { return i != 0 },
	)
	columnsSelect := generateColumnNames(filteredFieldNames)

	fillsArray := make([]string, len(filteredFieldNames))
	for i := 1; i < entity.NumField(); i++ {
		fillsArray[i-1] = fmt.Sprintf("$%d", i)
	}
	fills := generateColumnNames(fillsArray)

	return fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s) RETURNING id", tableName, columnsSelect, fills)
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

func stringArrayContains(stack []string, needle string) bool {
	for _, s := range stack {
		if s == needle {
			return true
		}
	}
	return false
}

func readEntityFieldsName(entity reflect.Type) []string {
	num := entity.NumField()
	result := make([]string, num)
	for i := 0; i < num; i++ {
		result[i] = entity.Field(i).Name
	}
	return result
}

func generateColumnNames(names []string) string {
	columnsSelect := ""
	for i, name := range names {
		name := ToSnakeCase(name)

		if len(name) < 1 {
			continue
		}

		if i == 0 {
			columnsSelect += name
		} else {
			columnsSelect += ", " + name
		}
	}
	return columnsSelect
}

func joinToString(fields []string) string {
	result := ""
	for i, field := range fields {
		if i == 0 {
			result += field
		} else {
			result += ", " + field
		}
	}
	return result
}

func filterFieldNames(fieldNames []string, without []string) []string {
	result := make([]string, len(fieldNames))
	i := 0
	for _, name := range fieldNames {
		if stringArrayContains(without, ToSnakeCase(name)) {
			continue
		}
		result[i] = name
		i += 1
	}
	return result
}

func filterFieldNamesWithFunc(
	fieldNames []string,
	filterFunc func(string, int) bool,
) []string {
	result := make([]string, len(fieldNames))
	i := 0
	for index, name := range fieldNames {
		if filterFunc(name, index) {
			result[i] = name
			i += 1
		}
	}
	return result
}
