package market

import (
	"math"
	"testing"

	"github.com/CrimsonCoder42/calcbench/num/stats"
)

// Kyle 1985 Theorem 1: Δp = λ·(order flow); fitted slope hand-checked: x = −2, 0, 2, Δp = −1, 0.5, 2 →
// Sxy = 6, Sxx = 8, λ = 0.75. Tolerance 1e-15.
func TestKyleLambda_handChecked(t *testing.T) {
	got, err := KyleLambda([]float64{-1, 0.5, 2}, []float64{-2, 0, 2})
	if err != nil || !near(got, 0.75, 1e-15) {
		t.Errorf("KyleLambda = %v, %v; want 0.75", got, err)
	}
}

func TestKyleLambda_exactLine(t *testing.T) {
	// Δp = 0.01 + 0.002·v exactly → λ = 0.002
	v := []float64{-300, -100, 50, 200, 400}
	dp := make([]float64, len(v))
	for i, x := range v {
		dp[i] = 0.01 + 0.002*x
	}
	got, err := KyleLambda(dp, v)
	if err != nil || !near(got, 0.002, 1e-15) {
		t.Errorf("exact line: %v, %v", got, err)
	}
}

func TestKyleLambda_errors(t *testing.T) {
	_, err := KyleLambda(nil, nil)
	wantErr(t, "empty", err, ErrEmpty)
	_, err = KyleLambda([]float64{1}, []float64{1, 2})
	wantErr(t, "lengths", err, ErrLengthMismatch)
	_, err = KyleLambda([]float64{math.NaN(), 1}, []float64{1, 2})
	wantErr(t, "NaN dp", err, ErrNotFinite)
	_, err = KyleLambda([]float64{1, 1}, []float64{1, math.Inf(1)})
	wantErr(t, "Inf v", err, ErrNotFinite)
	_, err = KyleLambda([]float64{1}, []float64{1})
	wantErr(t, "one point", err, ErrDomain)
	wantErr(t, "one point (stats)", err, stats.ErrTooFew)
	_, err = KyleLambda([]float64{1, 2}, []float64{3, 3})
	wantErr(t, "constant volume", err, ErrDomain)
	wantErr(t, "constant volume (stats)", err, stats.ErrZeroVariance)
	_, err = KyleLambda([]float64{-math.MaxFloat64, math.MaxFloat64}, []float64{-0.5, 0.5}) // slope 2·MaxFloat64
	wantErr(t, "overflow", err, ErrOverflow)
}

func BenchmarkKyleLambda(b *testing.B) {
	v, dp := make([]float64, 1000), make([]float64, 1000)
	for i := range v {
		v[i] = float64(i%21 - 10)
		dp[i] = 0.001*v[i] + float64(i%3-1)*0.0005
	}
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = KyleLambda(dp, v)
	}
}
