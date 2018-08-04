package dbutil

import (
	"fmt"
	"io/ioutil"
	"testing"
)

type EntityExamples2 struct {
	ID      int64
	Name    string
	Age     int
	Address string
}

func TestCreateFile(t *testing.T) {
	file, err := generateRepository("domain", EntityExamples2{})
	if err != nil {
		t.Error(err)
	}

	fileContents := fmt.Sprintf("%#v", file)

	expectedFileBytes, err := ioutil.ReadFile("entity_examples2.go.gen")
	if err != nil {
		t.Error(err)
	}

	expectedFileContents := string(expectedFileBytes)

	if fileContents != expectedFileContents {
		file.Save("/tmp/entity_examples2.go.gen")
		t.Error("Generated is not same as 'entity_examples2.go.gen'. Check with /tmp/entity_examples2.go.gen")
		t.Log(expectedFileContents)
		t.Log(fileContents)
	}
}
