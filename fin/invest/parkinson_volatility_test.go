package invest

import (
	"fmt"
	"math"
	"testing"
)

func TestParkinsonVolatility_handWorked(t *testing.T) {
	got, err := ParkinsonVolatility([]float64{110}, []float64{100})
	want := math.Log(1.1) / math.Sqrt(4*math.Ln2)
	if err != nil || !near(got, want, 1e-16) || !near(got, 0.0572395, 1e-6) {
		t.Errorf("ParkinsonVolatility = %v, %v; want %v", got, err, want)
	}
}

func TestParkinsonVolatility_cases(t *testing.T) {
	got, err := ParkinsonVolatility([]float64{100, 100}, []float64{100, 100})
	if err != nil || got != 0 {
		t.Errorf("flat bars: %v, %v; want 0", got, err)
	}
	// Two bars with ln H/L = a and b: √((a² + b²)/(8 ln 2)).
	a, b := math.Log(1.1), math.Log(1.2)
	got, err = ParkinsonVolatility([]float64{110, 120}, []float64{100, 100})
	if want := math.Sqrt((a*a + b*b) / (8 * math.Ln2)); err != nil || !near(got, want, 1e-16) {
		t.Errorf("two bars: %v, %v; want %v", got, err, want)
	}
}

func TestParkinsonVolatility_errors(t *testing.T) {
	_, err := ParkinsonVolatility(nil, nil)
	wantErr(t, "empty", err, ErrTooFew)
	_, err = ParkinsonVolatility([]float64{110, 120}, []float64{100})
	wantErr(t, "lengths", err, ErrLengthMismatch)
	_, err = ParkinsonVolatility([]float64{math.NaN()}, []float64{100})
	wantErr(t, "NaN", err, ErrNotFinite)
	_, err = ParkinsonVolatility([]float64{110}, []float64{0})
	wantErr(t, "zero low", err, ErrDomain)
	_, err = ParkinsonVolatility([]float64{100}, []float64{110})
	wantErr(t, "low > high", err, ErrDomain)
}

func TestParkinsonVolatility_noAlloc(t *testing.T) {
	h, l := []float64{110}, []float64{100}
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = ParkinsonVolatility(h, l) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkParkinsonVolatility(b *testing.B) {
	for _, n := range benchSizes {
		_, h, l, _ := genBars(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, sinkErr = ParkinsonVolatility(h, l)
			}
		})
	}
}
