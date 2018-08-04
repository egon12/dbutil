package dbutil

import (
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

	// TODO, check entity must struct
	// check entity must have fields

	file, err := createFile(packageName)
	if err != nil {
		return file, err
	}

	file, err = createStruct(file, entity)
	if err != nil {
		return file, err
	}

	file, err = createSelectFunc(file, entity)
	if err != nil {
		return file, err
	}

	file, err = createInsertFunc(file, entity)
	if err != nil {
		return file, err
	}

	file, err = createUpdateFunc(file, entity)
	if err != nil {
		return file, err
	}

	/*
		file, err = createDeleteFunc(file, entity)
		if err != nil {
			return "", err
		}
	*/

	return file, nil

}

func createFile(packageName string) (*File, error) {
	f := NewFile(packageName)

	return f, nil
}

func createStruct(file *File, entity reflect.Type) (*File, error) {

	file.Type().Id(structName(entity)).Struct(
		Op("*").Add().Qual("database/sql", "DB"),
	)

	return file, nil

}

func structName(entity reflect.Type) string {
	return "Postgre" + entity.Name() + "Repository"
}

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
func createSelectFunc(file *File, entity reflect.Type) (*File, error) {

	query, err := createSelectQuery(entity)
	if err != nil {
		return file, err
	}

	queryLit := Lit(query)

	fields, err := createSelectField(entity)
	if err != nil {
		return file, err
	}

	receiver := structName(entity)

	functionName := "Select"

	param1 := Id("where").String()
	param2 := Id("limit")
	param3 := Id("offset").Int()
	params := []Code{param1, param2, param3}

	returnType1 := Index().Qual(entity.PkgPath(), entity.Name())
	returnType2 := Error()
	returnType := []Code{returnType1, returnType2}

	theReturn := Return(Id("result"), Err())

	createRepoFunction(file, receiver, functionName, params, returnType).Block(
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

	return file, nil
}

func createSelectQuery(entity reflect.Type) (string, error) {

	columnsSelect := ""
	for i := 0; i < entity.NumField(); i++ {
		name := ToSnakeCase(entity.Field(i).Name)
		if i == 0 {
			columnsSelect += name
		} else {
			columnsSelect += ", " + name
		}
	}

	tableName := ToSnakeCase(entity.Name())

	query := fmt.Sprintf("SELECT %s FROM %s LIMIT $1 OFFSET $2", columnsSelect, tableName)

	return query, nil
}

func createSelectField(entity reflect.Type) (Code, error) {

	fields := []Code{}
	for i := 0; i < entity.NumField(); i++ {
		idField := Op("&").Add().Id("entity").Dot(entity.Field(i).Name)
		fields = append(fields, idField)
	}

	return List(fields...), nil
}

// createInsertFunc
//
// result:
//
// func Inssert(entity &s) {
// 	return sql.Exec(%s, %s)
// }
//
//
func createInsertFunc(file *File, entity reflect.Type) (*File, error) {

	query, err := createInsertQuery(entity)
	if err != nil {
		return file, err
	}

	queryLit := Lit(query)

	fields, err := createInsertField(entity)
	if err != nil {
		return file, err
	}

	receiver := structName(entity)

	functionName := "Insert"

	param1 := Id("entity").Qual(entity.PkgPath(), entity.Name())
	params := []Code{param1}

	returnType1 := Error()
	returnType := []Code{returnType1}

	createRepoFunction(file, receiver, functionName, params, returnType).Block(
		List(Id("_"), Err()).Op(":=").Id("r").Dot("Exec").Call(queryLit, fields),
		Return().Err(),
	)

	return file, nil
}

func createInsertQuery(entity reflect.Type) (string, error) {

	columnsSelect := ""
	fills := ""
	for i := 1; i < entity.NumField(); i++ {
		name := ToSnakeCase(entity.Field(i).Name)
		if i == 1 {
			columnsSelect += name
			fills += fmt.Sprintf("$%d", i)
		} else {
			columnsSelect += ", " + name
			fills += fmt.Sprintf(", $%d", i)
		}
	}

	tableName := ToSnakeCase(entity.Name())

	query := fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s)", tableName, columnsSelect, fills)

	return query, nil
}

func createInsertField(entity reflect.Type) (Code, error) {

	fields := []Code{}
	for i := 1; i < entity.NumField(); i++ {
		idField := Id("entity").Dot(entity.Field(i).Name)
		fields = append(fields, idField)
	}

	return List(fields...), nil
}

// createUpdateFunc
func createUpdateFunc(file *File, entity reflect.Type) (*File, error) {

	query, err := createUpdateQuery(entity)
	if err != nil {
		return file, err
	}

	queryLit := Lit(query)

	fields, err := createUpdateField(entity)
	if err != nil {
		return file, err
	}

	receiver := structName(entity)

	functionName := "Update"

	param1 := Id("entity").Qual(entity.PkgPath(), entity.Name())
	params := []Code{param1}

	returnType1 := Error()
	returnType := []Code{returnType1}

	createRepoFunction(file, receiver, functionName, params, returnType).Block(
		List(Id("_"), Err()).Op(":=").Id("r").Dot("Exec").Call(queryLit, fields),
		Return().Err(),
	)

	return file, nil
}

func createUpdateQuery(entity reflect.Type) (string, error) {

	sets := ""
	for i := 1; i < entity.NumField(); i++ {
		name := ToSnakeCase(entity.Field(i).Name)
		if i == 1 {
			sets += fmt.Sprintf("%s = $%d", name, i)
		} else {
			sets += fmt.Sprintf(", %s = $%d", name, i)
		}
	}

	tableName := ToSnakeCase(entity.Name())

	idFieldNumber := fmt.Sprintf("$%d", entity.NumField())

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = %s", tableName, sets, idFieldNumber)

	return query, nil
}

func createUpdateField(entity reflect.Type) (Code, error) {
	fields := []Code{}
	for i := 1; i < entity.NumField(); i++ {
		idField := Id("entity").Dot(entity.Field(i).Name)
		fields = append(fields, idField)
	}

	idField := Id("entity").Dot(entity.Field(0).Name)
	fields = append(fields, idField)

	return List(fields...), nil
}

func createRepoFunction(file *File, receiver, functionName string, params, returnType []Code) *Statement {

	return file.Func().Params(Id("r").Id(receiver)).Id(functionName).Params(params...).Params(returnType...)

}
