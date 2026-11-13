package personal

import (
	"errors"
	"testing"
)

var (
	sinkCents int64
	sinkErr   error
)

func wantErr(t *testing.T, name string, err, sentinel error) {
	t.Helper()
	if !errors.Is(err, sentinel) {
		t.Errorf("%s: err = %v, want %v", name, err, sentinel)
	}
}

func TestWeeklyPayEquivalent_source(t *testing.T) {
	// 29 CFR 778.113(b): monthly × 12 ÷ 52 and semimonthly × 24 ÷ 52 (hand computations).
	cases := []struct {
		name    string
		pay     int64
		periods int
		want    int64
	}{
		{"monthly 8,666.67 → 2,000.00 (866667×12/52 = 200000.08)", 866667, 12, 200000},
		{"semimonthly 4,333.33 → 2,000.00 (433333×24/52 = 199999.85)", 433333, 24, 200000},
		{"monthly 5,200.00 → exactly 1,200.00", 520000, 12, 120000},
	}
	for _, c := range cases {
		got, err := WeeklyPayEquivalent(c.pay, c.periods)
		if err != nil || got != c.want {
			t.Errorf("%s: got %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

func TestWeeklyPayEquivalent_cases(t *testing.T) {
	cases := []struct {
		name    string
		pay     int64
		periods int
		want    int64
	}{
		{"zero pay", 0, 12, 0},
		{"annual 104,000 → 2,000", 10400000, 1, 200000},
		{"biweekly 4,000 → 2,000", 400000, 26, 200000},
		{"weekly is unchanged", 123457, 52, 123457},
		{"half cent rounds to even (down)", 1, 26, 0}, // 26/52 = 0.5 → 0
		{"half cent rounds to even (up)", 3, 26, 2},   // 78/52 = 1.5 → 2
		{"just above half rounds up", 27, 1, 1},       // 27/52 = 0.519 → 1
		{"daily 366 periods", 100, 366, 704},          // 36600/52 = 703.85 → 704
		{"largest exact pay (want set below)", maxExactCents, 366, 0},
	}
	// 2⁵³ × 366 / 52 rounded half-even, by exact integer arithmetic (the product is below 2⁶³).
	big := int64(maxExactCents) * 366
	q, r := big/52, big%52
	if 2*r > 52 || (2*r == 52 && q%2 == 1) {
		q++
	}
	cases[len(cases)-1].want = q
	for _, c := range cases {
		got, err := WeeklyPayEquivalent(c.pay, c.periods)
		if err != nil || got != c.want {
			t.Errorf("%s: got %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

func TestWeeklyPayEquivalent_errors(t *testing.T) {
	_, err := WeeklyPayEquivalent(-1, 12)
	wantErr(t, "negative pay", err, ErrNegative)
	_, err = WeeklyPayEquivalent(maxExactCents+1, 12)
	wantErr(t, "pay above 2^53", err, ErrOutOfRange)
	_, err = WeeklyPayEquivalent(100, 0)
	wantErr(t, "zero periods", err, ErrOutOfRange)
	_, err = WeeklyPayEquivalent(100, 367)
	wantErr(t, "too many periods", err, ErrOutOfRange)
}

func TestWeeklyPayEquivalent_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkCents, sinkErr = WeeklyPayEquivalent(866667, 12) }); a != 0 {
		t.Errorf("allocates %v times, want 0", a)
	}
}

func BenchmarkWeeklyPayEquivalent(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, sinkErr = WeeklyPayEquivalent(866667, 12)
	}
}
