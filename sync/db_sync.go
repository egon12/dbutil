package sync

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
)

// To Use the sync
// s := SyncUtils{db}
// s.InteractiveSync(entityExamples{})

type (
	SyncUtils struct {
		Db *sql.DB
	}

	ISyncUtils interface {
		InteractiveSync(entity interface{})
		ForceSync(entity interface{})
		CheckTable(entity interface{}) error
		DropTable(entity interface{}) error
		CreateTable(entity interface{}) error
	}
)

func (s SyncUtils) InteractiveSync(entity interface{}) {
	var err error

	if err = s.CheckTable(entity); err == nil {
		return
	}

	fmt.Println(err)
	fmt.Print("Do you want to drop and recreate table entity? [y/n] :")

	var answer string
	if fmt.Scanln(&answer); answer != "y" {
		return
	}

	s.RecreateTable(entity)
}

func (s SyncUtils) ForceSync(entity interface{}) {
	var err error

	if err = s.CheckTable(entity); err == nil {
		return
	}

	fmt.Println(err)

	s.RecreateTable(entity)
}

// CheckTable CheckTable is exists and have columns that ready to filled with Field from struct
func (s SyncUtils) CheckTable(entity interface{}) error {

	entityType := reflect.TypeOf(entity)

	getter := DbColumnGetter{s.Db}

	tableName := ToSnakeCase(entityType.Name())

	columns, err := getter.GetColumns(tableName)
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
		err = isSame(columns[i], fields[i], i)
		if err != nil {
			return err
		}
	}

	return nil
}

// DropTable just drop the table
func (s SyncUtils) DropTable(entity interface{}) error {

	entityType := reflect.TypeOf(entity)

	tableName := ToSnakeCase(entityType.Name())

	query := fmt.Sprintf("DROP TABLE %s;", tableName)

	_, err := s.Db.Exec(query)

	return err
}

// CreateTable create table with columns as same as field
func (s SyncUtils) CreateTable(entity interface{}) error {
	entityType := reflect.TypeOf(entity)

	columns := []string{}
	for i := 0; i < entityType.NumField(); i++ {
		entityField := entityType.Field(i)
		column := processField(entityField)
		columns = append(columns, column)
	}

	tableName := ToSnakeCase(entityType.Name())

	query := getCreateSQL(tableName, columns)

	_, err := s.Db.Exec(query)

	if err != nil {
		return errors.New(err.Error() + " query: " + query)
	}

	return nil
}

func (s SyncUtils) RecreateTable(entity interface{}) {
	s.DropTable(entity)
	if err := s.CreateTable(entity); err != nil {
		fmt.Println(err)
	}
}

type DbColumnGetter struct {
	Db *sql.DB
}

func (d DbColumnGetter) GetColumns(tableName string) ([]*sql.ColumnType, error) {
	var err error
	var result []*sql.ColumnType

	query := fmt.Sprintf("SELECT * FROM %s LIMIT 1;", tableName)

	rows, err := d.Db.Query(query)
	if err != nil {
		return result, err
	}

	fmt.Printf("%+v\n", rows)

	result, err = rows.ColumnTypes()
	if err != nil {
		return result, err
	}

	fmt.Printf("%+v\n", result)
	return result, nil
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
	case reflect.Bool:
		return processBoolean(field)
	case reflect.Int:
		return processInteger(field)
	case reflect.Int32:
		return processInteger32(field)
	case reflect.Int64:
		return processInteger64(field)
	case reflect.String:
		return processString(field)
	default:
		log.Printf("Cannot process %s field", field.Type.Kind())

	}

	return ""
}

func processPrimaryKey(field reflect.StructField) string {
	if field.Type.Kind() != reflect.Int64 {
		log.Printf("ID should be use int64")
	}
	return "id SERIAL8"
}

func processBoolean(field reflect.StructField) string {
	return fmt.Sprintf("%s BOOL", ToSnakeCase(field.Name))
}

func processInteger(field reflect.StructField) string {
	return fmt.Sprintf("%s INT", ToSnakeCase(field.Name))
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
