package stats

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// Worked example: Baldi et al. 2000 eq. (11) / Chicco and Jurman 2020 eq. 2,
// MCC = (TP·TN − FP·FN)/√((TP+FP)(TP+FN)(TN+FP)(TN+FN)). TP = 90, TN = 80, FP = 20, FN = 10 →
// 7000/√99000000 = 0.7035264706814485 (hand-checked).
func TestMatthewsCorrelation_worked(t *testing.T) {
	// Tolerance 1e-15: a few roundings.
	if got, err := MatthewsCorrelation(90, 80, 20, 10); err != nil || !near(got, 7000/math.Sqrt(99000000), 1e-15) {
		t.Errorf("MatthewsCorrelation(90, 80, 20, 10) = %v, %v; want 0.70352647", got, err)
	}
}

func TestMatthewsCorrelation_cases(t *testing.T) {
	for _, tt := range []struct {
		name           string
		tp, tn, fp, fn int
		want           float64
	}{
		{"perfect", 10, 5, 0, 0, 1},
		{"inverse", 0, 0, 4, 6, -1},
		{"no association", 5, 5, 5, 5, 0},
		{"huge counts, exact numerator", math.MaxInt, math.MaxInt, math.MaxInt - 1, math.MaxInt, 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MatthewsCorrelation(tt.tp, tt.tn, tt.fp, tt.fn)
			// Tolerance 1e-15 absolute; the huge case's numerator is MaxInt, against a denominator near
			// 2^126, so r ≈ 5e-20.
			if err != nil || math.Abs(got-tt.want) > 1e-15 {
				t.Errorf("MatthewsCorrelation(%d, %d, %d, %d) = %v, %v; want %v", tt.tp, tt.tn, tt.fp, tt.fn, got, err, tt.want)
			}
		})
	}
}

func TestMatthewsCorrelation_errors(t *testing.T) {
	for _, tt := range []struct {
		name           string
		tp, tn, fp, fn int
		want           error
	}{
		{"negative TP", -1, 1, 1, 1, ErrNegative},
		{"negative FN", 1, 1, 1, -1, ErrNegative},
		{"no predicted positives", 0, 5, 0, 3, ErrZeroTotal},
		{"no actual negatives", 4, 0, 0, 2, ErrZeroTotal},
		{"all zero", 0, 0, 0, 0, ErrZeroTotal},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := MatthewsCorrelation(tt.tp, tt.tn, tt.fp, tt.fn); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("MatthewsCorrelation = %v, %v; want 0, %v", got, err, tt.want)
			}
		})
	}
}

// The oracle evaluates the formula in 256-bit floating point.
func TestMatthewsCorrelation_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	bi := func(v int) *big.Float { return new(big.Float).SetPrec(256).SetInt64(int64(v)) }
	for trial := range 400 {
		limit := int64(1) << (1 + rng.IntN(62))
		tp, tn := int(rng.Int64N(limit))+1, int(rng.Int64N(limit))+1
		fp, fn := int(rng.Int64N(limit))+1, int(rng.Int64N(limit))+1
		num := new(big.Float).SetPrec(256).Mul(bi(tp), bi(tn))
		num.Sub(num, new(big.Float).SetPrec(256).Mul(bi(fp), bi(fn)))
		den := new(big.Float).SetPrec(256).SetInt64(1)
		for _, s := range [][2]int{{tp, fp}, {tp, fn}, {tn, fp}, {tn, fn}} {
			m := new(big.Float).SetPrec(256).Add(bi(s[0]), bi(s[1]))
			den.Mul(den, m)
		}
		den.Sqrt(den)
		want, _ := new(big.Float).SetPrec(256).Quo(num, den).Float64()
		// Tolerance 1e-15 absolute: r is in [−1, 1]; the margins round once above 2^53.
		if got, err := MatthewsCorrelation(tp, tn, fp, fn); err != nil || math.Abs(got-want) > 1e-15 {
			t.Fatalf("trial %d: MatthewsCorrelation(%d, %d, %d, %d) = %v, %v; oracle %v", trial, tp, tn, fp, fn, got, err, want)
		}
	}
}

func TestMatthewsCorrelation_noAllocation(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = MatthewsCorrelation(90, 80, 20, 10) }); a != 0 {
		t.Errorf("MatthewsCorrelation allocates %v times per call, want 0", a)
	}
}

func BenchmarkMatthewsCorrelation(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchSink, _ = MatthewsCorrelation(90, 80, 20, 10)
	}
}
