package dbutil

import (
	. "github.com/dave/jennifer/jen"
	"reflect"
)

// createUpdateFunc
func createUpdateFunc(entity reflect.Type) Code {

	query := generateUpdateQuery(entity)

	queryLit := Lit(query)

	fields := createUpdateField(entity)

	functionName := "Update"

	param1 := Id("entity").Qual(entity.PkgPath(), entity.Name())
	params := []Code{param1}

	returnType1 := Error()
	returnType := []Code{returnType1}

	return createRepoFunction(entity, functionName, params, returnType).Block(
		List(Id("_"), Err()).Op(":=").Id("r").Dot("ReadWrite").Dot("Exec").Call(queryLit, fields),
		Return().Err(),
	)

}

func createUpdateField(entity reflect.Type) Code {
	fields := []Code{}
	for i := 1; i < entity.NumField(); i++ {
		idField := Id("entity").Dot(entity.Field(i).Name)
		fields = append(fields, idField)
	}

	idField := Id("entity").Dot(entity.Field(0).Name)
	fields = append(fields, idField)

	return List(fields...)
}
