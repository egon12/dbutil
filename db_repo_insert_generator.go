package dbutil

import (
	. "github.com/dave/jennifer/jen"
	"reflect"
)

// createInsertFunc
//
// result:
//
// func Inssert(entity &s) {
// 	return sql.Exec(%s, %s)
// }
//
//
func createInsertFunc(entity reflect.Type) Code {

  params := []Code{
    Id("entity").Qual(entity.PkgPath(), entity.Name()),
  }

	returnType := []Code{Error()}

	query := generateInsertQuery(entity)
	queryLit := Lit(query)
	fields := createInsertField(entity)

	return createRepoFunction(entity, "Insert", params, returnType).Block(
		List(Id("_"), Err()).Op(":=").Id("r").Dot("ReadWrite").Dot("Exec").Call(queryLit, fields),
		Return().Err(),
	)
}

func createInsertField(entity reflect.Type) Code {

	fields := []Code{}
	for i := 1; i < entity.NumField(); i++ {
		idField := Id("entity").Dot(entity.Field(i).Name)
		fields = append(fields, idField)
	}

	return List(fields...)
}
