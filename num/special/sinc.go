package special

import "math"

// Sinc returns the normalized sinc function sin(πx)/(πx), with Sinc(0) = 1. The argument is first
// reduced exactly to r = x − n with n = round(x), r ∈ [−½, ½], and sin(πx) = (−1)ⁿ·sin(πr); every nonzero
// integer gives r = 0, so Sinc of every nonzero integer is exactly (+)0.
// Near 0 (|πx| < 1e-4) the Taylor polynomial 1 − y²/6 + y⁴/120 (y = πx) is used. Sinc(±Inf) = 0
// (the limit) and Sinc(NaN) = NaN; it has no error result, as R4 of the build brief fixes.
//
// ID:      NUM-SPEC-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Oppenheim & Schafer, Discrete-Time Signal Processing, 3rd ed., Pearson, 2010, §2.9 (ideal low-pass impulse response); NIST DLMF §4.5 / eq. 4.6.1 (sin series, for the small-argument branch)
// Tags:    special-function, sinc, signal, interpolation, window
// Units:   dimensionless
// Example: Sinc(0.5) = 2/π = 0.6366197723675814
func Sinc(x float64) float64 {
	if math.IsNaN(x) {
		return x
	}
	if math.IsInf(x, 0) {
		return 0
	}
	y := math.Pi * x
	if math.Abs(y) < 1e-4 {
		y2 := y * y
		return 1 - y2/6 + y2*y2/120
	}
	n := math.Round(x)
	r := x - n // exact (Sterbenz)
	if r == 0 {
		return 0
	}
	s := math.Sin(math.Pi * r)
	if math.Mod(n, 2) != 0 {
		s = -s
	}
	return s / y
}
