package cron

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"
	"testing"
	"time"
	_ "time/tzdata" // zone rules embedded so America/New_York and Australia/Lord_Howe load on any host
)

func mustLoad(t testing.TB, name string) *time.Location {
	t.Helper()
	loc, err := time.LoadLocation(name)
	if err != nil {
		t.Skipf("zone %s unavailable: %v", name, err)
	}
	return loc
}

func utc(y int, mo time.Month, d, h, mi int) time.Time {
	return time.Date(y, mo, d, h, mi, 0, 0, time.UTC)
}

// checkNext asserts c.Next(after) is the instant want, in after's location, with the UTC offset wantOff.
func checkNext(t *testing.T, spec string, after, want time.Time, wantOff time.Duration) {
	t.Helper()
	got, err := mustParse(t, spec).Next(after)
	if err != nil {
		t.Fatalf("%q.Next(%v): %v", spec, after, err)
	}
	if !got.Equal(want) {
		t.Fatalf("%q.Next(%v) = %v, want %v", spec, after, got, want.In(after.Location()))
	}
	if got.Location().String() != after.Location().String() {
		t.Fatalf("%q.Next(%v) is in %v, want %v", spec, after, got.Location(), after.Location())
	}
	if _, off := got.Zone(); time.Duration(off)*time.Second != wantOff {
		t.Fatalf("%q.Next(%v) = %v has offset %v, want %v", spec, after, got, time.Duration(off)*time.Second, wantOff)
	}
	if got.Second() != 0 || got.Nanosecond() != 0 {
		t.Fatalf("%q.Next(%v) = %v has nonzero seconds", spec, after, got)
	}
}

const (
	est = -5 * time.Hour
	edt = -4 * time.Hour
)

// Worked examples of the BRIEF (doc/toolbench_requests_2026-11-13/BRIEF.md, B3). The 2026 US transitions are
// the second Sunday of March (8th) and the first Sunday of November (1st), 2:00 local (tz database, America/New_York).
func TestNext_workedExamples(t *testing.T) {
	ny := mustLoad(t, "America/New_York")
	// 02:30 does not exist on 2026-03-08: it fires at the transition instant, 03:00 EDT = 07:00Z.
	checkNext(t, "30 2 * * *", time.Date(2026, 3, 7, 3, 0, 0, 0, ny), utc(2026, 3, 8, 7, 0), edt)
	// 2026-11-13 is a Friday: the next weekday 10:00 is Monday the 16th.
	checkNext(t, "0 10 * * 1-5", utc(2026, 11, 13, 10, 0), utc(2026, 11, 16, 10, 0), 0)
}

func TestNext_springForwardGap(t *testing.T) {
	ny := mustLoad(t, "America/New_York")
	gapEnd := utc(2026, 3, 8, 7, 0) // 03:00 EDT, the first instant after 02:00–02:59 EST, which do not exist
	// After the transition the gap time fires again only the next day.
	checkNext(t, "30 2 * * *", gapEnd.In(ny), utc(2026, 3, 9, 6, 30), edt)
	// Sixty matching wall times in the gap fire once, at 03:00 EDT.
	checkNext(t, "* 2 * * *", utc(2026, 3, 8, 6, 59).In(ny), gapEnd, edt) // after 01:59 EST
	checkNext(t, "* 2 * * *", gapEnd.In(ny), utc(2026, 3, 9, 6, 0), edt)  // after 03:00 EDT → next day 02:00
	// A gap time and 03:00 itself both fire at the transition instant; the next firing is 03:01.
	checkNext(t, "* 2,3 * * *", utc(2026, 3, 8, 6, 59).In(ny), gapEnd, edt)
	checkNext(t, "* 2,3 * * *", gapEnd.In(ny), utc(2026, 3, 8, 7, 1), edt)
	// After carrying fractional seconds just before the gap.
	checkNext(t, "30 2 * * *", time.Date(2026, 3, 8, 1, 59, 59, 999_999_999, ny), gapEnd, edt)
	// 01:59 EST fires normally just before the gap.
	checkNext(t, "59 1 * * *", time.Date(2026, 3, 8, 0, 0, 0, 0, ny), utc(2026, 3, 8, 6, 59), est)
}

