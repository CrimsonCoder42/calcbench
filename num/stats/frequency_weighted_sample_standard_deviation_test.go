package stats

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// OpenStax Introductory Statistics 2e, Example 2.32: s² = 9.7375/19 = 0.5125, s = 0.715891 (printed to 6
// decimals). Tolerance 1e-14 relative against √0.5125, and the printed value to 6 decimals.
func TestFrequencyWeightedSampleStandardDeviation_sourceExample(t *testing.T) {
	got, err := FrequencyWeightedSampleStandardDeviation(fifthGradeAges, fifthGradeFreqs)
	if err != nil || !near(got, math.Sqrt(0.5125), 1e-14) || math.Abs(got-0.715891) > 5e-7 {
		t.Errorf("= (%v, %v), want √0.5125 ≈ 0.715891", got, err)
	}
}

func TestFrequencyWeightedSampleStandardDeviation_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		counts []int
		want   float64
	}{
		{"two equal", []float64{4}, []int{2}, 0},
		{"two values", []float64{1, 3}, []int{1, 1}, math.Sqrt2},
		{"large offset", []float64{1e9 + 4, 1e9 + 7, 1e9 + 13, 1e9 + 16}, []int{1, 1, 1, 1}, math.Sqrt(30)},
	}
	for _, tc := range tests {
		got, err := FrequencyWeightedSampleStandardDeviation(tc.values, tc.counts)
		if err != nil || !near(got, tc.want, 1e-12) {
			t.Errorf("%s: = (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// Oracle: NUM-STATS-019 on the expanded data (fixed seed, 300 tables). Tolerance 1e-12 relative.
func TestFrequencyWeightedSampleStandardDeviation_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		m := 1 + rng.IntN(30)
		values := randomFloats(rng, m)
		counts := make([]int, m)
		for i := range counts {
			counts[i] = rng.IntN(21)
		}
		counts[0] += 2
		want, _ := SampleStandardDeviation(expand(values, counts))
		got, err := FrequencyWeightedSampleStandardDeviation(values, counts)
		if err != nil || !near(got, want, 1e-12) {
			t.Fatalf("got (%v, %v), want %v", got, err, want)
		}
	}
}

func TestFrequencyWeightedSampleStandardDeviation_errors(t *testing.T) {
	f := FrequencyWeightedSampleStandardDeviation
	if _, err := f([]float64{5, 6}, []int{1, 0}); !errors.Is(err, ErrTooFew) {
		t.Errorf("total 1: %v", err)
	}
	if _, err := f(nil, nil); !errors.Is(err, ErrEmpty) {
		t.Errorf("empty: %v", err)
	}
	if _, err := f([]float64{math.NaN()}, []int{2}); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN: %v", err)
	}
	if _, err := f([]float64{1}, []int{-2}); !errors.Is(err, ErrNegative) {
		t.Errorf("negative: %v", err)
	}
	// s = MaxFloat64·√2 is beyond float64.
	if _, err := f([]float64{-math.MaxFloat64, math.MaxFloat64}, []int{1, 1}); !errors.Is(err, ErrOverflow) {
		t.Errorf("±MaxFloat64: %v, want ErrOverflow", err)
	}
}

func TestFrequencyWeightedSampleStandardDeviation_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() {
		benchSink, _ = FrequencyWeightedSampleStandardDeviation(fifthGradeAges, fifthGradeFreqs)
	}); a != 0 {
		t.Errorf("allocates %v times, want 0", a)
	}
}

func BenchmarkFrequencyWeightedSampleStandardDeviation(b *testing.B) {
	for _, size := range benchSizes {
		values, counts := frequencyBenchInput(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSink, _ = FrequencyWeightedSampleStandardDeviation(values, counts)
			}
		})
	}
}
