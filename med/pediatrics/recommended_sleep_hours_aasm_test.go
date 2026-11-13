package pediatrics

import (
	"errors"
	"math"
	"testing"
)

// Paruthi 2016: 4–12 months 12–16 h; 1–2 years 11–14; 3–5 years 10–13; 6–12 years 9–12; 13–18 years 8–10.
// Each age edge on both sides.
func TestRecommendedSleepHoursAASM_example(t *testing.T) {
	tests := []struct{ months, lo, hi float64 }{
		{4, 12, 16}, {11.99, 12, 16}, {12, 11, 14}, {35.99, 11, 14}, {36, 10, 13}, {71.99, 10, 13},
		{72, 9, 12}, {155.99, 9, 12}, {156, 8, 10}, {227.99, 8, 10},
	}
	for _, tc := range tests {
		lo, hi, err := RecommendedSleepHoursAASM(tc.months)
		if err != nil || lo != tc.lo || hi != tc.hi {
			t.Errorf("RecommendedSleepHoursAASM(%v) = (%v, %v, %v), want %v–%v", tc.months, lo, hi, err, tc.lo, tc.hi)
		}
	}
}

// parity: ledger pediatric_sleep_calculator, source cases 1–7 (min and max hours; the source's below_min and
// above_max are the caller's subtraction). Source case 1 (2 months → 14–17 h) is outside the AASM statement,
// which excludes infants under 4 months; the primary's answer is ErrOutOfRange.
func TestRecommendedSleepHoursAASM_parity(t *testing.T) {
	if _, _, err := RecommendedSleepHoursAASM(2); !errors.Is(err, ErrOutOfRange) { // source case 1: 14–17
		t.Errorf("2 months: error = %v, want ErrOutOfRange (source 14–17)", err)
	}
	tests := []struct{ months, lo, hi float64 }{
		{6, 12, 16},  // source case 2
		{24, 11, 14}, // source case 3
		{48, 10, 13}, // source case 4
		{100, 9, 12}, // source case 5
		{200, 8, 10}, // source case 6
		{12, 11, 14}, // source case 7
	}
	for _, tc := range tests {
		lo, hi, err := RecommendedSleepHoursAASM(tc.months)
		if err != nil || lo != tc.lo || hi != tc.hi {
			t.Errorf("RecommendedSleepHoursAASM(%v) = (%v, %v, %v), want %v–%v", tc.months, lo, hi, err, tc.lo, tc.hi)
		}
	}
}

func TestRecommendedSleepHoursAASM_errors(t *testing.T) {
	cases := []struct {
		m    float64
		want error
	}{
		{math.NaN(), ErrNotFinite}, {math.Inf(1), ErrNotFinite}, {math.Inf(-1), ErrNotFinite},
		{3.99, ErrOutOfRange}, {228, ErrOutOfRange}, {-1, ErrOutOfRange},
	}
	for _, c := range cases {
		if _, _, err := RecommendedSleepHoursAASM(c.m); !errors.Is(err, c.want) {
			t.Errorf("RecommendedSleepHoursAASM(%v) error = %v, want %v", c.m, err, c.want)
		}
	}
}

func TestRecommendedSleepHoursAASM_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _, _ = RecommendedSleepHoursAASM(48) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkRecommendedSleepHoursAASM(b *testing.B) {
	b.ReportAllocs()
	m := 48.0
	for b.Loop() {
		sinkFloat, _, _ = RecommendedSleepHoursAASM(m)
	}
}
