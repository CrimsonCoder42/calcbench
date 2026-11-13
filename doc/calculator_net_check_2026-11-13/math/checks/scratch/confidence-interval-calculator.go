//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/stats"
)

func main() {
	for _, c := range []struct {
		n          int
		mean, sd, cl float64
	}{{50, 20.6, 3.2, 0.95}, {10, 10, 2, 0.95}, {30, -4.5, 12.25, 0.99}} {
		fmt.Println(c)
		fmt.Println(" Z", fmt.Sprint(stats.MeanConfidenceIntervalZ(c.mean, c.sd, c.n, c.cl)))
		fmt.Println(" T", fmt.Sprint(stats.MeanConfidenceIntervalT(c.mean, c.sd, c.n, c.cl)))
	}
}
