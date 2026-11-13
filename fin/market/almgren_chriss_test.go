package market

import (
	"math"
	"testing"
)

var sinkSlice []float64

// Almgren & Chriss 2001 §2.2. Hand-checked: σ = 1, η = 1, γ = 0, λ = 0.5, τ = 1 → κ̃² = 0.5,
// cosh κ = 1 + ½·0.5 = 1.25, κ = ln 2; x_1 = sinh(ln 2)/sinh(2 ln 2)·100 = 0.75/1.875·100 = 40.
// Tolerance 1e-12: exp/expm1 roundings times X = 100.
func TestAlmgrenChrissTrajectory_handChecked(t *testing.T) {
	got, err := AlmgrenChrissTrajectory(nil, 100, 2, 2, 1, 1, 0, 0.5)
	want := []float64{100, 40, 0}
	if err != nil || len(got) != 3 {
		t.Fatalf("got %v, %v", got, err)
	}
	for i := range want {
		if !near(got[i], want[i], 1e-12) {
			t.Errorf("x[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestAlmgrenChrissTrajectory_riskNeutralIsLinear(t *testing.T) {
	got, err := AlmgrenChrissTrajectory(nil, 100, 1, 4, 0.3, 0.1, 0.01, 0)
	want := []float64{100, 75, 50, 25, 0}
	if err != nil || len(got) != len(want) {
		t.Fatalf("got %v, %v", got, err)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("x[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestAlmgrenChrissTrajectory_largeKappaNoOverflow(t *testing.T) {
	// κτ ≈ ln(1e12) ≈ 27.6, so κT ≈ 2.8e4 ≫ 710: sinh would overflow; the ratio must stay finite and non-increasing
	got, err := AlmgrenChrissTrajectory(nil, 1, 1, 1000, 1, 1e-6, 0, 1e12)
	if err != nil {
		t.Fatal(err)
	}
	if got[0] != 1 || got[1000] != 0 {
		t.Errorf("ends = %v, %v", got[0], got[1000])
	}
	for j := 1; j < len(got); j++ {
		if !(got[j] <= got[j-1]) || math.IsNaN(got[j]) {
			t.Fatalf("x[%d] = %v after %v", j, got[j], got[j-1])
		}
	}
}

func TestAlmgrenChrissTrajectory_reusesDst(t *testing.T) {
	dst := make([]float64, 0, 10)
	got, err := AlmgrenChrissTrajectory(dst, 100, 1, 4, 0, 1, 0, 0)
	if err != nil || &got[0] != &dst[:1][0] {
		t.Errorf("dst not reused: %v, %v", got, err)
	}
}

func TestAlmgrenChrissTrajectory_errors(t *testing.T) {
	_, err := AlmgrenChrissTrajectory(nil, 100, 1, 0, 1, 1, 0, 1)
	wantErr(t, "periods 0", err, ErrOutOfRange)
	_, err = AlmgrenChrissTrajectory(nil, math.NaN(), 1, 1, 1, 1, 0, 1)
	wantErr(t, "NaN shares", err, ErrNotFinite)
	_, err = AlmgrenChrissTrajectory(nil, 100, 0, 1, 1, 1, 0, 1)
	wantErr(t, "horizon 0", err, ErrOutOfRange)
	_, err = AlmgrenChrissTrajectory(nil, 100, 1, 1, 1, 1, 0, -1)
	wantErr(t, "lambda < 0", err, ErrOutOfRange)
	_, err = AlmgrenChrissTrajectory(nil, 100, 1, 1, -1, 1, 0, 1)
	wantErr(t, "sigma < 0", err, ErrOutOfRange)
	_, err = AlmgrenChrissTrajectory(nil, 100, 1, 1, 1, 0, 0, 1)
	wantErr(t, "eta 0", err, ErrOutOfRange)
	_, err = AlmgrenChrissTrajectory(nil, 100, 1, 1, 1, 1, math.Inf(1), 1)
	wantErr(t, "gamma Inf", err, ErrNotFinite)
	_, err = AlmgrenChrissTrajectory(nil, 100, 2, 1, 1, 1, 1, 1) // η̃ = 1 − 1·2/2 = 0
	wantErr(t, "eta tilde 0", err, ErrDomain)
	_, err = AlmgrenChrissTrajectory(nil, 100, 1, 1, math.MaxFloat64, 1, 0, math.MaxFloat64)
	wantErr(t, "kappa overflow", err, ErrOverflow)
}

func BenchmarkAlmgrenChrissTrajectory(b *testing.B) {
	dst := make([]float64, 1001)
	b.ReportAllocs()
	for b.Loop() {
		sinkSlice, sinkErr = AlmgrenChrissTrajectory(dst, 1e6, 1, 1000, 0.95, 2.5e-6, 2.5e-7, 1e-6)
	}
}

// Almgren & Chriss 2001 §1, linear impact. Hand-checked: x = 100, 40, 0, τ = 1, σ = 1, η = 1, γ = 0.1,
// ε = 0.5: η̃ = 0.95, E = ½·0.1·100² + 0.5·100 + 0.95·(60² + 40²) = 500 + 50 + 4940 = 5490;
// V = 1·1·(40² + 0²) = 1600. Tolerance 1e-9.
func TestAlmgrenChrissExpectedCost_handChecked(t *testing.T) {
	e, v, err := AlmgrenChrissExpectedCost([]float64{100, 40, 0}, 1, 1, 1, 0.1, 0.5)
	if err != nil || !near(e, 5490, 1e-9) || !near(v, 1600, 1e-9) {
		t.Errorf("got %v, %v, %v; want 5490, 1600", e, v, err)
	}
}

func TestAlmgrenChrissExpectedCost_optimalBeatsLinearOnUtility(t *testing.T) {
	// with λ = 0.5 the optimal schedule minimises E + λV among schedules, so it beats the straight line
	const lambda = 0.5
	opt, _ := AlmgrenChrissTrajectory(nil, 100, 2, 2, 1, 1, 0, lambda)
	eo, vo, err1 := AlmgrenChrissExpectedCost(opt, 1, 1, 1, 0, 0)
	el, vl, err2 := AlmgrenChrissExpectedCost([]float64{100, 50, 0}, 1, 1, 1, 0, 0)
	if err1 != nil || err2 != nil || !(eo+lambda*vo < el+lambda*vl) {
		t.Errorf("optimal %v + λ%v not below linear %v + λ%v (%v, %v)", eo, vo, el, vl, err1, err2)
	}
}

func TestAlmgrenChrissExpectedCost_errors(t *testing.T) {
	_, _, err := AlmgrenChrissExpectedCost([]float64{0}, 1, 1, 1, 0, 0)
	wantErr(t, "one holding", err, ErrTooFew)
	_, _, err = AlmgrenChrissExpectedCost([]float64{1, math.NaN()}, 1, 1, 1, 0, 0)
	wantErr(t, "NaN", err, ErrNotFinite)
	_, _, err = AlmgrenChrissExpectedCost([]float64{1, 0}, 1, 1, 1, 0, -1)
	wantErr(t, "epsilon < 0", err, ErrOutOfRange)
	_, _, err = AlmgrenChrissExpectedCost([]float64{1, 0}, 0, 1, 1, 0, 0)
	wantErr(t, "tau 0", err, ErrOutOfRange)
	_, _, err = AlmgrenChrissExpectedCost([]float64{1, 0}, 2, 1, 1, 1, 0)
	wantErr(t, "eta tilde 0", err, ErrDomain)
	_, _, err = AlmgrenChrissExpectedCost([]float64{1, 1}, 1, 1, 1, 0, 0)
	wantErr(t, "not liquidated", err, ErrDomain)
	_, _, err = AlmgrenChrissExpectedCost([]float64{1e200, 0}, 1, 1, 1, 0, 0)
	wantErr(t, "overflow", err, ErrOverflow)
}

func BenchmarkAlmgrenChrissExpectedCost(b *testing.B) {
	x, _ := AlmgrenChrissTrajectory(nil, 1e6, 1, 1000, 0.95, 2.5e-6, 2.5e-7, 1e-6)
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkFloat2, sinkErr = AlmgrenChrissExpectedCost(x, 0.001, 0.95, 2.5e-6, 2.5e-7, 0.0625)
	}
}
