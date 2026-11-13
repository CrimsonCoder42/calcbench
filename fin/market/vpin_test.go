package market

import (
	"math"
	"testing"
)

// ELO 2012 §2: VPIN = Σ|V^S − V^B|/(nV); hand-checked: V = 100, (60, 40), (30, 70) → 60/200 = 0.3.
// Tolerance 1e-15: one division rounding.
func TestVPIN_cases(t *testing.T) {
	for _, c := range []struct {
		name      string
		buy, sell []float64
		v, want   float64
	}{
		{"example", []float64{60, 30}, []float64{40, 70}, 100, 0.3},
		{"balanced", []float64{50, 50}, []float64{50, 50}, 100, 0},
		{"one-sided", []float64{100}, []float64{0}, 100, 1},
	} {
		got, err := VPIN(c.buy, c.sell, c.v)
		if err != nil || !near(got, c.want, 1e-15) {
			t.Errorf("%s: %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestVPIN_errors(t *testing.T) {
	_, err := VPIN(nil, nil, 1)
	wantErr(t, "empty", err, ErrEmpty)
	_, err = VPIN([]float64{1}, []float64{1, 2}, 1)
	wantErr(t, "lengths", err, ErrLengthMismatch)
	_, err = VPIN([]float64{1}, []float64{1}, 0)
	wantErr(t, "V = 0", err, ErrOutOfRange)
	_, err = VPIN([]float64{1}, []float64{1}, math.NaN())
	wantErr(t, "V NaN", err, ErrOutOfRange)
	_, err = VPIN([]float64{-1}, []float64{1}, 1)
	wantErr(t, "negative", err, ErrNegative)
	_, err = VPIN([]float64{1}, []float64{math.Inf(1)}, 1)
	wantErr(t, "Inf", err, ErrNotFinite)
	_, err = VPIN([]float64{math.MaxFloat64, math.MaxFloat64}, []float64{0, 0}, 1)
	wantErr(t, "overflow", err, ErrOverflow)
}

func BenchmarkVPIN(b *testing.B) {
	buy, sell := make([]float64, 50), make([]float64, 50)
	for i := range buy {
		buy[i], sell[i] = float64(40+i%20), float64(60-i%20)
	}
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = VPIN(buy, sell, 100)
	}
}
