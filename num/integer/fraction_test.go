package integer

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// Knuth §4.5.1 defines the operations; the cases below are hand computations, e.g.
// 2/7 + 3/8 = (16 + 21)/56 = 37/56; 2/7 − 3/8 = (16 − 21)/56 = −5/56; 2/7 × 3/8 = 6/56 = 3/28;
// 2/7 ÷ 3/8 = 16/21; 1/6 + 1/3 = 1/6 + 2/6 = 1/2 (exercises d₁ = 3, d₂ = 3).
func TestFraction_handExamples(t *testing.T) {
	type op func(a, b, c, d int64) (int64, int64, error)
	tests := []struct {
		name         string
		f            op
		a, b, c, d   int64
		wantN, wantD int64
	}{
		{"add", FractionAdd, 2, 7, 3, 8, 37, 56},
		{"sub", FractionSub, 2, 7, 3, 8, -5, 56},
		{"mul", FractionMul, 2, 7, 3, 8, 3, 28},
		{"div", FractionDiv, 2, 7, 3, 8, 16, 21},
		{"add common factor", FractionAdd, 1, 6, 1, 3, 1, 2},
		{"add to zero", FractionAdd, 1, 2, -2, 4, 0, 1},
		{"add negative dens", FractionAdd, 1, -2, 1, -3, -5, 6},
		{"mixed -2 3/4 + 3 5/7", FractionAdd, -11, 4, 26, 7, 27, 28},
		{"mul by zero", FractionMul, 0, 5, 7, 9, 0, 1},
		{"div negative", FractionDiv, 1, 2, -1, 4, -2, 1},
	}
	for _, tc := range tests {
		n, d, err := tc.f(tc.a, tc.b, tc.c, tc.d)
		if err != nil || n != tc.wantN || d != tc.wantD {
			t.Errorf("%s: got %d/%d, %v; want %d/%d", tc.name, n, d, err, tc.wantN, tc.wantD)
		}
	}
}

func TestFractionReduce_cases(t *testing.T) {
	tests := []struct{ n, d, wantN, wantD int64 }{
		{21, 98, 3, 14},
		{-6, -4, 3, 2},
		{6, -4, -3, 2},
		{0, -5, 0, 1},
		{math.MaxInt64, math.MaxInt64, 1, 1},
	}
	for _, tc := range tests {
		n, d, err := FractionReduce(tc.n, tc.d)
		if err != nil || n != tc.wantN || d != tc.wantD {
			t.Errorf("FractionReduce(%d, %d) = %d/%d, %v; want %d/%d", tc.n, tc.d, n, d, err, tc.wantN, tc.wantD)
		}
	}
}

func TestFraction_errors(t *testing.T) {
	if _, _, err := FractionReduce(1, 0); !errors.Is(err, ErrZeroDenominator) {
		t.Errorf("Reduce 1/0: %v", err)
	}
	if _, _, err := FractionReduce(math.MinInt64, 1); !errors.Is(err, ErrOverflow) {
		t.Errorf("Reduce MinInt64: %v", err)
	}
	if _, _, err := FractionAdd(1, 0, 1, 2); !errors.Is(err, ErrZeroDenominator) {
		t.Errorf("Add 1/0: %v", err)
	}
	if _, _, err := FractionSub(1, 2, 1, 0); !errors.Is(err, ErrZeroDenominator) {
		t.Errorf("Sub x/0: %v", err)
	}
	if _, _, err := FractionMul(1, 2, math.MinInt64, 3); !errors.Is(err, ErrOverflow) {
		t.Errorf("Mul MinInt64: %v", err)
	}
	if _, _, err := FractionDiv(1, 2, 0, 3); !errors.Is(err, ErrZeroDenominator) {
		t.Errorf("Div by 0/3: %v", err)
	}
	if _, _, err := FractionAdd(math.MaxInt64, 1, 1, 1); !errors.Is(err, ErrOverflow) {
		t.Errorf("Add overflow: %v", err)
	}
	if _, _, err := FractionMul(math.MaxInt64, 1, 2, 1); !errors.Is(err, ErrOverflow) {
		t.Errorf("Mul overflow: %v", err)
	}
	if _, _, err := FractionAdd(1, math.MaxInt64, 1, math.MaxInt64-1); !errors.Is(err, ErrOverflow) {
		t.Errorf("Add denominator overflow: %v", err)
	}
}

// The oracle is math/big.Rat, which keeps results in lowest terms.
func TestFraction_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	gen := func() int64 { return r.Int64N(1<<20) - 1<<19 }
	for range 500 {
		a, b, c, d := gen(), gen(), gen(), gen()
		if b == 0 || d == 0 {
			continue
		}
		x, y := big.NewRat(a, b), big.NewRat(c, d)
		check := func(name string, n, den int64, err error, want *big.Rat) {
			if err != nil {
				t.Fatalf("%s(%d/%d, %d/%d): %v", name, a, b, c, d, err)
			}
			if want.Num().Int64() != n || want.Denom().Int64() != den {
				t.Fatalf("%s(%d/%d, %d/%d) = %d/%d, want %s", name, a, b, c, d, n, den, want)
			}
		}
		n, den, err := FractionAdd(a, b, c, d)
		check("Add", n, den, err, new(big.Rat).Add(x, y))
		n, den, err = FractionSub(a, b, c, d)
		check("Sub", n, den, err, new(big.Rat).Sub(x, y))
		n, den, err = FractionMul(a, b, c, d)
		check("Mul", n, den, err, new(big.Rat).Mul(x, y))
		if c != 0 {
			n, den, err = FractionDiv(a, b, c, d)
			check("Div", n, den, err, new(big.Rat).Quo(x, y))
		}
	}
}

func TestFraction_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { _, _, _ = FractionAdd(2, 7, 3, 8) }); a != 0 {
		t.Errorf("FractionAdd allocates %v", a)
	}
}

func BenchmarkFractionAdd(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _, _ = FractionAdd(123456789, 987654321, 192837465, 564738291)
	}
}

func BenchmarkFractionMul(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _, _ = FractionMul(123456789, 987654321, 192837465, 564738291)
	}
}

func BenchmarkFractionReduce(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _, _ = FractionReduce(123456789, 987654321)
	}
}

func BenchmarkFractionSub(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _, _ = FractionSub(123456789, 987654321, 192837465, 564738291)
	}
}

func BenchmarkFractionDiv(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _, _ = FractionDiv(123456789, 987654321, 192837465, 564738291)
	}
}
