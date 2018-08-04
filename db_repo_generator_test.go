package dbutil

import (
	"fmt"
	. "github.com/egon12/dbutil/mydomain"
	"io/ioutil"
	"testing"
)

func TestCreateFile(t *testing.T) {
	file, err := generateRepository("generated", EntityExamples2{})
	if err != nil {
		t.Error(err)
	}

	fileContents := fmt.Sprintf("%#v", file)

	expectedFileBytes, err := ioutil.ReadFile("./generated/entity_examples2_repo.go")
	if err != nil {
		t.Error(err)
	}

	expectedFileContents := string(expectedFileBytes)

	if fileContents != expectedFileContents {
		file.Save("/tmp/entity_examples2_repo.go")
		t.Error("Generated is not same as 'generated/entity_examples2_repo.go'. Check with /tmp/entity_examples2_repo.go")
	}
}
