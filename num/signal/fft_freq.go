package signal

import (
	"fmt"
	"math"
)

// FFTFreq returns the frequency of each FFT bin for n samples taken sampleSpacing apart:
// f = [0, 1, …, ⌈n/2⌉−1, −⌊n/2⌋, …, −1] / (n·sampleSpacing), numpy's fftfreq order (non-negative
// frequencies first, then the negative ones). Bin k of an n-point DFT is the frequency k/(n·Δt), and bins
// above n/2 alias to k − n. The result goes into dst (grown to n when short).
// Errors: ErrDimension (n < 1 or n > 2⁴⁰), ErrNotFinite (sampleSpacing NaN or ±Inf), ErrOutOfRange
// (sampleSpacing = 0), ErrOverflow (a frequency not finite, from a subnormal spacing). The result is nil
// with any error.
//
// ID:      NUM-SIG-005
// Time:    O(n)   n = number of samples
// Space:   O(1) extra, O(n) for the result
// Source:  Oppenheim & Schafer, Discrete-Time Signal Processing, 3rd ed., Pearson, 2010, §8.1 and §8.4 (DFT bin k samples the DTFT at ω = 2πk/N), §4.2 (Ω = ω/T for sampling period T)
// Tags:    fftfreq, frequency, bins, fft, spectrum, signal
// Units:   sampleSpacing in any time unit (seconds, say); the result in cycles per that unit (hertz)
// Example: FFTFreq(n = 4, sampleSpacing = 0.1) = [0, 2.5, −5, −2.5]
func FFTFreq(dst []float64, n int, sampleSpacing float64) ([]float64, error) {
	if err := checkFreqArgs(n, sampleSpacing); err != nil {
		return nil, err
	}
	f := growFloat(dst, n)
	scale := float64(n) * sampleSpacing
	pos := (n + 1) / 2
	for i := range pos {
		f[i] = float64(i) / scale
	}
	for i := pos; i < n; i++ {
		f[i] = float64(i-n) / scale
	}
	return finishFreq(f)
}

// RFFTFreq returns the frequency of each RFFT bin for n samples taken sampleSpacing apart:
// f = [0, 1, …, ⌊n/2⌋] / (n·sampleSpacing), n/2+1 values (numpy's rfftfreq). The result goes into dst
// (grown to n/2+1 when short).
// Errors: as FFTFreq: ErrDimension, ErrNotFinite, ErrOutOfRange, ErrOverflow. The result is nil with any
// error.
//
// ID:      NUM-SIG-006
// Time:    O(n)   n = number of samples
// Space:   O(1) extra, O(n) for the result
// Source:  Oppenheim & Schafer, Discrete-Time Signal Processing, 3rd ed., Pearson, 2010, §8.1 and §8.4 (bin k at k/(N·T)); §8.6 (X[N−k] = conj X[k] for real x, so bins 0…N/2 suffice)
// Tags:    rfftfreq, frequency, bins, rfft, spectrum, signal
// Units:   sampleSpacing in any time unit; the result in cycles per that unit
// Example: RFFTFreq(n = 4, sampleSpacing = 0.1) = [0, 2.5, 5]
func RFFTFreq(dst []float64, n int, sampleSpacing float64) ([]float64, error) {
	if err := checkFreqArgs(n, sampleSpacing); err != nil {
		return nil, err
	}
	f := growFloat(dst, n/2+1)
	scale := float64(n) * sampleSpacing
	for i := range f {
		f[i] = float64(i) / scale
	}
	return finishFreq(f)
}

func checkFreqArgs(n int, sampleSpacing float64) error {
	if n < 1 || n > maxLength {
		return fmt.Errorf("%w: n = %d", ErrDimension, n)
	}
	if math.IsNaN(sampleSpacing) || math.IsInf(sampleSpacing, 0) {
		return fmt.Errorf("%w: sampleSpacing = %v", ErrNotFinite, sampleSpacing)
	}
	if sampleSpacing == 0 {
		return fmt.Errorf("%w: sampleSpacing = 0", ErrOutOfRange)
	}
	return nil
}

func finishFreq(f []float64) ([]float64, error) {
	for i, v := range f {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return nil, fmt.Errorf("%w: frequency[%d]", ErrOverflow, i)
		}
	}
	return f, nil
}
