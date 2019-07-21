package dbutil

import (
	"errors"
	"fmt"
	. "github.com/dave/jennifer/jen"
	"reflect"
)

func GenerateRepository(packageName string, realEntity interface{}, filename string) error {

	file, err := generateRepository(packageName, realEntity)
	if err != nil {
		return err
	}

	err = file.Save(filename)
	if err != nil {
		return err
	}

	return nil
}

func generateRepository(packageName string, realEntity interface{}) (*File, error) {

	entity := reflect.TypeOf(realEntity)

	if entity.Kind() != reflect.Struct {
		errMsg := fmt.Sprintf("%s is not a struct. We can only generate repository for a struct", entity.Name())
		return nil, errors.New(errMsg)
	}

	file := NewFile(packageName)
	file.Add(generateStruct(entity))
	file.Add(generateConstructor(entity))

	file.Add(createSelectFunc(entity))

	file.Add(createInsertFunc(entity))
	file.Add(createUpdateFunc(entity))
	file.Add(createDeleteFunc(entity))
	file.Add(generateWhereFactoryStruct(entity))

	for _, code := range generateWhereFactoryFunctions(entity) {
		file.Add(code)
	}

	for _, code := range generateWhereFactoryStandaloneFunctions(entity) {
		file.Add(code)
	}

	return file, nil

}

func generateStruct(entity reflect.Type) Code {
	return Type().Id(structName(entity)).Struct(
		Id("ReadWrite").Op("*").Add().Qual("database/sql", "DB"),
		Id("ReadOnly").Op("*").Add().Qual("database/sql", "DB"),
	)
}

func generateConstructor(entity reflect.Type) Code {

	name := structName(entity)

	params := []Code{
		Id("readWrite").Op("*").Add().Qual("database/sql", "DB"),
		Id("readOnly").Op("*").Add().Qual("database/sql", "DB"),
	}

	returnType := []Code{
		Id(name),
		Error(),
	}

	return Func().Id("New" + name).Params(params...).Params(returnType...).Block(
		If(Id("readOnly").Op("!=").Nil()).Block(
			Return().List(Id(name).Values(Id("readWrite"), Id("readOnly")), Nil()),
		).Else().Block(
			Return().List(Id(name).Values(Id("readWrite"), Id("readWrite")), Nil()),
		),
	)
}

func structName(entity reflect.Type) string {
	return "Postgre" + entity.Name() + "Repository"
}

func createRepoFunction(entity reflect.Type, functionName string, params, returnType []Code) *Statement {

	receiver := structName(entity)

	return Func().Params(Id("r").Id(receiver)).Id(functionName).Params(params...).Params(returnType...)

}
