package loan

import (
	"errors"
	"fmt"
	"math"
	"testing"
)

// Hand-checked (Python, same operations) on the OpenStax §8.3 mortgage ($140,000, 0.3 %/month, 240):
// $100 extra a month pays it off in period 204 and saves $9,257.23 of interest.
func TestExtraPaymentPayoff_openStaxMortgage(t *testing.T) {
	n, saved, err := ExtraPaymentPayoff(14_000_000, 0.003, 240, 10_000)
	if err != nil || n != 204 || saved != 925_723 {
		t.Fatalf("got (%d, %d, %v), want (204, 925723)", n, saved, err)
	}
}

// Oracle: the saving equals Σinterest of the schedule minus Σinterest of the prepaid run; with no
// extra the payoff is the full term and nothing is saved; paying the whole loan in period 1 saves all
// interest but the first period's.
func TestExtraPaymentPayoff_cases(t *testing.T) {
	rows, _ := AmortizationSchedule(3_200_000, 0.005, 36)
	var total int64
	for _, r := range rows {
		total += r.InterestCents
	}
	cases := []struct {
		name      string
		principal int64
		rate      float64
		periods   int
		extra     int64
		wantN     int
		wantSaved int64
	}{
		{"noExtra", 3_200_000, 0.005, 36, 0, 36, 0},
		{"payAllAtOnce", 3_200_000, 0.005, 36, 3_200_000, 1, total - 16_000},
		{"zeroPrincipal", 0, 0.005, 36, 100, 0, 0},
		{"zeroRate", 1200, 0, 12, 100, 6, 0},
		{"onePeriod", 100_000, 0.01, 1, 500, 1, 0},
		{"tinyLoanNoExtraEndsWithSchedule", 100, 0.005, 360, 0, 100, 0},
	}
	for _, c := range cases {
		n, saved, err := ExtraPaymentPayoff(c.principal, c.rate, c.periods, c.extra)
		if err != nil || n != c.wantN || saved != c.wantSaved {
			t.Errorf("%s: got (%d, %d, %v), want (%d, %d)", c.name, n, saved, err, c.wantN, c.wantSaved)
		}
	}
	if total != 304_610 {
		t.Fatalf("schedule interest = %d, want 304610", total)
	}
}

// More extra never takes longer and never saves less.
func TestExtraPaymentPayoff_monotone(t *testing.T) {
	prevN, prevSaved := 241, int64(-1)
	for extra := int64(0); extra <= 200_000; extra += 5_000 {
		n, saved, err := ExtraPaymentPayoff(14_000_000, 0.003, 240, extra)
		if err != nil || n > prevN || saved < prevSaved {
			t.Fatalf("extra %d: (%d, %d, %v) after (%d, %d)", extra, n, saved, err, prevN, prevSaved)
		}
		prevN, prevSaved = n, saved
	}
}

func TestExtraPaymentPayoff_errors(t *testing.T) {
	cases := []struct {
		name      string
		principal int64
		rate      float64
		periods   int
		extra     int64
		want      error
	}{
		{"principalAboveMax", maxExactCents + 1, 0.01, 12, 0, ErrOutOfRange},
		{"extraAboveMax", 1000, 0.01, 12, maxExactCents + 1, ErrOutOfRange},
		{"principalNegative", -1, 0.01, 12, 0, ErrNegative},
		{"extraNegative", 1000, 0.01, 12, -1, ErrNegative},
		{"rateNaN", 1000, math.NaN(), 12, 0, ErrNotFinite},
		{"rateInf", 1000, math.Inf(1), 12, 0, ErrRate},
		{"rateMinusOne", 1000, -1, 12, 0, ErrRate},
		{"zeroPeriods", 1000, 0.01, 0, 0, ErrPeriods},
		{"periodsAboveCap", 1, 2, 1 << 40, 0, ErrPeriods},
		{"periodsJustAboveCap", 1000, 0.01, maxSchedulePeriods + 1, 0, ErrPeriods},
		{"paymentPlusExtraOverflow", maxExactCents, 0, 1, 1, ErrOverflow},
		{"interestTotalOverflow", maxExactCents / 2, 1, 3, 0, ErrOverflow},
	}
	for _, c := range cases {
		if _, _, err := ExtraPaymentPayoff(c.principal, c.rate, c.periods, c.extra); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestExtraPaymentPayoff_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(20, func() { _, sinkInt, _ = ExtraPaymentPayoff(14_000_000, 0.003, 240, 10_000) }); a != 0 {
		t.Fatalf("allocs = %v, want 0", a)
	}
}

func BenchmarkExtraPaymentPayoff(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, sinkInt, _ = ExtraPaymentPayoff(100_000_000_000, 0.0001, n, 1)
			}
		})
	}
}
