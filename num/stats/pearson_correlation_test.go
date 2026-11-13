package stats

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// Worked example: NIST StRD Norris, certified R² = 0.999993745883712 with a positive slope, so
// r = +√R² = 0.9999968729369667 (hand-checked from the certified value).
func TestPearsonCorrelation_norris(t *testing.T) {
	// Tolerance 1e-14: the certified R² has 15 significant digits.
	if got, err := PearsonCorrelation(norrisX, norrisY); err != nil || !near(got, math.Sqrt(0.999993745883712), 1e-14) {
		t.Errorf("PearsonCorrelation(Norris) = %v, %v; want 0.9999968729369667", got, err)
	}
}

func TestPearsonCorrelation_cases(t *testing.T) {
	m := math.MaxFloat64
	for _, tt := range []struct {
		name string
		x, y []float64
		want float64
	}{
		{"two points", []float64{0, 1}, []float64{5, 7}, 1},
		{"perfect negative", []float64{1, 2, 3}, []float64{3, 2, 1}, -1},
		{"uncorrelated", []float64{-1, 0, 1, -1, 0, 1}, []float64{0, 1, 2, 2, 1, 0}, 0},
		{"hand-checked", []float64{1, 2, 3, 4}, []float64{2, 1, 4, 3}, 0.6},
		{"huge y, small x", []float64{0, 1e-100, 2e-100}, []float64{-m / 2, 0, m / 2}, 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Tolerance 1e-15: a few roundings (r is clamped to [−1, 1]).
			if got, err := PearsonCorrelation(tt.x, tt.y); err != nil || !near(got, tt.want, 1e-15) {
				t.Errorf("PearsonCorrelation(%v, %v) = %v, %v; want %v", tt.x, tt.y, got, err, tt.want)
			}
		})
	}
}

func TestPearsonCorrelation_errors(t *testing.T) {
	for _, tt := range []struct {
		name string
		x, y []float64
		want error
	}{
		{"empty", []float64{}, []float64{}, ErrEmpty},
		{"length mismatch", []float64{1}, []float64{1, 2}, ErrLengthMismatch},
		{"one pair", []float64{1}, []float64{2}, ErrTooFew},
		{"NaN", []float64{1, 2}, []float64{math.NaN(), 1}, ErrNotFinite},
		{"constant x", []float64{2, 2}, []float64{1, 2}, ErrZeroVariance},
		{"constant y", []float64{1, 2}, []float64{4, 4}, ErrZeroVariance},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := PearsonCorrelation(tt.x, tt.y); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("PearsonCorrelation = %v, %v; want 0, %v", got, err, tt.want)
			}
		})
	}
}

func TestPearsonCorrelation_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	bf := func(v float64) *big.Float { return new(big.Float).SetPrec(512).SetFloat64(v) }
	for trial := range 200 {
		n := 2 + rng.IntN(200)
		x, y := randomFloats(rng, n), randomFloats(rng, n)
		for i := range y {
			y[i] += rng.Float64() * x[i]
		}
		mx, my := bigMean(x), bigMean(y)
		sxx, syy, sxy := bf(0), bf(0), bf(0)
		for i := range x {
			dx := bf(x[i])
			dx.Sub(dx, mx)
			dy := bf(y[i])
			dy.Sub(dy, my)
			sxx.Add(sxx, new(big.Float).SetPrec(512).Mul(dx, dx))
			syy.Add(syy, new(big.Float).SetPrec(512).Mul(dy, dy))
			sxy.Add(sxy, new(big.Float).SetPrec(512).Mul(dx, dy))
		}
		den := new(big.Float).SetPrec(512).Mul(sxx, syy)
		den.Sqrt(den)
		want, _ := new(big.Float).SetPrec(512).Quo(sxy, den).Float64()
		// Tolerance 1e-12 absolute: r is in [−1, 1]; O(n·u) in the centred sums.
		if got, err := PearsonCorrelation(x, y); err != nil || math.Abs(got-want) > 1e-12 {
			t.Fatalf("trial %d: PearsonCorrelation(n=%d) = %v, %v; oracle %v", trial, n, got, err, want)
		}
	}
}

func TestPearsonCorrelation_noAllocation(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = PearsonCorrelation(norrisX, norrisY) }); a != 0 {
		t.Errorf("PearsonCorrelation allocates %v times per call, want 0", a)
	}
}

func BenchmarkPearsonCorrelation(b *testing.B) {
	for _, size := range benchSizes {
		xs, ys := benchInputs(size.n), benchInputs(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				benchSink, _ = PearsonCorrelation(xs[i&7], ys[(i+1)&7])
			}
		})
	}
}

// Tiny spread in x. x = 2^-700·(1, 2, 3) or 1e-200·(1, 2, 3) with y = 2, 1, 3: centred x ∝ (−1, 0, 1),
// centred y = (0, −1, 1), so S_xy ∝ 1, S_xx ∝ 2, S_yy = 2 and r = 1/√(2·2) = 1/2 (hand); the decimal case
// is also checked against exact rational arithmetic on the float64 inputs. Before the underflow fix both
// returned ErrZeroVariance claiming every x was equal. Two smallest subnormals: two points, r = 1.
func TestPearsonCorrelation_tinySpread(t *testing.T) {
	y := []float64{2, 1, 3}
	// Tolerance 4e-16 relative: r = sxy/(√sxx·√syy) rounds the two roots and the product and quotient.
	if got, err := PearsonCorrelation([]float64{0x1p-700, 0x2p-700, 0x3p-700}, y); err != nil || !near(got, 0.5, 4e-16) {
		t.Errorf("PearsonCorrelation(2^-700·(1 2 3), 2 1 3) = %v, %v; want 0.5", got, err)
	}
	x := []float64{1e-200, 2e-200, 3e-200}
	_, _, sxx, syy, sxy := exactCentred(x, y)
	f := func(r *big.Rat) *big.Float { return new(big.Float).SetPrec(512).SetRat(r) }
	den := new(big.Float).SetPrec(512).Mul(f(sxx), f(syy))
	wantF, _ := new(big.Float).SetPrec(512).Quo(f(sxy), den.Sqrt(den)).Float64()
	// Tolerance 4e-16 relative for the same roundings; and the exact r is 1/2 to within the inputs' rounding.
	if got, err := PearsonCorrelation(x, y); err != nil || !near(got, wantF, 4e-16) || !near(got, 0.5, 1e-15) {
		t.Errorf("PearsonCorrelation(1e-200·(1 2 3), 2 1 3) = %v, %v; want %v", got, err, wantF)
	}
	// Two points: r = 1 exactly in theory; 4e-16 covers the rounded roots (√2-type factors).
	if got, err := PearsonCorrelation([]float64{0x1p-1074, 0x1p-1073}, []float64{1, 2}); err != nil || !near(got, 1, 4e-16) {
		t.Errorf("PearsonCorrelation(2^-1074 2^-1073, 1 2) = %v, %v; want 1", got, err)
	}
	if _, err := PearsonCorrelation([]float64{1e-300, 1e-300, 1e-300}, y); !errors.Is(err, ErrZeroVariance) {
		t.Errorf("PearsonCorrelation(equal tiny x) error = %v; want ErrZeroVariance", err)
	}
}
