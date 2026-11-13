package reduce

import (
	"math"
	"strconv"
	"testing"
)

// IEEE 754-2019 §5.7.2: isFinite is false exactly for NaN and the infinities.
func TestAllFinite_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   bool
	}{
		{"none", nil, true},
		{"zeros", []float64{0, math.Copysign(0, -1)}, true},
		{"extremes", []float64{math.MaxFloat64, -math.MaxFloat64, math.SmallestNonzeroFloat64}, true},
		{"NaN", []float64{1, math.NaN()}, false},
		{"+Inf", []float64{math.Inf(1)}, false},
		{"-Inf", []float64{2, math.Inf(-1)}, false},
	}
	for _, tc := range tests {
		if got := AllFinite(tc.values...); got != tc.want {
			t.Errorf("%s: AllFinite = %v, want %v", tc.name, got, tc.want)
		}
	}
	if !AllFinite() || !AllFinite(1, 2) || AllFinite(1, math.NaN()) {
		t.Errorf("scalar form wrong")
	}
}

func TestAllFinite_noAlloc(t *testing.T) {
	xs := []float64{1, 2, 3}
	if a := testing.AllocsPerRun(100, func() { _ = AllFinite(xs...); _ = AllFinite(1, 2) }); a != 0 {
		t.Errorf("AllFinite allocates %v times, want 0", a)
	}
}

func BenchmarkAllFinite(b *testing.B) {
	for _, n := range []int{1e3, 1e4, 1e5} {
		values := benchValues(n)
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = AllFinite(values...)
			}
		})
	}
}
