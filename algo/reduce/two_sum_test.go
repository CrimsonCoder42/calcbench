package reduce

import (
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// Knuth TAOCP vol. 2 §4.2.2 Theorem B: s + e = a + b exactly. 1 + 2⁻⁶⁰ rounds to 1 with error 2⁻⁶⁰.
func TestTwoSum_example(t *testing.T) {
	s, e := TwoSum(1, 0x1p-60)
	if s != 1 || e != 0x1p-60 {
		t.Errorf("TwoSum(1, 2^-60) = (%v, %v), want (1, 2^-60)", s, e)
	}
	s, e = TwoSum(0x1p-60, 1) // order does not matter
	if s != 1 || e != 0x1p-60 {
		t.Errorf("TwoSum(2^-60, 1) = (%v, %v), want (1, 2^-60)", s, e)
	}
}

func TestTwoSum_cases(t *testing.T) {
	tests := []struct {
		name       string
		a, b, s, e float64
	}{
		{"zeros", 0, 0, 0, 0},
		{"exact", 1.5, 2.25, 3.75, 0},
		{"cancel", 1e300, -1e300, 0, 0},
		{"max", math.MaxFloat64, -math.MaxFloat64, 0, 0},
	}
	for _, tc := range tests {
		s, e := TwoSum(tc.a, tc.b)
		if s != tc.s || e != tc.e {
			t.Errorf("%s: TwoSum = (%v, %v), want (%v, %v)", tc.name, s, e, tc.s, tc.e)
		}
	}
}

func TestTwoSum_overflow(t *testing.T) {
	s, e := TwoSum(math.MaxFloat64, math.MaxFloat64)
	if !math.IsInf(s, 1) || !math.IsNaN(e) {
		t.Errorf("TwoSum(Max, Max) = (%v, %v), want (+Inf, NaN)", s, e)
	}
}

// Oracle: s + e equals a + b exactly (big.Float) on 500 random pairs over wide exponents.
func TestTwoSum_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 500 {
		a := (r.Float64()*2 - 1) * math.Ldexp(1, r.IntN(200)-100)
		b := (r.Float64()*2 - 1) * math.Ldexp(1, r.IntN(200)-100)
		s, e := TwoSum(a, b)
		lhs := new(big.Float).SetPrec(1000).SetFloat64(s)
		lhs.Add(lhs, new(big.Float).SetPrec(1000).SetFloat64(e))
		rhs := new(big.Float).SetPrec(1000).SetFloat64(a)
		rhs.Add(rhs, new(big.Float).SetPrec(1000).SetFloat64(b))
		if lhs.Cmp(rhs) != 0 {
			t.Fatalf("TwoSum(%v, %v) = (%v, %v): s + e != a + b", a, b, s, e)
		}
	}
}

func BenchmarkTwoSum(b *testing.B) {
	b.ReportAllocs()
	x := 0.1
	for b.Loop() {
		s, e := TwoSum(x, 0x1p-60)
		x = s + e
	}
}
