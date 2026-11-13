package loan

import (
	"errors"
	"math"
	"testing"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

// Hand-checked (Python, same float operations): income $6,000/month, housing ratio 0.28, taxes and
// insurance $400 → P&I allowance $1,280; at 0.5 %/month over 360 the loan is $213,493.27. Oracle: the
// level payment on that loan (FIN-INT-018) is the allowance to within a cent.
func TestAffordableLoanAmount_example(t *testing.T) {
	got, err := AffordableLoanAmount(600_000, 0.28, 40_000, 0.005, 360)
	if err != nil || got != 21_349_327 {
		t.Fatalf("got %d, %v; want 21349327", got, err)
	}
	pay, err := interest.AnnuityPaymentFromPresentValue(got, 0.005, 360)
	if err != nil || pay < 127_999 || pay > 128_001 {
		t.Fatalf("payment on the loan = %d, %v; want 128000 ± 1", pay, err)
	}
}

func TestAffordableLoanAmount_cases(t *testing.T) {
	cases := []struct {
		name    string
		income  int64
		ratio   float64
		ti      int64
		rate    float64
		periods int
		want    int64
	}{
		{"zeroRate", 600_000, 0.25, 0, 0, 12, 1_800_000},
		{"taxesExceedAllowance", 100_000, 0.1, 20_000, 0.005, 360, 0},
		{"taxesEqualAllowance", 100_000, 0.25, 25_000, 0.005, 360, 0},
		{"zeroRatio", 600_000, 0, 0, 0.005, 360, 0},
		{"zeroIncome", 0, 0.28, 0, 0.005, 360, 0},
		{"onePeriod", 100_000, 1, 0, 0.0, 1, 100_000},
	}
	for _, c := range cases {
		got, err := AffordableLoanAmount(c.income, c.ratio, c.ti, c.rate, c.periods)
		if err != nil || got != c.want {
			t.Errorf("%s: got %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

func TestAffordableLoanAmount_errors(t *testing.T) {
	cases := []struct {
		name    string
		income  int64
		ratio   float64
		ti      int64
		rate    float64
		periods int
		want    error
	}{
		{"incomeAboveMax", maxExactCents + 1, 0.28, 0, 0.005, 360, ErrOutOfRange},
		{"incomeNegative", -1, 0.28, 0, 0.005, 360, ErrNegative},
		{"ratioNaN", 600_000, math.NaN(), 0, 0.005, 360, ErrNotFinite},
		{"ratioInf", 600_000, math.Inf(1), 0, 0.005, 360, ErrNotFinite},
		{"ratioNegative", 600_000, -0.1, 0, 0.005, 360, ErrOutOfRange},
		{"ratioAboveOne", 600_000, 1.01, 0, 0.005, 360, ErrOutOfRange},
		{"tiAboveMax", 600_000, 0.28, maxExactCents + 1, 0.005, 360, ErrOutOfRange},
		{"tiNegative", 600_000, 0.28, -1, 0.005, 360, ErrNegative},
		{"rateNaN", 600_000, 0.28, 0, math.NaN(), 360, ErrNotFinite},
		{"rateInf", 600_000, 0.28, 0, math.Inf(1), 360, ErrRate},
		{"rateMinusOne", 600_000, 0.28, 0, -1, 360, ErrRate},
		{"zeroPeriods", 600_000, 0.28, 0, 0.005, 0, ErrPeriods},
		{"loanOverflow", maxExactCents, 1, 0, 0, 2, ErrOverflow},
	}
	for _, c := range cases {
		if _, err := AffordableLoanAmount(c.income, c.ratio, c.ti, c.rate, c.periods); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestAffordableLoanAmount_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkInt, _ = AffordableLoanAmount(600_000, 0.28, 40_000, 0.005, 360) }); a != 0 {
		t.Fatalf("allocs = %v, want 0", a)
	}
}

func BenchmarkAffordableLoanAmount(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = AffordableLoanAmount(600_000, 0.28, 40_000, 0.005, 360)
	}
}
