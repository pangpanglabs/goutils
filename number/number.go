package number

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

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
		f, _ := strconv.ParseFloat(strconv.FormatFloat(num*output, 'f', 6, 64), 64)
		return math.Ceil(f) / output
	case "floor", "Floor":
		output := math.Pow(10, float64(setting.RoundDigit))
		f, _ := strconv.ParseFloat(strconv.FormatFloat(num*output, 'f', 6, 64), 64)
		return math.Floor(f) / output
	case "round", "Round":
		output := math.Pow(10, float64(setting.RoundDigit))
		return float64(Round(num*output)) / output
	default: // BankRound
		s := fmt.Sprint(num)                      // "2.1965"
		pointPos := strings.Index(s, ".")         // 1
		roundPos := pointPos + setting.RoundDigit // 3

		if pointPos < 0 || roundPos+1 >= len(s) {
			return num
		}

		intNum, _ := strconv.ParseInt(s[:pointPos]+s[pointPos+1:roundPos+1], 10, 64) // 219

		switch s[roundPos] {
		case '0', '2', '4', '6', '8':
			if s[roundPos+1] >= '5' {
				intNum += 1
			}
		default:
			if s[roundPos+1] > '5' {
				intNum += 1
			}
		}

		return float64(intNum) / math.Pow(10, float64(setting.RoundDigit))
	}
}

func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
