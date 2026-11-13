package signal

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// ConvolutionMode selects which part of the full linear convolution Convolve and Correlate return
// (numpy's mode argument).
type ConvolutionMode int

const (
	// ModeFull returns every output where the sequences overlap at all: n+m−1 values.
	ModeFull ConvolutionMode = iota
	// ModeSame returns max(n, m) values centred on the full output (offset ⌊(min(n, m)−1)/2⌋).
	ModeSame
	// ModeValid returns only the outputs where the sequences overlap completely: max(n, m)−min(n, m)+1 values.
	ModeValid
)

// convolveFFTMin is the length of the shorter input at and above which Convolve and Correlate switch
// from the direct O(n·m) sum to the FFT path. Below it the direct sum does fewer flops (64·64 = 4096
// multiply-adds against two complex FFTs of length 128 plus allocation) and is exact up to one rounding
// per term; the reviewer's benchmark confirms or moves the threshold.
const convolveFFTMin = 64

// Convolve returns the discrete linear convolution (a ∗ v)[k] = Σᵢ a[i]·v[k−i] in the requested mode:
// ModeFull (n+m−1 values, equal to the product of the polynomials with coefficients a and v),
// ModeSame (max(n, m) values, the centre of the full result) or ModeValid (max−min+1 values). The
// arguments commute. When min(n, m) < 64 the sum is direct (one fused multiply-add per term); otherwise
// both inputs go into one complex FFT of length L = 2^⌈log₂(n+m−1)⌉ (a in the real part, v in the
// imaginary part), the spectra are separated and multiplied, and one inverse FFT gives the result, whose
// error is then about ε·log₂L·‖a‖₂‖v‖₂ rather than per-term. The result goes into dst (grown when short);
// dst must not overlap a or v.
// Errors: ErrEmpty (a or v empty), ErrDimension (n+m−1 > 2⁴⁰), ErrMode (mode unknown), ErrNotFinite (an
// element of a or v), ErrOverflow (an output not finite). The result is nil with any error.
//
// ID:      NUM-SIG-014
// Time:    O(n·m) direct when min(n, m) < 64; O(L log L) on the FFT path, L < 2(n+m)   n = len(a), m = len(v)
// Space:   O(1) extra direct; O(L) extra on the FFT path; O(n+m) for the result
// Source:  Oppenheim & Schafer, Discrete-Time Signal Processing, 3rd ed., Pearson, 2010, §2.3 (convolution sum y[n] = Σ x[k]·h[n−k]); §8.7 (linear convolution by a DFT of length ≥ n+m−1)
// Tags:    convolve, convolution, filter, fir, polynomial-multiply, fft, signal
// Units:   the result has the units of a times the units of v
// Example: Convolve([1, 2, 3], [0, 1, 0.5], ModeFull) = [0, 1, 2.5, 4, 1.5]; ModeSame = [1, 2.5, 4]; ModeValid = [2.5]
func Convolve(dst, a, v []float64, mode ConvolutionMode) ([]float64, error) {
	return linearConvolve(dst, a, v, mode, false)
}

// Correlate returns the discrete cross-correlation c[k] = Σₙ a[n+k]·v[n] (numpy's correlate on real
// input), which is Convolve(a, reverse(v)) in the same mode: ModeFull lists lags k = −(m−1)…n−1, ModeSame
// and ModeValid take the same slices as Convolve. The arguments do not commute: swapping a and v reverses
// the ModeFull result. Paths (direct below min(n, m) = 64, FFT above), dst rules and errors as for
// Convolve.
// Errors: ErrEmpty, ErrDimension, ErrMode, ErrNotFinite, ErrOverflow. The result is nil with any error.
//
// ID:      NUM-SIG-015
// Time:    O(n·m) direct when min(n, m) < 64; O(L log L) on the FFT path, L < 2(n+m)   n = len(a), m = len(v)
// Space:   O(1) extra direct; O(L) extra on the FFT path; O(n+m) for the result
// Source:  Proakis & Manolakis, Digital Signal Processing: Principles, Algorithms, and Applications, 4th ed., Pearson, 2007, §2.6.1 (crosscorrelation r_xy(l) = Σ x(n)·y(n−l) = x(l) ∗ y(−l))
// Source:  Oppenheim & Schafer, Discrete-Time Signal Processing, 3rd ed., Pearson, 2010, §2.3 (convolution sum); §8.7 (linear convolution by the DFT)
// Tags:    correlate, cross-correlation, correlation, matched-filter, lag, fft, signal
// Units:   the result has the units of a times the units of v
// Example: Correlate([1, 2, 3], [0, 1, 0.5], ModeFull) = [0.5, 2, 3.5, 3, 0]; ModeSame = [2, 3.5, 3]; ModeValid = [3.5]
func Correlate(dst, a, v []float64, mode ConvolutionMode) ([]float64, error) {
	return linearConvolve(dst, a, v, mode, true)
}

