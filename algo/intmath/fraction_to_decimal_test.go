package intmath

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"strings"
	"testing"
)

// Hardy & Wright §9.2: 1/3 = .333…, 3/7 = .142857142857…; a mixed recurring decimal such as 1/6 = .1666….
func TestFractionToDecimal_sourceExample(t *testing.T) {
	tests := []struct {
		num, den int64
		want     string
	}{{1, 3, "0.(3)"}, {3, 7, "0.(428571)"}, {1, 6, "0.1(6)"}, {1, 7, "0.(142857)"}}
	for _, tt := range tests {
		if got, err := FractionToDecimal(tt.num, tt.den, 100); err != nil || got != tt.want {
			t.Errorf("FractionToDecimal(%d, %d) = (%q, %v), want %q", tt.num, tt.den, got, err, tt.want)
		}
	}
}

func TestFractionToDecimal_cases(t *testing.T) {
	tests := []struct {
		num, den int64
		max      int
		want     string
	}{
		{0, 5, 0, "0"}, {0, -5, 0, "0"}, {2, 1, 0, "2"}, {1, 2, 1, "0.5"}, {-1, 2, 5, "-0.5"},
		{1, -2, 5, "-0.5"}, {-1, -2, 5, "0.5"}, {-1, 6, 2, "-0.1(6)"}, {22, 7, 6, "3.(142857)"},
		{1, 8, 3, "0.125"}, {4, 333, 3, "0.(012)"}, {1, 12, 3, "0.08(3)"}, {2, 6, 1, "0.(3)"},
		{math.MinInt64, -1, 0, "9223372036854775808"}, {math.MinInt64, 1, 0, "-9223372036854775808"},
		{math.MaxInt64, math.MinInt64, 63, "-0.999999999999999999891579782751449556599254719913005828857421875"},
		{1, math.MinInt64, 63, "-0.000000000000000000108420217248550443400745280086994171142578125"},
		{1, 3, 1, "0.(3)"},
	}
	for _, tt := range tests {
		if got, err := FractionToDecimal(tt.num, tt.den, tt.max); err != nil || got != tt.want {
			t.Errorf("FractionToDecimal(%d, %d, %d) = (%q, %v), want %q", tt.num, tt.den, tt.max, got, err, tt.want)
		}
	}
}

func TestFractionToDecimal_errors(t *testing.T) {
	tests := []struct {
		num, den int64
		max      int
		wantErr  error
	}{
		{1, 0, 10, ErrDivideByZero}, {0, 0, 10, ErrDivideByZero}, {1, 3, -1, ErrNegative},
		{1, 7, 5, ErrTooLong}, {1, 6, 1, ErrTooLong}, {1, 2, 0, ErrTooLong},
		{1, 1 << 62, 61, ErrTooLong}, // 2^−62 has 62 fractional digits
	}
	for _, tt := range tests {
		if got, err := FractionToDecimal(tt.num, tt.den, tt.max); !errors.Is(err, tt.wantErr) || got != "" {
			t.Errorf("FractionToDecimal(%d, %d, %d) = (%q, %v), want %v", tt.num, tt.den, tt.max, got, err, tt.wantErr)
		}
	}
}

// Oracle: parse the result back to an exact rational (finite part plus repetend/(10^p − 1)·10^−k) and
// compare it with numerator/denominator; and the repetend must be minimal (it starts at the first
// repeated remainder, so it cannot be shortened or started earlier).
func TestFractionToDecimal_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 500 {
		num := rng.Int64N(2_000_001) - 1_000_000
		den := rng.Int64N(2_001) - 1_000
		if den == 0 {
			continue
		}
		got, err := FractionToDecimal(num, den, 10_000)
		if err != nil {
			t.Fatalf("FractionToDecimal(%d, %d) error %v", num, den, err)
		}
		if parsed := parseRepeating(t, got); parsed.Cmp(big.NewRat(num, den)) != 0 {
			t.Fatalf("FractionToDecimal(%d, %d) = %q = %s, want %d/%d", num, den, got, parsed, num, den)
		}
	}
}

// parseRepeating reads "[-]I[.F[(R)]]" as an exact rational.
func parseRepeating(t *testing.T, s string) *big.Rat {
	t.Helper()
	neg := strings.HasPrefix(s, "-")
	s = strings.TrimPrefix(s, "-")
	intPart, frac, _ := strings.Cut(s, ".")
	fixed, rep, hasRep := strings.Cut(frac, "(")
	rep = strings.TrimSuffix(rep, ")")
	r, ok := new(big.Rat).SetString(intPart)
	if !ok {
		t.Fatalf("bad integer part in %q", s)
	}
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(len(fixed))), nil)
	if fixed != "" {
		f, _ := new(big.Int).SetString(fixed, 10)
		r.Add(r, new(big.Rat).SetFrac(f, scale))
	}
	if hasRep {
		p, _ := new(big.Int).SetString(rep, 10)
		den := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(len(rep))), nil)
		den.Sub(den, big.NewInt(1))
		den.Mul(den, scale)
		r.Add(r, new(big.Rat).SetFrac(p, den))
	}
	if neg {
		r.Neg(r)
	}
	return r
}

func BenchmarkFractionToDecimal(b *testing.B) {
	// Primes p for which 10 has order p − 1 (982, 9966, 99988 checked), so the expansion has p − 1 digits.
	for _, c := range []struct {
		name string
		den  int64
	}{{"1e3", 983}, {"1e4", 9_967}, {"1e5", 99_989}} {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				s, _ := FractionToDecimal(1, c.den, 200_000)
				benchSink = uint64(len(s))
			}
		})
	}
}
