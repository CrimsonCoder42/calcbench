package invest

import (
	"fmt"
	"math"
	"testing"
)

func TestTrackingError_handWorked(t *testing.T) {
	got, err := TrackingError([]float64{0.02, 0.03, 0.01}, []float64{0.01, 0.01, 0.01})
	if err != nil || !near(got, 0.01, 1e-16) {
		t.Errorf("TrackingError = %v, %v; want 0.01", got, err)
	}
}

func TestTrackingError_cases(t *testing.T) {
	got, err := TrackingError([]float64{0.5, 0.75}, []float64{0.25, 0.5})
	if err != nil || got != 0 {
		t.Errorf("constant active return: %v, %v; want 0", got, err)
	}
	// Active returns 1, 2, 3, 4: s = √(5/3).
	got, err = TrackingError([]float64{1, 2, 3, 4}, []float64{0, 0, 0, 0})
	if err != nil || !near(got, math.Sqrt(5.0/3), 1e-15) {
		t.Errorf("1..4: %v, %v; want √(5/3)", got, err)
	}
}

func TestTrackingError_errors(t *testing.T) {
	_, err := TrackingError([]float64{0.01}, []float64{0.01})
	wantErr(t, "one", err, ErrTooFew)
	_, err = TrackingError([]float64{0.01, 0.02}, []float64{0.01})
	wantErr(t, "lengths", err, ErrLengthMismatch)
	_, err = TrackingError([]float64{0.01, 0.02}, []float64{0.01, math.Inf(1)})
	wantErr(t, "Inf", err, ErrNotFinite)
	_, err = TrackingError([]float64{1.7e308, 0}, []float64{-1.7e308, 0})
	wantErr(t, "active overflow", err, ErrOverflow)
	_, err = TrackingError([]float64{1e300, -1e300}, []float64{0, 0})
	wantErr(t, "variance overflow", err, ErrOverflow)
}

func TestTrackingError_noAlloc(t *testing.T) {
	p, bm := []float64{0.02, 0.03, 0.01}, []float64{0.01, 0.01, 0.01}
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = TrackingError(p, bm) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkTrackingError(b *testing.B) {
	for _, n := range benchSizes {
		p := genReturns(n)
		bm := genReturns(n + 1)[1:]
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, sinkErr = TrackingError(p, bm)
			}
		})
	}
}
