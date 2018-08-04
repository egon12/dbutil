package dbutil

import (
	. "github.com/dave/jennifer/jen"
	"reflect"
)

func createDeleteFunc(entity reflect.Type) Code {

	params := []Code{
		Id("entity").Qual(entity.PkgPath(), entity.Name()),
	}

	returnType := []Code{Error()}

	query := Lit(generateDeleteQuery(entity))

	field := Id("entity").Dot("ID")

	return createRepoFunction(entity, "Insert", params, returnType).Block(
		List(Id("_"), Err()).Op(":=").Id("r").Dot("ReadWrite").Dot("Exec").Call(query, field),
		Return().Err(),
	)
}