func TestNext_fallBackOverlap(t *testing.T) {
	ny := mustLoad(t, "America/New_York")
	// 2026-11-01: 01:00–01:59 occurs at -04:00 (05:00Z–05:59Z) and again at -05:00 (06:00Z–06:59Z).
	checkNext(t, "30 1 * * *", time.Date(2026, 11, 1, 0, 0, 0, 0, ny), utc(2026, 11, 1, 5, 30), edt) // first occurrence
	checkNext(t, "30 1 * * *", utc(2026, 11, 1, 5, 30).In(ny), utc(2026, 11, 2, 6, 30), est)         // second never fires
	checkNext(t, "30 1 * * *", utc(2026, 11, 1, 6, 10).In(ny), utc(2026, 11, 2, 6, 30), est)         // after 01:10 EST
	// Every minute: the second pass 01:00–01:59 EST never fires; 02:00 EST follows 01:59 EDT.
	checkNext(t, "* * * * *", utc(2026, 11, 1, 5, 58).In(ny), utc(2026, 11, 1, 5, 59), edt)
	checkNext(t, "* * * * *", utc(2026, 11, 1, 5, 59).In(ny), utc(2026, 11, 1, 7, 0), est)
	checkNext(t, "* * * * *", utc(2026, 11, 1, 6, 0).In(ny), utc(2026, 11, 1, 7, 0), est)
	checkNext(t, "* * * * *", utc(2026, 11, 1, 6, 59).In(ny), utc(2026, 11, 1, 7, 0), est)
	checkNext(t, "0 2 * * *", utc(2026, 11, 1, 5, 0).In(ny), utc(2026, 11, 1, 7, 0), est)
}

// Australia/Lord_Howe shifts by 30 minutes: +10:30 standard, +11:00 summer; summer time starts the first
// Sunday of October at 02:00 and ends the first Sunday of April at 02:00 (tz database, rule LH from 2008).
func TestNext_halfHourDST(t *testing.T) {
	lh := mustLoad(t, "Australia/Lord_Howe")
	std, dst := 10*time.Hour+30*time.Minute, 11*time.Hour
	// 2026-10-04 02:15 does not exist: it fires at 02:30 +11:00 = 2026-10-03T15:30Z.
	checkNext(t, "15 2 * * *", time.Date(2026, 10, 3, 12, 0, 0, 0, lh), utc(2026, 10, 3, 15, 30), dst)
	// 2026-04-05 01:45 occurs at +11:00 (14:45Z) and +10:30 (15:15Z): only the first fires.
	checkNext(t, "45 1 * * *", time.Date(2026, 4, 5, 0, 0, 0, 0, lh), utc(2026, 4, 4, 14, 45), dst)
	checkNext(t, "45 1 * * *", utc(2026, 4, 4, 14, 45).In(lh), utc(2026, 4, 5, 15, 15), std)
}

