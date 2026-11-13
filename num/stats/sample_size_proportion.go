package stats

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/num/prob"
)

// twoSidedZ returns the standard normal quantile z with P(−z < Z < z) = confidence.
func twoSidedZ(confidence float64) (float64, error) {
	if math.IsNaN(confidence) || math.IsInf(confidence, 0) {
		return 0, fmt.Errorf("%w: confidence = %v", ErrNotFinite, confidence)
	}
	if confidence <= 0 || confidence >= 1 {
		return 0, fmt.Errorf("%w: confidence = %v, want 0 < confidence < 1", ErrOutOfRange, confidence)
	}
	return prob.NormalQuantile(0.5+confidence/2, 0, 1)
}

// SampleSizeProportion returns the simple-random-sample size n needed to estimate a population proportion
// to within ±marginOfError at the given two-sided confidence: n₀ = z²·p(1 − p)/d², then, for a finite
// population of size population, n = n₀/(1 + (n₀ − 1)/N). population = 0 means an infinite population
// (n = n₀). The result is not rounded; the sample is the next whole number up, ⌈n⌉.
// Errors: ErrNotFinite for NaN or ±Inf input; ErrOutOfRange unless 0 < confidence < 1, 0 ≤ proportion ≤ 1,
// 0 < marginOfError ≤ 1 and population ≥ 0.
//
// ID:      NUM-STATS-047
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Cochran, Sampling Techniques, 3rd ed., Wiley, 1977, §4.4 "The Formula for n in Sampling for Proportions" (n₀ = t²PQ/d², n = n₀/(1 + (n₀ − 1)/N))
// Tags:    sample-size, proportion, margin-of-error, survey, finite-population-correction, cochran
// Units:   confidence, proportion and marginOfError as fractions (0.95, 0.5, 0.05); population a count; the result a count (unrounded)
// Example: 95 %, p = 0.5, d = 0.05, infinite → 1.959964²·0.25/0.0025 = 384.146 → 385 (hand-checked)
func SampleSizeProportion(confidence, proportion, marginOfError float64, population int64) (float64, error) {
	z, err := twoSidedZ(confidence)
	if err != nil {
		return 0, err
	}
	if math.IsNaN(proportion) || math.IsNaN(marginOfError) || math.IsInf(marginOfError, 0) {
		return 0, fmt.Errorf("%w: proportion = %v, marginOfError = %v", ErrNotFinite, proportion, marginOfError)
	}
	if proportion < 0 || proportion > 1 || marginOfError <= 0 || marginOfError > 1 || population < 0 {
		return 0, fmt.Errorf("%w: proportion = %v, marginOfError = %v, population = %d", ErrOutOfRange, proportion, marginOfError, population)
	}
	n0 := z * z * proportion * (1 - proportion) / (marginOfError * marginOfError)
	if population == 0 {
		return n0, nil
	}
	return n0 / (1 + (n0-1)/float64(population)), nil
}

// MarginOfErrorProportion returns the half-width d of the two-sided normal-approximation confidence
// interval for a proportion estimated from a simple random sample: d = z·√(p(1 − p)/n)·√((N − n)/(N − 1)),
// the last factor (the finite population correction) omitted when population = 0 (infinite).
// Errors: ErrNotFinite for NaN or ±Inf input; ErrOutOfRange unless 0 < confidence < 1, 0 ≤ proportion ≤ 1,
// sampleSize ≥ 1, population ≥ 0, and, when population ≠ 0, population ≥ 2 and sampleSize ≤ population.
//
// ID:      NUM-STATS-048
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Cochran, Sampling Techniques, 3rd ed., Wiley, 1977, §4.4 "The Formula for n in Sampling for Proportions" (d = t·√((N − n)/(N − 1))·√(PQ/n))
// Tags:    margin-of-error, proportion, confidence-interval, survey, finite-population-correction, cochran
// Units:   confidence and proportion as fractions; sampleSize and population counts; the result a fraction (0.05 = ±5 percentage points)
// Example: 95 %, p = 0.6, n = 100, infinite → 1.959964·√(0.24/100) = 0.0960182 (hand-checked)
func MarginOfErrorProportion(confidence, proportion float64, sampleSize, population int64) (float64, error) {
	z, err := twoSidedZ(confidence)
	if err != nil {
		return 0, err
	}
	if math.IsNaN(proportion) {
		return 0, fmt.Errorf("%w: proportion = %v", ErrNotFinite, proportion)
	}
	if proportion < 0 || proportion > 1 || sampleSize < 1 || population < 0 ||
		(population != 0 && (population < 2 || sampleSize > population)) {
		return 0, fmt.Errorf("%w: proportion = %v, sampleSize = %d, population = %d", ErrOutOfRange, proportion, sampleSize, population)
	}
	n := float64(sampleSize)
	d := z * math.Sqrt(proportion*(1-proportion)/n)
	if population != 0 {
		N := float64(population)
		d *= math.Sqrt((N - n) / (N - 1))
	}
	return d, nil
}
