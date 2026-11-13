package respiratory

import (
	"errors"
	"math"
	"testing"
)

// HAND-DERIVED from Lim 2003's criteria (PMID 12728155): confused, urea 8, RR 24, BP 110/70, age 70 → 3.
func TestCURB65Score_handExample(t *testing.T) {
	if got, err := CURB65Score(true, 8, 24, 110, 70, 70); err != nil || got != 3 {
		t.Errorf("got (%d, %v), want 3", got, err)
	}
	if got, err := CURB65Score(true, 7.1, 30, 89, 90, 65); err != nil || got != 5 {
		t.Errorf("all five: got (%d, %v), want 5", got, err)
	}
	if got, err := CURB65Score(false, 5, 16, 120, 80, 40); err != nil || got != 0 {
		t.Errorf("none: got (%d, %v), want 0", got, err)
	}
}

// Both sides of each criterion's edge: urea > 7, RR ≥ 30, SBP < 90 or DBP ≤ 60, age ≥ 65.
func TestCURB65Score_edges(t *testing.T) {
	base := func(urea, rr, sbp, dbp, age float64) int {
		n, err := CURB65Score(false, urea, rr, sbp, dbp, age)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
		return n
	}
	cases := []struct {
		name                    string
		urea, rr, sbp, dbp, age float64
		want                    int
	}{
		{"urea 7", 7, 20, 120, 80, 50, 0}, {"urea 7.01", 7.01, 20, 120, 80, 50, 1},
		{"RR 29", 5, 29, 120, 80, 50, 0}, {"RR 30", 5, 30, 120, 80, 50, 1},
		{"SBP 90", 5, 20, 90, 80, 50, 0}, {"SBP 89.9", 5, 20, 89.9, 80, 50, 1},
		{"DBP 61", 5, 20, 120, 61, 50, 0}, {"DBP 60", 5, 20, 120, 60, 50, 1},
		{"SBP and DBP low count once", 5, 20, 80, 50, 50, 1},
		{"age 64.9", 5, 20, 120, 80, 64.9, 0}, {"age 65", 5, 20, 120, 80, 65, 1},
	}
	for _, c := range cases {
		if got := base(c.urea, c.rr, c.sbp, c.dbp, c.age); got != c.want {
			t.Errorf("%s: got %d, want %d", c.name, got, c.want)
		}
	}
}

func TestCURB65Score_errors(t *testing.T) {
	tests := []struct {
		name                    string
		urea, rr, sbp, dbp, age float64
		wantErr                 error
	}{
		{"negative urea", -1, 20, 120, 80, 50, ErrOutOfRange},
		{"negative RR", 5, -1, 120, 80, 50, ErrOutOfRange},
		{"zero SBP", 5, 20, 0, 80, 50, ErrOutOfRange},
		{"zero DBP", 5, 20, 120, 0, 50, ErrOutOfRange},
		{"negative age", 5, 20, 120, 80, -1, ErrOutOfRange},
		{"NaN urea", math.NaN(), 20, 120, 80, 50, ErrNotFinite},
		{"+Inf RR", 5, math.Inf(1), 120, 80, 50, ErrNotFinite},
		{"-Inf SBP", 5, 20, math.Inf(-1), 80, 50, ErrNotFinite},
		{"NaN DBP", 5, 20, 120, math.NaN(), 50, ErrNotFinite},
		{"+Inf age", 5, 20, 120, 80, math.Inf(1), ErrNotFinite},
	}
	for _, tc := range tests {
		if _, err := CURB65Score(false, tc.urea, tc.rr, tc.sbp, tc.dbp, tc.age); !errors.Is(err, tc.wantErr) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.wantErr)
		}
	}
}

func TestCURB65Score_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = CURB65Score(true, 8, 24, 110, 70, 70) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkCURB65Score(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = CURB65Score(true, 8, 24, 110, 70, 70)
	}
}
