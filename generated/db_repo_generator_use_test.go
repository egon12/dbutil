package generated

import (
	"database/sql"
	"fmt"
	"github.com/egon12/dbutil"
	. "github.com/egon12/dbutil/mydomain"
	_ "github.com/lib/pq"
	"testing"
)

func init() {
	sql.Register("mockDriver", dbutil.MockDb)
}

func TestUseFile(t *testing.T) {

	db, err := sql.Open("mockDriver", "")
	defer db.Close()
	if err != nil {
		t.Error(err)
	}

	dbutil.Db = db
	dbutil.ForceSync(EntityExamples2{})

	repo, err := NewPostgreEntityExamples2Repository(db, nil)
	if err != nil {
		t.Error(err)
	}

	entity := EntityExamples2{
		Name:    "Egon",
		Age:     34,
		Address: "Jalan something",
	}

	e1, err := repo.Insert(entity)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", e1)
	e2, err := repo.Insert(entity)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", e2)
	e3, err := repo.Insert(entity)
	if err != nil {
		t.Error(err)
	}

	e3.Age = 35
	e3.Name = "Nova"
	err = repo.Update(e3)
	if err != nil {
		t.Error(err)
	}

	tenEntity, err := repo.Select(10, 0)
	if err != nil {
		t.Error(err)
	}

	for _, e := range tenEntity {
		fmt.Printf("%+v\n", e)
		repo.Delete(e)
	}
}
