# goutils/number

## Getting Started

main.go:
```golang
package main

import (
	"fmt"

	"github.com/pangpanglabs/goutils/number"
)

func main() {
	s := number.Setting{
		RoundDigit:    3,
		RoundStrategy: "ceil",
	}
	r := number.ToFixed(1/3.0, &s)
	fmt.Println("The result of 1/3.0 with 3 decimal places: ", r)
}
```