package apperrors

import "fmt"

type ErrorStrategyNotReady struct {
	StrategyName string
}

func (err ErrorStrategyNotReady) Error() string {
	return fmt.Sprintf(
		"strategy '%s' not ready",
		err.StrategyName,
	)
}
