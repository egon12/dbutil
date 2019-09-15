package sync

import (
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	camel_case := ToSnakeCase("camelCase")
	if camel_case != "camel_case" {
		t.Error("'" + camel_case + "' is not same as expected camel_case")
	}

	name_with_snake_case := ToSnakeCase("nameWithSnakeCase")
	if name_with_snake_case != "name_with_snake_case" {
		t.Error("'" + name_with_snake_case + "' is not same as expected name_with_snake_case")
	}

	camel_case_with_hypen := ToSnakeCase("Camel_CaseWith_Hyphen")
	if camel_case_with_hypen != "camel_case_with_hyphen" {
		t.Error("'" + camel_case_with_hypen + "' is not same as expected camel_case_with_hyphen")
	}

}
