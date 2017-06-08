package number

import "math"

type Setting struct {
	RoundDigit    int
	RoundStrategy string
}

const (
	DefaultRoundDigit    = 2
	DefaultRoundStrategy = "round"
)

func ToFixed(num float64, setting *Setting) float64 {
	if setting == nil {
		setting = &Setting{
			RoundDigit:    DefaultRoundDigit,
			RoundStrategy: DefaultRoundStrategy,
		}
	}
	switch setting.RoundStrategy {
	case "ceil", "Ceil":
		output := math.Pow(10, float64(setting.RoundDigit))
		return math.Ceil(num*output) / output
	case "floor", "Floor":
		output := math.Pow(10, float64(setting.RoundDigit))
		return math.Floor(num*output) / output
	default:
		output := math.Pow(10, float64(setting.RoundDigit))
		return float64(Round(num*output)) / output
	}
}

func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
