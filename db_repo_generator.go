package dbutil

import (
	"errors"
	"fmt"
	. "github.com/dave/jennifer/jen"
	"reflect"
)

func GenerateRepository(packageName string, realEntity interface{}, filename string) error {

	generator, err := NewRepoGenerator(packageName, realEntity, nil)
	if err != nil {
		return err
	}

	generator.Generate()

	err = generator.Save(filename)
	if err != nil {
		return err
	}

	return nil
}

type IRepoGenerator interface {
	Generate()
	GoString()
	Save(string)
}

type RepoGenerator struct {
	entity          reflect.Type
	file            *File
	fieldFilterFunc func(reflect.StructField) bool
}

func NewRepoGenerator(
	packageName string,
	realEntity interface{},
	fieldFilterFunc func(reflect.StructField) bool,
) (*RepoGenerator, error) {

	repoGenerator := RepoGenerator{}

	repoGenerator.entity = reflect.TypeOf(realEntity)
	if repoGenerator.entity.Kind() != reflect.Struct {
		errMsg := fmt.Sprintf("%s is not a struct. We can only generate repository for a struct", repoGenerator.entity.Name())
		return &repoGenerator, errors.New(errMsg)
	}

	repoGenerator.file = NewFile(packageName)

	repoGenerator.fieldFilterFunc = fieldFilterFunc

	return &repoGenerator, nil
}

func (repoGenerator *RepoGenerator) Generate() {
	file := repoGenerator.file
	entity := repoGenerator.entity

	file.Add(generateRepoStruct(entity))
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
}

func (repoGenerator *RepoGenerator) GoString() string {
	return repoGenerator.file.GoString()
}

func (repoGenerator *RepoGenerator) Save(filename string) error {
	return repoGenerator.file.Save(filename)
}

func (repoGenerator *RepoGenerator) getRepoName() string {
	return "Postgre" + repoGenerator.entity.Name() + "Repository"
}

func generateRepoStruct(entity reflect.Type) Code {

	repoName := getRepoName(entity)

	return Type().Id(repoName).Struct(
		Id("ReadWrite").Op("*").Add().Qual("database/sql", "DB"),
		Id("ReadOnly").Op("*").Add().Qual("database/sql", "DB"),
	)
}

func generateConstructor(entity reflect.Type) Code {

	name := getRepoName(entity)

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

func createRepoFunction(entity reflect.Type, functionName string, params, returnType []Code) *Statement {
	receiver := getRepoName(entity)
	return Func().Params(Id("r").Id(receiver)).Id(functionName).Params(params...).Params(returnType...)
}

func getRepoName(entity reflect.Type) string {
	return "Postgre" + entity.Name() + "Repository"
}
