package apperrors

import (
	"fmt"
	"testing"
)

func TestErrorInvalidInputs(t *testing.T) {
	err := ErrorInvalidInputs{
		InputsName: []string{
			"a",
			"b",
		},
	}

	fmt.Println(err)
}
