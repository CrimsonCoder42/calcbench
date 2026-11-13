//go:build ignore

package main

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/num/stats"
)

func main() {
	cases := []struct {
		v    []float64
		want float64
	}{
		{[]float64{2, 8}, 4}, {[]float64{1, 3, 9}, 3}, {[]float64{1e300, 1e-300}, 1},
		{[]float64{0x1p600, 0x1p600}, 0x1p600},
		{[]float64{0x1p-1000, 0x1p-1060, 0x1p-1070, 0x1p1000}, math.Ldexp(math.Sqrt2, -533)},
		{[]float64{0x1p-3, 0x1p-4}, math.Sqrt2 / 16},
		{[]float64{10, 2, 38, 23, 38, 23, 21}, 16.412764443111},
	}
	for _, c := range cases {
		g, err := stats.GeometricMean(c.v)
		fmt.Printf("%v rel %.3g %v\n", g, math.Abs(g-c.want)/c.want, err)
	}
}
