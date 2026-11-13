package cron

import (
	"errors"
	"strings"
	"testing"
)

func mustParse(t testing.TB, spec string) CronSpec {
	t.Helper()
	c, err := ParseCron(spec)
	if err != nil {
		t.Fatalf("ParseCron(%q): %v", spec, err)
	}
	return c
}

// Worked example of the BRIEF (doc/toolbench_requests_2026-11-13/BRIEF.md, B3): "0 10 * * 1-5" is 10:00 on
// Monday to Friday; the field meanings and ranges are POSIX crontab, INPUT FILES.
func TestParseCron_workedExample(t *testing.T) {
	c := mustParse(t, "0 10 * * 1-5")
	want := CronSpec{
		minute:  1,
		hour:    1 << 10,
		dom:     0xFFFFFFFE, // 1..31
		month:   0x1FFE,     // 1..12
		dow:     0b0111110,  // Monday..Friday
		domStar: true,
		dowStar: false,
	}
	if c != want {
		t.Fatalf("got %+v, want %+v", c, want)
	}
}

func TestParseCron_fields(t *testing.T) {
	cases := []struct {
		name string
		spec string
		want CronSpec
	}{
		{"all stars", "* * * * *", CronSpec{minute: 1<<60 - 1, hour: 1<<24 - 1, dom: 0xFFFFFFFE, month: 0x1FFE, dow: 0x7F, domStar: true, dowStar: true}},
		{"lists and ranges", "1,3-5,59 0,23 1,31 2,12 0,6", CronSpec{minute: 1<<1 | 1<<3 | 1<<4 | 1<<5 | 1<<59, hour: 1 | 1<<23, dom: 1<<1 | 1<<31, month: 1<<2 | 1<<12, dow: 1 | 1<<6}},
		{"single-value range", "5-5 * * * *", CronSpec{minute: 1 << 5, hour: 1<<24 - 1, dom: 0xFFFFFFFE, month: 0x1FFE, dow: 0x7F, domStar: true, dowStar: true}},
		{"leading zeros", "05 007 * * *", CronSpec{minute: 1 << 5, hour: 1 << 7, dom: 0xFFFFFFFE, month: 0x1FFE, dow: 0x7F, domStar: true, dowStar: true}},
		{"tabs and extra blanks", " \t0\t\t0  * *   * \t", CronSpec{minute: 1, hour: 1, dom: 0xFFFFFFFE, month: 0x1FFE, dow: 0x7F, domStar: true, dowStar: true}},
		{"full range is not a star", "* * 1-31 * 0-6", CronSpec{minute: 1<<60 - 1, hour: 1<<24 - 1, dom: 0xFFFFFFFE, month: 0x1FFE, dow: 0x7F}},
		{"overlapping list", "1-10,5-6,10 * * * *", CronSpec{minute: 0x7FE, hour: 1<<24 - 1, dom: 0xFFFFFFFE, month: 0x1FFE, dow: 0x7F, domStar: true, dowStar: true}},
		{"Feb 29 alone", "0 0 29 2 *", CronSpec{minute: 1, hour: 1, dom: 1 << 29, month: 1 << 2, dow: 0x7F, dowStar: true}},
		{"Feb 31 with a weekday (OR) fires", "0 0 31 2 0-6", CronSpec{minute: 1, hour: 1, dom: 1 << 31, month: 1 << 2, dow: 0x7F}},
		{"31 in a list with a 31-day month", "0 0 31 2,3 *", CronSpec{minute: 1, hour: 1, dom: 1 << 31, month: 1<<2 | 1<<3, dow: 0x7F, dowStar: true}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := mustParse(t, tc.spec); got != tc.want {
				t.Fatalf("ParseCron(%q) = %+v, want %+v", tc.spec, got, tc.want)
			}
		})
	}
}

func TestParseCron_syntaxErrors(t *testing.T) {
	for _, spec := range []string{
		"", "   ", "\t",
		"* * * *",       // 4 fields
		"* * * * * *",   // 6 fields (a seconds or year field)
		"0 0 0 * * * *", // 7 fields
		"*/5 * * * *",   // step
		"0-30/5 * * * *",
		"1,,2 * * * *",
		",1 * * * *",
		"1, * * * *",
		"-1 * * * *",
		"1- * * * *",
		"1-2-3 * * * *",
		"a * * * *",
		"*,1 * * * *",
		"** * * * *",
		"0 0 * JAN *",
		"0 0 * * MON",
		"0 0 ? * *",
		"0 0 L * *",
		"0 0 1W * *",
		"0 0 * * 5L",
		"0 0 * * 5#2",
		"@hourly",
		"@daily extra words here x",
		"0\n0 * * *",
		"0 0 * * +1",
		"0 0 * * 1.0",
		"0 0 * *  1", // a non-blank space joins the field
	} {
		if _, err := ParseCron(spec); !errors.Is(err, ErrSyntax) {
			t.Errorf("ParseCron(%q) error = %v, want ErrSyntax", spec, err)
		}
	}
}

func TestParseCron_rangeErrors(t *testing.T) {
	cases := []struct{ spec, field string }{
		{"60 * * * *", "minute"},
		{"0-60 * * * *", "minute"},
		{"* 24 * * *", "hour"},
		{"* * 0 * *", "day of month"},
		{"* * 32 * *", "day of month"},
		{"* * * 0 *", "month"},
		{"* * * 13 *", "month"},
		{"* * * * 7", "day of week"},
		{"5-3 * * * *", "minute"},
		{"99999999999999999999999 * * * *", "minute"},
		{"* * * * 0,7", "day of week"},
	}
	for _, tc := range cases {
		_, err := ParseCron(tc.spec)
		if !errors.Is(err, ErrRange) {
			t.Errorf("ParseCron(%q) error = %v, want ErrRange", tc.spec, err)
			continue
		}
		if !strings.Contains(err.Error(), tc.field+" field") {
			t.Errorf("ParseCron(%q) error %q does not name the %s field", tc.spec, err, tc.field)
		}
	}
}

// Worked example of the BRIEF: "0 0 31 2 *" never fires (February has at most 29 days).
func TestParseCron_neverFires(t *testing.T) {
	for _, spec := range []string{
		"0 0 31 2 *",
		"0 0 30 2 *",
		"0 0 30,31 2 *",
		"0 0 31 4,6,9,11 *",
		"0 0 31 2,4 *",
		"* * 30-31 2 *",
	} {
		if _, err := ParseCron(spec); !errors.Is(err, ErrNeverFires) {
			t.Errorf("ParseCron(%q) error = %v, want ErrNeverFires", spec, err)
		}
	}
}

func TestParseCron_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		if _, err := ParseCron("0,15,30,45 9-17 1,15 1-12 1-5"); err != nil {
			t.Fatal(err)
		}
	}); n != 0 {
		t.Fatalf("ParseCron allocates %v times per call, want 0", n)
	}
}

func BenchmarkParseCron(b *testing.B) {
	for _, bc := range []struct{ name, spec string }{
		{"stars", "* * * * *"},
		{"lists", "0,15,30,45 9-17 1,15 1-12 1-5"},
	} {
		b.Run(bc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := ParseCron(bc.spec); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
