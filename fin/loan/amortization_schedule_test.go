package loan

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

var (
	sinkInt   int64
	sinkFloat float64
	sinkRows  []AmortizationRow
)

// OpenStax, Principles of Finance, §8.3 Table 8.4: $32,000 at 6 %/yr monthly (0.5 %), 36 payments,
// payment $973.50; rows 1–3 and row 36's interest ($4.84) are compared exactly. The book's row 36 shows
// principal 968.66 because it carries the unrounded payment; here the last payment absorbs the
// rounding of the level payment (97360 = 484 + 96876).
func TestAmortizationSchedule_openStaxAutoLoan(t *testing.T) {
	rows, err := AmortizationSchedule(3_200_000, 0.005, 36)
	if err != nil {
		t.Fatal(err)
	}
	want := []AmortizationRow{
		{1, 97350, 16000, 81350, 3118650},
		{2, 97350, 15593, 81757, 3036893},
		{3, 97350, 15184, 82166, 2954727},
	}
	for k, w := range want {
		if rows[k] != w {
			t.Errorf("row %d = %+v, want %+v", k+1, rows[k], w)
		}
	}
	if last := rows[35]; last != (AmortizationRow{36, 97360, 484, 96876, 0}) {
		t.Errorf("row 36 = %+v", last)
	}
}

// OpenStax, Principles of Finance, §8.3: $140,000 at 3.6 %/yr monthly (0.3 %), 240 payments, payment
// $819.16; row 239 interest 4.89 and principal 814.27, row 240 interest 2.45.
func TestAmortizationSchedule_openStaxMortgage(t *testing.T) {
	rows, err := AmortizationSchedule(14_000_000, 0.003, 240)
	if err != nil {
		t.Fatal(err)
	}
	if rows[0].PaymentCents != 81916 {
		t.Fatalf("payment = %d, want 81916", rows[0].PaymentCents)
	}
	if r := rows[238]; r.InterestCents != 489 || r.PrincipalCents != 81427 {
		t.Errorf("row 239 = %+v, want interest 489, principal 81427", r)
	}
	if r := rows[239]; r.InterestCents != 245 || r.BalanceCents != 0 {
		t.Errorf("row 240 = %+v, want interest 245, balance 0", r)
	}
}

// reconcile checks the schedule identities: at most periods rows, Σprincipal = loan, final balance 0,
// no negative balance, each payment = interest + principal, each balance = previous − principal, every
// row but the last pays the level payment, each interest = round_half_even(exact previous balance ×
// rate) (big.Rat oracle), and RemainingBalance equals every row and is 0 after the last row.
func reconcile(t *testing.T, principal int64, rate float64, periods int, rows []AmortizationRow) {
	t.Helper()
	n := len(rows)
	if n > periods || (principal > 0 && n == 0) {
		t.Fatalf("len = %d for periods %d", n, periods)
	}
	level, err := interest.AnnuityPaymentFromPresentValue(principal, rate, periods)
	if err != nil {
		t.Fatal(err)
	}
	rateRat := new(big.Rat).SetFloat64(rate)
	var sumPrincipal int64
	prev := principal
	for k, r := range rows {
		if r.Period != k+1 {
			t.Fatalf("row %d: Period = %d", k+1, r.Period)
		}
		if r.PaymentCents != r.InterestCents+r.PrincipalCents {
			t.Fatalf("row %d: payment %d ≠ interest + principal", k+1, r.PaymentCents)
		}
		if r.BalanceCents != prev-r.PrincipalCents || r.BalanceCents < 0 {
			t.Fatalf("row %d: balance %d, previous %d, principal %d", k+1, r.BalanceCents, prev, r.PrincipalCents)
		}
		if k < n-1 && r.PaymentCents != level {
			t.Fatalf("row %d: payment %d ≠ level %d", k+1, r.PaymentCents, level)
		}
		exact := new(big.Rat).Mul(new(big.Rat).SetInt64(prev), rateRat)
		f, _ := exact.Float64()
		if want := math.RoundToEven(f); float64(r.InterestCents) != want {
			t.Fatalf("row %d: interest %d, exact %v", k+1, r.InterestCents, exact.FloatString(6))
		}
		if rb, err := RemainingBalance(principal, rate, periods, k+1); err != nil || rb != r.BalanceCents {
			t.Fatalf("row %d: RemainingBalance = %d, %v; want %d", k+1, rb, err, r.BalanceCents)
		}
		sumPrincipal += r.PrincipalCents
		prev = r.BalanceCents
	}
	if sumPrincipal != principal {
		t.Fatalf("Σprincipal = %d, want %d", sumPrincipal, principal)
	}
	if prev != 0 {
		t.Fatalf("final balance = %d, want 0", prev)
	}
	if n < periods {
		if rb, err := RemainingBalance(principal, rate, periods, periods); err != nil || rb != 0 {
			t.Fatalf("RemainingBalance after payoff = %d, %v; want 0", rb, err)
		}
	}
}

