package dbutil

import (
	"fmt"
	. "github.com/egon12/dbutil/mydomain"
	"io/ioutil"
	"strings"
	"testing"
)

func Test_Generate_With_Generator(t *testing.T) {

	generator, _ := NewRepoGenerator("generated", EntityExamples2{}, nil)

	generator.Generate()

	fileContents := fmt.Sprintf("%#v", generator)

	expectedFileBytes, err := ioutil.ReadFile("./generated/entity_examples2_repo.go")
	if err != nil {
		t.Error(err)
	}

	expectedFileContents := string(expectedFileBytes)

	if fileContents != expectedFileContents {
		generator.Save("/tmp/entity_examples2_repo.go")
		t.Error("Generated is not same as 'generated/entity_examples2_repo.go'. Check with /tmp/entity_examples2_repo.go")
	}
}

func Test_Entity_IsNotStruct(t *testing.T) {
	falseEntity := ""

	_, err := NewRepoGenerator("generated", falseEntity, nil)

	if err == nil {
		t.Error("It should be some error")
	}
}

func Test_Error_Message_Should_Tell_Type(t *testing.T) {

	type SomeType string

	falseEntity := new(SomeType)
	_, err := NewRepoGenerator("generated", *falseEntity, nil)

	if !strings.Contains(err.Error(), "SomeType") {
		t.Error("Error should tell the type of entity. Got: " + err.Error())
	}

}
