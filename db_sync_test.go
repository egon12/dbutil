package utils

import (
	"database/sql"
	"reflect"
	"testing"
)

type EntityExamples1 struct {
	ID   int64
	Name string
}

func init() {
	sql.Register("mockDriver", MockDb)
}

func TestGetDbColumns(t *testing.T) {
	// Register the DB
	Db, _ = sql.Open("mockDriver", "")

	ex := EntityExamples1{}
	entity := reflect.TypeOf(ex)

	MockDb.ColumnsValue = []string{"id", "name"}

	columns, _ := getDbColumns(entity)

	if columns[0].Name() != "id" && columns[1].Name() != "name" {
		t.Error("The Column should be 'id, name' got:", columns)
	}

	if MockDb.QueryValue != "SELECT * FROM entity_examples1 LIMIT 1;" {
		t.Error("Query is not same got:", MockDb.QueryValue)
	}
}

func TestGetFields(t *testing.T) {

	ex := EntityExamples1{}
	entity := reflect.TypeOf(ex)

	fields, _ := getFields(entity)

	if fields[0].Name != "ID" {
		t.Error("First field is not id. got:", fields[0])
	}

	if fields[1].Name != "Name" {
		t.Error("First field is not name. got:", fields[1])
	}
}

type EntityExamples1WithTrait struct {
	EntityExamples1
	Address string
}

func TestGetFieldsNestedClass(t *testing.T) {
	t.Skip("for now it doesn't yet support nested class")

	ex := EntityExamples1WithTrait{}
	entity := reflect.TypeOf(ex)

	fields, _ := getFields(entity)

	if fields[0].Name != "id" {
		t.Error("First field is not id. got:", fields[0])
	}

	if fields[1].Name != "name" {
		t.Error("First field is not name. got:", fields[1])
	}

	if fields[2].Name != "address" {
		t.Error("First field is not name. got:", fields[2])
	}

}

func TestCheckTable(t *testing.T) {
	Db, _ = sql.Open("mockDriver", "")

	ex := EntityExamples1{}

	MockDb.ColumnsValue = []string{"id", "name"}
	MockDb.ColumnsTypeValue = []reflect.Type{reflect.TypeOf(ex.ID), reflect.TypeOf(ex.Name)}

	err := CheckTable(ex)
	if err != nil {
		t.Error("it should be not error", err)
	}
}

func TestCheckTableDbColumnsNotComplete(t *testing.T) {
	Db, _ = sql.Open("mockDriver", "")

	ex := EntityExamples1{}

	MockDb.ColumnsValue = []string{"id"}
	MockDb.ColumnsTypeValue = []reflect.Type{reflect.TypeOf(ex.ID)}

	err := CheckTable(ex)
	if err == nil {
		t.Error("It should be error")
	}
}

func TestCreateTable(t *testing.T) {

	Db, _ = sql.Open("mockDriver", "")

	ex := EntityExamples1{}

	MockDb.ColumnsValue = []string{"id"}
	MockDb.ColumnsTypeValue = []reflect.Type{reflect.TypeOf(ex.ID)}

	err := CreateTable(ex)
	if err != nil {
		t.Error("It should be not error")
	}

	if MockDb.QueryValue != "CREATE TABLE entity_examples1 (id SERIAL, name VARCHAR(255));" {
		t.Error("Wrong Query :", MockDb.QueryValue)
	}

}

func TestDropTable(t *testing.T) {

	Db, _ = sql.Open("mockDriver", "")

	ex := EntityExamples1{}

	MockDb.ColumnsValue = []string{"id"}
	MockDb.ColumnsTypeValue = []reflect.Type{reflect.TypeOf(ex.ID)}

	err := CreateTable(ex)
	if err != nil {
		t.Error("It should be not error")
	}

	if MockDb.QueryValue != "CREATE TABLE entity_examples1 (id SERIAL, name VARCHAR(255));" {
		t.Error("Wrong Query :", MockDb.QueryValue)
	}

}
