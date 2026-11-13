package personal

import (
	"math"
	"testing"
)

func TestHourlyRateFromWeeklyPay_source(t *testing.T) {
	// 29 CFR 778.113(a): weekly salary ÷ hours it is meant to cover (hand computation).
	got, err := HourlyRateFromWeeklyPay(200000, 40)
	if err != nil || got != 5000 {
		t.Errorf("2,000/40 h: got %d, %v; want 5000", got, err)
	}
}

func TestHourlyRateFromWeeklyPay_cases(t *testing.T) {
	cases := []struct {
		name  string
		pay   int64
		hours float64
		want  int64
	}{
		{"zero pay", 0, 40, 0},
		{"1,000 over 37.5 h = 26.666… → 26.67", 100000, 37.5, 2667},
		{"half cent rounds to even (down)", 5, 2, 2}, // 2.5 → 2
		{"half cent rounds to even (up)", 7, 2, 4},   // 3.5 → 4
		{"all 168 hours", 168000, 168, 1000},
		{"fractional hour", 1000, 0.5, 2000},
	}
	for _, c := range cases {
		got, err := HourlyRateFromWeeklyPay(c.pay, c.hours)
		if err != nil || got != c.want {
			t.Errorf("%s: got %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

func TestHourlyRateFromWeeklyPay_errors(t *testing.T) {
	_, err := HourlyRateFromWeeklyPay(-1, 40)
	wantErr(t, "negative pay", err, ErrNegative)
	_, err = HourlyRateFromWeeklyPay(maxExactCents+1, 40)
	wantErr(t, "pay above 2^53", err, ErrOutOfRange)
	for _, h := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		_, err = HourlyRateFromWeeklyPay(100, h)
		wantErr(t, "non-finite hours", err, ErrNotFinite)
	}
	for _, h := range []float64{0, -1, 168.01} {
		_, err = HourlyRateFromWeeklyPay(100, h)
		wantErr(t, "hours outside (0, 168]", err, ErrDomain)
	}
	_, err = HourlyRateFromWeeklyPay(maxExactCents, 1e-300)
	wantErr(t, "tiny hours overflow", err, ErrOutOfRange)
}

func TestHourlyRateFromWeeklyPay_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkCents, sinkErr = HourlyRateFromWeeklyPay(200000, 40) }); a != 0 {
		t.Errorf("allocates %v times, want 0", a)
	}
}

func BenchmarkHourlyRateFromWeeklyPay(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, sinkErr = HourlyRateFromWeeklyPay(200000, 40)
	}
}
