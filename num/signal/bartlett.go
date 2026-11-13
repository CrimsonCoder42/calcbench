package signal

// Bartlett returns the symmetric Bartlett (triangular) window of m points,
// w[n] = 1 − |2n/(m−1) − 1|, n = 0…m−1, which is 0 at both ends and 1 at the centre (numpy's bartlett).
// m = 0 gives an empty slice and m = 1 gives [1]. The result goes into dst (grown to m when short).
// Errors: ErrDimension (m < 0 or m > 2⁴⁰). The result is nil with any error.
//
// ID:      NUM-SIG-012
// Time:    O(m)   m = number of points
// Space:   O(1) extra, O(m) for the result
// Source:  Harris, "On the use of windows for harmonic analysis with the discrete Fourier transform", Proc. IEEE 66(1) (1978) 51–83, §V.B, eq. (20), doi:10.1109/PROC.1978.10837
// Tags:    window, bartlett, triangular, fejer, taper, spectrum, signal
// Units:   dimensionless weights in [0, 1]
// Example: Bartlett(5) = [0, 0.5, 1, 0.5, 0]
func Bartlett(dst []float64, m int) ([]float64, error) {
	w, ok, err := checkWindowLength(dst, m)
	if err != nil || !ok {
		return w, err
	}
	den := float64(m - 1)
	for n := 0; 2*n < m; n++ {
		v := 2 * float64(n) / den
		if v > 1 {
			v = 2 - v
		}
		w[n] = v
		w[m-1-n] = v
	}
	return w, nil
}
