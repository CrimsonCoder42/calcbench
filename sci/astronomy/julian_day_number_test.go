package astronomy

import (
	"errors"
	"math/rand/v2"
	"testing"
	"time"
)

// Meeus, Astronomical Algorithms, 2nd ed., ch. 7: Example 7.a, 1957 October 4.81 is JD 2436116.31, so the
// day's noon is JDN 2436116; J2000.0 (2000 January 1.5) is JD 2451545.0.
func TestGregorianToJulianDayNumber_meeusExamples(t *testing.T) {
	cases := []struct {
		name string
		y    int64
		m, d int
		want int64
	}{
		{"example 7.a sputnik", 1957, 10, 4, 2436116},
		{"J2000", 2000, 1, 1, 2451545},
		{"example 7.e 1954-06-30 (JD 2434923.5 at 0h)", 1954, 6, 30, 2434924},
		{"first supported day", -4712, 1, 1, 38},
		{"last supported day", 1000000, 12, 31, 366963925},
		{"leap day 2000", 2000, 2, 29, 2451604},
		{"gregorian reform 1582-10-15 (JD 2299160.5 at 0h)", 1582, 10, 15, 2299161},
	}
	for _, c := range cases {
		got, err := GregorianToJulianDayNumber(c.y, c.m, c.d)
		if err != nil || got != c.want {
			t.Errorf("%s: got %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

func TestGregorianToJulianDayNumber_errors(t *testing.T) {
	cases := []struct {
		y    int64
		m, d int
		want error
	}{
		{2023, 2, 29, ErrInvalidDate},
		{1900, 2, 29, ErrInvalidDate},
		{2024, 13, 1, ErrInvalidDate},
		{2024, 0, 1, ErrInvalidDate},
		{2024, 4, 31, ErrInvalidDate},
		{2024, 1, 0, ErrInvalidDate},
		{-4713, 12, 31, ErrOutOfRange},
		{1000001, 1, 1, ErrOutOfRange},
	}
	for _, c := range cases {
		if _, err := GregorianToJulianDayNumber(c.y, c.m, c.d); !errors.Is(err, c.want) {
			t.Errorf("%d-%d-%d: err %v, want %v", c.y, c.m, c.d, err, c.want)
		}
	}
}

// Oracle: Go's time package uses the proleptic Gregorian calendar; JDN(1970-01-01) = 2440588.
func TestJulianDayNumber_oracleAgainstTimePackage(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < 500; i++ {
		z := minJDN + r.Int64N(maxJDN-minJDN+1)
		y, m, d, err := JulianDayNumberToGregorian(z)
		if err != nil {
			t.Fatalf("jdn %d: %v", z, err)
		}
		tm := time.Unix((z-2440588)*86400, 0).UTC()
		if int(y) != tm.Year() || m != int(tm.Month()) || d != tm.Day() {
			t.Fatalf("jdn %d: got %d-%d-%d, time says %v", z, y, m, d, tm)
		}
		if back, _ := GregorianToJulianDayNumber(y, m, d); back != z {
			t.Fatalf("round trip %d → %d-%d-%d → %d", z, y, m, d, back)
		}
	}
}

func TestJulianDayNumberToGregorian_examplesAndRange(t *testing.T) {
	if y, m, d, err := JulianDayNumberToGregorian(2436116); err != nil || y != 1957 || m != 10 || d != 4 {
		t.Errorf("Meeus 7.c: got %d-%d-%d %v", y, m, d, err)
	}
	for _, z := range []int64{37, -1, 366963926} {
		if _, _, _, err := JulianDayNumberToGregorian(z); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("jdn %d: err %v", z, err)
		}
	}
}

func BenchmarkGregorianToJulianDayNumber(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = GregorianToJulianDayNumber(1957, 10, 4)
	}
}

func BenchmarkJulianDayNumberToGregorian(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _, _, _ = JulianDayNumberToGregorian(2436116)
	}
}