func TestNext_calendar(t *testing.T) {
	cases := []struct {
		name        string
		spec        string
		after, want time.Time
	}{
		{"31st skips 30-day and short months", "0 0 31 * *", utc(2026, 1, 31, 0, 0), utc(2026, 3, 31, 0, 0)},
		{"31st skips April", "0 0 31 * *", utc(2026, 3, 31, 0, 0), utc(2026, 5, 31, 0, 0)},
		{"Feb 29 across 2100", "0 0 29 2 *", utc(2097, 1, 1, 0, 0), utc(2104, 2, 29, 0, 0)},
		{"Feb 29 window edge", "0 0 29 2 *", utc(2096, 2, 29, 0, 0), utc(2104, 2, 29, 0, 0)},
		{"Feb 29 ordinary", "0 0 29 2 *", utc(2104, 2, 29, 0, 0), utc(2108, 2, 29, 0, 0)},
		{"13th or Friday: a Friday", "0 0 13 * 5", utc(2026, 11, 13, 0, 0), utc(2026, 11, 20, 0, 0)},
		{"13th or Friday: a Sunday the 13th", "0 0 13 * 5", utc(2026, 12, 11, 0, 0), utc(2026, 12, 13, 0, 0)},
		{"13th or Friday: Friday after the 13th", "0 0 13 * 5", utc(2026, 12, 13, 0, 0), utc(2026, 12, 18, 0, 0)},
		{"Feb 31 or any weekday", "0 0 31 2 0-6", utc(2026, 1, 15, 0, 0), utc(2026, 2, 1, 0, 0)},
		{"month only", "0 0 * 2 *", utc(2026, 3, 1, 0, 0), utc(2027, 2, 1, 0, 0)},
		{"month only next day", "0 0 * 2 *", utc(2027, 2, 1, 0, 0), utc(2027, 2, 2, 0, 0)},
		{"dow only Sunday", "0 12 * * 0", utc(2026, 11, 13, 0, 0), utc(2026, 11, 15, 12, 0)},
		{"after is a match", "0 10 * * *", utc(2026, 11, 13, 10, 0), utc(2026, 11, 14, 10, 0)},
		{"after carries seconds", "* * * * *", time.Date(2026, 11, 13, 10, 0, 30, 500_000_000, time.UTC), utc(2026, 11, 13, 10, 1)},
		{"after one nanosecond past a match", "0 * * * *", time.Date(2026, 11, 13, 10, 0, 0, 1, time.UTC), utc(2026, 11, 13, 11, 0)},
		{"after one nanosecond before a match", "0 * * * *", time.Date(2026, 11, 13, 9, 59, 59, 999_999_999, time.UTC), utc(2026, 11, 13, 10, 0)},
		{"year-end carry", "0 0 1 1 *", time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC), utc(2027, 1, 1, 0, 0)},
		{"last minute of the year", "59 23 31 12 *", utc(2026, 12, 31, 23, 59), utc(2027, 12, 31, 23, 59)},
		{"hour carry", "0 * * * *", utc(2026, 11, 13, 23, 0), utc(2026, 11, 14, 0, 0)},
		{"minute list", "5,50 * * * *", utc(2026, 11, 13, 10, 5), utc(2026, 11, 13, 10, 50)},
		{"year 1", "* * * * *", time.Time{}, utc(1, 1, 1, 0, 1)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) { checkNext(t, tc.spec, tc.after, tc.want, 0) })
	}
}

func TestNext_fixedZones(t *testing.T) {
	ist := time.FixedZone("IST", 5*3600+1800)
	npt := time.FixedZone("NPT", 5*3600+45*60)
	nfz := time.FixedZone("M0330", -(3*3600 + 30*60))
	checkNext(t, "0 9 * * *", time.Date(2026, 11, 13, 9, 0, 0, 0, ist), time.Date(2026, 11, 14, 9, 0, 0, 0, ist), 5*time.Hour+30*time.Minute)
	checkNext(t, "0 0 * * *", time.Date(2026, 11, 13, 23, 59, 59, 0, npt), time.Date(2026, 11, 14, 0, 0, 0, 0, npt), 5*time.Hour+45*time.Minute)
	checkNext(t, "30 23 31 12 *", time.Date(2026, 12, 31, 23, 30, 0, 0, nfz), time.Date(2027, 12, 31, 23, 30, 0, 0, nfz), -(3*time.Hour + 30*time.Minute))
	// The same instant read in two locations: the wall clock of after's location decides.
	at := utc(2026, 11, 13, 0, 0)
	checkNext(t, "0 9 * * *", at, utc(2026, 11, 13, 9, 0), 0)
	checkNext(t, "0 9 * * *", at.In(ist), utc(2026, 11, 13, 3, 30), 5*time.Hour+30*time.Minute)
}

