package reduce

import (
	"errors"
	"math"
	"strconv"
	"testing"
)

// Knuth §1.2.10 Algorithm M: 3, 7, 7, 2 has its first maximum at index 1.
func TestArgMax_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   int
	}{
		{"knuth style", []float64{3, 7, 7, 2}, 1},
		{"one", []float64{7}, 0},
		{"sorted", []float64{1, 2, 3}, 2},
		{"reversed", []float64{3, 2, 1}, 0},
		{"+Inf", []float64{0, math.Inf(1), 1}, 1},
		{"signed zeros equal", []float64{math.Copysign(0, -1), 0}, 0},
	}
	for _, tc := range tests {
		got, err := ArgMax(tc.values)
		if err != nil || got != tc.want {
			t.Errorf("%s: ArgMax = (%v, %v), want %d", tc.name, got, err, tc.want)
		}
	}
	if got, _ := ArgMax([]int{math.MinInt, math.MaxInt, math.MaxInt}); got != 1 {
		t.Errorf("int extremes: ArgMax = %d, want 1", got)
	}
}

func TestArgMax_errors(t *testing.T) {
	if _, err := ArgMax([]int(nil)); !errors.Is(err, ErrEmpty) {
		t.Errorf("empty: err = %v", err)
	}
	if _, err := ArgMax([]float32{float32(math.NaN())}); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN: err = %v", err)
	}
}

func BenchmarkArgMax(b *testing.B) {
	for _, n := range []int{1e3, 1e4, 1e5} {
		values := benchValues(n)
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = ArgMax(values)
			}
		})
	}
}