func TestAmortizationSchedule_reconciles(t *testing.T) {
	cases := []struct {
		name      string
		principal int64
		rate      float64
		periods   int
		wantRows  int
	}{
		{"auto", 3_200_000, 0.005, 36, 36},
		{"mortgage", 14_000_000, 0.003, 240, 240},
		{"thirtyYear", 50_000_000, 0.065 / 12, 360, 360},
		{"zeroRate", 1000, 0, 3, 3},
		{"onePeriod", 100_000, 0.01, 1, 1},
		{"negativeRate", 1_000_000, -0.001, 24, 24},
		{"oneCentPaymentRoundsToZero", 1, 0.01, 12, 12},
		{"zeroPrincipal", 0, 0.01, 5, 0},
		{"highRate", 1_000_000, 0.5, 40, 40},
		{"tinyLoanEarlyPayoff", 100, 0.005, 360, 100},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rows, err := AmortizationSchedule(c.principal, c.rate, c.periods)
			if err != nil {
				t.Fatal(err)
			}
			if rows == nil || len(rows) != c.wantRows {
				t.Fatalf("len = %d (nil %v), want %d", len(rows), rows == nil, c.wantRows)
			}
			reconcile(t, c.principal, c.rate, c.periods, rows)
		})
	}
}

// Review case: $1.00 at 0.5 %/month over 360 has a level payment of 1 cent (0.5996 rounded); the
// interest rounds to 0 every row, so the loan is repaid in row 100 and the schedule ends there.
func TestAmortizationSchedule_tinyLoan(t *testing.T) {
	rows, err := AmortizationSchedule(100, 0.005, 360)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 100 || cap(rows) != 360 {
		t.Fatalf("len %d cap %d, want 100, 360", len(rows), cap(rows))
	}
	if rows[0] != (AmortizationRow{1, 1, 0, 1, 99}) || rows[99] != (AmortizationRow{100, 1, 0, 1, 0}) {
		t.Fatalf("rows[0] = %+v, rows[99] = %+v", rows[0], rows[99])
	}
}

// Generated small loans (1…1000 cents, fixed seed), where the rounded payment often repays early: every
// schedule must reconcile, never go negative, and end at 0.
func TestAmortizationSchedule_randomSmallLoans(t *testing.T) {
	rng := rand.New(rand.NewPCG(3, 4))
	for trial := 0; trial < 300; trial++ {
		principal := 1 + rng.Int64N(1000)
		rate := rng.Float64() * 0.02
		periods := 1 + rng.IntN(400)
		rows, err := AmortizationSchedule(principal, rate, periods)
		if err != nil {
			t.Fatalf("trial %d: %v", trial, err)
		}
		reconcile(t, principal, rate, periods, rows)
	}
}

