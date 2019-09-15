package sync

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type mysqlGormColumn struct {
	Name string
	Type string
}

type SyncUtilGorm struct {
	Db *sql.DB
}

// CheckTable CheckTable is exists and have columns that ready to filled with Field from struct
func (s SyncUtilGorm) CheckTableGorm(entity interface{}) error {

	entityType := reflect.TypeOf(entity)

	columns, err := s.getDbColumnsGorm(entityType)
	if err != nil {
		return err
	}

	fields, err := getFieldsGorm(entityType)
	if err != nil {
		return err
	}

	if len(columns) != len(fields) {
		errorMsg := "Different fields between Struct and DB"
		errorMsg += " Fields [" + joinFieldsToString(fields) + "]"
		errorMsg += " Column [" + joinMysqlGormColumnsToString(columns) + "]"
		return errors.New(errorMsg)
	}

	for i := range fields {
		column := columns[i]
		field := fields[i]
		err = isSameGorm(column, field, i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s SyncUtilGorm) getDbColumnsGorm(entity reflect.Type) ([]mysqlGormColumn, error) {

	var err error
	var result []mysqlGormColumn

	query := queryDescribeGorm(entity)

	rows, err := s.Db.Query(query)
	if err != nil {
		return result, err
	}

	var unusedColumn1, unusedColumn2, unusedColumn4 string
	var unusedColumn3 interface{}

	for rows.Next() {
		item := mysqlGormColumn{}
		err = rows.Scan(&item.Name, &item.Type, &unusedColumn1, &unusedColumn2, &unusedColumn3, &unusedColumn4)
		if err != nil {
			return result, err
		}
		result = append(result, item)

	}

	return result, nil
}

func queryDescribeGorm(entity reflect.Type) string {

	tableName := ToSnakeCase(entity.Name()) + "s"

	return fmt.Sprintf("DESCRIBE %s;", tableName)
}

func getFieldsGorm(entity reflect.Type) ([]reflect.StructField, error) {
	var result []reflect.StructField

	if entity.Kind() != reflect.Struct {
		errorMessage := fmt.Sprintf("\"%s\" is not a struct. Do you really want to add DB for this type?", entity.Name())
		return result, errors.New(errorMessage)
	}

	numFieldsNotStruct := entity.NumField()
	if numFieldsNotStruct == 0 {
		errorMessage := fmt.Sprintf("\"%s\" doesnt have field. Do you really want to add DB for this type?", entity.Name())
		return result, errors.New(errorMessage)
	}

	additionNumFields := 0
	for i := 0; i < numFieldsNotStruct; i++ {
		// only support one nested only
		if entity.Field(i).Type.Kind() == reflect.Struct {
			additionNumFields += entity.Field(i).Type.NumField() - 1
		}
	}

	numFields := numFieldsNotStruct + additionNumFields

	result = make([]reflect.StructField, numFields)
	k := 0
	for i := 0; i < numFieldsNotStruct; i++ {
		// only support one nested only
		if entity.Field(i).Type.Kind() == reflect.Struct {
			for j := 0; j < entity.Field(i).Type.NumField(); j++ {
				result[k] = entity.Field(i).Type.Field(j)
				k += 1
			}
		} else {
			result[k] = entity.Field(i)
			k += 1
		}
	}

	return result, nil
}

func isSameGorm(column mysqlGormColumn, field reflect.StructField, index int) error {

	fieldName := ToSnakeCase(field.Name)

	// skip
	if column.Name == "id" {
		return nil
	}

	if column.Name == "created_at" {
		return nil
	}

	if column.Name == "updated_at" {
		return nil
	}

	if column.Name == "deleted_at" {
		return nil
	}

	if column.Name != fieldName {
		errorMsg := fmt.Sprintf("Different name in column number \"%d\". DB : %s, Struct : %s\n", index, column.Name, fieldName)
		return errors.New(errorMsg)
	}

	fieldType := field.Tag.Get("gorm")

	if strings.Contains(fieldType, "type") {
		columnTypeToCheck := "type:" + column.Type
		if fieldType != "" && columnTypeToCheck != fieldType {
			errorMsg := fmt.Sprintf("Different type in column \"%s\". DB : %s, Struct : %s\n", fieldName, columnTypeToCheck, fieldType)
			return errors.New(errorMsg)
		}
	} else if field.Type.Kind() == reflect.String {
		if !strings.Contains(column.Type, "varchar") {
			errorMsg := fmt.Sprintf("Different type in column \"%s\". DB : %s, Struct : %s\n", fieldName, column.Type, "varchar")
			return errors.New(errorMsg)
		}
	} else if field.Type.Kind() == reflect.Uint64 {
		if column.Type != "bigint(20) unsigned" {
			errorMsg := fmt.Sprintf("Different type integer in column \"%s\". DB : %s, Struct : %s\n", fieldName, column.Type, reflect.Uint64)
			return errors.New(errorMsg)
		}
	} else if field.Type.Kind() == reflect.Float64 {
		if column.Type != "double" {
			errorMsg := fmt.Sprintf("Different type integer in column \"%s\". DB : %s, Struct : %s\n", fieldName, column.Type, reflect.Uint64)
			return errors.New(errorMsg)
		}
	} else {
		errorMsg := fmt.Sprintf("Cannot decide type in column \"%s\". DB : %s, Struct : %s\n", fieldName, column.Type, field.Type)
		return errors.New(errorMsg)
	}

	return nil
}

func joinMysqlGormColumnsToString(items []mysqlGormColumn) string {
	if len(items) == 0 {
		return ""
	}

	result := items[0].Name
	for i := 1; i < len(items); i++ {
		result += ", " + items[i].Name
	}

	return result
}
