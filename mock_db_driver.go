package dbutil

import (
	"database/sql/driver"
	"errors"
	"reflect"
)

var (
	MockDb = MockDbDriver{}
)

type MockDbDriver struct {
	QueryValue        string
	NumInputValue     int
	LastInsertIDValue int64
	RowsAffectedValue int64
	ColumnsValue      []string
	ColumnsTypeValue  []reflect.Type
	ColumnsLength     []int64
	Rows              [][]interface{}
	Cursor            int64
}

/**
MockDbDriver as Driver
**/
func (m MockDbDriver) Open(conn string) (driver.Conn, error) {
	return m, nil
}

/**
MockDbDriver as Connection
**/
func (m *MockDbDriver) Driver() *MockDbDriver {
	return m
}

func (m MockDbDriver) Prepare(query string) (driver.Stmt, error) {
	MockDb.QueryValue = query
	return m, nil
}

func (m MockDbDriver) Close() error {
	return nil
}

func (m MockDbDriver) Begin() (driver.Tx, error) {
	return m, nil
}

/**
MockDbDriver as Stmt
**/
func (m MockDbDriver) NumInput() int {
	return MockDb.NumInputValue
}

func (m MockDbDriver) Exec(args []driver.Value) (driver.Result, error) {
	return m, nil
}

/**
MockDbDriver as Result
**/
func (m MockDbDriver) LastInsertId() (int64, error) {
	return MockDb.LastInsertIDValue, nil
}

func (m MockDbDriver) RowsAffected() (int64, error) {
	return MockDb.RowsAffectedValue, nil
}

func (m MockDbDriver) Query(args []driver.Value) (driver.Rows, error) {
	return m, nil
}

func (m MockDbDriver) Columns() []string {
	return MockDb.ColumnsValue
}

func (m MockDbDriver) Next(dest []driver.Value) error {

	if len(MockDb.Rows) == 0 {
		return errors.New("Empty Rows, forgot to set the Rows?")
	}

	result := MockDb.Rows[MockDb.Cursor]

	if len(result) != len(dest) {
		return errors.New("Destination and Row value is different. Maybe column length is not same as values in rows")
	}

	for index, item := range result {
		dest[index] = item
	}
	return nil
}

func (m MockDbDriver) Commit() error {
	return nil
}

func (m MockDbDriver) Rollback() error {
	return nil
}

/**
to fullfill RowsColumnTypeScanType
**/
func (m MockDbDriver) ColumnTypeScanType(i int) reflect.Type {
	if i >= len(MockDb.ColumnsTypeValue) {
		return reflect.TypeOf(new(interface{})).Elem()
	}
	return MockDb.ColumnsTypeValue[i]
}

func (m MockDbDriver) ColumnTypeLength(i int) (int64, bool) {
	if i >= len(MockDb.ColumnsLength) {
		return 255, true
	}
	return MockDb.ColumnsLength[i], true
}
