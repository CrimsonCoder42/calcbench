package cron

import (
	"errors"
	"fmt"
)

// Errors returned by this package. They are wrapped with the offending field; test with errors.Is.
var (
	// ErrSyntax: the specification is not five blank-separated POSIX fields of numbers, ranges, commas or `*`.
	ErrSyntax = errors.New("cron: syntax error")
	// ErrRange: a number is outside its field's range, or a range has lo > hi.
	ErrRange = errors.New("cron: value out of range")
	// ErrNeverFires: no instant can ever match the specification.
	ErrNeverFires = errors.New("cron: specification never fires")
)

// CronSpec is a parsed crontab time specification: one bitset per field and whether each day field was
// literally `*`. It is a comparable value; build it with ParseCron. The zero value never fires.
type CronSpec struct {
	minute  uint64 // bits 0–59
	hour    uint32 // bits 0–23
	dom     uint32 // bits 1–31
	month   uint16 // bits 1–12
	dow     uint8  // bits 0–6, 0 = Sunday
	domStar bool   // day-of-month field was literally "*"
	dowStar bool   // day-of-week field was literally "*"
}

// field describes one of the five positions.
type field struct {
	name   string
	lo, hi int
}

var fields = [5]field{
	{"minute", 0, 59},
	{"hour", 0, 23},
	{"day of month", 1, 31},
	{"month", 1, 12},
	{"day of week", 0, 6},
}

// maxDays is the longest length of each month (February 29), index 1–12.
var maxDays = [13]int{0, 31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

// ParseCron returns the CronSpec of a five-field POSIX crontab time specification such as "30 2 * * 1-5".
// Fields are separated by spaces or tabs; leading and trailing blanks are ignored.
// Errors: ErrSyntax when the field count is not 5 or a field is not `*` or a comma list of numbers and
// lo-hi ranges (steps, names, `?`, `L`, macros and empty elements included); ErrRange when a number is
// outside its field (minute 0–59, hour 0–23, day of month 1–31, month 1–12, day of week 0–6) or lo > hi;
// ErrNeverFires when day of week is `*`, day of month is restricted, and no allowed month has an allowed day
// (February counted as 29 days). Each error is wrapped with the field's name and text.
//
// ID:      ALG-CRON-001
// Time:    O(n)   n = len(spec)
// Space:   O(1), no allocation on success
// Source:  IEEE Std 1003.1-2024 (POSIX.1-2024), Shell & Utilities, "crontab", INPUT FILES (five fields; `*`, numbers, lo-hi ranges, comma lists; ranges 0–59, 0–23, 1–31, 1–12, 0–6 with 0 = Sunday), https://pubs.opengroup.org/onlinepubs/9799919799/utilities/crontab.html (read 2026-11-13)
// Tags:    cron, crontab, schedule, parse, posix
// Units:   spec is text; the CronSpec holds minutes, hours, days, months and weekdays as bitsets
// Example: ParseCron("0 10 * * 1-5") fires at 10:00 Monday to Friday; ParseCron("0 0 31 2 *") is ErrNeverFires
func ParseCron(spec string) (CronSpec, error) {
	var parts [5]string
	n := 0
	for i := 0; i < len(spec); {
		if spec[i] == ' ' || spec[i] == '\t' {
			i++
			continue
		}
		j := i
		for j < len(spec) && spec[j] != ' ' && spec[j] != '\t' {
			j++
		}
		if n < len(parts) {
			parts[n] = spec[i:j]
		}
		n++
		i = j
	}
	if n != len(parts) {
		return CronSpec{}, fmt.Errorf("%w: %d fields, want 5 (minute hour day-of-month month day-of-week)", ErrSyntax, n)
	}
	var set [5]uint64
	var star [5]bool
	for k, f := range fields {
		s, st, err := parseField(parts[k], f.lo, f.hi)
		if err != nil {
			return CronSpec{}, fmt.Errorf("%w: %s field %q", err, f.name, parts[k])
		}
		set[k], star[k] = s, st
	}
	c := CronSpec{
		minute:  set[0],
		hour:    uint32(set[1]),
		dom:     uint32(set[2]),
		month:   uint16(set[3]),
		dow:     uint8(set[4]),
		domStar: star[2],
		dowStar: star[4],
	}
	if c.dowStar && !c.domStar {
		fires := false
		for m := 1; m <= 12; m++ {
			if c.month&(1<<m) != 0 && c.dom&(uint32(1)<<(maxDays[m]+1)-2) != 0 {
				fires = true
				break
			}
		}
		if !fires {
			return CronSpec{}, fmt.Errorf("%w: no allowed month has an allowed day of month in %q", ErrNeverFires, spec)
		}
	}
	return c, nil
}

// parseField returns the bitset of one field (bit v set when value v is allowed) and whether it was "*".
// The returned error is a bare sentinel; the caller wraps it.
func parseField(s string, lo, hi int) (uint64, bool, error) {
	if s == "*" {
		return (uint64(1)<<(hi+1) - 1) &^ (uint64(1)<<lo - 1), true, nil
	}
	var set uint64
	for i := 0; ; {
		a, j := number(s, i)
		if j == i {
			return 0, false, ErrSyntax // empty element, or not a digit
		}
		b := a
		if j < len(s) && s[j] == '-' {
			k := j + 1
			b, j = number(s, k)
			if j == k {
				return 0, false, ErrSyntax
			}
		}
		if j < len(s) && s[j] != ',' {
			return 0, false, ErrSyntax
		}
		if a < lo || a > hi || b < lo || b > hi || a > b {
			return 0, false, ErrRange
		}
		set |= (uint64(1)<<(b+1) - 1) &^ (uint64(1)<<a - 1)
		if j == len(s) {
			return set, false, nil
		}
		i = j + 1 // past the comma; an empty next element is caught above
	}
}

// number reads decimal digits of s from i and returns their value (saturated at 1000, above every field
// maximum) and the index after them. No digits: it returns (0, i).
func number(s string, i int) (int, int) {
	v := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		if v < 1000 {
			v = v*10 + int(s[i]-'0')
		}
		i++
	}
	return v, i
}
