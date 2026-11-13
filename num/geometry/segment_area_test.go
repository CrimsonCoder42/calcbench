package geometry

import (
	"math"
	"math/big"
	"testing"
)

// Worked examples: CRC 31st ed. §4.6.6 p. 334 (area of segment = ½R²(θ − sin θ)), hand-derived:
// r 2, θ π/3 → 2(π/3 − √3/2) ≈ 0.362344294824; r 1, θ π → π/2 (half the disc).
func TestSegmentArea_example(t *testing.T) {
	runValCases(t, []valCase{
		{"SegmentArea(2, π/3)", func() (float64, error) { return SegmentArea(2, math.Pi/3) }, 2 * (math.Pi/3 - math.Sqrt(3)/2), tolSome},
		{"SegmentArea(1, π)", func() (float64, error) { return SegmentArea(1, math.Pi) }, math.Pi / 2, tolFew},
	})
}

func TestSegmentArea_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"full circle is the disc", func() (float64, error) { return SegmentArea(1, 2*math.Pi) }, math.Pi, tolFew},
		// θ − sin θ = θ³/6 − θ⁵/120 + θ⁷/5040 − …; at θ = 10⁻⁶ the third term is below an ulp of the sum.
		{"tiny angle, series branch", func() (float64, error) { return SegmentArea(1, 1e-6) }, 0.5 * (1e-18/6 - 1e-30/120), tolFew},
	})
}

func TestSegmentArea_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { return e1(SegmentArea(nan, 1)) }, ErrNotFinite},
		{"Inf angle", func() error { return e1(SegmentArea(1, inf)) }, ErrNotFinite},
		{"zero radius", func() error { return e1(SegmentArea(0, 1)) }, ErrNonPositive},
		{"zero angle", func() error { return e1(SegmentArea(1, 0)) }, ErrOutOfRange},
		{"beyond a turn", func() error { return e1(SegmentArea(1, 6.5)) }, ErrOutOfRange},
		{"overflow", func() error { return e1(SegmentArea(1e160, 3)) }, ErrOverflow},
	})
}

// thetaMinusSinExact returns θ − sin θ from its Maclaurin series (DLMF 4.19.1) in 256-bit big.Float
// arithmetic, rounded once.
func thetaMinusSinExact(theta float64) float64 {
	const prec = 256
	x := new(big.Float).SetPrec(prec).SetFloat64(theta)
	x2 := new(big.Float).SetPrec(prec).Mul(x, x)
	term := new(big.Float).SetPrec(prec).Mul(x2, x)
	term.Quo(term, big.NewFloat(6))
	sum := new(big.Float).SetPrec(prec).Set(term)
	for k := 4; k < 400; k += 2 {
		term.Mul(term, x2)
		term.Quo(term, big.NewFloat(float64(k*(k+1))))
		term.Neg(term)
		sum.Add(sum, term)
		if term.Sign() == 0 || term.MantExp(nil) < sum.MantExp(nil)-prec {
			break
		}
	}
	f, _ := sum.Float64()
	return f
}

// Oracle: the 256-bit series over θ from 10⁻⁸ to 2π, log-uniform. Just above the switch at θ = 0.5 the
// direct form multiplies sin's one-ulp error by θ/(θ − sin θ) ≈ 24, so 5e-15 (about 22 ulps) is allowed;
// the series branch is far better than that.
func TestSegmentArea_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		theta := math.Exp(math.Log(1e-8) + rng.Float64()*(math.Log(2*math.Pi)-math.Log(1e-8)))
		want := 0.5 * thetaMinusSinExact(theta)
		if got, err := SegmentArea(1, theta); err != nil || !relNear(got, want, 5e-15) {
			t.Fatalf("SegmentArea(1, %v) = (%v, %v), want %v", theta, got, err, want)
		}
	}
	for _, theta := range []float64{math.Nextafter(segmentSeriesBelowRad, 0), segmentSeriesBelowRad} {
		want := 0.5 * thetaMinusSinExact(theta)
		if got, err := SegmentArea(1, theta); err != nil || !relNear(got, want, 5e-15) {
			t.Errorf("at the switch, SegmentArea(1, %v) = (%v, %v), want %v", theta, got, err, want)
		}
	}
}

func TestSegmentArea_allocs(t *testing.T) {
	assertNoAllocs(t, "SegmentArea", func() { sink, _ = SegmentArea(2, 0.3) })
}

func BenchmarkSegmentArea(b *testing.B) {
	b.ReportAllocs()
	r, theta := 2.0, 0.3
	for b.Loop() {
		sink, _ = SegmentArea(r, theta)
	}
}
