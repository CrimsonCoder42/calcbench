package stats

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// Worked example: NIST StRD Linear Least Squares, Norris (36 points): certified B0 = −0.262323073774029,
// B1 = 1.00211681802045.
func TestLinearRegression_norris(t *testing.T) {
	intercept, slope, err := LinearRegression(norrisX, norrisY)
	// Tolerance: the certified values have 15 significant digits; the intercept is a small difference
	// of quantities near 400, so it is checked to 1e-11 relative, the slope to 1e-14.
	if err != nil || !near(intercept, -0.262323073774029, 1e-11) || !near(slope, 1.00211681802045, 1e-14) {
		t.Errorf("LinearRegression(Norris) = %v, %v, %v; want -0.262323073774029, 1.00211681802045", intercept, slope, err)
	}
}

func TestLinearRegression_cases(t *testing.T) {
	for _, tt := range []struct {
		name             string
		x, y             []float64
		intercept, slope float64
	}{
		{"two points", []float64{0, 1}, []float64{1, 3}, 1, 2},
		{"exact line", []float64{1, 2, 3, 4}, []float64{-1, -4, -7, -10}, 2, -3},
		{"horizontal", []float64{1, 2, 3}, []float64{5, 5, 5}, 5, 0},
		{"large offset in x", []float64{1e9, 1e9 + 1, 1e9 + 2}, []float64{0, 1, 2}, -1e9, 1},
		{"y squares overflow, slope finite", []float64{0, 1, 2}, []float64{-math.MaxFloat64 / 2, 0, math.MaxFloat64 / 2}, -math.MaxFloat64 / 2, math.MaxFloat64 / 2},
		{"symmetric scatter", []float64{-1, 0, 1, -1, 0, 1}, []float64{0, 1, 2, 2, 1, 0}, 1, 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			intercept, slope, err := LinearRegression(tt.x, tt.y)
			// Tolerance 1e-15 relative (absolute at 0): exact data, a few roundings.
			if err != nil || !near(intercept, tt.intercept, 1e-15) || !near(slope, tt.slope, 1e-15) {
				t.Errorf("LinearRegression(%v, %v) = %v, %v, %v; want %v, %v", tt.x, tt.y, intercept, slope, err, tt.intercept, tt.slope)
			}
		})
	}
}

func TestLinearRegression_errors(t *testing.T) {
	m := math.MaxFloat64
	for _, tt := range []struct {
		name string
		x, y []float64
		want error
	}{
		{"empty", nil, nil, ErrEmpty},
		{"length mismatch", []float64{1, 2}, []float64{1}, ErrLengthMismatch},
		{"one point", []float64{1}, []float64{1}, ErrTooFew},
		{"NaN x", []float64{1, math.NaN()}, []float64{1, 2}, ErrNotFinite},
		{"Inf y", []float64{1, 2}, []float64{1, math.Inf(1)}, ErrNotFinite},
		{"all x equal", []float64{3, 3, 3}, []float64{1, 2, 3}, ErrZeroVariance},
		{"slope overflows", []float64{0, 1e-10}, []float64{-m / 2, m / 2}, ErrOverflow},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if a, b, err := LinearRegression(tt.x, tt.y); !errors.Is(err, tt.want) || a != 0 || b != 0 {
				t.Errorf("LinearRegression = %v, %v, %v; want 0, 0, %v", a, b, err, tt.want)
			}
		})
	}
}

