package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

// Db place to inject the Db..
// To use it, we need to inject the Db first
// for ex:
//
// ```
// utils.Db = db
// utils.CheckTable(Entity{})
// ```
var Db *sql.DB

func InteractiveSync(entity interface{}) {
	err := CheckTable(entity)
	if err != nil {
		fmt.Print("Do you want to drop and recreate table entity? [y/n] :")
		answer := ""
		fmt.Scanln(&answer)
		if answer == "y" {
			DropTable(entity)

			err2 := CreateTable(entity)
			if err2 != nil {
				fmt.Println(err2)
			}
		}
	}
}

// CheckTable CheckTable is exists and have columns that ready to filled with Field from struct
func CheckTable(entity interface{}) error {

	entityType := reflect.TypeOf(entity)

	columns, err := getDbColumns(entityType)
	if err != nil {
		return err
	}

	fields, err := getFields(entityType)
	if err != nil {
		return err
	}

	if len(columns) != len(fields) {
		errorMsg := "Different fields between Struct and DB"
		errorMsg += " Fields [" + joinFieldsToString(fields) + "]"
		errorMsg += " Column [" + joinColumnsToString(columns) + "]"
		return errors.New(errorMsg)
	}

	for i := range fields {
		column := columns[i]
		field := fields[i]
		err = isSame(column, field, i)
		if err != nil {
			return err
		}
	}

	return nil
}

// DropTable just drop the table
func DropTable(entity interface{}) error {

	entityType := reflect.TypeOf(entity)

	tableName := ToSnakeCase(entityType.Name())

	query := fmt.Sprintf("DROP TABLE %s;", tableName)

	_, err := Db.Exec(query)

	return err
}

// CreateTable create table with columns as same as field
func CreateTable(entity interface{}) error {
	entityType := reflect.TypeOf(entity)

	columns := []string{}
	for i := 0; i < entityType.NumField(); i++ {
		entityField := entityType.Field(i)
		column := processField(entityField)
		columns = append(columns, column)
	}

	tableName := ToSnakeCase(entityType.Name())

	query := getCreateSQL(tableName, columns)

	_, err := Db.Exec(query)

	if err != nil {
		return errors.New(err.Error() + " query: " + query)
	}

	return nil
}

func getCreateSQL(tableName string, columns []string) string {

	columnsCombined := ""

	for i, c := range columns {
		if i == 0 {
			columnsCombined += c
		} else {
			columnsCombined += ", " + c
		}
	}

	return fmt.Sprintf("CREATE TABLE %s (%s);", tableName, columnsCombined)
}

func processField(field reflect.StructField) string {

	if field.Name == "ID" {
		return processPrimaryKey(field)
	}

	switch field.Type.Kind() {
	case reflect.String:
		return processString(field)
	case reflect.Int32:
		return processInteger32(field)
	case reflect.Int64:
		return processInteger64(field)

	}

	return ""
}

func processPrimaryKey(field reflect.StructField) string {
	return "id SERIAL"
}

func processInteger32(field reflect.StructField) string {
	return fmt.Sprintf("%s INT", ToSnakeCase(field.Name))
}

func processInteger64(field reflect.StructField) string {
	return fmt.Sprintf("%s INT8", ToSnakeCase(field.Name))
}

func processString(field reflect.StructField) string {
	length := "255"

	newlength, ok := field.Tag.Lookup("varchar")
	if ok {
		length = newlength
	}

	return fmt.Sprintf("%s VARCHAR(%s)", ToSnakeCase(field.Name), length)
}

func getDbColumns(entity reflect.Type) ([]*sql.ColumnType, error) {

	var err error
	var result []*sql.ColumnType

	query := querySelect(entity)

	rows, err := Db.Query(query)
	if err != nil {
		return result, err
	}

	result, err = rows.ColumnTypes()
	if err != nil {
		return result, err
	}

	return result, nil
}

func querySelect(entity reflect.Type) string {

	tableName := ToSnakeCase(entity.Name())

	return fmt.Sprintf("SELECT * FROM %s LIMIT 1;", tableName)
}

func getFields(entity reflect.Type) ([]reflect.StructField, error) {
	var result []reflect.StructField

	if entity.Kind() != reflect.Struct {
		errorMessage := fmt.Sprintf("\"%s\" is not a struct. Do you really want to add DB for this type?", entity.Name())
		return result, errors.New(errorMessage)
	}

	numFields := entity.NumField()
	if numFields == 0 {
		errorMessage := fmt.Sprintf("\"%s\" doesnt have field. Do you really want to add DB for this type?", entity.Name())
		return result, errors.New(errorMessage)
	}

	result = make([]reflect.StructField, numFields)
	for i := 0; i < numFields; i++ {
		result[i] = entity.Field(i)
	}

	return result, nil
}

func isSame(column *sql.ColumnType, field reflect.StructField, index int) error {

	fieldName := ToSnakeCase(field.Name)

	if column.Name() != fieldName {
		errorMsg := fmt.Sprintf("Different name in column number %d. DB : %s, Struct : %s\n", index, column.Name(), fieldName)
		return errors.New(errorMsg)
	}

	if column.ScanType() != field.Type {
		errorMsg := fmt.Sprintf("Different type in column %s. DB : %s, Struct : %s\n", fieldName, column.ScanType(), field.Type)
		return errors.New(errorMsg)
	}

	if field.Type.Kind() == reflect.String {
		length, ok := field.Tag.Lookup("varchar")
		if !ok {
			length = "255"
		}

		columnLengthInt, ok := column.Length()
		if !ok {
			errorMsg := fmt.Sprintf("Driver doesn't support get column length. Please update the driver!")
			return errors.New(errorMsg)
		}

		columnLength := fmt.Sprintf("%d", columnLengthInt)
		if length != columnLength {
			errorMsg := fmt.Sprintf("Different length in column %s. DB : %s, Struct : %s\n", fieldName, columnLength, length)
			return errors.New(errorMsg)
		}
	}

	return nil
}

func joinColumnsToString(columns []*sql.ColumnType) string {
	result := ""
	for i, column := range columns {
		if i == 0 {
			result += column.Name()
		} else {
			result += ", " + column.Name()
		}
	}
	return result
}

func joinFieldsToString(fields []reflect.StructField) string {
	result := ""
	for i, field := range fields {
		if i == 0 {
			result += field.Name
		} else {
			result += ", " + field.Name
		}
	}
	return result
}
