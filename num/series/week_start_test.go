package series

import (
	"errors"
	"testing"
	"time"
)

func TestWeekStart_cases(t *testing.T) {
	loc := newYork(t)
	sun := time.Date(2026, 10, 4, 15, 0, 0, 0, loc) // a Sunday
	cases := []struct {
		name  string
		at    time.Time
		start time.Weekday
		want  time.Time
	}{
		// Doc example: a Sunday in a Monday week belongs to the week of the Monday before.
		{"sunday, monday weeks", sun, time.Monday, time.Date(2026, 9, 28, 0, 0, 0, 0, loc)},
		{"sunday, sunday weeks", sun, time.Sunday, time.Date(2026, 10, 4, 0, 0, 0, 0, loc)},
		{"saturday weeks", sun, time.Saturday, time.Date(2026, 10, 3, 0, 0, 0, 0, loc)},
		{"local date, not UTC", time.Date(2026, 10, 5, 2, 0, 0, 0, time.UTC), time.Monday, time.Date(2026, 9, 28, 0, 0, 0, 0, loc)},
		{"DST week", time.Date(2026, 11, 7, 23, 0, 0, 0, loc), time.Sunday, time.Date(2026, 11, 1, 0, 0, 0, 0, loc)},
		{"year boundary", time.Date(2027, 1, 1, 12, 0, 0, 0, loc), time.Monday, time.Date(2026, 12, 28, 0, 0, 0, 0, loc)},
	}
	for _, c := range cases {
		got, err := WeekStart(c.at, loc, c.start)
		if err != nil || !got.Equal(c.want) {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestWeekStart_errors(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	if _, err := WeekStart(now, nil, time.Sunday); !errors.Is(err, ErrNilLocation) {
		t.Errorf("nil loc: %v", err)
	}
	for _, d := range []time.Weekday{-1, 7} {
		if _, err := WeekStart(now, time.UTC, d); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("weekday %d: %v", d, err)
		}
	}
}

func TestWeekStart_noAlloc(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	if a := testing.AllocsPerRun(100, func() { _, _ = WeekStart(now, time.UTC, time.Monday) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkWeekStart(b *testing.B) {
	loc := newYork(b)
	now := time.Date(2026, 11, 4, 12, 0, 0, 0, loc)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = WeekStart(now, loc, time.Monday)
	}
}
