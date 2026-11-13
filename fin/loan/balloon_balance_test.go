package loan

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

// Hand-checked (Python, same float operations): a 30-year $100,000 loan at 0.5 %/month has the level
// payment $599.55; after 60 payments the balloon is $93,054.37.
func TestBalloonBalance_thirtyYearFiveYearBalloon(t *testing.T) {
	payment, err := interest.AnnuityPaymentFromPresentValue(10_000_000, 0.005, 360)
	if err != nil || payment != 59_955 {
		t.Fatalf("payment = %d, %v; want 59955", payment, err)
	}
	got, err := BalloonBalance(10_000_000, payment, 0.005, 60)
	if err != nil || got != 9_305_437 {
		t.Fatalf("got %d, %v; want 9305437", got, err)
	}
}

// Oracle: before the last period, the balloon balance under the level payment equals the schedule's
// remaining balance (same recursion, same rounding).
func TestBalloonBalance_matchesRemainingBalance(t *testing.T) {
	const principal, rate, periods = 14_000_000, 0.003, 240
	payment, _ := interest.AnnuityPaymentFromPresentValue(principal, rate, periods)
	for _, k := range []int{0, 1, 36, 120, 239} {
		bb, err1 := BalloonBalance(principal, payment, rate, k)
		rb, err2 := RemainingBalance(principal, rate, periods, k)
		if err1 != nil || err2 != nil || bb != rb {
			t.Errorf("k=%d: BalloonBalance %d (%v), RemainingBalance %d (%v)", k, bb, err1, rb, err2)
		}
	}
}

func TestBalloonBalance_cases(t *testing.T) {
	cases := []struct {
		name               string
		principal, payment int64
		rate               float64
		made               int
		want               int64
	}{
		{"noPayments", 1000, 100, 0.01, 0, 1000},
		{"zeroRate", 1000, 100, 0, 3, 700},
		{"paidOffIsZero", 100, 60, 0, 2, 0},
		{"staysZeroAfterPayoff", 100, 60, 0, 50, 0},
		{"exactPayoff", 100, 50, 0, 2, 0},
		{"zeroPrincipal", 0, 10, 0.01, 5, 0},
		{"interestOnly", 100_000, 1000, 0.01, 12, 100_000},
		{"negativeAmortization", 100_000, 0, 0.01, 1, 101_000},
		{"negativeRate", 100_000, 0, -0.5, 1, 50_000},
		{"maxExactNoPayments", maxExactCents, 0, 0.01, 0, maxExactCents},
	}
	for _, c := range cases {
		got, err := BalloonBalance(c.principal, c.payment, c.rate, c.made)
		if err != nil || got != c.want {
			t.Errorf("%s: got %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

func TestBalloonBalance_errors(t *testing.T) {
	cases := []struct {
		name               string
		principal, payment int64
		rate               float64
		made               int
		want               error
	}{
		{"principalAboveMax", maxExactCents + 1, 0, 0.01, 1, ErrOutOfRange},
		{"paymentAboveMax", 100, maxExactCents + 1, 0.01, 1, ErrOutOfRange},
		{"principalNegative", -1, 0, 0.01, 1, ErrNegative},
		{"paymentNegative", 100, -1, 0.01, 1, ErrNegative},
		{"rateNaN", 100, 0, math.NaN(), 1, ErrNotFinite},
		{"rateInf", 100, 0, math.Inf(1), 1, ErrRate},
		{"rateMinusOne", 100, 0, -1, 1, ErrRate},
		{"madeNegative", 100, 0, 0.01, -1, ErrPeriods},
		{"madeAboveCap", 10_000_000, 50_000, 0.005, math.MaxInt, ErrPeriods},
		{"madeJustAboveCap", 100, 0, 0.01, maxSchedulePeriods + 1, ErrPeriods},
		{"balanceOverflow", maxExactCents, 0, 1, 1, ErrOverflow},
		{"interestOverflow", maxExactCents, 0, 3, 1, ErrOverflow},
	}
	for _, c := range cases {
		if _, err := BalloonBalance(c.principal, c.payment, c.rate, c.made); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestBalloonBalance_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(20, func() { sinkInt, _ = BalloonBalance(10_000_000, 59_955, 0.005, 60) }); a != 0 {
		t.Fatalf("allocs = %v, want 0", a)
	}
}

func BenchmarkBalloonBalance(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkInt, _ = BalloonBalance(100_000_000_000, 10_000_000, 0.0001, n)
			}
		})
	}
}
