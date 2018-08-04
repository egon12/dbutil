package dbutil

import (
	. "github.com/dave/jennifer/jen"
	"reflect"
)

// createSelectFunc
//
// result:
//
// func Select() {
// 	result := %s{}
// 	rows, err := sql.Query(%s)
// 	if err != nil {
// 		return result, err
// 	}
// 	if rows.HasNext() {
// 		rows.Scan(%s)
// 	}
// 	return result, nil
// }
//
//
func createSelectFunc(entity reflect.Type) Code {

	queryLit := Lit(generateSelectQuery(entity))

	fields := createSelectField(entity)

	functionName := "Select"

	param1 := Id("limit")
	param2 := Id("offset").Int()
	params := []Code{param1, param2}

	returnType1 := Index().Qual(entity.PkgPath(), entity.Name())
	returnType2 := Error()
	returnType := []Code{returnType1, returnType2}

	theReturn := Return(Id("result"), Err())

	return createRepoFunction(entity, functionName, params, returnType).Block(
		Id("result").Op(":=").Index().Qual(entity.PkgPath(), entity.Name()).Values(),

		List(Id("rows"), Err()).Op(":=").Id("r").Dot("ReadOnly").Dot("Query").Call(queryLit, Id("limit"), Id("offset")),
		If(Id("err").Op("!=").Nil()).Block(
			theReturn,
		),

		For(Id("rows").Dot("Next").Call()).Block(
			Id("entity").Op(":=").Qual(entity.PkgPath(), entity.Name()).Values(),
			Id("rows").Dot("Scan").Call(fields),
			Id("result").Op("=").Append(Id("result"), Id("entity")),
		),
		theReturn,
	)

}

func createSelectField(entity reflect.Type) Code {

	fields := []Code{}
	for i := 0; i < entity.NumField(); i++ {
		idField := Op("&").Add().Id("entity").Dot(entity.Field(i).Name)
		fields = append(fields, idField)
	}

	return List(fields...)
}
