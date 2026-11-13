//go:build ignore

package main

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/num/stats"
)

func main() {
	for _, c := range []struct {
		conf, p, d float64
		N          int64
	}{{0.95, 0.5, 0.05, 0}, {0.95, 0.5, 0.05, 1000}, {0.99, 0.3, 0.02, 0}, {0.90, 0.5, 0.05, 10000}} {
		n, err := stats.SampleSizeProportion(c.conf, c.p, c.d, c.N)
		fmt.Println(c, n, math.Ceil(n), err)
	}
	fmt.Println(stats.MarginOfErrorProportion(0.95, 0.6, 100, 0))
	fmt.Println(stats.MarginOfErrorProportion(0.95, 0.6, 100, 1000))
	n, _ := stats.SampleSizeProportion(0.95, 0.5, 0.05, 10000)
	fmt.Println(stats.MarginOfErrorProportion(0.95, 0.5, int64(math.Round(n)), 10000))
}
