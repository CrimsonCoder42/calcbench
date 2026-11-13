//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/geometry"
)

func main() {
	fmt.Println(geometry.RightTriangleHypotenuse(3, 4))
	fmt.Println(geometry.RightTriangleHypotenuse(1e-5, 2.5e-5))
	fmt.Println(geometry.RightTriangleLeg(13, 5))
	fmt.Println(geometry.RightTriangleLeg(100000.001, 100000))
}
