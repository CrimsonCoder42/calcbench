package loan

import (
	"errors"
	"math"
	"testing"
)

// Hand-checked from the Fannie Mae Selling Guide B2-1.2-01 definition (loan ÷ property value):
// $360,000 on $400,000 → 0.9. A correctly rounded division gives the double nearest 0.9 exactly.
func TestLoanToValueRatio_example(t *testing.T) {
	got, err := LoanToValueRatio(36_000_000, 40_000_000)
	if err != nil || got != 0.9 {
		t.Fatalf("got %v, %v; want 0.9", got, err)
	}
}

func TestLoanToValueRatio_cases(t *testing.T) {
	cases := []struct {
		name        string
		loan, value int64
		want        float64
	}{
		{"zeroLoan", 0, 100, 0},
		{"equal", 100, 100, 1},
		{"underwater", 150, 100, 1.5},
		{"maxExact", maxExactCents, maxExactCents, 1},
		{"oneCentOfMax", 1, maxExactCents, 1.0 / maxExactCents},
	}
	for _, c := range cases {
		got, err := LoanToValueRatio(c.loan, c.value)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestLoanToValueRatio_errors(t *testing.T) {
	cases := []struct {
		name        string
		loan, value int64
		want        error
	}{
		{"loanAboveMax", maxExactCents + 1, 100, ErrOutOfRange},
		{"valueBelowMin", 100, math.MinInt64, ErrOutOfRange},
		{"loanNegative", -1, 100, ErrNegative},
		{"valueNegative", 100, -1, ErrNegative},
		{"valueZero", 100, 0, ErrZeroDenominator},
	}
	for _, c := range cases {
		if _, err := LoanToValueRatio(c.loan, c.value); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestLoanToValueRatio_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, _ = LoanToValueRatio(36_000_000, 40_000_000) }); a != 0 {
		t.Fatalf("allocs = %v, want 0", a)
	}
}

func BenchmarkLoanToValueRatio(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = LoanToValueRatio(36_000_000, 40_000_000)
	}
}
