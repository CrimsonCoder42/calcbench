//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/dp"
	"github.com/CrimsonCoder42/calcbench/num/series"
)

func main() {
	fmt.Println(series.ArithmeticSequenceTerm(2, 5, 20))
	fmt.Println(series.ArithmeticSeriesSum(2, 5, 20))
	fmt.Println(series.ArithmeticSequenceTerm(-1.5, 0.25, 7))
	fmt.Println(series.ArithmeticSeriesSum(-1.5, 0.25, 7))
	fmt.Println(series.GeometricSequenceTerm(2, 5, 12))
	fmt.Println(series.GeometricSeriesSum(2, 5, 12))
	fmt.Println(series.GeometricSequenceTerm(3, 0.5, 10))
	fmt.Println(series.GeometricSeriesSum(3, 0.5, 10))
	fmt.Println(dp.Fibonacci(10))
	fmt.Println(dp.Fibonacci(50))
}
