package trig

import "math"

// piLo is π − math.Pi rounded to float64: math.Pi + piLo carries π to about 107 bits.
const piLo = 1.2246467991473532e-16

// thirdAngle returns the third angle z = π − x − y of a triangle with angles x and y (x + y < π) and
// sin z, both accurate to a few ulps relative even when z is tiny. x + y is formed exactly as s + e
// (Knuth's TwoSum), so z = ((math.Pi − s) + piLo) − e, where math.Pi − s is exact (Sterbenz) whenever
// s ≥ π/2, which is the only case in which z can be small. sin z is taken as sin s when s < π/2 (z is
// then near π, where sin z would lose digits) and as sin z otherwise. A z that rounds below zero is
// taken as zero.
func thirdAngle(x, y float64) (z, sinZ float64) {
	s := x + y
	yy := s - x
	e := (x - (s - yy)) + (y - yy)
	z = ((math.Pi - s) + piLo) - e
	if z <= 0 {
		return 0, 0
	}
	if s < math.Pi/2 {
		return z, math.Sin(s)
	}
	return z, math.Sin(z)
}
