package combinatorics

import (
	"errors"
	"math"
	"math/big"
	"testing"
)

// n = 10: √(20π)·(10/e)¹⁰ = 3 598 695.6187, against 10! = 3 628 800, a ratio of 1.008365 inside Robbins's
// (e^(1/121), e^(1/120)) = (1.008299, 1.008368) (hand-checked). Tolerance 1e-3 absolute: the value is
// quoted to four decimals.
func TestStirlingApproximation_sourceExample(t *testing.T) {
	got, err := StirlingApproximation(10)
	if err != nil || math.Abs(got-3598695.6187) > 1e-3 {
		t.Errorf("StirlingApproximation(10) = (%v, %v), want (3598695.6187, nil)", got, err)
	}
}

// Robbins 1955: for every n ≥ 1, n!/approximation lies strictly in (e^(1/(12n+1)), e^(1/(12n))). Checked
// for n = 1…170 with n! exact and the ratio in 200-bit floats. At n = 170 the true ratio is about 5.6e-10 (relative) below
// the upper bound and 2.4e-7 above the lower; the computed approximation is within about 2e-13 of
// exact, so the check is decisive.
func TestStirlingApproximation_robbinsBound(t *testing.T) {
	for n := 1; n <= stirlingApproximationMax; n++ {
		got, err := StirlingApproximation(n)
		if err != nil || math.IsInf(got, 0) || got <= 0 {
			t.Fatalf("StirlingApproximation(%d) = (%v, %v), want finite positive", n, got, err)
		}
		exact := new(big.Float).SetPrec(200).SetInt(bigFactorial(n))
		ratio, _ := new(big.Float).SetPrec(200).Quo(exact, new(big.Float).SetPrec(200).SetFloat64(got)).Float64()
		lo, hi := math.Exp(1/float64(12*n+1)), math.Exp(1/float64(12*n))
		if ratio <= lo || ratio >= hi {
			t.Errorf("n = %d: n!/approximation = %.15f, want in (%.15f, %.15f)", n, ratio, lo, hi)
		}
	}
}

func TestStirlingApproximation_errors(t *testing.T) {
	for _, n := range []int{0, -1, math.MinInt} {
		if got, err := StirlingApproximation(n); !errors.Is(err, ErrOutOfRange) || got != 0 {
			t.Errorf("StirlingApproximation(%d) = (%v, %v), want (0, ErrOutOfRange)", n, got, err)
		}
	}
	for _, n := range []int{171, math.MaxInt} {
		if got, err := StirlingApproximation(n); !errors.Is(err, ErrOverflow) || got != 0 {
			t.Errorf("StirlingApproximation(%d) = (%v, %v), want (0, ErrOverflow)", n, got, err)
		}
	}
}

func TestStirlingApproximation_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkF, _ = StirlingApproximation(100) }); a != 0 {
		t.Errorf("StirlingApproximation allocates %v times, want 0", a)
	}
}

func BenchmarkStirlingApproximation(b *testing.B) {
	b.Run("n=100", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sinkF, _ = StirlingApproximation(100)
		}
	})
}
