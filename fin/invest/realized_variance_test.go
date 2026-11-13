package invest

import (
	"fmt"
	"math"
	"testing"
)

func TestRealizedVariance_handWorked(t *testing.T) {
	r := []float64{0.01, -0.02, 0.02}
	v, err := RealizedVariance(r)
	if err != nil || !near(v, 0.0009, 1e-18) {
		t.Errorf("RealizedVariance = %v, %v; want 0.0009", v, err)
	}
	s, err := RealizedVolatility(r)
	if err != nil || !near(s, 0.03, 1e-16) {
		t.Errorf("RealizedVolatility = %v, %v; want 0.03", s, err)
	}
}

func TestRealizedVariance_cases(t *testing.T) {
	for _, c := range []struct {
		name string
		r    []float64
		want float64
	}{
		{"one", []float64{-0.5}, 0.25},
		{"zeros", []float64{0, 0}, 0},
		{"no mean removed", []float64{0.1, 0.1, 0.1}, 0.03},
	} {
		got, err := RealizedVariance(c.r)
		if err != nil || !near(got, c.want, 1e-16) {
			t.Errorf("%s: %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestRealizedVariance_errors(t *testing.T) {
	_, err := RealizedVariance(nil)
	wantErr(t, "empty", err, ErrTooFew)
	_, err = RealizedVolatility(nil)
	wantErr(t, "vol empty", err, ErrTooFew)
	_, err = RealizedVariance([]float64{math.NaN()})
	wantErr(t, "NaN", err, ErrNotFinite)
	_, err = RealizedVolatility([]float64{math.Inf(-1)})
	wantErr(t, "Inf", err, ErrNotFinite)
	_, err = RealizedVariance([]float64{1e200})
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestRealizedVariance_noAlloc(t *testing.T) {
	r := []float64{0.01, -0.02, 0.02}
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = RealizedVolatility(r) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkRealizedVariance(b *testing.B) {
	for _, n := range benchSizes {
		r := genReturns(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, sinkErr = RealizedVariance(r)
			}
		})
	}
}

func BenchmarkRealizedVolatility(b *testing.B) {
	for _, n := range benchSizes {
		r := genReturns(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, sinkErr = RealizedVolatility(r)
			}
		})
	}
}
