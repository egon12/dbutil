package dbutil

import (
	. "github.com/egon12/dbutil/mydomain"
	"reflect"
	"testing"
)

func TestGenerateSelectQuery(t *testing.T) {
	entity := reflect.TypeOf(EntityExamples2{})
	query := generateSelectQuery(entity)

	expectedQuery := "SELECT id, name, age, address FROM entity_examples2 LIMIT $1 OFFSET $2"

	if query != expectedQuery {
		t.Error("Wrong query:", query)
	}
}

func TestGenerateInsertQuery(t *testing.T) {
	entity := reflect.TypeOf(EntityExamples2{})
	query := generateInsertQuery(entity)

	expectedQuery := "INSERT INTO entity_examples2(name, age, address) VALUES ($1, $2, $3) RETURNING id"

	if query != expectedQuery {
		t.Error("Wrong query:", query)
	}
}

func TestGenerateUpdateQuery(t *testing.T) {
	entity := reflect.TypeOf(EntityExamples2{})
	query := generateUpdateQuery(entity)

	expectedQuery := "UPDATE entity_examples2 SET name = $1, age = $2, address = $3 WHERE id = $4"

	if query != expectedQuery {
		t.Error("Wrong query:", query)
	}
}

func TestGenerateDeleteQuery(t *testing.T) {
	entity := reflect.TypeOf(EntityExamples2{})
	query := generateDeleteQuery(entity)

	expectedQuery := "DELETE FROM entity_examples2 WHERE id = $1"

	if query != expectedQuery {
		t.Error("Wrong query:", query)
	}
}
