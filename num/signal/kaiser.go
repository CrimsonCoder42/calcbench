package signal

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/num/special"
)

// Kaiser returns the symmetric Kaiser–Bessel window of m points,
// w[n] = I₀(β·√(1 − (2n/(m−1) − 1)²)) / I₀(β), n = 0…m−1 (numpy's kaiser), with I₀ the modified
// Bessel function of order zero (NUM-SPEC-006). β = 0 is the rectangular window; larger β narrows the
// main lobe's taper (β ≈ 8.6 resembles Blackman). I₀ is even, so −β gives the same window as β. m = 0
// gives an empty slice and m = 1 gives [1]. The result goes into dst (grown to m when short).
// Errors: ErrDimension (m < 0 or m > 2⁴⁰), ErrNotFinite (beta NaN or ±Inf), ErrOverflow (|beta| above
// about 713.98, where I₀(β) overflows). The result is nil with any error.
//
// ID:      NUM-SIG-013
// Time:    O(m·c)   m = number of points, c ≤ 60 series terms per I₀ (|β| ≤ 20) or ≤ 30 asymptotic terms
// Space:   O(1) extra, O(m) for the result
// Source:  Kaiser, "Nonrecursive digital filter design using the I₀-sinh window function", Proc. IEEE Int. Symp. Circuits and Systems, 1974, 20–23
// Source:  Harris, "On the use of windows for harmonic analysis with the discrete Fourier transform", Proc. IEEE 66(1) (1978) 51–83, §V.H, eq. (46), doi:10.1109/PROC.1978.10837
// Tags:    window, kaiser, kaiser-bessel, bessel, taper, spectrum, filter-design, signal
// Units:   beta dimensionless; weights dimensionless in (0, 1]
// Range:   |beta| ≤ about 713.98
// Example: Kaiser(3, β = 1) = [1/I₀(1), 1, 1/I₀(1)] = [0.7898574933...,1, 0.7898574933...] (I₀(1) = 1.2660658777520082, A&S Table 9.8)
func Kaiser(dst []float64, m int, beta float64) ([]float64, error) {
	if math.IsNaN(beta) || math.IsInf(beta, 0) {
		return nil, fmt.Errorf("%w: beta = %v", ErrNotFinite, beta)
	}
	w, ok, err := checkWindowLength(dst, m)
	if err != nil || !ok {
		return w, err
	}
	den, err := special.BesselI0(beta)
	if err != nil {
		return nil, besselError(err, beta)
	}
	mm := float64(m - 1)
	for n := 0; 2*n < m; n++ {
		r := 2*float64(n)/mm - 1
		num, err := special.BesselI0(beta * math.Sqrt(math.Max(0, 1-r*r)))
		if err != nil {
			return nil, besselError(err, beta)
		}
		v := num / den
		w[n] = v
		w[m-1-n] = v
	}
	return w, nil
}

func besselError(err error, beta float64) error {
	if errors.Is(err, special.ErrOverflow) {
		return fmt.Errorf("%w: I0(beta), beta = %v", ErrOverflow, beta)
	}
	return fmt.Errorf("%w: beta = %v (%v)", ErrNotFinite, beta, err)
}
