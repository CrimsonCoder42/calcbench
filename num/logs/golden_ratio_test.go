package logs

import (
	"math"
	"math/big"
	"testing"
)

// The digits are OEIS A001622 (read 2026-11-12). Two independent checks: the defining identity
// φ² = φ + 1 to 200 bits of big.Float precision (the constant has 63 digits, about 209 bits, so the
// residual must be below 2^-190), and (1+√5)/2 computed in big.Float.
func TestGoldenRatio_identity(t *testing.T) {
	const prec = 256
	phi := new(big.Float).SetPrec(prec)
	if _, ok := phi.SetString("1.61803398874989484820458683436563811772030917980576286213544862"); !ok {
		t.Fatal("parse")
	}
	sq := new(big.Float).SetPrec(prec).Mul(phi, phi)
	rhs := new(big.Float).SetPrec(prec).Add(phi, big.NewFloat(1).SetPrec(prec))
	diff := new(big.Float).SetPrec(prec).Sub(sq, rhs)
	diff.Abs(diff)
	limit := new(big.Float).SetPrec(prec).SetMantExp(big.NewFloat(1), -190)
	if diff.Cmp(limit) > 0 {
		t.Fatalf("φ² − (φ+1) = %g, want < 2^-190", diff)
	}
}

func TestGoldenRatio_closedForm(t *testing.T) {
	const prec = 256
	five := new(big.Float).SetPrec(prec).SetInt64(5)
	root := new(big.Float).SetPrec(prec).Sqrt(five)
	want := new(big.Float).SetPrec(prec).Add(root, big.NewFloat(1).SetPrec(prec))
	want.Quo(want, big.NewFloat(2).SetPrec(prec))
	got, _ := want.Float64()
	if got != GoldenRatio {
		t.Fatalf("float64((1+√5)/2) = %v, GoldenRatio = %v", got, float64(GoldenRatio))
	}
}

// The float64 value must be the same double as the standard library's math.Phi (exact comparison:
// both are the round-to-nearest double of the same real number).
func TestGoldenRatio_matchesMathPhi(t *testing.T) {
	if float64(GoldenRatio) != math.Phi {
		t.Fatalf("GoldenRatio = %v, math.Phi = %v", float64(GoldenRatio), math.Phi)
	}
}

var sinkFloat float64

func BenchmarkGoldenRatio(b *testing.B) {
	b.ReportAllocs()
	x := 1.0
	for b.Loop() {
		sinkFloat = x * GoldenRatio
	}
}
