package apperrors

import "fmt"

type ErrValidation struct {
	Issue  error
	Caller string
}

const areaErrServiceValidation = "Validation"

func (e ErrValidation) Error() string {
	var res [3]string

	res[0] = fmt.Sprintf("Area: %s", areaErrServiceValidation)
	res[1] = fmt.Sprintf("Caller: %s", e.Caller)

	res[2] = ""
	if e.Issue != nil {
		res[2] = fmt.Sprintf("Issue: %s", e.Issue.Error())
	}

	return res[0] + _space + res[1] + _space + res[2]
}

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
