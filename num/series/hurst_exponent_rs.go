package series

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
	"github.com/CrimsonCoder42/calcbench/num/stats"
)

// HurstExponentRS returns the Hurst exponent H of a series estimated by rescaled-range (R/S) analysis.
// For each chunk size s in chunkSizes the series is cut into ⌊n/s⌋ consecutive non-overlapping chunks
// (a remainder at the end is not used). In each chunk, with mean x̄ and cumulative deviations
// Y_j = Σ_{i≤j}(x_i − x̄), j = 1…s, the range is R = max(0, Y_1…Y_s) − min(0, Y_1…Y_s) and S is the
// population standard deviation of the chunk; chunks with S = 0 are skipped. The average R/S over the
// chunks of size s is one point (ln s, ln R/S), and H is the least-squares slope of ln R/S on ln s
// (NUM-STATS-009 LinearRegression). H ≈ 0.5 for independent increments, > 0.5 for persistence, < 0.5 for
// anti-persistence. The estimator is biased upward for small s (no Anis–Lloyd correction is applied).
// values and chunkSizes are read, never modified.
// Errors: ErrNotFinite when a value is NaN or ±Inf; ErrTooFew when fewer than 2 chunk sizes are given;
// ErrOutOfRange when a chunk size is < 2 or > len(values), or two chunk sizes are equal; ErrZeroVariance
// when every chunk of some size is constant; ErrOverflow when a deviation or a sum exceeds the float64 range.
//
// ID:      NUM-SERIES-017
// Time:    O(n·c)   n = len(values), c = len(chunkSizes); each size reads the series twice
// Space:   O(c): two slices of c logarithms (two allocations)
// Source:  Hurst, "Long-Term Storage Capacity of Reservoirs", Transactions of the American Society of Civil Engineers, 1951, 116:770–808 (the adjusted range R of cumulative departures from the mean, scaled by the standard deviation; R/S ∝ N^K)
// Source:  Mandelbrot & Wallis, "Robustness of the Rescaled Range R/S in the Measurement of Noncyclic Long Run Statistical Dependence", Water Resources Research, 1969, 5(5):967–988, doi:10.1029/WR005i005p00967 (R/S over sub-records of lag s, the slope of log R/S against log s, the pox diagram)
// Tags:    time-series, hurst-exponent, rescaled-range, r-s-analysis, long-memory, persistence, fractal
// Units:   values in any one unit; chunkSizes in samples; result dimensionless
// Example: 1, 2, 3, 4, chunkSizes 2, 4 → R/S(2) = 0.5/0.5 = 1, R/S(4) = 2/√1.25 → H = ln(4/√5)/ln 2 = 2 − ½·log₂5 ≈ 0.839036 (hand-checked)
func HurstExponentRS(values []float64, chunkSizes []int) (float64, error) {
	n := len(values)
	for i, v := range values {
		if !reduce.AllFinite(v) {
			return 0, fmt.Errorf("%w: values[%d] = %v", ErrNotFinite, i, v)
		}
	}
	c := len(chunkSizes)
	if c < 2 {
		return 0, fmt.Errorf("%w: %d chunk sizes, want ≥ 2", ErrTooFew, c)
	}
	for i, s := range chunkSizes {
		if s < 2 || s > n {
			return 0, fmt.Errorf("%w: chunkSizes[%d] = %d, want 2 ≤ s ≤ %d", ErrOutOfRange, i, s, n)
		}
		for j := 0; j < i; j++ {
			if chunkSizes[j] == s {
				return 0, fmt.Errorf("%w: chunkSizes[%d] = chunkSizes[%d] = %d", ErrOutOfRange, j, i, s)
			}
		}
	}
	logS := make([]float64, c)
	logRS := make([]float64, c)
	for i, s := range chunkSizes {
		rs, err := meanRescaledRange(values, s)
		if err != nil {
			return 0, err
		}
		logS[i] = math.Log(float64(s))
		logRS[i] = math.Log(rs)
	}
	_, h, err := stats.LinearRegression(logS, logRS)
	if err != nil { // inputs are finite and distinct, so only an overflow can remain
		return 0, fmt.Errorf("%w: slope: %v", ErrOverflow, err)
	}
	return h, nil
}

// meanRescaledRange returns the mean R/S over the non-overlapping chunks of size s of values.
func meanRescaledRange(values []float64, s int) (float64, error) {
	var sum, comp float64
	used := 0
	for start := 0; start+s <= len(values); start += s {
		chunk := values[start : start+s]
		var ms, mc float64
		for _, v := range chunk {
			ms, mc = reduce.NeumaierAdd(ms, mc, v)
		}
		mean := (ms + mc) / float64(s)
		var y, yc, ss, sc float64
		lo, hi := 0.0, 0.0
		for _, v := range chunk {
			d := v - mean
			y, yc = reduce.NeumaierAdd(y, yc, d)
			cum := y + yc
			lo, hi = math.Min(lo, cum), math.Max(hi, cum)
			ss, sc = reduce.NeumaierAdd(ss, sc, d*d)
		}
		sd := math.Sqrt((ss + sc) / float64(s))
		r := hi - lo
		if !reduce.AllFinite(mean, sd, r) {
			return 0, fmt.Errorf("%w: chunk of size %d at index %d", ErrOverflow, s, start)
		}
		if sd == 0 {
			continue
		}
		sum, comp = reduce.NeumaierAdd(sum, comp, r/sd)
		used++
	}
	if used == 0 {
		return 0, fmt.Errorf("%w: every chunk of size %d is constant", ErrZeroVariance, s)
	}
	return (sum + comp) / float64(used), nil
}
