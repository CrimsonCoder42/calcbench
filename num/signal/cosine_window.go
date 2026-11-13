package signal

import (
	"fmt"
	"math"
)

// checkWindowLength validates a window length m and sizes dst. ok is false when m is 0 or 1, in which
// case w is already the full answer (empty, or [1]).
func checkWindowLength(dst []float64, m int) (w []float64, ok bool, err error) {
	if m < 0 || m > maxLength {
		return nil, false, fmt.Errorf("%w: m = %d", ErrDimension, m)
	}
	w = growFloat(dst, m)
	if m == 1 {
		w[0] = 1
	}
	return w, m > 1, nil
}

// cosineWindow fills w with a₀ − a₁cos(2πn/(M−1)) + a₂cos(4πn/(M−1)), n = 0…M−1 (M ≥ 2), filling the
// second half from the first so the window is exactly symmetric.
func cosineWindow(w []float64, a0, a1, a2 float64) {
	m := len(w)
	den := float64(m - 1)
	for n := 0; 2*n < m; n++ {
		t := 2 * math.Pi * float64(n) / den
		v := a0 - a1*math.Cos(t)
		if a2 != 0 {
			v += a2 * math.Cos(2*t)
		}
		w[n] = v
		w[m-1-n] = v
	}
}

// Hanning returns the symmetric Hann ("hanning") window of m points,
// w[n] = 0.5 − 0.5·cos(2πn/(m−1)), n = 0…m−1 (numpy's hanning). m = 0 gives an empty slice and m = 1
// gives [1]. The result goes into dst (grown to m when short).
// Errors: ErrDimension (m < 0 or m > 2⁴⁰). The result is nil with any error.
//
// ID:      NUM-SIG-009
// Time:    O(m)   m = number of points; m/2 cosines
// Space:   O(1) extra, O(m) for the result
// Source:  Harris, "On the use of windows for harmonic analysis with the discrete Fourier transform", Proc. IEEE 66(1) (1978) 51–83, §V.C, eq. (27a), doi:10.1109/PROC.1978.10837
// Tags:    window, hann, hanning, raised-cosine, taper, spectrum, signal
// Units:   dimensionless weights in [0, 1]
// Example: Hanning(5) = [0, 0.5, 1, 0.5, 0]
func Hanning(dst []float64, m int) ([]float64, error) {
	w, ok, err := checkWindowLength(dst, m)
	if err != nil || !ok {
		return w, err
	}
	cosineWindow(w, 0.5, 0.5, 0)
	return w, nil
}

// Hamming returns the symmetric Hamming window of m points, w[n] = 0.54 − 0.46·cos(2πn/(m−1)),
// n = 0…m−1 (numpy's hamming). m = 0 gives an empty slice and m = 1 gives [1]. The result goes into dst
// (grown to m when short).
// Errors: ErrDimension (m < 0 or m > 2⁴⁰). The result is nil with any error.
//
// ID:      NUM-SIG-010
// Time:    O(m)   m = number of points
// Space:   O(1) extra, O(m) for the result
// Source:  Harris, "On the use of windows for harmonic analysis with the discrete Fourier transform", Proc. IEEE 66(1) (1978) 51–83, §V.D, eq. (30a), doi:10.1109/PROC.1978.10837
// Tags:    window, hamming, raised-cosine, taper, spectrum, signal
// Units:   dimensionless weights in [0.08, 1]
// Example: Hamming(5) = [0.08, 0.54, 1, 0.54, 0.08]
func Hamming(dst []float64, m int) ([]float64, error) {
	w, ok, err := checkWindowLength(dst, m)
	if err != nil || !ok {
		return w, err
	}
	cosineWindow(w, 0.54, 0.46, 0)
	return w, nil
}

// Blackman returns the symmetric three-term Blackman window of m points,
// w[n] = 0.42 − 0.5·cos(2πn/(m−1)) + 0.08·cos(4πn/(m−1)), n = 0…m−1 (numpy's blackman). The end
// points are 0 up to rounding (about ±1e-17, as in numpy). m = 0 gives an empty slice and m = 1 gives [1].
// The result goes into dst (grown to m when short).
// Errors: ErrDimension (m < 0 or m > 2⁴⁰). The result is nil with any error.
//
// ID:      NUM-SIG-011
// Time:    O(m)   m = number of points
// Space:   O(1) extra, O(m) for the result
// Source:  Harris, "On the use of windows for harmonic analysis with the discrete Fourier transform", Proc. IEEE 66(1) (1978) 51–83, §V.E, eq. (32), doi:10.1109/PROC.1978.10837
// Tags:    window, blackman, cosine-sum, taper, spectrum, signal
// Units:   dimensionless weights in [0, 1]
// Example: Blackman(5) = [0, 0.34, 1, 0.34, 0]
func Blackman(dst []float64, m int) ([]float64, error) {
	w, ok, err := checkWindowLength(dst, m)
	if err != nil || !ok {
		return w, err
	}
	cosineWindow(w, 0.42, 0.5, 0.08)
	return w, nil
}
