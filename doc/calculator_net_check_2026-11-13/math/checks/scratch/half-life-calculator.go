//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/logs"
)

func main() {
	fmt.Println(logs.ExponentialDecayHalfLife(100, 10, 50))
	fmt.Println(logs.ExponentialDecayElapsed(1, 0.2, 5730))
	fmt.Println(logs.ExponentialDecayRemaining(250, 1000, 5730))
	fmt.Println(logs.MeanLifetimeFromHalfLife(5730))
	fmt.Println(logs.DecayConstantFromHalfLife(5730))
}
