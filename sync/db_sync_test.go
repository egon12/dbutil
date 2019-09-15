package sync

import (
	"database/sql"
	"github.com/egon12/dbutil/mock"
	"reflect"
	"testing"
)

type EntityExamples1 struct {
	ID   int64
	Name string
}

var (
	mockDb *mock.MockDbDriver
	db     *sql.DB
)

func init() {
	mockDb = &mock.MockDbDriver{}

	sql.Register("mockDriver", mockDb)

	db, _ = sql.Open("mockDriver", "")
}

func TestGetDbColumns(t *testing.T) {
	getter := DbColumnGetter{db}

	mockDb.ColumnsValue = []string{"id", "name"}

	columns, _ := getter.GetColumns("entity_examples1")

	if columns[0].Name() != "id" && columns[1].Name() != "name" {
		t.Error("The Column should be 'id, name' got:", columns)
	}

	if mockDb.QueryValue != "SELECT * FROM entity_examples1 LIMIT 1;" {
		t.Error("Query is not same got:", mockDb.QueryValue)
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

	ex := EntityExamples1{}

	s := SyncUtils{db}

	mockDb.ColumnsValue = []string{"id", "name"}
	mockDb.ColumnsTypeValue = []reflect.Type{reflect.TypeOf(ex.ID), reflect.TypeOf(ex.Name)}

	err := s.CheckTable(ex)
	if err != nil {
		t.Error("It should be not error:", err)
	}
}

func TestCheckTableDbColumnsNotComplete(t *testing.T) {

	ex := EntityExamples1{}

	s := SyncUtils{db}

	mockDb.ColumnsValue = []string{"id"}
	mockDb.ColumnsTypeValue = []reflect.Type{reflect.TypeOf(ex.ID)}

	err := s.CheckTable(ex)
	if err == nil {
		t.Error("It should be error")
	}
}

func TestCreateTable(t *testing.T) {

	s := SyncUtils{db}

	ex := EntityExamples1{}

	mockDb.ColumnsValue = []string{"id"}
	mockDb.ColumnsTypeValue = []reflect.Type{reflect.TypeOf(ex.ID)}

	err := s.CreateTable(ex)
	if err != nil {
		t.Error("It should be not error")
	}

	if mockDb.QueryValue != "CREATE TABLE entity_examples1 (id SERIAL8, name VARCHAR(255));" {
		t.Error("Wrong Query :", mockDb.QueryValue)
	}

}

func TestDropTable(t *testing.T) {

	s := SyncUtils{db}

	ex := EntityExamples1{}

	mockDb.ColumnsValue = []string{"id"}
	mockDb.ColumnsTypeValue = []reflect.Type{reflect.TypeOf(ex.ID)}

	err := s.CreateTable(ex)
	if err != nil {
		t.Error("It should be not error")
	}

	if mockDb.QueryValue != "CREATE TABLE entity_examples1 (id SERIAL8, name VARCHAR(255));" {
		t.Error("Wrong Query :", mockDb.QueryValue)
	}

}
