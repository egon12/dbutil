package dbutil

import (
	. "github.com/dave/jennifer/jen"
	"reflect"
)

func generateWhereFactoryStruct(entity reflect.Type) Code {

	structName := entity.Name() + "WhereFactory"

	codes := []Code{}
	for i := 0; i < entity.NumField(); i++ {
		field := entity.Field(i)
		fieldName := "where" + field.Name + "Value"
		code := Id(fieldName).Op("*").Add().Id(field.Type.String())
		codes = append(codes, code)
	}

	return Type().Id(structName).Struct(codes...)

}

func generateWhereFactoryFunctions(entity reflect.Type) []Code {

	structName := entity.Name() + "WhereFactory"

	codes := []Code{}
	for i := 0; i < entity.NumField(); i++ {

		field := entity.Field(i)
		funcName := "Where" + field.Name
		fieldName := "where" + field.Name + "Value"
		fieldType := field.Type.String()

		code := Func().Params(Id("w").Id(structName)).Id(funcName).Params(Id("value").Id(fieldType)).Block(
			If(Id("w").Dot(fieldName).Op("==").Nil()).Block(Id("w").Dot(fieldName).Op("=").New(Id(fieldType))),
			Op("*").Add().Id("w").Dot(fieldName).Op("=").Id("value"),
		)

		codes = append(codes, code)
	}

	return codes

}

func generateWhereFactoryStandaloneFunctions(entity reflect.Type) []Code {

	repoName := getRepoName(entity)
	structName := entity.Name() + "WhereFactory"

	codes := []Code{}
	for i := 0; i < entity.NumField(); i++ {

		field := entity.Field(i)
		funcName := "Where" + field.Name
		valueType := field.Type.String()
		code := Func().Params(Id("r").Id(repoName)).Id(funcName).Params(Id("value").Id(valueType)).Id(structName).Block(
			Id("w").Op(":=").Id(structName).Values(),
			Id("w").Dot(funcName).Call(Id("value")),
			Return(Id("w")),
		)

		codes = append(codes, code)
	}

	return codes

}

/*
func generateSelectWhereFunc(entity reflect.Type) Code {

	query := generateSelectQuery(entity)

	queryLit := Lit(query)

	fields, _ := createSelectField(entity)

	receiver := structName(entity)

	functionName := "Select"

	param1 := Id("limit")
	param2 := Id("offset").Int()
	params := []Code{param1, param2}

	returnType1 := Index().Qual(entity.PkgPath(), entity.Name())
	returnType2 := Error()
	returnType := []Code{returnType1, returnType2}

	theReturn := Return(Id("result"), Err())

	return createRepoFunction(file, receiver, functionName, params, returnType).Block(
		Id("result").Op(":=").Index().Qual(entity.PkgPath(), entity.Name()).Values(),

		List(Id("rows"), Err()).Op(":=").Id("r").Dot("Query").Call(queryLit, Id("limit"), Id("offset")),
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
*/
