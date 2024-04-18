package coin

import (
	"strings"

	"github.com/govalues/decimal"
)

type DecimalValues []decimal.Decimal

func (d DecimalValues) String() string {
	result := make([]string, len(d), len(d))

	for ix, value := range d {
		result[ix] = value.String()
	}

	return strings.Join(result, ",")
}
