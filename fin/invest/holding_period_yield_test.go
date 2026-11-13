package invest

import "testing"

// OpenStax Principles of Finance §15.1: Target bought at $128.74, sold at $176.53, $2.68 dividends →
// total dollar return $50.47 and 39.20 % (displayed to 2 decimals of a percent, so ±5e-5); 3M per
// share: $175.00 → $169.72 with $5.91 dividends → 0.36 %.
func TestHoldingPeriodYield_openStax(t *testing.T) {
	got, err := HoldingPeriodYield(12874, 17653, 268)
	if err != nil || got != 5047.0/12874.0 || !near(got, 0.3920, 5e-5) {
		t.Errorf("Target: %v, %v; want 5047/12874 ≈ 0.3920", got, err)
	}
	got, err = HoldingPeriodYield(17500, 16972, 591)
	if err != nil || !near(got, 0.0036, 5e-5) {
		t.Errorf("3M: %v, %v; want ≈ 0.0036", got, err)
	}
	// Identity: HPY = capital gains yield + dividend yield.
	cg, _ := CapitalGainsYield(12874, 17653)
	dy, _ := DividendYield(268, 12874)
	if hpy, _ := HoldingPeriodYield(12874, 17653, 268); !near(hpy, cg+dy, 1e-15) {
		t.Errorf("HPY %v ≠ CGY + DY %v", hpy, cg+dy)
	}
}

func TestHoldingPeriodYield_cases(t *testing.T) {
	cases := []struct {
		name               string
		begin, end, income int64
		want               float64
	}{
		{"no change", 1000, 1000, 0, 0},
		{"total loss", 1000, 0, 0, -1},
		{"negative income", 1000, 1000, -10, -0.01},
		{"extremes", maxCents, maxCents, maxCents, 1},
	}
	for _, c := range cases {
		got, err := HoldingPeriodYield(c.begin, c.end, c.income)
		if err != nil || got != c.want {
			t.Errorf("%s: %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestHoldingPeriodYield_errors(t *testing.T) {
	_, err := HoldingPeriodYield(0, 100, 0)
	wantErr(t, "begin 0", err, ErrZeroDenominator)
	_, err = HoldingPeriodYield(-1, 100, 0)
	wantErr(t, "begin < 0", err, ErrNegative)
	_, err = HoldingPeriodYield(100, -1, 0)
	wantErr(t, "end < 0", err, ErrNegative)
	_, err = HoldingPeriodYield(100, 100, maxCents+1)
	wantErr(t, "income > 2^53", err, ErrOutOfRange)
	_, err = HoldingPeriodYield(100, 100, -maxCents-1)
	wantErr(t, "income < −2^53", err, ErrOutOfRange)
	_, err = HoldingPeriodYield(maxCents+1, 100, 0)
	wantErr(t, "begin > 2^53", err, ErrOutOfRange)
}

func TestHoldingPeriodYield_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = HoldingPeriodYield(12874, 17653, 268) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkHoldingPeriodYield(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = HoldingPeriodYield(12874, 17653, 268)
	}
}
