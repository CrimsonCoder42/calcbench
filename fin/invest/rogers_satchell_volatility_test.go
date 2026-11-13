package invest

import (
	"fmt"
	"math"
	"testing"
)

func TestRogersSatchellVolatility_handWorked(t *testing.T) {
	got, err := RogersSatchellVolatility([]float64{100}, []float64{110}, []float64{95}, []float64{105})
	want := math.Sqrt(math.Log(110.0/105)*math.Log(1.1) + math.Log(95.0/105)*math.Log(0.95))
	if err != nil || !near(got, want, 1e-16) || !near(got, 0.097813, 1e-5) {
		t.Errorf("RogersSatchellVolatility = %v, %v; want %v", got, err, want)
	}
}

// A pure trend bar (open = low, close = high) has zero Rogers–Satchell variance: the drift is removed.
func TestRogersSatchellVolatility_trendBar(t *testing.T) {
	got, err := RogersSatchellVolatility([]float64{100, 110}, []float64{110, 121}, []float64{100, 110}, []float64{110, 121})
	if err != nil || got != 0 {
		t.Errorf("trend bars: %v, %v; want 0", got, err)
	}
}

func TestRogersSatchellVolatility_generated(t *testing.T) {
	o, h, l, c := genBars(500)
	got, err := RogersSatchellVolatility(o, h, l, c)
	if err != nil || got <= 0 || got > 0.1 {
		t.Errorf("generated: %v, %v; want in (0, 0.1]", got, err)
	}
}

func TestRogersSatchellVolatility_errors(t *testing.T) {
	_, err := RogersSatchellVolatility(nil, nil, nil, nil)
	wantErr(t, "empty", err, ErrTooFew)
	_, err = RogersSatchellVolatility([]float64{100, 100}, []float64{110}, []float64{100}, []float64{105})
	wantErr(t, "lengths", err, ErrLengthMismatch)
	_, err = RogersSatchellVolatility([]float64{100}, []float64{math.Inf(1)}, []float64{100}, []float64{105})
	wantErr(t, "Inf", err, ErrNotFinite)
	_, err = RogersSatchellVolatility([]float64{-1}, []float64{110}, []float64{-2}, []float64{105})
	wantErr(t, "negative", err, ErrDomain)
}

func TestRogersSatchellVolatility_noAlloc(t *testing.T) {
	o, h, l, c := []float64{100}, []float64{110}, []float64{95}, []float64{105}
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = RogersSatchellVolatility(o, h, l, c) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkRogersSatchellVolatility(b *testing.B) {
	for _, n := range benchSizes {
		o, h, l, c := genBars(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, sinkErr = RogersSatchellVolatility(o, h, l, c)
			}
		})
	}
}