// linearConvolve computes a ∗ v (reverseV false) or a ∗ reverse(v) (reverseV true) and slices it by mode.
func linearConvolve(dst, a, v []float64, mode ConvolutionMode, reverseV bool) ([]float64, error) {
	n, m := len(a), len(v)
	if n == 0 {
		return nil, fmt.Errorf("%w: a", ErrEmpty)
	}
	if m == 0 {
		return nil, fmt.Errorf("%w: v", ErrEmpty)
	}
	if n > maxLength || m > maxLength-n {
		return nil, fmt.Errorf("%w: len(a)+len(v)-1 > 2^40", ErrDimension)
	}
	full := n + m - 1
	short, long := min(n, m), max(n, m)
	var start, size int
	switch mode {
	case ModeFull:
		start, size = 0, full
	case ModeSame:
		start, size = (short-1)/2, long
	case ModeValid:
		start, size = short-1, long-short+1
	default:
		return nil, fmt.Errorf("%w: %d", ErrMode, int(mode))
	}
	if err := firstNotFinite("a", a); err != nil {
		return nil, err
	}
	if err := firstNotFinite("v", v); err != nil {
		return nil, err
	}
	out := growFloat(dst, size)
	if short < convolveFFTMin {
		// direct: out[k] = Σ a[i]·v'[start+k−i] over the overlap, v' = v or reverse(v).
		for k := range size {
			t := start + k
			lo, hi := max(0, t-m+1), min(n-1, t)
			s := 0.0
			if reverseV {
				for i := lo; i <= hi; i++ {
					s = math.FMA(a[i], v[m-1-t+i], s)
				}
			} else {
				for i := lo; i <= hi; i++ {
					s = math.FMA(a[i], v[t-i], s)
				}
			}
			out[k] = s
		}
	} else {
		l := nextPow2(full)
		z := make([]complex128, l)
		for i := range n {
			z[i] = complex(a[i], 0)
		}
		for j := range m {
			vj := v[j]
			if reverseV {
				vj = v[m-1-j]
			}
			z[j] = complex(real(z[j]), vj)
		}
		transform(z, z, -1)
		// Z = A + iV; A[k] = (Z[k] + conj Z[−k])/2, V[k] = (Z[k] − conj Z[−k])/(2i);
		// A·V = (Z[k]² − conj(Z[−k])²)/(4i). Pairs k and −k are updated together.
		for k := 0; k <= l/2; k++ {
			mk := (l - k) % l
			zk, zm := z[k], z[mk]
			pk := (zk*zk - conj(zm)*conj(zm)) / complex(0, 4)
			pm := (zm*zm - conj(zk)*conj(zk)) / complex(0, 4)
			z[k], z[mk] = pk, pm
		}
		transform(z, z, +1)
		inv := 1 / float64(l)
		for k := range size {
			out[k] = real(z[start+k]) * inv
		}
	}
	if !reduce.AllFinite(out...) {
		for i, x := range out {
			if math.IsNaN(x) || math.IsInf(x, 0) {
				return nil, fmt.Errorf("%w: output[%d]", ErrOverflow, i)
			}
		}
	}
	return out, nil
}

// firstNotFinite returns ErrNotFinite naming the first NaN or ±Inf element of x, or nil.
func firstNotFinite(name string, x []float64) error {
	if reduce.AllFinite(x...) {
		return nil
	}
	for i, v := range x {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return fmt.Errorf("%w: %s[%d] = %v", ErrNotFinite, name, i, v)
		}
	}
	return nil
}