// The oracle computes the textbook two-pass sums in 512-bit floating point.
func TestLinearRegression_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	bf := func(v float64) *big.Float { return new(big.Float).SetPrec(512).SetFloat64(v) }
	for trial := range 200 {
		n := 2 + rng.IntN(200)
		x, y := randomFloats(rng, n), randomFloats(rng, n)
		for i := range y {
			y[i] += 0.5 * x[i]
		}
		mx, my := bigMean(x), bigMean(y)
		sxx, sxy := bf(0), bf(0)
		for i := range x {
			dx := bf(x[i])
			dx.Sub(dx, mx)
			dy := bf(y[i])
			dy.Sub(dy, my)
			sxx.Add(sxx, new(big.Float).SetPrec(512).Mul(dx, dx))
			sxy.Add(sxy, new(big.Float).SetPrec(512).Mul(dx, dy))
		}
		bSlope := new(big.Float).SetPrec(512).Quo(sxy, sxx)
		bIntercept := new(big.Float).SetPrec(512).Mul(bSlope, mx)
		bIntercept.Sub(my, bIntercept)
		wantSlope, _ := bSlope.Float64()
		wantIntercept, _ := bIntercept.Float64()
		intercept, slope, err := LinearRegression(x, y)
		// Tolerance 1e-11 relative, plus 1e-11 absolute: the slope's condition grows when x has small
		// spread; the intercept subtracts slope·x̄ ≈ 50 from ȳ ≈ 150.
		if err != nil || math.Abs(slope-wantSlope) > 1e-11*(1+math.Abs(wantSlope)) ||
			math.Abs(intercept-wantIntercept) > 1e-11*(1+math.Abs(wantIntercept)) {
			t.Fatalf("trial %d: LinearRegression(n=%d) = %v, %v, %v; oracle %v, %v", trial, n, intercept, slope, err, wantIntercept, wantSlope)
		}
	}
}

func TestLinearRegression_noAllocation(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { benchSink, _, _ = LinearRegression(norrisX, norrisY) }); a != 0 {
		t.Errorf("LinearRegression allocates %v times per call, want 0", a)
	}
}

func BenchmarkLinearRegression(b *testing.B) {
	for _, size := range benchSizes {
		xs, ys := benchInputs(size.n), benchInputs(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				_, benchSink, _ = LinearRegression(xs[i&7], ys[(i+1)&7])
			}
		})
	}
}

// Tiny spread in x. x = 2^-700·(1, 2, 3), y = 2, 1, 3: S_xy = 2^-700, S_xx = 2·2^-1400, so slope = 2^699
// and intercept = ȳ − slope·x̄ = 2 − 2^699·2^-699 = 1, both exact (hand). x = 1e-200·(1, 2, 3): slope
// 1e-200/(2e-400) = 5e199 and intercept 1 (hand), checked also against exact rational arithmetic on the
// float64 inputs. Before the underflow fix both returned ErrZeroVariance.
func TestLinearRegression_tinySpread(t *testing.T) {
	y := []float64{2, 1, 3}
	if b0, b1, err := LinearRegression([]float64{0x1p-700, 0x2p-700, 0x3p-700}, y); err != nil || b0 != 1 || b1 != 0x1p699 {
		t.Errorf("LinearRegression(2^-700·(1 2 3), 2 1 3) = %v, %v, %v; want 1, 2^699", b0, b1, err)
	}
	x := []float64{1e-200, 2e-200, 3e-200}
	xMean, yMean, sxx, _, sxy := exactCentred(x, y)
	slope := new(big.Rat).Quo(sxy, sxx)
	wantB1 := ratFloat(slope)
	wantB0 := ratFloat(new(big.Rat).Sub(yMean, new(big.Rat).Mul(slope, xMean)))
	// Tolerance 4e-16 relative on the slope (a quotient of two rounded scaled sums) and 1e-15 absolute
	// on the intercept (ȳ − slope·x̄ with ȳ = 2 and slope·x̄ ≈ 1).
	b0, b1, err := LinearRegression(x, y)
	if err != nil || !near(b1, wantB1, 4e-16) || math.Abs(b0-wantB0) > 1e-15 || !near(b1, 5e199, 1e-15) || math.Abs(b0-1) > 1e-15 {
		t.Errorf("LinearRegression(1e-200·(1 2 3), 2 1 3) = %v, %v, %v; want %v, %v", b0, b1, err, wantB0, wantB1)
	}
}