func TestNext_neverFires(t *testing.T) {
	if _, err := (CronSpec{}).Next(utc(2026, 11, 13, 0, 0)); !errors.Is(err, ErrNeverFires) {
		t.Fatalf("zero CronSpec Next error = %v, want ErrNeverFires", err)
	}
	// Every field restricted to nothing but the day fields: a spec ParseCron cannot return, built here to
	// prove the search window ends.
	c := CronSpec{minute: 1, hour: 1, month: 1 << 2, dom: 1 << 30, dowStar: true}
	if _, err := c.Next(utc(2026, 11, 13, 0, 0)); !errors.Is(err, ErrNeverFires) {
		t.Fatalf("Feb 30 Next error = %v, want ErrNeverFires", err)
	}
}

func TestNext_noAllocation(t *testing.T) {
	ny := mustLoad(t, "America/New_York")
	c := mustParse(t, "30 2 * * 1-5")
	for _, after := range []time.Time{utc(2026, 11, 13, 0, 0), time.Date(2026, 3, 6, 3, 0, 0, 0, ny)} {
		if n := testing.AllocsPerRun(100, func() {
			if _, err := c.Next(after); err != nil {
				t.Fatal(err)
			}
		}); n != 0 {
			t.Fatalf("Next(%v) allocates %v times per call, want 0", after, n)
		}
	}
}

// oracleSpec is a randomly generated specification: its text and, independently, the allowed values.
type oracleSpec struct {
	text  string
	allow [5][64]bool
	star  [5]bool
}

// matches is the definition, written out: month, hour and minute restrict; the day rule of the package doc.
func (o *oracleSpec) matches(t time.Time) bool {
	if !o.allow[0][t.Minute()] || !o.allow[1][t.Hour()] || !o.allow[3][int(t.Month())] {
		return false
	}
	dom, dow := o.allow[2][t.Day()], o.allow[4][int(t.Weekday())]
	switch {
	case o.star[2] && o.star[4]:
		return true
	case o.star[2]:
		return dow
	case o.star[4]:
		return dom
	default:
		return dom || dow
	}
}

func randomSpec(r *rand.Rand) *oracleSpec {
	o := &oracleSpec{}
	var parts [5]string
	for k, f := range fields {
		if r.IntN(100) < 40 {
			parts[k] = "*"
			o.star[k] = true
			for v := f.lo; v <= f.hi; v++ {
				o.allow[k][v] = true
			}
			continue
		}
		var elems []string
		for range 1 + r.IntN(3) {
			a := f.lo + r.IntN(f.hi-f.lo+1)
			b := a
			if r.IntN(2) == 0 {
				b = a + r.IntN(f.hi-a+1)
				elems = append(elems, fmt.Sprintf("%d-%d", a, b))
			} else {
				elems = append(elems, fmt.Sprintf("%d", a))
			}
			for v := a; v <= b; v++ {
				o.allow[k][v] = true
			}
		}
		parts[k] = strings.Join(elems, ",")
	}
	o.text = strings.Join(parts[:], " ")
	return o
}

// TestNext_oracle compares Next in UTC with a minute-by-minute search over at most two years after a random
// instant, for 300 random specifications (fixed seed). UTC only: no daylight saving.
func TestNext_oracle(t *testing.T) {
	const window = 2 * 366 * 24 * time.Hour
	r := rand.New(rand.NewPCG(1, 2))
	lo, hi := utc(2000, 1, 1, 0, 0).Unix(), utc(2031, 1, 1, 0, 0).Unix()
	for i := range 300 {
		o := randomSpec(r)
		after := time.Unix(lo+r.Int64N(hi-lo), r.Int64N(1e9)).UTC()
		c, err := ParseCron(o.text)
		if errors.Is(err, ErrNeverFires) {
			continue // proved by TestParseCron_neverFires; the window here could not show it
		}
		if err != nil {
			t.Fatalf("#%d ParseCron(%q): %v", i, o.text, err)
		}
		got, err := c.Next(after)
		if err != nil {
			t.Fatalf("#%d %q.Next(%v): %v", i, o.text, after, err)
		}
		end := after.Add(window)
		want := time.Time{}
		for m := after.Truncate(time.Minute).Add(time.Minute); !m.After(end); m = m.Add(time.Minute) {
			if o.matches(m) {
				want = m
				break
			}
		}
		switch {
		case want.IsZero() && !got.After(end):
			t.Fatalf("#%d %q.Next(%v) = %v, oracle finds nothing before %v", i, o.text, after, got, end)
		case !want.IsZero() && !got.Equal(want):
			t.Fatalf("#%d %q.Next(%v) = %v, oracle %v", i, o.text, after, got, want)
		}
	}
}