// The payoff payment differs from the level payment by at most (F/A,i,n) cents for i ≥ 0 (payment
// rounding and per-row interest rounding, each ≤ ½ cent, compounded to the end).
func TestAmortizationSchedule_lastPaymentBound(t *testing.T) {
	rng := rand.New(rand.NewPCG(5, 6))
	for trial := 0; trial < 300; trial++ {
		principal := 10_000 + rng.Int64N(100_000_000_000)
		rate := rng.Float64() * 0.01
		periods := 1 + rng.IntN(480)
		rows, err := AmortizationSchedule(principal, rate, periods)
		if err != nil {
			t.Fatal(err)
		}
		if len(rows) != periods {
			continue
		}
		fa, _ := interest.UniformSeriesCompoundAmount(rate, periods)
		if d := math.Abs(float64(rows[periods-1].PaymentCents - rows[0].PaymentCents)); periods > 1 && d > fa {
			t.Fatalf("trial %d: last − level = %v cents > (F/A) = %v", trial, d, fa)
		}
	}
}

func TestAmortizationSchedule_randomReconcile(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := 0; trial < 200; trial++ {
		principal := rng.Int64N(100_000_000_000)
		rate := rng.Float64() * 0.02
		periods := 1 + rng.IntN(400)
		rows, err := AmortizationSchedule(principal, rate, periods)
		if err != nil {
			t.Fatalf("trial %d: %v", trial, err)
		}
		reconcile(t, principal, rate, periods, rows)
	}
}

func TestAmortizationSchedule_zeroRate(t *testing.T) {
	rows, err := AmortizationSchedule(1000, 0, 3)
	if err != nil {
		t.Fatal(err)
	}
	want := []AmortizationRow{{1, 333, 0, 333, 667}, {2, 333, 0, 333, 334}, {3, 334, 0, 334, 0}}
	for k := range want {
		if rows[k] != want[k] {
			t.Errorf("row %d = %+v, want %+v", k+1, rows[k], want[k])
		}
	}
}

func TestAmortizationSchedule_maxExactPrincipal(t *testing.T) {
	rows, err := AmortizationSchedule(maxExactCents, 0, 2)
	if err != nil {
		t.Fatal(err)
	}
	reconcile(t, maxExactCents, 0, 2, rows)
}

func TestAmortizationSchedule_errors(t *testing.T) {
	cases := []struct {
		name      string
		principal int64
		rate      float64
		periods   int
		want      error
	}{
		{"principalAboveMax", maxExactCents + 1, 0.01, 12, ErrOutOfRange},
		{"principalNegative", -1, 0.01, 12, ErrNegative},
		{"rateNaN", 1000, math.NaN(), 12, ErrNotFinite},
		{"rateInf", 1000, math.Inf(1), 12, ErrRate},
		{"rateNegInf", 1000, math.Inf(-1), 12, ErrRate},
		{"rateMinusOne", 1000, -1, 12, ErrRate},
		{"zeroPeriods", 1000, 0.01, 0, ErrPeriods},
		{"negativePeriods", 1000, 0.01, -1, ErrPeriods},
		{"tooManyPeriods", 1000, 0.01, maxSchedulePeriods + 1, ErrPeriods},
		{"paymentOverflow", maxExactCents, 1, 1, ErrOverflow},
	}
	for _, c := range cases {
		rows, err := AmortizationSchedule(c.principal, c.rate, c.periods)
		if !errors.Is(err, c.want) || rows != nil {
			t.Errorf("%s: rows %v, err = %v, want %v", c.name, rows != nil, err, c.want)
		}
	}
}

// A payment overflow comes from fin/interest; it must match both packages' sentinels.
func TestAmortizationSchedule_interestErrorMapped(t *testing.T) {
	_, err := AmortizationSchedule(maxExactCents, 1, 1)
	if !errors.Is(err, ErrOverflow) || !errors.Is(err, interest.ErrOverflow) {
		t.Fatalf("err = %v, want loan.ErrOverflow and interest.ErrOverflow", err)
	}
}

func TestAmortizationSchedule_oneAllocation(t *testing.T) {
	if a := testing.AllocsPerRun(20, func() { sinkRows, _ = AmortizationSchedule(14_000_000, 0.003, 240) }); a != 1 {
		t.Fatalf("allocs = %v, want 1", a)
	}
}

func BenchmarkAmortizationSchedule(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkRows, _ = AmortizationSchedule(100_000_000_000, 0.0001, n)
			}
		})
	}
}
