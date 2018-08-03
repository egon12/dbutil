# dbutil

This is experimental for make my life easier, without dependend to any ORM/driver besides
what that provide by go. 

Right now, it only support Postgres (because of different syntax between SQL)

## Database Synchronization

Put this script into file like `sync.go`

``` Golang
package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"

	"yourdomain"
	"github.com/egon12/dbutil"
)

func main() {

	db, err := sql.Open("postgres", "user=user dbname=db_to_sync password=password sslmode=disable")
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}
	dbutil.Db = db

	emptyEntity := yourdomain.Entity{}
	dbutil.InteractiveSync(emptyEntity)
}
```

and then execute

```
go run sync.go
```

## Repository Generator

Put this script into file like `generate.go`

``` Golang
package main

import (
	"yourdomain"
	"github.com/egon12/dbutil"
	"log"
)

func main() {

	emptyEntity := yourdomain.Entity{}
	err := dbutil.GenerateRepository("gateway", emptyEntity, "./gateway/repository.go")
	if err != nil {
		log.Fatal(err)
	}
}
```

and then execute

```
go run generate.go
```

