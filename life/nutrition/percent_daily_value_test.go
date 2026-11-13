package nutrition

import (
	"math"
	"testing"
)

// Results are whole numbers, compared exactly.

// 21 CFR 101.9(d)(7)(ii) with DRVs from 101.9(c)(9) (adults and children ≥ 4 years); HAND-DERIVED:
// sodium 1150 / 2300 × 100 = 50; total carbohydrate 30 / 275 × 100 = 10.909… → 11.
func TestPercentDailyValue_example(t *testing.T) {
	if got, err := PercentDailyValue(1150, 2300); err != nil || got != 50 {
		t.Errorf("sodium: got (%v, %v), want 50", got, err)
	}
	if got, err := PercentDailyValue(30, 275); err != nil || got != 11 {
		t.Errorf("carbohydrate: got (%v, %v), want 11", got, err)
	}
}

// parity: ledger percent_daily_value, source cases 1–7. Case 3 (2.5 g of 20 g = 12.5 %): the FDA Food
// Labeling Guide (2013), Appendix H, p. 129, rounds an exact half up, so 13 (agrees with the source; D-073).
func TestPercentDailyValue_parity(t *testing.T) {
	tests := []struct {
		name             string
		amount, dv, want float64
	}{
		{"source case 1", 25, 50, 50},
		{"source case 2", 30, 275, 11},
		{"source case 3 (exact half rounds up)", 2.5, 20, 13},
		{"source case 4", 3, 28, 11},
		{"source case 5", 1150, 2300, 50},
		{"source case 6", 2400, 2300, 104},
		{"source case 7", 0, 300, 0},
	}
	for _, tc := range tests {
		if got, err := PercentDailyValue(tc.amount, tc.dv); err != nil || got != tc.want {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger percent_daily_value, source refusal: a DV of 10 g (the 1–3 years saturated-fat DRV) was
// refused by the source's adults-column bound; the regulation's formula applies to any DRV, so it is accepted:
// 5 / 10 × 100 = 50.
func TestPercentDailyValue_parityChildDRV(t *testing.T) {
	if got, err := PercentDailyValue(5, 10); err != nil || got != 50 {
		t.Errorf("got (%v, %v), want 50", got, err)
	}
}

// Rounding edges (FDA Food Labeling Guide, Appendix H: "2.5 to 2.99" round up, "2.01 to 2.49" round down).
// Exact ties need amount/DV exact in binary: 1/8 = 12.5 % → 13, 3/8 = 37.5 % → 38, 5/8 = 62.5 % → 63,
// 7/8 = 87.5 % → 88 (half up); 0.49 % → 0; 12.49 % → 12; 12.51 % → 13.
func TestPercentDailyValue_rounding(t *testing.T) {
	tests := []struct{ amount, dv, want float64 }{
		{1, 8, 13}, {3, 8, 38}, {5, 8, 63}, {7, 8, 88}, {0.49, 100, 0}, {12.49, 100, 12}, {12.51, 100, 13},
	}
	for _, tc := range tests {
		if got, err := PercentDailyValue(tc.amount, tc.dv); err != nil || got != tc.want {
			t.Errorf("(%v, %v) = (%v, %v), want %v", tc.amount, tc.dv, got, err, tc.want)
		}
	}
}

func TestPercentDailyValue_errors(t *testing.T) {
	f := func(amount, dv float64) func() error {
		return func() error { _, err := PercentDailyValue(amount, dv); return err }
	}
	runErrCases(t, []errCase{
		{"amount NaN", f(nan, 50), ErrNotFinite},
		{"amount +Inf", f(pInf, 50), ErrNotFinite},
		{"dv -Inf", f(10, nInf), ErrNotFinite},
		{"dv NaN", f(10, nan), ErrNotFinite},
		{"dv 0", f(10, 0), ErrOutOfRange},
		{"dv negative", f(10, -50), ErrOutOfRange},
		{"amount negative", f(-0.1, 50), ErrOutOfRange},
		{"overflow", f(math.MaxFloat64, 5e-324), ErrDomain},
	})
}

func TestPercentDailyValue_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = PercentDailyValue(30, 275) }); n != 0 {
		t.Errorf("PercentDailyValue allocates %v times, want 0", n)
	}
}

func BenchmarkPercentDailyValue(b *testing.B) {
	b.ReportAllocs()
	amount, dv := 30.0, 275.0
	for b.Loop() {
		sink, _ = PercentDailyValue(amount, dv)
	}
}
