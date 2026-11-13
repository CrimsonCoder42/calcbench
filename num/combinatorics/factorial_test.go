package combinatorics

import (
	"errors"
	"math"
	"math/big"
	"testing"
)

// Benchmark sinks keep results observable. Shared by the package's benchmarks.
var (
	sinkU   uint64
	sinkF   float64
	sinkBig *big.Int
)

// maxU64 is 2⁶⁴ − 1 as a big.Int, the overflow threshold of every uint64 function.
var maxU64 = new(big.Int).SetUint64(math.MaxUint64)

// wantU64 checks a uint64 result against an exact big.Int: the value when it fits, else ErrOverflow.
func wantU64(t *testing.T, name string, got uint64, err error, exact *big.Int) {
	t.Helper()
	if exact.Cmp(maxU64) > 0 {
		if !errors.Is(err, ErrOverflow) || got != 0 {
			t.Errorf("%s = (%d, %v), want (0, ErrOverflow); exact %s", name, got, err, exact)
		}
		return
	}
	if err != nil || got != exact.Uint64() {
		t.Errorf("%s = (%d, %v), want (%s, nil)", name, got, err, exact)
	}
}

// bigLog returns the natural log of a positive big.Int to float64 accuracy.
func bigLog(x *big.Int) float64 {
	var mant big.Float
	exp := new(big.Float).SetInt(x).MantExp(&mant)
	m, _ := mant.Float64()
	return math.Log(m) + float64(exp)*math.Ln2
}

// bigFactorial returns n! exactly.
func bigFactorial(n int) *big.Int {
	if n < 2 {
		return big.NewInt(1)
	}
	return new(big.Int).MulRange(1, int64(n))
}

// OEIS A000142: 10! = 3 628 800; 20! = 2 432 902 008 176 640 000, the last term below 2⁶⁴.
func TestFactorial_sourceExample(t *testing.T) {
	if got, err := Factorial(10); err != nil || got != 3628800 {
		t.Errorf("Factorial(10) = (%d, %v), want (3628800, nil)", got, err)
	}
	if got, err := Factorial(20); err != nil || got != 2432902008176640000 {
		t.Errorf("Factorial(20) = (%d, %v), want (2432902008176640000, nil)", got, err)
	}
}

func TestFactorial_cases(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want uint64
	}{
		{"zero", 0, 1},
		{"one", 1, 1},
		{"two", 2, 2},
		{"five", 5, 120},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Factorial(tt.n); err != nil || got != tt.want {
				t.Errorf("Factorial(%d) = (%d, %v), want (%d, nil)", tt.n, got, err, tt.want)
			}
		})
	}
}

func TestFactorial_errors(t *testing.T) {
	for _, n := range []int{-1, math.MinInt} {
		if got, err := Factorial(n); !errors.Is(err, ErrNegative) || got != 0 {
			t.Errorf("Factorial(%d) = (%d, %v), want (0, ErrNegative)", n, got, err)
		}
	}
	for _, n := range []int{21, 22, math.MaxInt} {
		if got, err := Factorial(n); !errors.Is(err, ErrOverflow) || got != 0 {
			t.Errorf("Factorial(%d) = (%d, %v), want (0, ErrOverflow)", n, got, err)
		}
	}
}

func TestFactorial_oracle(t *testing.T) {
	for n := range 30 {
		got, err := Factorial(n)
		wantU64(t, "Factorial", got, err, bigFactorial(n))
	}
}

func TestFactorial_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkU, _ = Factorial(20) }); a != 0 {
		t.Errorf("Factorial allocates %v times, want 0", a)
	}
}

func BenchmarkFactorial(b *testing.B) {
	b.Run("n=20", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sinkU, _ = Factorial(20)
		}
	})
}
