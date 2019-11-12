package number_test

import (
	"testing"

	"github.com/hillfolk/goutils/number"
	"github.com/hillfolk/goutils/test"
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
func TestBankRound(t *testing.T) {
	testdata := []struct{ a, b float64 }{
		{2.1965, 2.20},
		{2, 2},
		{2.2, 2.2},
		{2.20, 2.20},
		{2.27, 2.27},
		{2.195, 2.19},
		{2.185, 2.19},
		{2.175, 2.17},
		{2.165, 2.17},
		{100, 100},
		{.1, .1},
		{.175, .17},
		{.171, .17},
		{.166, .17},
		{.165, .17},
		{.17, .17},
	}

	for _, d := range testdata {
		fixedNumber := number.ToFixed(d.a, &number.Setting{
			RoundDigit:    2,
			RoundStrategy: "BankRound",
		})
		test.Equals(t, d.b, fixedNumber)
	}
}
