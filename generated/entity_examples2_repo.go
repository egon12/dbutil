package generated

import (
	"database/sql"
	mydomain "github.com/egon12/dbutil/mydomain"
)

type PostgreEntityExamples2Repository struct {
	ReadWrite *sql.DB
	ReadOnly  *sql.DB
}

func NewPostgreEntityExamples2Repository(readWrite *sql.DB, readOnly *sql.DB) (PostgreEntityExamples2Repository, error) {
	if readOnly != nil {
		return PostgreEntityExamples2Repository{readWrite, readOnly}, nil
	} else {
		return PostgreEntityExamples2Repository{readWrite, readWrite}, nil
	}
}
func (r PostgreEntityExamples2Repository) Select(limit, offset int) ([]mydomain.EntityExamples2, error) {
	result := []mydomain.EntityExamples2{}
	rows, err := r.ReadOnly.Query("SELECT id, name, age, address FROM entity_examples2 LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		entity := mydomain.EntityExamples2{}
		rows.Scan(&entity.ID, &entity.Name, &entity.Age, &entity.Address)
		result = append(result, entity)
	}
	return result, err
}
func (r PostgreEntityExamples2Repository) Insert(entity mydomain.EntityExamples2) (mydomain.EntityExamples2, error) {
	err := r.ReadWrite.QueryRow("INSERT INTO entity_examples2(name, age, address) VALUES ($1, $2, $3) RETURNING id", entity.Name, entity.Age, entity.Address).Scan(&entity.ID)
	return entity, err
}
func (r PostgreEntityExamples2Repository) Update(entity mydomain.EntityExamples2) error {
	_, err := r.ReadWrite.Exec("UPDATE entity_examples2 SET name = $1, age = $2, address = $3 WHERE id = $4", entity.Name, entity.Age, entity.Address, entity.ID)
	return err
}
func (r PostgreEntityExamples2Repository) Delete(entity mydomain.EntityExamples2) error {
	_, err := r.ReadWrite.Exec("DELETE FROM entity_examples2 WHERE id = $1", entity.ID)
	return err
}

type EntityExamples2WhereFactory struct {
	whereIDValue      *int64
	whereNameValue    *string
	whereAgeValue     *int32
	whereAddressValue *string
}

func (w EntityExamples2WhereFactory) WhereID(value int64) {
	if w.whereIDValue == nil {
		w.whereIDValue = new(int64)
	}
	*w.whereIDValue = value
}
func (w EntityExamples2WhereFactory) WhereName(value string) {
	if w.whereNameValue == nil {
		w.whereNameValue = new(string)
	}
	*w.whereNameValue = value
}
func (w EntityExamples2WhereFactory) WhereAge(value int32) {
	if w.whereAgeValue == nil {
		w.whereAgeValue = new(int32)
	}
	*w.whereAgeValue = value
}
func (w EntityExamples2WhereFactory) WhereAddress(value string) {
	if w.whereAddressValue == nil {
		w.whereAddressValue = new(string)
	}
	*w.whereAddressValue = value
}
func (r PostgreEntityExamples2Repository) WhereID(value int64) EntityExamples2WhereFactory {
	w := EntityExamples2WhereFactory{}
	w.WhereID(value)
	return w
}
func (r PostgreEntityExamples2Repository) WhereName(value string) EntityExamples2WhereFactory {
	w := EntityExamples2WhereFactory{}
	w.WhereName(value)
	return w
}
func (r PostgreEntityExamples2Repository) WhereAge(value int32) EntityExamples2WhereFactory {
	w := EntityExamples2WhereFactory{}
	w.WhereAge(value)
	return w
}
func (r PostgreEntityExamples2Repository) WhereAddress(value string) EntityExamples2WhereFactory {
	w := EntityExamples2WhereFactory{}
	w.WhereAddress(value)
	return w
}
