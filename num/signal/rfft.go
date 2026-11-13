package signal

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// RFFT returns the first n/2+1 bins X[0…⌊n/2⌋] of the DFT of a real sequence x (the rest follow from
// X[n−k] = conj X[k]). For even n it packs z[j] = x[2j] + i·x[2j+1], takes one complex FFT of length n/2
// (radix-2 or Bluestein, NUM-SIG-001) and separates the even and odd halves,
// X[k] = Fe[k] + e^(−2πik/n)·Fo[k] with Fe[k] = (Z[k] + conj Z[n/2−k])/2, Fo[k] = (Z[k] − conj Z[n/2−k])/(2i);
// odd n takes a complex FFT of length n. The result goes into dst (grown to n/2+1 when short).
// Errors: ErrEmpty (n = 0), ErrDimension (n > 2⁴⁰), ErrNotFinite (an element of x), ErrOverflow. The
// result is nil with any error.
//
// ID:      NUM-SIG-003
// Time:    O(n log n)   n = len(x); even n costs about half a complex FFT of length n
// Space:   O(n) extra (the inner FFT's buffers; odd n also an n-element complex copy), O(n) for the result
// Source:  Press, Teukolsky, Vetterling, Flannery, Numerical Recipes, 3rd ed., Cambridge University Press, 2007, §12.3.2, eqs. (12.3.5)–(12.3.6) (FFT of a single real function)
// Tags:    rfft, real-fft, fft, dft, fourier, spectrum, signal
// Units:   X has the units of x
// Example: RFFT([1, 2, 3, 4]) = [10, −2+2i, −2]
func RFFT(dst []complex128, x []float64) ([]complex128, error) {
	n := len(x)
	if n == 0 {
		return nil, fmt.Errorf("%w: x", ErrEmpty)
	}
	if n > maxLength {
		return nil, fmt.Errorf("%w: len(x) = %d > 2^40", ErrDimension, n)
	}
	if !reduce.AllFinite(x...) {
		for i, v := range x {
			if math.IsNaN(v) || math.IsInf(v, 0) {
				return nil, fmt.Errorf("%w: x[%d] = %v", ErrNotFinite, i, v)
			}
		}
	}
	h := n / 2
	y := growComplex(dst, h+1)
	if n%2 == 1 {
		z := make([]complex128, n)
		for j, v := range x {
			z[j] = complex(v, 0)
		}
		transform(z, z, -1)
		copy(y, z[:h+1])
	} else {
		z := y[:h]
		for j := range h {
			z[j] = complex(x[2*j], x[2*j+1])
		}
		transform(z, z, -1)
		z0 := z[0]
		y[0] = complex(real(z0)+imag(z0), 0)
		y[h] = complex(real(z0)-imag(z0), 0)
		for k := 1; 2*k <= h; k++ {
			zk, zm := z[k], z[h-k]
			fe := (zk + conj(zm)) / 2
			fo := (zk - conj(zm)) / complex(0, 2)
			s, c := math.Sincos(-2 * math.Pi * float64(k) / float64(n))
			w := complex(c, s)
			y[k] = fe + w*fo
			if 2*k != h {
				// index h−k: Fe' = conj Fe, Fo' = (zm − conj zk)/(2i) = conj Fo, w^(h−k) = −conj w.
				y[h-k] = conj(fe - w*fo)
			}
		}
	}
	if err := checkResult("RFFT", y); err != nil {
		return nil, err
	}
	return y, nil
}

// IRFFT returns the real sequence x of length n whose RFFT is spectrum, which must hold n/2+1 bins
// (numpy's irfft with an explicit n). The imaginary parts of spectrum[0] and, for even n,
// spectrum[n/2] are ignored, as a real signal's spectrum has them zero. Even n undoes RFFT's packing
// with one inverse complex FFT of length n/2; odd n rebuilds the Hermitian spectrum of length n and takes
// one inverse FFT. Normalised by 1/n, so IRFFT(RFFT(x), len(x)) = x up to rounding. The result goes into
// dst (grown to n when short).
// Errors: ErrDimension (n < 1, n > 2⁴⁰ or len(spectrum) ≠ n/2+1), ErrNotFinite, ErrOverflow. The result
// is nil with any error.
//
// ID:      NUM-SIG-004
// Time:    O(n log n)   n = output length
// Space:   O(n) extra (one complex buffer of n/2 for even n, n for odd n, plus the inner FFT's), O(n) for the result
// Source:  Press, Teukolsky, Vetterling, Flannery, Numerical Recipes, 3rd ed., Cambridge University Press, 2007, §12.3.2, eqs. (12.3.5)–(12.3.7) (inverse of the real FFT)
// Tags:    irfft, inverse-real-fft, fft, dft, fourier, signal
// Units:   x has the units of the spectrum
// Example: IRFFT([10, −2+2i, −2], 4) = [1, 2, 3, 4]
func IRFFT(dst []float64, spectrum []complex128, n int) ([]float64, error) {
	if n < 1 || n > maxLength {
		return nil, fmt.Errorf("%w: n = %d", ErrDimension, n)
	}
	h := n / 2
	if len(spectrum) != h+1 {
		return nil, fmt.Errorf("%w: len(spectrum) = %d, want n/2+1 = %d", ErrDimension, len(spectrum), h+1)
	}
	if err := checkComplex("spectrum", spectrum); err != nil {
		return nil, err
	}
	out := growFloat(dst, n)
	if n%2 == 1 {
		full := make([]complex128, n)
		full[0] = complex(real(spectrum[0]), 0)
		for k := 1; k <= h; k++ {
			full[k] = spectrum[k]
			full[n-k] = conj(spectrum[k])
		}
		transform(full, full, +1)
		inv := 1 / float64(n)
		for j := range out {
			out[j] = real(full[j]) * inv
		}
	} else {
		z := make([]complex128, h)
		x0, xh := real(spectrum[0]), real(spectrum[h])
		z[0] = complex((x0+xh)/2, (x0-xh)/2)
		for k := 1; k < h; k++ {
			xk, xm := spectrum[k], spectrum[h-k]
			fe := (xk + conj(xm)) / 2
			s, c := math.Sincos(2 * math.Pi * float64(k) / float64(n)) // conj of e^(−2πik/n)
			fo := (xk - conj(xm)) * complex(c, s) / 2
			z[k] = fe + complex(0, 1)*fo
		}
		transform(z, z, +1)
		inv := 1 / float64(h)
		for j := range h {
			out[2*j] = real(z[j]) * inv
			out[2*j+1] = imag(z[j]) * inv
		}
	}
	for i, v := range out {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return nil, fmt.Errorf("%w: IRFFT[%d]", ErrOverflow, i)
		}
	}
	return out, nil
}
