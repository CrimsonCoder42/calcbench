package meteorology

import (
	"errors"
	"math"
	"testing"
)

// Hand computations from the Rothfusz regression and NWS WPC adjustments (see HeatIndexF Source):
// 85 °F 70 % → 92.7021; 100 °F 10 % (low-RH adjustment) → 94.1225; 84 °F 90 % (high-RH adjustment) → 98.3425;
// 70 °F 50 % (simple form, mean with T < 80) → 0.5·(70 + 61 + 2.4 + 4.7) = 69.05. Tolerance 1e-3 °F
// (hand arithmetic to 4 decimals).
func TestHeatIndexF_cases(t *testing.T) {
	cases := []struct{ t, rh, want float64 }{
		{85, 70, 92.7021},
		{100, 10, 94.1225},
		{84, 90, 98.3425},
		{70, 50, 69.05},
	}
	for _, c := range cases {
		got, err := HeatIndexF(c.t, c.rh)
		if err != nil || math.Abs(got-c.want) > 1e-3 {
			t.Errorf("%v °F %v %%: %v %v, want %v", c.t, c.rh, got, err, c.want)
		}
	}
}

func TestHeatIndexF_errors(t *testing.T) {
	cases := []struct {
		t, rh float64
		want  error
	}{
		{85, -1, ErrOutOfRange},
		{85, 101, ErrOutOfRange},
		{201, 50, ErrOutOfRange},
		{math.NaN(), 50, ErrNotFinite},
		{85, math.Inf(-1), ErrNotFinite},
	}
	for _, c := range cases {
		if _, err := HeatIndexF(c.t, c.rh); !errors.Is(err, c.want) {
			t.Errorf("(%v, %v): %v, want %v", c.t, c.rh, err, c.want)
		}
	}
}

func BenchmarkHeatIndexF(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = HeatIndexF(85, 70)
	}
}
