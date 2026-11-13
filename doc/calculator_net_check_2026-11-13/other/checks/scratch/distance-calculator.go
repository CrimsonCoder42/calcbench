//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/sci/geoscience"
)

func main() {
	d, err := geoscience.VincentyInverseDistance(55.75, 0, -(33 + 26.0/60), 108+13.0/60, 6377397.155, 1/299.1528128)
	fmt.Printf("vincenty line a %.4f %v\n", d, err)
	for _, p := range [][4]float64{{38.8976, -77.0366, 39.9496, -75.1503}, {51.5074, -0.1278, 40.7128, -74.0060}, {-33.8688, 151.2093, 35.6762, 139.6503}, {0, 0, 0, 1}, {0, 0, 0.5, 179.7}} {
		w, err := geoscience.WGS84Distance(p[0], p[1], p[2], p[3])
		g, _ := geoscience.GreatCircleDistance(p[0], p[1], p[2], p[3], geoscience.MeanEarthRadiusMeters)
		fmt.Printf("%v wgs84 %.3f m (%.4f km, %.4f mi) sphere %.3f %v\n", p, w, w/1000, w/1609.344, g, err)
	}
}
