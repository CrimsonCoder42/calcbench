//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	v, _ := interest.FutureValue(1000000, 0.05-0.01, 10) // their convention: return net of the expense ratio
	fmt.Println(v)
}
