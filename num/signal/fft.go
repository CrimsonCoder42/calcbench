package signal

import (
	"errors"
	"fmt"
	"math"
	"math/bits"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Sentinel errors. Every error a function in this package returns wraps one of these.
var (
	// ErrEmpty is returned when an input that must hold at least one value is empty.
	ErrEmpty = errors.New("signal: empty input")
	// ErrNotFinite is returned when an input value is NaN or ±Inf.
	ErrNotFinite = errors.New("signal: value is NaN or infinite")
	// ErrDimension is returned when a length or size argument is negative, too large or inconsistent.
	ErrDimension = errors.New("signal: invalid length")
	// ErrOverflow is returned when a finite input produces a non-finite result.
	ErrOverflow = errors.New("signal: result overflows float64")
	// ErrOutOfRange is returned when a scalar parameter is outside its documented range.
	ErrOutOfRange = errors.New("signal: argument out of range")
	// ErrMode is returned when a ConvolutionMode is not one of ModeFull, ModeSame, ModeValid.
	ErrMode = errors.New("signal: unknown convolution mode")
)

// maxLength caps every transform length so the padded Bluestein length 2·n−1 rounded up to a power of
// two, and every index product, stays far inside int on 64-bit platforms.
const maxLength = 1 << 40

// growComplex returns dst resliced to length n when its capacity allows, else a new slice of length n.
func growComplex(dst []complex128, n int) []complex128 {
	if cap(dst) >= n {
		return dst[:n]
	}
	return make([]complex128, n)
}

// growFloat returns dst resliced to length n when its capacity allows, else a new slice of length n.
func growFloat(dst []float64, n int) []float64 {
	if cap(dst) >= n {
		return dst[:n]
	}
	return make([]float64, n)
}

// cmplxFinite reports whether both parts of z are finite.
func cmplxFinite(z complex128) bool {
	return reduce.AllFinite(real(z), imag(z))
}

// checkComplex returns ErrNotFinite naming the first non-finite element of x.
func checkComplex(name string, x []complex128) error {
	for i, z := range x {
		if !cmplxFinite(z) {
			return fmt.Errorf("%w: %s[%d] = %v", ErrNotFinite, name, i, z)
		}
	}
	return nil
}

// checkResult returns ErrOverflow naming the first non-finite element of a result computed from finite
// inputs.
func checkResult(name string, y []complex128) error {
	for i, z := range y {
		if !cmplxFinite(z) {
			return fmt.Errorf("%w: %s[%d]", ErrOverflow, name, i)
		}
	}
	return nil
}

// twiddles returns w[j] = exp(sign·2πi·j/n) for 0 ≤ j < n/2, each from its own Sincos call (no
// recurrence, so no error growth across the table). n is a power of two ≥ 2.
func twiddles(n int, sign float64) []complex128 {
	w := make([]complex128, n/2)
	for j := range w {
		s, c := math.Sincos(sign * 2 * math.Pi * float64(j) / float64(n))
		w[j] = complex(c, s)
	}
	return w
}

// radix2 transforms x (length n, a power of two) in place with the iterative decimation-in-time
// Cooley–Tukey algorithm: bit-reversal permutation, then log₂ n butterfly stages. tw is twiddles(n, ±1);
// its sign picks the direction. No normalisation is applied.
func radix2(x []complex128, tw []complex128) {
	n := len(x)
	if n < 2 {
		return
	}
	shift := uint(64 - bits.TrailingZeros(uint(n)))
	for i := range n {
		j := int(bits.Reverse64(uint64(i)) >> shift)
		if j > i {
			x[i], x[j] = x[j], x[i]
		}
	}
	for size := 2; size <= n; size <<= 1 {
		half := size >> 1
		step := n / size
		for start := 0; start < n; start += size {
			for j := range half {
				t := tw[j*step] * x[start+j+half]
				u := x[start+j]
				x[start+j] = u + t
				x[start+j+half] = u - t
			}
		}
	}
}

// isPow2 reports whether n ≥ 1 is a power of two.
func isPow2(n int) bool { return n > 0 && n&(n-1) == 0 }

// nextPow2 returns the smallest power of two ≥ n, for 1 ≤ n ≤ maxLength·2.
func nextPow2(n int) int {
	if n <= 1 {
		return 1
	}
	return 1 << bits.Len(uint(n-1))
}

// transform writes the unnormalised DFT of x into y (len(y) == len(x) ≥ 1, finite input). sign = −1
// is the forward transform e^(−2πijk/n), +1 the inverse kernel. y may be x itself.
func transform(y, x []complex128, sign float64) {
	n := len(x)
	if isPow2(n) {
		copy(y, x)
		radix2(y, twiddles(n, sign))
		return
	}
	bluestein(y, x, sign)
}

// bluestein writes the length-n DFT of x into y with Bluestein's chirp-z identity
// jk = (j² + k² − (k−j)²)/2: X[k] = w[k]·Σ_j (x[j]·w[j])·conj(w[k−j]), w[k] = e^(sign·πi·k²/n), a
// circular convolution of length L = 2^⌈log₂(2n−1)⌉ done with three radix-2 transforms. k² is reduced
// mod 2n incrementally so the chirp angle stays exact for every n.
func bluestein(y, x []complex128, sign float64) {
	n := len(x)
	l := nextPow2(2*n - 1)
	chirp := make([]complex128, n)
	var k2 uint64 // k² mod 2n
	m := uint64(2 * n)
	for k := range n {
		s, c := math.Sincos(sign * math.Pi * float64(k2) / float64(n))
		chirp[k] = complex(c, s)
		k2 = (k2 + 2*uint64(k) + 1) % m
	}
	a := make([]complex128, l)
	b := make([]complex128, l)
	for k := range n {
		a[k] = x[k] * chirp[k]
	}
	b[0] = conj(chirp[0])
	for k := 1; k < n; k++ {
		b[k] = conj(chirp[k])
		b[l-k] = b[k]
	}
	fwd := twiddles(l, -1)
	radix2(a, fwd)
	radix2(b, fwd)
	for k := range a {
		a[k] *= b[k]
	}
	// inverse length-l transform by conjugation: IDFT(z) = conj(DFT(conj z))/l.
	for k := range a {
		a[k] = conj(a[k])
	}
	radix2(a, fwd)
	inv := 1 / float64(l)
	for k := range n {
		y[k] = chirp[k] * complex(real(a[k])*inv, -imag(a[k])*inv)
	}
}

func conj(z complex128) complex128 { return complex(real(z), -imag(z)) }

// FFT returns the discrete Fourier transform X[k] = Σⱼ x[j]·e^(−2πi·jk/n), k = 0…n−1, unnormalised
// (numpy's default "backward" norm). Powers of two use the iterative radix-2 Cooley–Tukey algorithm; any
// other length uses Bluestein's chirp-z algorithm on a radix-2 transform of length 2^⌈log₂(2n−1)⌉, so
// every n is O(n log n). The result goes into dst (grown to n when its capacity is short); dst may be x
// itself but must not otherwise overlap it. x is not modified unless it is dst.
// Errors: ErrEmpty (n = 0), ErrDimension (n > 2⁴⁰), ErrNotFinite (an element of x), ErrOverflow (a
// non-finite output from finite input). The result is nil with any error.
//
// ID:      NUM-SIG-001
// Time:    O(n log n)   n = len(x); about 5·n·log₂ n flops for powers of two, about 3 times that via Bluestein
// Space:   O(n) extra (twiddle table n/2; Bluestein three buffers of length L < 4n), O(n) for the result
// Source:  Cooley & Tukey, "An algorithm for the machine calculation of complex Fourier series", Math. Comp. 19 (1965) 297–301, doi:10.1090/S0025-5718-1965-0178586-1
// Source:  Bluestein, "A linear filtering approach to the computation of discrete Fourier transform", IEEE Trans. Audio Electroacoust. AU-18(4) (1970) 451–455, doi:10.1109/TAU.1970.1162132
// Tags:    fft, dft, fourier, spectrum, radix-2, bluestein, chirp-z, signal
// Units:   X has the units of x; index k is the frequency bin k/(n·Δt) (see FFTFreq)
// Example: FFT([1, 2, 3, 4]) = [10, −2+2i, −2, −2−2i] (direct DFT sum, hand-checked)
func FFT(dst, x []complex128) ([]complex128, error) {
	if err := checkTransformInput(x); err != nil {
		return nil, err
	}
	y := growComplex(dst, len(x))
	transform(y, x, -1)
	if err := checkResult("FFT", y); err != nil {
		return nil, err
	}
	return y, nil
}

// IFFT returns the inverse discrete Fourier transform x[j] = (1/n)·Σₖ X[k]·e^(+2πi·jk/n), j = 0…n−1, so
// IFFT(FFT(x)) = x up to rounding. Algorithms, dst rules and lengths as for FFT.
// Errors: ErrEmpty, ErrDimension (n > 2⁴⁰), ErrNotFinite, ErrOverflow. The result is nil with any error.
//
// ID:      NUM-SIG-002
// Time:    O(n log n)   n = len(spectrum)
// Space:   O(n) extra, O(n) for the result
// Source:  Oppenheim & Schafer, Discrete-Time Signal Processing, 3rd ed., Pearson, 2010, §8.2, eqs. (8.67)–(8.68) (DFT pair)
// Source:  Cooley & Tukey 1965 (radix-2); Bluestein 1970 (other lengths), as for FFT (NUM-SIG-001)
// Tags:    ifft, inverse-fft, dft, fourier, radix-2, bluestein, signal
// Units:   x has the units of the spectrum
// Example: IFFT([10, −2+2i, −2, −2−2i]) = [1, 2, 3, 4]
func IFFT(dst, spectrum []complex128) ([]complex128, error) {
	if err := checkTransformInput(spectrum); err != nil {
		return nil, err
	}
	n := len(spectrum)
	y := growComplex(dst, n)
	transform(y, spectrum, +1)
	inv := 1 / float64(n)
	for k := range y {
		y[k] = complex(real(y[k])*inv, imag(y[k])*inv)
	}
	if err := checkResult("IFFT", y); err != nil {
		return nil, err
	}
	return y, nil
}

func checkTransformInput(x []complex128) error {
	if len(x) == 0 {
		return fmt.Errorf("%w: x", ErrEmpty)
	}
	if len(x) > maxLength {
		return fmt.Errorf("%w: len(x) = %d > 2^40", ErrDimension, len(x))
	}
	return checkComplex("x", x)
}
