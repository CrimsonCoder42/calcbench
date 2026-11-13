package stats

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// fifthGrade is OpenStax Introductory Statistics 2e, §2.7, Example 2.32: ages of 20 fifth graders as a
// frequency table. Σf(x − x̄)² = 9.7375, x̄ = 10.525, s² = 0.5125, s = 0.715891.
var (
	fifthGradeAges  = []float64{9, 9.5, 10, 10.5, 11, 11.5}
	fifthGradeFreqs = []int{1, 2, 4, 4, 6, 3}
)

// expand repeats values[i] counts[i] times (the oracle's input).
func expand(values []float64, counts []int) []float64 {
	var out []float64
	for i, v := range values {
		for range counts[i] {
			out = append(out, v)
		}
	}
	return out
}

// OpenStax Example 2.32: population form σ = √(9.7375/20). Tolerance 1e-14 relative.
func TestFrequencyWeightedPopulationStandardDeviation_sourceExample(t *testing.T) {
	got, err := FrequencyWeightedPopulationStandardDeviation(fifthGradeAges, fifthGradeFreqs)
	if want := math.Sqrt(9.7375 / 20); err != nil || !near(got, want, 1e-14) {
		t.Errorf("= (%v, %v), want %v", got, err, want)
	}
}

// parity: ledger C9 (weighted.go:44), source WeightedSD: population SD over (value, count) pairs, rounded
// to 2 decimals by the source (rounding not ported). Case built from the source's documented promise "a
// value counted N times gives what N copies in a slice would": 120×3, 130×1 → mean 122.5, σ = √(75/4) ≈ 4.33.
func TestFrequencyWeightedPopulationStandardDeviation_parity(t *testing.T) {
	got, err := FrequencyWeightedPopulationStandardDeviation([]float64{120, 130}, []int{3, 1})
	if err != nil || math.Round(got*100)/100 != 4.33 || !near(got, math.Sqrt(18.75), 1e-14) {
		t.Errorf("= (%v, %v), want √18.75 ≈ 4.33", got, err)
	}
}

func TestFrequencyWeightedPopulationStandardDeviation_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		counts []int
		want   float64
	}{
		{"one value", []float64{7}, []int{5}, 0},
		{"zero count ignored", []float64{1, 1000, 3}, []int{1, 0, 1}, 1},
		{"unsorted duplicates", []float64{3, 1, 3}, []int{1, 1, 1}, math.Sqrt(8.0 / 9)},
		{"negative values", []float64{-2, 2}, []int{2, 2}, 2},
		{"large offset", []float64{1e9 + 4, 1e9 + 7, 1e9 + 13, 1e9 + 16}, []int{1, 1, 1, 1}, math.Sqrt(22.5)},
		{"huge values", []float64{-1e308, 1e308}, []int{1, 1}, 1e308},
		{"subnormal values", []float64{0, 4e-323}, []int{1, 1}, 2e-323},
	}
	for _, tc := range tests {
		got, err := FrequencyWeightedPopulationStandardDeviation(tc.values, tc.counts)
		if err != nil || !near(got, tc.want, 1e-12) {
			t.Errorf("%s: = (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// Oracle: NUM-STATS-022 on the expanded data, 300 random tables (fixed seed, m ≤ 30, counts ≤ 20).
// Tolerance 1e-12 relative (both are compensated two-pass forms).
func TestFrequencyWeightedPopulationStandardDeviation_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		m := 1 + rng.IntN(30)
		values := randomFloats(rng, m)
		counts := make([]int, m)
		for i := range counts {
			counts[i] = rng.IntN(21)
		}
		counts[0]++ // at least one count
		want, _ := PopulationStandardDeviation(expand(values, counts))
		got, err := FrequencyWeightedPopulationStandardDeviation(values, counts)
		if err != nil || !near(got, want, 1e-12) {
			t.Fatalf("got (%v, %v), want %v", got, err, want)
		}
	}
}

func TestFrequencyWeightedPopulationStandardDeviation_errors(t *testing.T) {
	f := FrequencyWeightedPopulationStandardDeviation
	if _, err := f(nil, nil); !errors.Is(err, ErrEmpty) {
		t.Errorf("empty: %v", err)
	}
	if _, err := f([]float64{1}, []int{1, 2}); !errors.Is(err, ErrLengthMismatch) {
		t.Errorf("mismatch: %v", err)
	}
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := f([]float64{1, v}, []int{1, 1}); !errors.Is(err, ErrNotFinite) {
			t.Errorf("value %v: %v", v, err)
		}
	}
	if _, err := f([]float64{1, 2}, []int{1, -1}); !errors.Is(err, ErrNegative) {
		t.Errorf("negative count: %v", err)
	}
	if _, err := f([]float64{1, 2}, []int{0, 0}); !errors.Is(err, ErrZeroTotal) {
		t.Errorf("zero total: %v", err)
	}
	if _, err := f([]float64{1, 2}, []int{math.MaxInt, 1}); !errors.Is(err, ErrOverflow) {
		t.Errorf("total overflow: %v", err)
	}
	if got, err := f([]float64{-math.MaxFloat64, math.MaxFloat64}, []int{1, 1}); err != nil && !errors.Is(err, ErrOverflow) {
		t.Errorf("±MaxFloat64: (%v, %v), want MaxFloat64 or ErrOverflow", got, err)
	}
}

func TestFrequencyWeightedPopulationStandardDeviation_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() {
		benchSink, _ = FrequencyWeightedPopulationStandardDeviation(fifthGradeAges, fifthGradeFreqs)
	}); a != 0 {
		t.Errorf("allocates %v times, want 0", a)
	}
}

// frequencyBenchInput returns m random values with counts 0…20 (fixed seed).
func frequencyBenchInput(m int) ([]float64, []int) {
	rng := rand.New(rand.NewPCG(1, 2))
	values := randomFloats(rng, m)
	counts := make([]int, m)
	for i := range counts {
		counts[i] = 1 + rng.IntN(20)
	}
	return values, counts
}

func BenchmarkFrequencyWeightedPopulationStandardDeviation(b *testing.B) {
	for _, size := range benchSizes {
		values, counts := frequencyBenchInput(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSink, _ = FrequencyWeightedPopulationStandardDeviation(values, counts)
			}
		})
	}
}
