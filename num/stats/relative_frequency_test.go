package stats

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Worked example: e-Handbook §1.3.3.14, "the count in a class divided by the total number of
// observations". Counts 1, 3, 1 → 0.2, 0.6, 0.2 (hand-checked).
func TestRelativeFrequency_worked(t *testing.T) {
	got, err := RelativeFrequency([]int{1, 3, 1})
	if err != nil || len(got) != 3 || got[0] != 0.2 || got[1] != 0.6 || got[2] != 0.2 {
		t.Errorf("RelativeFrequency(1 3 1) = %v, %v; want [0.2 0.6 0.2]", got, err)
	}
}

func TestRelativeFrequency_cases(t *testing.T) {
	for _, tt := range []struct {
		name   string
		counts []int
		want   []float64
	}{
		{"one class", []int{5}, []float64{1}},
		{"zeros kept", []int{0, 2, 0, 2}, []float64{0, 0.5, 0, 0.5}},
		{"largest total", []int{math.MaxInt - 1, 1}, []float64{1, 1 / float64(math.MaxInt)}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RelativeFrequency(tt.counts)
			if err != nil || len(got) != len(tt.want) {
				t.Fatalf("RelativeFrequency(%v) = %v, %v; want %v", tt.counts, got, err, tt.want)
			}
			for i := range got {
				// Tolerance 1e-15: one division each.
				if !near(got[i], tt.want[i], 1e-15) {
					t.Errorf("RelativeFrequency(%v)[%d] = %v; want %v", tt.counts, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestRelativeFrequency_errors(t *testing.T) {
	for _, tt := range []struct {
		name   string
		counts []int
		want   error
	}{
		{"empty", nil, ErrEmpty},
		{"negative", []int{1, -1}, ErrNegative},
		{"all zero", []int{0, 0}, ErrZeroTotal},
		{"total overflows", []int{math.MaxInt, 1}, ErrOverflow},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := RelativeFrequency(tt.counts); !errors.Is(err, tt.want) || got != nil {
				t.Errorf("RelativeFrequency(%v) = %v, %v; want nil, %v", tt.counts, got, err, tt.want)
			}
		})
	}
}

func TestRelativeFrequency_sumsToOne(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		counts := make([]int, 1+rng.IntN(50))
		for i := range counts {
			counts[i] = rng.IntN(1000)
		}
		counts[0]++ // total > 0
		got, err := RelativeFrequency(counts)
		var s float64
		for _, f := range got {
			s += f
		}
		// Tolerance 1e-13: b ≤ 50 roundings in the check's own sum.
		if err != nil || math.Abs(s-1) > 1e-13 {
			t.Fatalf("trial %d: RelativeFrequency sums to %v, %v; want 1", trial, s, err)
		}
	}
}

func BenchmarkRelativeFrequency(b *testing.B) {
	for _, size := range benchSizes {
		counts := make([]int, size.n)
		for i := range counts {
			counts[i] = i % 97
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				f, _ := RelativeFrequency(counts)
				benchSink = f[1]
			}
		})
	}
}
