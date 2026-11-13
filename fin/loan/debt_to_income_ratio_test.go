package loan

import (
	"errors"
	"testing"
)

// CFPB, "What is a debt-to-income ratio?": $1,500 + $100 + $400 = $2,000 of monthly debt on $6,000 gross
// monthly income is "33 percent". Exact comparison with 1.0/3: both are the correctly rounded quotient.
func TestDebtToIncomeRatio_cfpbExample(t *testing.T) {
	got, err := DebtToIncomeRatio(150_000+10_000+40_000, 600_000)
	if err != nil || got != 1.0/3 {
		t.Fatalf("got %v, %v; want 1/3", got, err)
	}
}

func TestDebtToIncomeRatio_cases(t *testing.T) {
	cases := []struct {
		name         string
		debt, income int64
		want         float64
	}{
		{"noDebt", 0, 600_000, 0},
		{"allIncome", 600_000, 600_000, 1},
		{"aboveOne", 900_000, 600_000, 1.5},
		{"maxExact", maxExactCents, maxExactCents, 1},
	}
	for _, c := range cases {
		got, err := DebtToIncomeRatio(c.debt, c.income)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestDebtToIncomeRatio_errors(t *testing.T) {
	cases := []struct {
		name         string
		debt, income int64
		want         error
	}{
		{"debtAboveMax", maxExactCents + 1, 100, ErrOutOfRange},
		{"incomeAboveMax", 100, maxExactCents + 1, ErrOutOfRange},
		{"debtNegative", -1, 100, ErrNegative},
		{"incomeNegative", 100, -1, ErrNegative},
		{"incomeZero", 100, 0, ErrZeroDenominator},
	}
	for _, c := range cases {
		if _, err := DebtToIncomeRatio(c.debt, c.income); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestDebtToIncomeRatio_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, _ = DebtToIncomeRatio(200_000, 600_000) }); a != 0 {
		t.Fatalf("allocs = %v, want 0", a)
	}
}

func BenchmarkDebtToIncomeRatio(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = DebtToIncomeRatio(200_000, 600_000)
	}
}
