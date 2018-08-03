package utils

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

type EntityExamples2 struct {
	ID      int64
	Name    string
	Age     int
	Address string
}

func TestCreateFile(t *testing.T) {
	file, err := generateRepository("domain", EntityExamples2{})
	if err != nil {
		t.Error(err)
	}

	fileContents := fmt.Sprintf("%#v", file)

	expectedFileBytes, err := ioutil.ReadFile("entity_examples2.go.gen")
	if err != nil {
		t.Error(err)
	}

	expectedFileContents := string(expectedFileBytes)

	if fileContents != expectedFileContents {
		file.Save("/tmp/entity_examples2.go.gen")
		t.Error("Generated is not same as 'entity_examples2.go.gen'. Check with /tmp/entity_examples2.go.gen")
		t.Log(expectedFileContents)
		t.Log(fileContents)
	}
}

func TestGenerateSelectQuery(t *testing.T) {
	entity := reflect.TypeOf(EntityExamples2{})
	query, _ := createSelectQuery(entity)

	expectedQuery := "SELECT id, name, age, address FROM entity_examples2 LIMIT $1 OFFSET $2"

	if query != expectedQuery {
		t.Error("Wrong query:", query)
	}
}

func TestGenerateInsertQuery(t *testing.T) {
	entity := reflect.TypeOf(EntityExamples2{})
	query, _ := createInsertQuery(entity)

	expectedQuery := "INSERT INTO entity_examples2(name, age, address) VALUES ($1, $2, $3)"

	if query != expectedQuery {
		t.Error("Wrong query:", query)
	}
}

func TestGenerateUpdateQuery(t *testing.T) {
	entity := reflect.TypeOf(EntityExamples2{})
	query, _ := createUpdateQuery(entity)

	expectedQuery := "UPDATE entity_examples2 SET name = $1, age = $2, address = $3 WHERE id = $4"

	if query != expectedQuery {
		t.Error("Wrong query:", query)
	}
}
