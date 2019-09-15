package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
)

var (
	generateFileName string = "./generate.go"
)

func main() {

	importName, fullEntityName, packageName, fileName, customArguments := parseArguments()

	if fullEntityName == "" {
		fmt.Printf("GenRepo need 'entity' real: %s \n\n", fullEntityName)
		flag.Usage()
		return
	}

	if importName != "" {
		importName = fmt.Sprintf("\"%s\"", importName)
	}

	compiledScript := fmt.Sprintf(script,
		importName,
		fullEntityName,
		packageName,
		fileName,
	)

	createFile(compiledScript)
	defer removeFile()

	runCommand := exec.Command("go", "run", "generate.go", customArguments)
	result, err := runCommand.Output()
	if err != nil {
		fmt.Printf("Cannot get Output: %s\n%s\n", err.Error(), string(result))
		return
	}
}

func parseArguments() (string, string, string, string, string) {

	var importName string
	flag.StringVar(&importName, "import", "", "Package to be import")

	var fullEntityName string
	flag.StringVar(&fullEntityName, "entity", "", "Full Entity Name")

	var packageName string
	flag.StringVar(&packageName, "package", "main", "Package Name for the new Repository")

	var fileName string
	flag.StringVar(&fileName, "output", "generated.go", "FileName of generated Repo")

	var customArguments string
	flag.StringVar(&customArguments, "arg", "", "Arguments add when run generate")

	flag.Parse()

	return importName, fullEntityName, packageName, fileName, customArguments
}

func createFile(compiledScript string) {
	err := ioutil.WriteFile(generateFileName, []byte(compiledScript), 0666)
	if err != nil {
		fmt.Errorf("Cannot write file:" + err.Error())
		panic(err)
	}
}

func removeFile() {
	rmCommand := exec.Command("rm", generateFileName)
	err := rmCommand.Start()
	if err != nil {
		fmt.Errorf("Cannot remove file:" + err.Error())
		panic(err)
	}
}

var script1 = `
package main

func main() {
	println("Hello World!!!\n")
}

`

var script = `

package main

import (
	%s
	"github.com/egon12/dbutil"
	"log"
)

func main() {

	emptyEntity := %s{}
	err := dbutil.GenerateRepository("%s", emptyEntity, "%s")
	if err != nil {
		log.Fatal(err)
	}
}
`
