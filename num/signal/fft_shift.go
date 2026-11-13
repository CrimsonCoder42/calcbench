package signal

// FFTShift returns x cyclically shifted right by ⌊n/2⌋, which moves the zero-frequency bin of an FFT
// (or of FFTFreq) to the centre: out[(i + ⌊n/2⌋) mod n] = x[i]. The result goes into dst (grown to n when
// short). dst may be x itself (rotated in place by three reversals) but must not otherwise overlap it.
// An empty x gives an empty result. There are no errors.
//
// ID:      NUM-SIG-007
// Time:    O(n)   n = len(x)
// Space:   O(1) extra (no allocation when cap(dst) ≥ n), O(n) for the result
// Source:  Oppenheim & Schafer, Discrete-Time Signal Processing, 3rd ed., Pearson, 2010, §8.1 (the DFT is periodic in k with period N, so bins N/2…N−1 are the negative frequencies −N/2…−1)
// Tags:    fftshift, shift, rotate, fft, spectrum, centre, signal
// Example: FFTShift([0, 1, 2, 3, 4]) = [3, 4, 0, 1, 2]
func FFTShift[T any](dst, x []T) []T {
	n := len(x)
	return rotateRight(dst, x, n/2)
}

// IFFTShift returns x cyclically shifted right by ⌈n/2⌉ (left by ⌊n/2⌋), the inverse of FFTShift:
// out[(i + ⌈n/2⌉) mod n] = x[i]. For even n it equals FFTShift. dst rules as for FFTShift. An empty x
// gives an empty result. There are no errors.
//
// ID:      NUM-SIG-008
// Time:    O(n)   n = len(x)
// Space:   O(1) extra (no allocation when cap(dst) ≥ n), O(n) for the result
// Source:  Oppenheim & Schafer, Discrete-Time Signal Processing, 3rd ed., Pearson, 2010, §8.1 (periodicity of the DFT in k)
// Tags:    ifftshift, shift, rotate, fft, spectrum, signal
// Example: IFFTShift([3, 4, 0, 1, 2]) = [0, 1, 2, 3, 4]
func IFFTShift[T any](dst, x []T) []T {
	n := len(x)
	return rotateRight(dst, x, n-n/2)
}

// rotateRight writes x rotated right by r (0 ≤ r ≤ n) into dst.
func rotateRight[T any](dst, x []T, r int) []T {
	n := len(x)
	var out []T
	if cap(dst) >= n {
		out = dst[:n]
	} else {
		out = make([]T, n)
	}
	if n == 0 {
		return out
	}
	r %= n
	if &out[0] == &x[0] {
		reverse(out)
		reverse(out[:r])
		reverse(out[r:])
		return out
	}
	copy(out[r:], x[:n-r])
	copy(out[:r], x[n-r:])
	return out
}

func reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
