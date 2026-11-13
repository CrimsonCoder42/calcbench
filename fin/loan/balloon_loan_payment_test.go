package loan

import (
	"errors"
	"math"
	"testing"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

// Hand-checked (Python, same float operations): $100,000, balloon $50,000, 0.5 %/month, 60 payments →
// (P − B·1.005⁻⁶⁰)·(A/P) = $1,216.64. Oracle: running BalloonBalance with that payment leaves the
// balloon within one cent per period (the payment is rounded to the cent once; here 3 cents).
func TestBalloonLoanPayment_example(t *testing.T) {
	got, err := BalloonLoanPayment(10_000_000, 5_000_000, 0.005, 60)
	if err != nil || got != 121_664 {
		t.Fatalf("got %d, %v; want 121664", got, err)
	}
	left, err := BalloonBalance(10_000_000, got, 0.005, 60)
	if err != nil || left != 5_000_003 {
		t.Fatalf("balloon left = %d, %v; want 5000003", left, err)
	}
}

// balloon = 0 is the fully amortizing payment (FIN-INT-018); balloon = principal is interest only.
func TestBalloonLoanPayment_limits(t *testing.T) {
	for _, c := range []struct {
		p int64
		r float64
		n int
	}{{3_200_000, 0.005, 36}, {14_000_000, 0.003, 240}, {1000, 0, 3}, {12345, 0.07, 1}} {
		got, err := BalloonLoanPayment(c.p, 0, c.r, c.n)
		want, _ := interest.AnnuityPaymentFromPresentValue(c.p, c.r, c.n)
		if err != nil || got != want {
			t.Errorf("%+v balloon 0: got %d, %v; want %d", c, got, err, want)
		}
	}
	if got, err := BalloonLoanPayment(10_000_000, 10_000_000, 0.005, 60); err != nil || got != 50_000 {
		t.Errorf("interest only: got %d, %v; want 50000", got, err)
	}
}

func TestBalloonLoanPayment_cases(t *testing.T) {
	cases := []struct {
		name               string
		principal, balloon int64
		rate               float64
		periods            int
		want               int64
	}{
		{"zeroRate", 1000, 400, 0, 3, 200},
		{"zeroPrincipalZeroBalloon", 0, 0, 0.01, 12, 0},
		{"balloonAtFullGrowth", 10_000, 10_100, 0.01, 1, 0},
		{"onePeriod", 100_000, 50_000, 0.01, 1, 51_000},
	}
	for _, c := range cases {
		got, err := BalloonLoanPayment(c.principal, c.balloon, c.rate, c.periods)
		if err != nil || got != c.want {
			t.Errorf("%s: got %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

func TestBalloonLoanPayment_errors(t *testing.T) {
	cases := []struct {
		name               string
		principal, balloon int64
		rate               float64
		periods            int
		want               error
	}{
		{"principalAboveMax", maxExactCents + 1, 0, 0.01, 12, ErrOutOfRange},
		{"balloonAboveMax", 100, maxExactCents + 1, 0.01, 12, ErrOutOfRange},
		{"principalNegative", -1, 0, 0.01, 12, ErrNegative},
		{"balloonNegative", 100, -1, 0.01, 12, ErrNegative},
		{"rateNaN", 100, 0, math.NaN(), 12, ErrNotFinite},
		{"rateInf", 100, 0, math.Inf(-1), 12, ErrRate},
		{"rateMinusOne", 100, 0, -1, 12, ErrRate},
		{"zeroPeriods", 100, 0, 0.01, 0, ErrPeriods},
		{"balloonTooLarge", 100, 1000, 0.01, 1, ErrDomain},
		{"paymentOverflow", maxExactCents, 0, 1, 1, ErrOverflow},
	}
	for _, c := range cases {
		if _, err := BalloonLoanPayment(c.principal, c.balloon, c.rate, c.periods); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestBalloonLoanPayment_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkInt, _ = BalloonLoanPayment(10_000_000, 5_000_000, 0.005, 60) }); a != 0 {
		t.Fatalf("allocs = %v, want 0", a)
	}
}

func BenchmarkBalloonLoanPayment(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = BalloonLoanPayment(10_000_000, 5_000_000, 0.005, 60)
	}
}
