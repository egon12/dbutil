package main

import (
	"fmt"
	//"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	fset := token.NewFileSet() // positions are relative to fset

	src := `package foo

import (
	"fmt"
	"time"
)

func bar() {
	fmt.Println(time.Now())
}`

	// Parse src but stop after processing the imports.
	f, err := parser.ParseFile(fset, "", src, parser.Trace)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the imports from the file's AST.
	for _, s := range f.Imports {
		fmt.Println(s.Path.Value)
	}

	/*
		for _, s := range f.Decls {
			switch so := s.(type) {
			case *ast.GenDecl:
				//fmt.Printf("gen decl \n")
			case *ast.FuncDecl:
				//fmt.Printf("func decl %s\n", so.Name.Name)

			}

		}
	*/

}
