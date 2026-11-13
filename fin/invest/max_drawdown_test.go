package invest

import (
	"fmt"
	"math"
	"testing"
)

func TestMaxDrawdown_handWorked(t *testing.T) {
	dd, peak, trough, err := MaxDrawdown([]float64{100, 120, 90, 110, 80, 130, 125})
	if err != nil || !near(dd, 1.0/3, 1e-15) || peak != 1 || trough != 4 {
		t.Errorf("MaxDrawdown = %v, %d, %d, %v; want 1/3, 1, 4", dd, peak, trough, err)
	}
}

func TestMaxDrawdown_cases(t *testing.T) {
	for _, c := range []struct {
		name          string
		w             []float64
		dd            float64
		peak, trought int
	}{
		{"one value", []float64{100}, 0, 0, 0},
		{"rising", []float64{1, 2, 3}, 0, 0, 0},
		{"flat", []float64{5, 5, 5}, 0, 0, 0},
		{"falling", []float64{100, 50, 25}, 0.75, 0, 2},
		{"later deeper", []float64{100, 90, 200, 100}, 0.5, 2, 3},
		{"equal declines keep the first", []float64{100, 50, 100, 50}, 0.5, 0, 1},
	} {
		dd, p, tr, err := MaxDrawdown(c.w)
		if err != nil || !near(dd, c.dd, 1e-15) || p != c.peak || tr != c.trought {
			t.Errorf("%s: %v, %d, %d, %v; want %v, %d, %d", c.name, dd, p, tr, err, c.dd, c.peak, c.trought)
		}
	}
}

// Oracle: the O(n²) definition max over i ≤ j of (Wᵢ − Wⱼ)/Wᵢ.
func TestMaxDrawdown_oracle(t *testing.T) {
	w := genPrices(200)
	for n := 1; n <= len(w); n += 11 {
		want := 0.0
		for i := 0; i < n; i++ {
			for j := i; j < n; j++ {
				want = math.Max(want, (w[i]-w[j])/w[i])
			}
		}
		got, _, _, err := MaxDrawdown(w[:n])
		if err != nil || !near(got, want, 1e-15) {
			t.Errorf("n=%d: %v, %v; want %v", n, got, err, want)
		}
	}
}

func TestMaxDrawdown_errors(t *testing.T) {
	_, _, _, err := MaxDrawdown(nil)
	wantErr(t, "empty", err, ErrTooFew)
	_, _, _, err = MaxDrawdown([]float64{100, math.NaN()})
	wantErr(t, "NaN", err, ErrNotFinite)
	_, _, _, err = MaxDrawdown([]float64{100, 0})
	wantErr(t, "zero", err, ErrDomain)
}

func TestMaxDrawdown_noAlloc(t *testing.T) {
	w := []float64{100, 120, 90, 110, 80, 130, 125}
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkInt, sinkInt, sinkErr = MaxDrawdown(w) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkMaxDrawdown(b *testing.B) {
	for _, n := range benchSizes {
		w := genPrices(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, sinkInt, sinkInt, sinkErr = MaxDrawdown(w)
			}
		})
	}
}
