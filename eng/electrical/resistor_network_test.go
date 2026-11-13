package electrical

import (
	"errors"
	"math"
	"testing"
)

// Hand computations from OpenStax University Physics Vol. 2, §10.2, eqs. (10.2)-(10.3):
// 10+2+38+23+38+23+21 = 155 Ω; parallel: Σ1/R = 0.1+0.5+2/38+2/23+1/21 = 0.78720715..., R = 1.27031366795 Ω;
// 6 ∥ 3 = 2 Ω. Tolerance 1e-12 relative.
func TestSeriesResistance_cases(t *testing.T) {
	got, err := SeriesResistance([]float64{10, 2, 38, 23, 38, 23, 21})
	if err != nil || got != 155 {
		t.Errorf("got %v %v", got, err)
	}
	if got, err := SeriesResistance([]float64{1e16, 1, 1, 1}); err != nil || got != 1e16+3 {
		t.Errorf("compensated: %v %v", got, err)
	}
	for _, c := range []struct {
		in   []float64
		want error
	}{
		{nil, ErrEmpty},
		{[]float64{1, -1}, ErrNegative},
		{[]float64{math.NaN()}, ErrNotFinite},
		{[]float64{math.MaxFloat64, math.MaxFloat64}, ErrOverflow},
	} {
		if _, err := SeriesResistance(c.in); !errors.Is(err, c.want) {
			t.Errorf("%v: %v, want %v", c.in, err, c.want)
		}
	}
}

func TestParallelResistance_cases(t *testing.T) {
	want := 1 / (0.1 + 0.5 + 2.0/38 + 2.0/23 + 1.0/21)
	got, err := ParallelResistance([]float64{10, 2, 38, 23, 38, 23, 21})
	if err != nil || !near(got, want) {
		t.Errorf("got %v %v, want %v", got, err, want)
	}
	if got, err := ParallelResistance([]float64{6, 3}); err != nil || !near(got, 2) {
		t.Errorf("6||3: %v %v", got, err)
	}
	if got, err := ParallelResistance([]float64{5}); err != nil || got != 5 {
		t.Errorf("one: %v %v", got, err)
	}
	if got, err := ParallelResistance([]float64{6, 0, 3}); err != nil || got != 0 {
		t.Errorf("short: %v %v", got, err)
	}
	if got, err := ParallelResistance([]float64{5e-324}); err != nil || got != 0 {
		t.Errorf("subnormal: %v %v", got, err)
	}
	if _, err := ParallelResistance(nil); !errors.Is(err, ErrEmpty) {
		t.Errorf("empty: %v", err)
	}
	if _, err := ParallelResistance([]float64{1, math.Inf(1)}); !errors.Is(err, ErrNotFinite) {
		t.Errorf("Inf: %v", err)
	}
}

func benchResistors(n int) []float64 {
	rs := make([]float64, n)
	for i := range rs {
		rs[i] = float64(i%97 + 1)
	}
	return rs
}

func BenchmarkSeriesResistance(b *testing.B) {
	for _, n := range []int{1000, 10000, 100000} {
		rs := benchResistors(n)
		b.Run(itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = SeriesResistance(rs)
			}
		})
	}
}

func BenchmarkParallelResistance(b *testing.B) {
	for _, n := range []int{1000, 10000, 100000} {
		rs := benchResistors(n)
		b.Run(itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = ParallelResistance(rs)
			}
		})
	}
}

func itoa(n int) string {
	switch n {
	case 1000:
		return "1e3"
	case 10000:
		return "1e4"
	default:
		return "1e5"
	}
}
