//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/geometry"
)

func main() {
	fmt.Println(geometry.SphereVolume(2))
	fmt.Println(geometry.ConeVolume(3, 7.5))
	fmt.Println(geometry.ConeFrustumVolume(2, 5, 4))
	fmt.Println(geometry.EllipsoidVolume(1, 2, 3.5))
	fmt.Println(geometry.SphericalCapVolume(5, 2))
}
