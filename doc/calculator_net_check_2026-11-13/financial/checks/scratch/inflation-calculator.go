//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	f, _ := interest.FutureValue(250000, 0.035, 12)
	p, _ := interest.PresentValue(250000, 0.035, 12)
	fmt.Println("forward", f, "backward", p)
}
