package series

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"testing"
	"time"
)

func checkWeeks(t *testing.T, name string, got []WeekSum, want []WeekSum) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s: got %v, want %v", name, got, want)
	}
	for i := range want {
		if !got[i].Start.Equal(want[i].Start) || got[i].Sum != want[i].Sum || got[i].Count != want[i].Count {
			t.Errorf("%s: week %d = %v, want %v", name, i, got[i], want[i])
		}
	}
}

// parity: ledger WeeklyBuckets, source cases 1–3 (Sunday weeks, Monday weeks, a DST week).
func TestWeeklySums_parity(t *testing.T) {
	loc := newYork(t)
	at := func(m time.Month, d, h, min int) time.Time { return time.Date(2026, m, d, h, min, 0, 0, loc) }
	day := func(m time.Month, d int) time.Time { return at(m, d, 0, 0) }
	pts := []Point{{at(10, 5, 10, 0), 2}, {at(10, 10, 23, 30), 3}, {at(10, 11, 0, 30), 1}}
	from, to := day(10, 4), day(10, 25)
	got, err := WeeklySums(pts, from, to, loc, time.Sunday)
	if err != nil {
		t.Fatal(err)
	}
	checkWeeks(t, "sunday", got, []WeekSum{{day(10, 4), 5, 2}, {day(10, 11), 1, 1}, {day(10, 18), 0, 0}})
	got, _ = WeeklySums(pts, from, to, loc, time.Monday)
	checkWeeks(t, "monday", got, []WeekSum{{day(9, 28), 0, 0}, {day(10, 5), 6, 3}, {day(10, 12), 0, 0}, {day(10, 19), 0, 0}})
	got, _ = WeeklySums([]Point{{at(11, 7, 23, 0), 1}}, day(11, 1), day(11, 8), loc, time.Sunday)
	checkWeeks(t, "dst", got, []WeekSum{{day(11, 1), 1, 1}})
	if got[0].Start.AddDate(0, 0, 7).Sub(got[0].Start) != 169*time.Hour {
		t.Errorf("DST week is not 169 h")
	}
	if _, err := WeeklySums(pts, from, to, nil, time.Sunday); !errors.Is(err, ErrNilLocation) {
		t.Errorf("nil loc: %v", err)
	}
}

func TestWeeklySums_cases(t *testing.T) {
	t0 := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC) // a Monday
	// Points outside [from, to) are not summed; the window end is open.
	pts := []Point{{t0.Add(-time.Hour), 100}, {t0, 1}, {t0.Add(7 * 24 * time.Hour), 100}}
	got, err := WeeklySums(pts, t0, t0.Add(7*24*time.Hour), time.UTC, time.Monday)
	if err != nil {
		t.Fatal(err)
	}
	checkWeeks(t, "window", got, []WeekSum{{t0, 1, 1}})
	// No points: every week is present and empty.
	got, _ = WeeklySums(nil, t0, t0.Add(15*24*time.Hour), time.UTC, time.Monday)
	if len(got) != 3 || got[2].Count != 0 {
		t.Errorf("empty: %v", got)
	}
}

func TestWeeklySums_errors(t *testing.T) {
	t0 := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	t1 := t0.Add(24 * time.Hour)
	cases := []struct {
		name     string
		pts      []Point
		from, to time.Time
		start    time.Weekday
		want     error
	}{
		{"to equals from", nil, t0, t0, time.Monday, ErrOutOfRange},
		{"to before from", nil, t1, t0, time.Monday, ErrOutOfRange},
		{"bad weekday", nil, t0, t1, 9, ErrOutOfRange},
		{"unsorted", []Point{{t1, 1}, {t0, 1}}, t0, t1, time.Monday, ErrUnsorted},
		{"NaN", []Point{{t0, math.NaN()}}, t0, t1, time.Monday, ErrNotFinite},
		{"overflow", []Point{{t0, math.MaxFloat64}, {t0, math.MaxFloat64}}, t0, t1, time.Monday, ErrOverflow},
	}
	for _, c := range cases {
		if _, err := WeeklySums(c.pts, c.from, c.to, time.UTC, c.start); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// Oracle: each point inside the window is counted in the week WeekStart gives it (fixed seed).
func TestWeeklySums_oracle(t *testing.T) {
	loc := newYork(t)
	r := rand.New(rand.NewPCG(1, 2))
	t0 := time.Date(2026, 2, 20, 0, 0, 0, 0, loc)
	for trial := 0; trial < 200; trial++ {
		pts := make([]Point, r.IntN(60))
		at := t0
		for i := range pts {
			at = at.Add(time.Duration(r.IntN(72)) * time.Hour)
			pts[i] = Point{At: at, Value: float64(r.IntN(10))}
		}
		from := t0.Add(time.Duration(r.IntN(30*24)) * time.Hour)
		to := from.Add(time.Duration(1+r.IntN(60*24)) * time.Hour)
		ws := time.Weekday(r.IntN(7))
		got, err := WeeklySums(pts, from, to, loc, ws)
		if err != nil {
			t.Fatal(err)
		}
		want := map[int64]WeekSum{}
		for _, p := range pts {
			if p.At.Before(from) || !p.At.Before(to) {
				continue
			}
			s, _ := WeekStart(p.At, loc, ws)
			w := want[s.Unix()]
			w.Sum += p.Value
			w.Count++
			want[s.Unix()] = w
		}
		total := 0
		for _, g := range got {
			w := want[g.Start.Unix()]
			if g.Sum != w.Sum || g.Count != w.Count {
				t.Fatalf("trial %d: week %v = %v, want %v", trial, g.Start, g, w)
			}
			total += g.Count
		}
		n := 0
		for _, w := range want {
			n += w.Count
		}
		if total != n {
			t.Fatalf("trial %d: counted %d of %d", trial, total, n)
		}
	}
}

func BenchmarkWeeklySums(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		pts := benchPoints(n)
		from, to := pts[0].At, pts[n-1].At.Add(time.Second)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := WeeklySums(pts, from, to, time.UTC, time.Monday); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
