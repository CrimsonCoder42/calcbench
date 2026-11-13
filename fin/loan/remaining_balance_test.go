package loan

import (
	"errors"
	"fmt"
	"math"
	"testing"
)

// OpenStax, Principles of Finance, §8.3 Table 8.4: balance after payment 3 is $29,547.27; the mortgage
// rows 12 and 120 equal the schedule (see amortization_schedule_test.go for the full reconciliation).
func TestRemainingBalance_openStax(t *testing.T) {
	cases := []struct {
		principal int64
		rate      float64
		periods   int
		made      int
		want      int64
	}{
		{3_200_000, 0.005, 36, 3, 2_954_727},
		{3_200_000, 0.005, 36, 0, 3_200_000},
		{3_200_000, 0.005, 36, 36, 0},
		{3_200_000, 0.005, 36, 12, 2_196_502},
		{14_000_000, 0.003, 240, 120, 8_244_692},
		{14_000_000, 0.003, 240, 239, 81_545},
	}
	for _, c := range cases {
		got, err := RemainingBalance(c.principal, c.rate, c.periods, c.made)
		if err != nil || got != c.want {
			t.Errorf("%+v: got %d, %v", c, got, err)
		}
	}
}

func TestRemainingBalance_cases(t *testing.T) {
	cases := []struct {
		name      string
		principal int64
		rate      float64
		periods   int
		made      int
		want      int64
	}{
		{"zeroRate", 1000, 0, 3, 2, 334},
		{"onePeriodPaid", 100_000, 0.01, 1, 1, 0},
		{"zeroPrincipal", 0, 0.01, 12, 6, 0},
		{"maxExactNoneMade", maxExactCents, 0.001, 12, 0, maxExactCents},
		{"afterEarlyPayoff", 100, 0.005, 360, 200, 0},
		{"atEarlyPayoff", 100, 0.005, 360, 100, 0},
		{"beforeEarlyPayoff", 100, 0.005, 360, 99, 1},
	}
	for _, c := range cases {
		got, err := RemainingBalance(c.principal, c.rate, c.periods, c.made)
		if err != nil || got != c.want {
			t.Errorf("%s: got %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

func TestRemainingBalance_errors(t *testing.T) {
	cases := []struct {
		name      string
		principal int64
		rate      float64
		periods   int
		made      int
		want      error
	}{
		{"principalAboveMax", maxExactCents + 1, 0.01, 12, 1, ErrOutOfRange},
		{"principalNegative", -5, 0.01, 12, 1, ErrNegative},
		{"rateNaN", 1000, math.NaN(), 12, 1, ErrNotFinite},
		{"rateInf", 1000, math.Inf(1), 12, 1, ErrRate},
		{"rateBelowMinusOne", 1000, -1.5, 12, 1, ErrRate},
		{"zeroPeriods", 1000, 0.01, 0, 0, ErrPeriods},
		{"madeNegative", 1000, 0.01, 12, -1, ErrPeriods},
		{"madeAbovePeriods", 1000, 0.01, 12, 13, ErrPeriods},
		{"paymentOverflow", maxExactCents, 1, 1, 1, ErrOverflow},
		{"periodsAboveCap", 10_000_000, 0.005, 10_000_000_000, 10_000_000_000, ErrPeriods},
		{"periodsJustAboveCap", 1000, 0.005, maxSchedulePeriods + 1, 1, ErrPeriods},
	}
	for _, c := range cases {
		if _, err := RemainingBalance(c.principal, c.rate, c.periods, c.made); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestRemainingBalance_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(20, func() { sinkInt, _ = RemainingBalance(14_000_000, 0.003, 240, 120) }); a != 0 {
		t.Fatalf("allocs = %v, want 0", a)
	}
}

func BenchmarkRemainingBalance(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkInt, _ = RemainingBalance(100_000_000_000, 0.0001, n, n-1)
			}
		})
	}
}
