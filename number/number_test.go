package number_test

import (
	"testing"

	"github.com/pangpanglabs/goutils/number"
	"github.com/pangpanglabs/goutils/test"
)

func TestNumber(t *testing.T) {
	r := number.ToFixed(10/3.0, nil)
	test.Equals(t, r, 3.33)

	s := number.Setting{
		RoundDigit:    3,
		RoundStrategy: "ceil",
	}
	r = number.ToFixed(1/3.0, &s)
	test.Equals(t, r, 0.334)

	s = number.Setting{
		RoundDigit:    1,
		RoundStrategy: "floor",
	}
	r = number.ToFixed(2/3.0, &s)
	test.Equals(t, r, 0.6)
}