// FuzzParseCron: no input panics; an accepted specification fires strictly after t, at a minute that matches
// its fields, with no matching minute in between (checked for up to 3000 minutes).
func FuzzParseCron(f *testing.F) {
	for _, s := range []string{"* * * * *", "0 10 * * 1-5", "0 0 29 2 *", "0 0 13 * 5", "59 23 31 12 *",
		"1,,2 * * * *", "*/5 * * * *", "0 0 31 2 *", "5-3 * * * *", "0 0 31 2 0-6"} {
		f.Add(s, int64(1_789_000_000))
	}
	f.Fuzz(func(t *testing.T, spec string, sec int64) {
		c, err := ParseCron(spec)
		if err != nil {
			if !errors.Is(err, ErrSyntax) && !errors.Is(err, ErrRange) && !errors.Is(err, ErrNeverFires) {
				t.Fatalf("ParseCron(%q) unexpected error %v", spec, err)
			}
			return
		}
		if sec < 0 {
			sec = -(sec + 1)
		}
		after := time.Unix(946_684_800+sec%2_000_000_000, 0).UTC() // 2000 … 2063
		got, err := c.Next(after)
		if err != nil {
			t.Fatalf("%q.Next(%v): %v", spec, after, err)
		}
		if !got.After(after) || got.Second() != 0 || got.Nanosecond() != 0 || !c.matchesUTC(got) {
			t.Fatalf("%q.Next(%v) = %v: not a later matching minute", spec, after, got)
		}
		m := after.Truncate(time.Minute).Add(time.Minute)
		for i := 0; i < 3000 && m.Before(got); i, m = i+1, m.Add(time.Minute) {
			if c.matchesUTC(m) {
				t.Fatalf("%q.Next(%v) = %v, but %v matches earlier", spec, after, got, m)
			}
		}
	})
}

// matchesUTC reports whether the UTC wall clock of t matches c's fields (the day rule of the package doc).
func (c CronSpec) matchesUTC(t time.Time) bool {
	t = t.UTC()
	if c.minute&(1<<t.Minute()) == 0 || c.hour&(1<<t.Hour()) == 0 || c.month&(1<<int(t.Month())) == 0 {
		return false
	}
	dom, dow := c.dom&(1<<t.Day()) != 0, c.dow&(1<<int(t.Weekday())) != 0
	switch {
	case c.domStar && c.dowStar:
		return true
	case c.domStar:
		return dow
	case c.dowStar:
		return dom
	default:
		return dom || dow
	}
}

var sinkTime time.Time

func BenchmarkNext(b *testing.B) {
	ny := mustLoad(b, "America/New_York")
	for _, bc := range []struct {
		name, spec string
		after      time.Time
	}{
		{"UTC/every-minute", "* * * * *", utc(2026, 11, 13, 10, 0)},
		{"UTC/weekdays", "0 10 * * 1-5", utc(2026, 11, 13, 10, 0)},
		{"UTC/feb29-2104", "0 0 29 2 *", utc(2097, 1, 1, 0, 0)},
		{"NewYork/gap", "30 2 * * *", time.Date(2026, 3, 7, 3, 0, 0, 0, ny)},
		{"NewYork/overlap", "* * * * *", utc(2026, 11, 1, 6, 0).In(ny)},
	} {
		c := mustParse(b, bc.spec)
		b.Run(bc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				t, err := c.Next(bc.after)
				if err != nil {
					b.Fatal(err)
				}
				sinkTime = t
			}
		})
	}
}
