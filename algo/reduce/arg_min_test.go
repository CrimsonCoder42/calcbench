package reduce

import (
	"errors"
	"math"
	"strconv"
	"testing"
)

// Knuth §1.2.10 Algorithm M scans once; on ties we keep the first index.
func TestArgMin_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   int
	}{
		{"one", []float64{7}, 0},
		{"two", []float64{2, 1}, 1},
		{"ties keep first", []float64{3, 1, 2, 1}, 1},
		{"sorted", []float64{1, 2, 3}, 0},
		{"reversed", []float64{3, 2, 1}, 2},
		{"negative", []float64{-1, -5, 0}, 1},
		{"-Inf", []float64{0, math.Inf(-1), -1}, 1},
		{"signed zeros equal", []float64{0, math.Copysign(0, -1)}, 0},
	}
	for _, tc := range tests {
		got, err := ArgMin(tc.values)
		if err != nil || got != tc.want {
			t.Errorf("%s: ArgMin = (%v, %v), want %d", tc.name, got, err, tc.want)
		}
	}
	if got, _ := ArgMin([]string{"b", "a", "c"}); got != 1 {
		t.Errorf("strings: ArgMin = %d, want 1", got)
	}
	if got, _ := ArgMin([]int{math.MinInt, math.MaxInt}); got != 0 {
		t.Errorf("int extremes: ArgMin = %d, want 0", got)
	}
}

func TestArgMin_errors(t *testing.T) {
	if _, err := ArgMin([]float64{}); !errors.Is(err, ErrEmpty) {
		t.Errorf("empty: err = %v", err)
	}
	if _, err := ArgMin([]float64{1, math.NaN()}); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN: err = %v", err)
	}
}

func BenchmarkArgMin(b *testing.B) {
	for _, n := range []int{1e3, 1e4, 1e5} {
		values := benchValues(n)
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = ArgMin(values)
			}
		})
	}
}
