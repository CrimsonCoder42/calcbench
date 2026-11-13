package business

import (
	"errors"
	"math"
	"testing"
)

// Principles of Accounting Vol. 1 §9.3, BWW: 2017 average receivables $80,000, credit sales $400,000,
// 365 ÷ 5 = 73 days (exact: 365 × 8,000,000 ÷ 40,000,000 = 73). 2018: $87,500 and $450,000, printed
// 71.01 days from the rounded turnover 5.14; the unrounded value is 70.97…, so that row is checked against
// 365 × 87,500 ÷ 450,000 to 1e-12 relative (two float roundings) and against the printed value to 0.05
// (the source rounded the turnover before dividing).
func TestAverageCollectionPeriod_sourceExample(t *testing.T) {
	got, err := AverageCollectionPeriod(8000000, 40000000, 365)
	if err != nil || got != 73 {
		t.Fatalf("2017: got %v, %v; want 73", got, err)
	}
	got, err = AverageCollectionPeriod(8750000, 45000000, 365)
	if err != nil {
		t.Fatal(err)
	}
	want := 365.0 * 87500 / 450000
	if math.Abs(got-want) > 1e-12*want || math.Abs(got-71.01) > 0.05 {
		t.Fatalf("2018: got %v, want %v (printed 71.01)", got, want)
	}
}

func TestAverageCollectionPeriod_table(t *testing.T) {
	cases := []struct {
		name      string
		ar, sales int64
		days      int
		want      float64
	}{
		{"zero receivables", 0, 100, 365, 0},
		{"one day", 50, 100, 1, 0.5},
		{"360-day year", 100, 100, 360, 360},
		{"max receivables", maxExactCents, maxExactCents, 365, 365},
	}
	for _, c := range cases {
		got, err := AverageCollectionPeriod(c.ar, c.sales, c.days)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestAverageCollectionPeriod_errors(t *testing.T) {
	cases := []struct {
		name      string
		ar, sales int64
		days      int
		want      error
	}{
		{"zero sales", 1, 0, 365, ErrZeroDenominator},
		{"zero days", 1, 1, 0, ErrPeriods},
		{"negative days", 1, 1, -1, ErrPeriods},
		{"receivables above 2^53", maxExactCents + 1, 1, 365, ErrOutOfRange},
		{"sales above 2^53", 1, maxExactCents + 1, 365, ErrOutOfRange},
		{"negative receivables", -1, 1, 365, ErrNegative},
		{"negative sales", 1, -1, 365, ErrNegative},
	}
	for _, c := range cases {
		if _, err := AverageCollectionPeriod(c.ar, c.sales, c.days); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// Identity with FIN-BIZ-014: days ÷ receivables turnover, to 4 ulps (two to three roundings each side).
func TestAverageCollectionPeriod_turnoverIdentity(t *testing.T) {
	for _, c := range [][2]int64{{8000000, 40000000}, {12345, 987654}, {1, 3}} {
		turnover, _ := ReceivablesTurnover(c[1], c[0])
		got, _ := AverageCollectionPeriod(c[0], c[1], 365)
		want := 365 / turnover
		if math.Abs(got-want) > 4*math.Abs(want)*0x1p-52 {
			t.Errorf("%v: got %v, want %v", c, got, want)
		}
	}
}

func TestAverageCollectionPeriod_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = AverageCollectionPeriod(8000000, 40000000, 365) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkAverageCollectionPeriod(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = AverageCollectionPeriod(8000000, 40000000, 365)
	}
}
