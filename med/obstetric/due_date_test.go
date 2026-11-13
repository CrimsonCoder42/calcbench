package obstetric

import (
	"errors"
	"testing"
	"time"
)

func date(y int, m time.Month, d int) time.Time { return time.Date(y, m, d, 0, 0, 0, 0, time.UTC) }

// Worked examples, HAND-DERIVED from ACOG CO 700 (EDD = LMP + 280 days) on the Gregorian calendar:
// 2026-01-01 + 280 d = 2026-10-08 (2026 not leap); 2024-02-29 + 280 d = 2024-12-05; 2023-12-31 + 280 d = 2024-10-06
// (crosses the 2024 leap day); 1999-04-10 + 280 d = 2000-01-15 (crosses a year end).
func TestDueDateFromLMP_example(t *testing.T) {
	tests := []struct{ lmp, want time.Time }{
		{date(2026, 1, 1), date(2026, 10, 8)},
		{date(2024, 2, 29), date(2024, 12, 5)},
		{date(2023, 12, 31), date(2024, 10, 6)},
		{date(1999, 4, 10), date(2000, 1, 15)},
	}
	for _, tc := range tests {
		if got := DueDateFromLMP(tc.lmp); !got.Equal(tc.want) {
			t.Errorf("DueDateFromLMP(%s) = %s, want %s", tc.lmp.Format(time.DateOnly), got.Format(time.DateOnly), tc.want.Format(time.DateOnly))
		}
	}
}

// The time of day and the location are kept as the caller's calendar date: 23:30 in UTC−10 on 2026-01-01 is
// still LMP 2026-01-01, and the EDD is midnight 2026-10-08 in that location.
func TestDueDateFromLMP_location(t *testing.T) {
	loc := time.FixedZone("UTC-10", -10*3600)
	got := DueDateFromLMP(time.Date(2026, 1, 1, 23, 30, 0, 0, loc))
	if want := time.Date(2026, 10, 8, 0, 0, 0, 0, loc); !got.Equal(want) || got.Location() != loc {
		t.Errorf("got %v, want %v", got, want)
	}
}

// Gestational age, HAND-DERIVED: LMP 2026-01-01 → on 2026-01-01 0 w 0 d; on 2026-01-07 0 w 6 d; on 2026-01-08 1 w 0 d;
// on 2026-03-15 (73 days) 10 w 3 d; on 2026-10-08 (280 days) 40 w 0 d.
func TestGestationalAgeFromLMP_example(t *testing.T) {
	lmp := date(2026, 1, 1)
	tests := []struct {
		on          time.Time
		weeks, days int
	}{
		{date(2026, 1, 1), 0, 0},
		{date(2026, 1, 7), 0, 6},
		{date(2026, 1, 8), 1, 0},
		{date(2026, 3, 15), 10, 3},
		{date(2026, 10, 8), 40, 0},
	}
	for _, tc := range tests {
		w, d, err := GestationalAgeFromLMP(lmp, tc.on)
		if err != nil || w != tc.weeks || d != tc.days {
			t.Errorf("on %s: (%d, %d, %v), want (%d, %d)", tc.on.Format(time.DateOnly), w, d, err, tc.weeks, tc.days)
		}
	}
}

// Day 280 of GestationalAgeFromLMP is the date DueDateFromLMP returns, for every LMP in 2023–2025.
func TestGestationalAgeFromLMP_agreesWithDueDate(t *testing.T) {
	for lmp := date(2023, 1, 1); lmp.Before(date(2026, 1, 1)); lmp = lmp.AddDate(0, 0, 1) {
		w, d, err := GestationalAgeFromLMP(lmp, DueDateFromLMP(lmp))
		if err != nil || w != 40 || d != 0 {
			t.Fatalf("lmp %s: (%d, %d, %v), want (40, 0)", lmp.Format(time.DateOnly), w, d, err)
		}
	}
}

// Dates are compared as each shows in its own location: 2026-01-02 01:00 in UTC+14 is still the date 2026-01-02.
func TestGestationalAgeFromLMP_location(t *testing.T) {
	on := time.Date(2026, 1, 2, 1, 0, 0, 0, time.FixedZone("UTC+14", 14*3600))
	if w, d, err := GestationalAgeFromLMP(date(2026, 1, 1), on); err != nil || w != 0 || d != 1 {
		t.Errorf("got (%d, %d, %v), want (0, 1)", w, d, err)
	}
}

func TestGestationalAgeFromLMP_errors(t *testing.T) {
	if _, _, err := GestationalAgeFromLMP(date(2026, 1, 2), date(2026, 1, 1)); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("on before lmp: err = %v, want ErrOutOfRange", err)
	}
}

func TestDueDate_noAllocation(t *testing.T) {
	lmp, on := date(2026, 1, 1), date(2026, 3, 15)
	if n := testing.AllocsPerRun(100, func() {
		_ = DueDateFromLMP(lmp)
		_, _, _ = GestationalAgeFromLMP(lmp, on)
	}); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkDueDateFromLMP(b *testing.B) {
	b.ReportAllocs()
	lmp := date(2026, 1, 1)
	var s time.Time
	for b.Loop() {
		s = DueDateFromLMP(lmp)
	}
	_ = s
}

func BenchmarkGestationalAgeFromLMP(b *testing.B) {
	b.ReportAllocs()
	lmp, on := date(2026, 1, 1), date(2026, 3, 15)
	var w int
	for b.Loop() {
		w, _, _ = GestationalAgeFromLMP(lmp, on)
	}
	_ = w
}
