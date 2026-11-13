package signal

import (
	"errors"
	"math"
	"testing"
)

// Hand-worked from the scalar loop (Brown & Hwang 4th ed. §4.2 with φ = 1, H = 1): x̂₀ = 0, P₀ = 1, Q = 0,
// R = 1. Predict: P⁻ = 1. Update(2): K = 1/2, x̂ = 1, P = 1/2. Update(4) (Q = 0, Predict leaves P):
// K = 0.5/1.5 = 1/3, x̂ = 1 + 3/3 = 2, P = 0.5/1.5 = 1/3. Tolerance 1e-15.
func TestKalmanFilter1D_handWorked(t *testing.T) {
	k, err := NewKalmanFilter1D(0, 1, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	if err := k.Predict(); err != nil || k.Variance() != 1 || k.State() != 0 {
		t.Fatalf("predict: (%v, %v, %v)", k.State(), k.Variance(), err)
	}
	if err := k.Update(2); err != nil || k.Gain() != 0.5 || k.State() != 1 || k.Variance() != 0.5 {
		t.Errorf("update 2: K %v x %v P %v (%v), want 0.5, 1, 0.5", k.Gain(), k.State(), k.Variance(), err)
	}
	_ = k.Predict()
	if err := k.Update(4); err != nil || math.Abs(k.Gain()-1.0/3) > 1e-15 || math.Abs(k.State()-2) > 1e-15 || math.Abs(k.Variance()-1.0/3) > 1e-15 {
		t.Errorf("update 4: K %v x %v P %v (%v), want 1/3, 2, 1/3", k.Gain(), k.State(), k.Variance(), err)
	}
}

// With Q = 0, R = 1 and prior (0, P₀ = 1), the filter is the Bayesian mean of a constant: after
// measurements z_1…z_n, x̂ = Σz/(n + 1) and P = 1/(n + 1) (the prior counts as one measurement of 0).
// Oracle on 200 generated measurements; tolerance 1e-12 relative (n roundings).
func TestKalmanFilter1D_constantIsRunningMean(t *testing.T) {
	k, _ := NewKalmanFilter1D(0, 1, 0, 1)
	sum := 0.0
	for i := 1; i <= 200; i++ {
		z := math.Sin(float64(i)) * 10
		sum += z
		_ = k.Predict()
		if err := k.Update(z); err != nil {
			t.Fatal(err)
		}
		want := sum / float64(i+1)
		if math.Abs(k.State()-want) > 1e-12*math.Max(1, math.Abs(want)) || math.Abs(k.Variance()-1/float64(i+1)) > 1e-15 {
			t.Fatalf("n=%d: x %v P %v, want %v, %v", i, k.State(), k.Variance(), want, 1/float64(i+1))
		}
	}
}

// Steady state of the scalar Riccati recursion: P⁻ = (Q + √(Q² + 4QR))/2, K = P⁻/(P⁻ + R). Q = R = 1 gives
// P⁻ = (1 + √5)/2 and K = (√5 − 1)/2 ≈ 0.618 (hand-worked). Converges geometrically; 100 steps, 1e-12.
func TestKalmanFilter1D_steadyStateGain(t *testing.T) {
	k, _ := NewKalmanFilter1D(0, 100, 1, 1)
	for i := 0; i < 100; i++ {
		_ = k.Predict()
		_ = k.Update(0)
	}
	if want := (math.Sqrt(5) - 1) / 2; math.Abs(k.Gain()-want) > 1e-12 {
		t.Errorf("gain = %v, want %v", k.Gain(), want)
	}
}

func TestKalmanFilter1D_reset(t *testing.T) {
	k, _ := NewKalmanFilter1D(0, 1, 0.5, 2)
	_ = k.Predict()
	_ = k.Update(3)
	if err := k.Reset(7, 0.25); err != nil || k.State() != 7 || k.Variance() != 0.25 || k.Gain() != 0 {
		t.Errorf("reset: (%v, %v, %v, %v)", k.State(), k.Variance(), k.Gain(), err)
	}
	if err := k.Reset(1, -1); !errors.Is(err, ErrOutOfRange) || k.State() != 7 {
		t.Errorf("reset negative variance: %v, state %v (must be unchanged)", err, k.State())
	}
	if err := k.Reset(math.NaN(), 1); !errors.Is(err, ErrNotFinite) {
		t.Errorf("reset NaN: %v", err)
	}
	// A zero prior variance is accepted: the first Update then has gain 0.
	_ = k.Reset(5, 0)
	if err := k.Update(9); err != nil || k.State() != 5 || k.Gain() != 0 {
		t.Errorf("zero prior: x %v K %v (%v)", k.State(), k.Gain(), err)
	}
}

func TestKalmanFilter1D_errors(t *testing.T) {
	ctor := []struct {
		name       string
		x, p, q, r float64
		want       error
	}{
		{"NaN state", math.NaN(), 1, 1, 1, ErrNotFinite},
		{"Inf variance", 0, math.Inf(1), 1, 1, ErrNotFinite},
		{"Inf Q", 0, 1, math.Inf(1), 1, ErrNotFinite},
		{"NaN R", 0, 1, 1, math.NaN(), ErrNotFinite},
		{"negative P", 0, -1, 1, 1, ErrOutOfRange},
		{"negative Q", 0, 1, -1, 1, ErrOutOfRange},
		{"zero R", 0, 1, 1, 0, ErrOutOfRange},
		{"negative R", 0, 1, 1, -1, ErrOutOfRange},
	}
	for _, tc := range ctor {
		if k, err := NewKalmanFilter1D(tc.x, tc.p, tc.q, tc.r); !errors.Is(err, tc.want) || k != nil {
			t.Errorf("%s: (%v, %v), want nil, %v", tc.name, k, err, tc.want)
		}
	}
	k, _ := NewKalmanFilter1D(0, 1, 1, 1)
	if err := k.Update(math.Inf(-1)); !errors.Is(err, ErrNotFinite) || k.State() != 0 {
		t.Errorf("update Inf: %v", err)
	}
	// Predict overflow: P = Q = MaxFloat64 → P + Q = +Inf; the filter is left unchanged.
	big, _ := NewKalmanFilter1D(0, math.MaxFloat64, math.MaxFloat64, 1)
	if err := big.Predict(); !errors.Is(err, ErrOverflow) || big.Variance() != math.MaxFloat64 {
		t.Errorf("predict overflow: %v, P %v", err, big.Variance())
	}
	// Update overflow: P⁻ + R = +Inf.
	big2, _ := NewKalmanFilter1D(0, math.MaxFloat64, 0, math.MaxFloat64)
	if err := big2.Update(1); !errors.Is(err, ErrOverflow) {
		t.Errorf("update overflow (variance): %v", err)
	}
	// Update overflow in the state: innovation z − x̂ = MaxFloat64 − (−MaxFloat64).
	big3, _ := NewKalmanFilter1D(-math.MaxFloat64, 1, 0, 1)
	if err := big3.Update(math.MaxFloat64); !errors.Is(err, ErrOverflow) || big3.State() != -math.MaxFloat64 {
		t.Errorf("update overflow (state): %v, x %v", err, big3.State())
	}
}

func TestKalmanFilter1D_noAlloc(t *testing.T) {
	k, _ := NewKalmanFilter1D(0, 1, 0.01, 1)
	if a := testing.AllocsPerRun(100, func() { _ = k.Predict(); _ = k.Update(1) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

var sinkKalman float64

func BenchmarkKalmanFilter1D(b *testing.B) {
	k, _ := NewKalmanFilter1D(0, 1, 0.01, 1)
	b.ReportAllocs()
	z := 0.0
	for b.Loop() {
		_ = k.Predict()
		_ = k.Update(z)
		z += 0.001
	}
	sinkKalman = k.State()
}
