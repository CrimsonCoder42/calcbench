//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/algebra"
)

func main() {
	fmt.Println(algebra.RoundDecimalsMode(-2.5, 0, algebra.RoundTiesToEven))
	fmt.Println(algebra.RoundDecimalsMode(-2.5, 0, algebra.RoundTiesToAway))
	fmt.Println(algebra.RoundDecimalsMode(-2.5, 0, algebra.RoundTiesTowardZero))
	fmt.Println(algebra.RoundDecimalsMode(1.125, 2, algebra.RoundTiesToEven))
	fmt.Println(algebra.RoundDecimalsMode(-1.21, 1, algebra.RoundTowardNegative))
	fmt.Println(algebra.RoundDecimalsMode(2.675, 2, algebra.RoundTiesToAway))
	fmt.Println(algebra.RoundDecimalsMode(1250, -2, algebra.RoundTowardPositive))
}
