package series

import (
	"errors"
	"testing"
	"time"
)

// parity: ledger LocalDays, source case 1 (2026-11-01 New York is 25 h long).
func TestLocalDayWindow_parity(t *testing.T) {
	loc := newYork(t)
	from, to, err := LocalDayWindow(time.Date(2026, 11, 1, 12, 0, 0, 0, loc), loc, 1)
	if err != nil {
		t.Fatal(err)
	}
	if to.Sub(from) != 25*time.Hour {
		t.Errorf("DST end day length %v, want 25h", to.Sub(from))
	}
}

// parity: ledger LocalDays, source case 2 (7 days ending 2026-10-07 UTC).
func TestLocalDayWindow_cases(t *testing.T) {
	t0 := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	from, to, err := LocalDayWindow(t0.Add(6*24*time.Hour), time.UTC, 7)
	if err != nil || !from.Equal(t0) || !to.Equal(t0.Add(7*24*time.Hour)) {
		t.Errorf("7 days: [%v, %v) %v", from, to, err)
	}
	loc := newYork(t)
	// Spring forward: 2026-03-08 is 23 h.
	from, to, _ = LocalDayWindow(time.Date(2026, 3, 8, 23, 59, 0, 0, loc), loc, 1)
	if to.Sub(from) != 23*time.Hour {
		t.Errorf("DST start day length %v, want 23h", to.Sub(from))
	}
	// asOf is read in loc: 02:00 UTC on 11-01 is 10-31 22:00 in New York.
	from, _, _ = LocalDayWindow(time.Date(2026, 11, 1, 2, 0, 0, 0, time.UTC), loc, 1)
	if y, m, d := from.Date(); y != 2026 || m != 10 || d != 31 {
		t.Errorf("local date: %v", from)
	}
	// Across a month and year end.
	from, to, _ = LocalDayWindow(time.Date(2027, 1, 2, 5, 0, 0, 0, time.UTC), time.UTC, 3)
	if !from.Equal(time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)) || !to.Equal(time.Date(2027, 1, 3, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("year end: [%v, %v)", from, to)
	}
}

func TestLocalDayWindow_errors(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	if _, _, err := LocalDayWindow(now, nil, 7); !errors.Is(err, ErrNilLocation) {
		t.Errorf("nil loc: %v", err)
	}
	for _, d := range []int{0, -1} {
		if _, _, err := LocalDayWindow(now, time.UTC, d); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("days %d: %v", d, err)
		}
	}
}

func TestLocalDayWindow_noAlloc(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	if a := testing.AllocsPerRun(100, func() { _, _, _ = LocalDayWindow(now, time.UTC, 7) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkLocalDayWindow(b *testing.B) {
	loc := newYork(b)
	now := time.Date(2026, 11, 1, 12, 0, 0, 0, loc)
	b.ReportAllocs()
	for b.Loop() {
		_, _, _ = LocalDayWindow(now, loc, 7)
	}
}
