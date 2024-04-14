package apperrors

import "fmt"

type ErrorInvalidInput struct {
	InputName string
}

func (err ErrorInvalidInput) Error() string {
	return fmt.Sprintf(
		"invalid input name: '%s'",
		err.InputName,
	)
}

type ErrorInvalidInputs struct {
	InputsName []string
}

func (err ErrorInvalidInputs) Error() string {
	return fmt.Sprintf(
		"invalid inputs names: '%s'",
		err.InputsName,
	)
}
