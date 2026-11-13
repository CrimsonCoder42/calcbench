package stats

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// OpenStax Introductory Statistics 2e, Example 2.32 table, median by NIST §1.3.5.1 (N = 20 even: mean of
// the 10th and 11th values). Cumulative frequencies 1, 3, 7, 11, 17, 20: the 10th and 11th are 10.5, so
// the median is 10.5. Exact.
func TestFrequencyWeightedMedian_sourceExample(t *testing.T) {
	got, err := FrequencyWeightedMedian(fifthGradeAges, fifthGradeFreqs)
	if err != nil || got != 10.5 {
		t.Errorf("= (%v, %v), want 10.5", got, err)
	}
}

// parity: ledger C10 (weighted.go:107), the source's weightedMedian = median() over counted values.
// Cases built from its algorithm (it has no test file): 1×3, 5×1 → N = 4, ranks 1 and 2 (0-based) → 1, 1 →
// 1; 2×1, 8×1 → 5; 3×2, 4×1 → N = 3 odd → 3.
func TestFrequencyWeightedMedian_parity(t *testing.T) {
	for _, c := range []struct {
		values []float64
		counts []int
		want   float64
	}{
		{[]float64{1, 5}, []int{3, 1}, 1},
		{[]float64{2, 8}, []int{1, 1}, 5},
		{[]float64{3, 4}, []int{2, 1}, 3},
	} {
		got, err := FrequencyWeightedMedian(c.values, c.counts)
		if err != nil || got != c.want {
			t.Errorf("%v × %v: = (%v, %v), want %v", c.values, c.counts, got, err, c.want)
		}
	}
}

func TestFrequencyWeightedMedian_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		counts []int
		want   float64
	}{
		{"one value", []float64{4}, []int{1}, 4},
		{"unsorted", []float64{9, 1, 5}, []int{1, 1, 1}, 5},
		{"split between values", []float64{9, 1}, []int{2, 2}, 5},
		{"zero counts ignored", []float64{100, 1, 2}, []int{0, 1, 1}, 1.5},
		{"duplicate values", []float64{3, 3, 7}, []int{1, 1, 1}, 3},
		{"negative", []float64{-5, -1}, []int{1, 3}, -1},
		{"opposite huge", []float64{-math.MaxFloat64, math.MaxFloat64}, []int{1, 1}, 0},
	}
	for _, tc := range tests {
		got, err := FrequencyWeightedMedian(tc.values, tc.counts)
		if err != nil || got != tc.want {
			t.Errorf("%s: = (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// Oracle: NUM-STATS-012 on the expanded data, 400 random tables (fixed seed), many ties (values from a
// small set), m up to 200 so the selection loop runs; inputs unchanged. Exact comparison (both pick the
// same order statistics; the even case is the same midpoint).
func TestFrequencyWeightedMedian_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 400 {
		m := 1 + rng.IntN(200)
		values := make([]float64, m)
		counts := make([]int, m)
		for i := range values {
			values[i] = float64(rng.IntN(25))
			counts[i] = rng.IntN(6)
		}
		counts[0]++
		vb, cb := slices.Clone(values), slices.Clone(counts)
		want, _ := Median(expand(values, counts))
		got, err := FrequencyWeightedMedian(values, counts)
		if err != nil || got != want {
			t.Fatalf("got (%v, %v), want %v", got, err, want)
		}
		if !slices.Equal(values, vb) || !slices.Equal(counts, cb) {
			t.Fatal("input modified")
		}
	}
}

// The sort fallback: an adversarial order for median-of-3 still gives the right answer (the work budget
// switches to sorting).
func TestFrequencyWeightedMedian_sortedAndReversed(t *testing.T) {
	const m = 5000
	values := make([]float64, m)
	counts := make([]int, m)
	for i := range values {
		values[i] = float64(i)
		counts[i] = 1
	}
	if got, _ := FrequencyWeightedMedian(values, counts); got != 2499.5 {
		t.Errorf("sorted: %v, want 2499.5", got)
	}
	slices.Reverse(values)
	if got, _ := FrequencyWeightedMedian(values, counts); got != 2499.5 {
		t.Errorf("reversed: %v, want 2499.5", got)
	}
}

func TestFrequencyWeightedMedian_errors(t *testing.T) {
	f := FrequencyWeightedMedian
	if _, err := f(nil, nil); !errors.Is(err, ErrEmpty) {
		t.Errorf("empty: %v", err)
	}
	if _, err := f([]float64{1, 2}, []int{1}); !errors.Is(err, ErrLengthMismatch) {
		t.Errorf("mismatch: %v", err)
	}
	if _, err := f([]float64{math.Inf(-1)}, []int{1}); !errors.Is(err, ErrNotFinite) {
		t.Errorf("Inf: %v", err)
	}
	if _, err := f([]float64{1}, []int{-1}); !errors.Is(err, ErrNegative) {
		t.Errorf("negative: %v", err)
	}
	if _, err := f([]float64{1}, []int{0}); !errors.Is(err, ErrZeroTotal) {
		t.Errorf("zero: %v", err)
	}
	if _, err := f([]float64{1, 2}, []int{math.MaxInt, math.MaxInt}); !errors.Is(err, ErrOverflow) {
		t.Errorf("overflow: %v", err)
	}
}

func BenchmarkFrequencyWeightedMedian(b *testing.B) {
	for _, size := range benchSizes {
		values, counts := frequencyBenchInput(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSink, _ = FrequencyWeightedMedian(values, counts)
			}
		})
	}
}
