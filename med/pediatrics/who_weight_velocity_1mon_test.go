package pediatrics

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"errors"
	"strings"
	"testing"
)

//go:embed data/README.md
var pediatricsDataReadme string

var sinkWHOIncrement WHOIncrementPercentiles

// The embedded CSV must be the file whose SHA-256 data/README.md records.
func TestWHOWeightIncrement1MonthPercentiles_sha256(t *testing.T) {
	sum := sha256.Sum256([]byte(whoWeightVelocity1MonCSV))
	got := hex.EncodeToString(sum[:])
	if !strings.Contains(pediatricsDataReadme, "| who_weight_velocity_1mon_p.csv (embedded) | derived from the two files above | 2026-11-12 | "+got+" |") {
		t.Errorf("embedded CSV SHA-256 %s is not the one recorded in data/README.md", got)
	}
}

// Published values, read from the WHO workbooks ttt-weight-boys-1mon-p.xlsx and ttt-weight-girls-1mon-p.xlsx
// (WHO Child Growth Standards, weight velocity, 1-month increments, percentiles). Exact: the file's integers.
func TestWHOWeightIncrement1MonthPercentiles_published(t *testing.T) {
	type cell struct {
		female bool
		k      int
		label  string
		pick   func(WHOIncrementPercentiles) float64
		want   float64
	}
	p1 := func(r WHOIncrementPercentiles) float64 { return r.P1 }
	p3 := func(r WHOIncrementPercentiles) float64 { return r.P3 }
	p15 := func(r WHOIncrementPercentiles) float64 { return r.P15 }
	p50 := func(r WHOIncrementPercentiles) float64 { return r.P50 }
	p85 := func(r WHOIncrementPercentiles) float64 { return r.P85 }
	p97 := func(r WHOIncrementPercentiles) float64 { return r.P97 }
	p99 := func(r WHOIncrementPercentiles) float64 { return r.P99 }
	cases := []cell{
		// boys
		{false, 0, "0 – 4 wks", p50, 1023}, {false, 0, "0 – 4 wks", p3, 369}, {false, 0, "0 – 4 wks", p97, 1575},
		{false, 1, "4 wks – 2 mo", p99, 1955}, {false, 5, "5 – 6 mo", p1, -17}, {false, 6, "6 – 7 mo", p3, 0},
		{false, 11, "11 – 12 mo", p50, 239}, {false, 9, "9 – 10 mo", p85, 486},
		// girls
		{true, 0, "0 – 4 wks", p50, 879}, {true, 0, "0 – 4 wks", p1, 280}, {true, 2, "2 – 3 mo", p15, 494},
		{true, 7, "7 – 8 mo", p97, 711}, {true, 11, "11 – 12 mo", p50, 232}, {true, 11, "11 – 12 mo", p1, -222},
	}
	for _, c := range cases {
		r, err := WHOWeightIncrement1MonthPercentiles(c.female, c.k)
		if err != nil || r.Interval != c.label || c.pick(r) != c.want {
			t.Errorf("WHOWeightIncrement1MonthPercentiles(%v, %d) = (%+v, %v), want interval %q value %v", c.female, c.k, r, err, c.label, c.want)
		}
	}
	// One full published row: boys "2 – 3 mo".
	want := WHOIncrementPercentiles{"2 – 3 mo", 307, 397, 446, 577, 658, 815, 980, 1071, 1228, 1290, 1410}
	if r, _ := WHOWeightIncrement1MonthPercentiles(false, 2); r != want {
		t.Errorf("boys 2 – 3 mo = %+v, want %+v", r, want)
	}
}

// Every row is ordered P1 < P3 < … < P99 and labelled per the documented index mapping.
func TestWHOWeightIncrement1MonthPercentiles_rows(t *testing.T) {
	labels := [12]string{"0 – 4 wks", "4 wks – 2 mo", "2 – 3 mo", "3 – 4 mo", "4 – 5 mo", "5 – 6 mo",
		"6 – 7 mo", "7 – 8 mo", "8 – 9 mo", "9 – 10 mo", "10 – 11 mo", "11 – 12 mo"}
	for _, female := range []bool{false, true} {
		for k := 0; k < 12; k++ {
			r, err := WHOWeightIncrement1MonthPercentiles(female, k)
			if err != nil || r.Interval != labels[k] {
				t.Fatalf("(%v, %d) = (%+v, %v), want interval %q", female, k, r, err, labels[k])
			}
			p := []float64{r.P1, r.P3, r.P5, r.P15, r.P25, r.P50, r.P75, r.P85, r.P95, r.P97, r.P99}
			for i := 1; i < len(p); i++ {
				if !(p[i] > p[i-1]) {
					t.Errorf("(%v, %d) percentiles not increasing at %d: %v", female, k, i, p)
				}
			}
		}
	}
}

func TestWHOWeightIncrement1MonthPercentiles_errors(t *testing.T) {
	for _, k := range []int{-1, 12, 100} {
		if _, err := WHOWeightIncrement1MonthPercentiles(false, k); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("intervalIndex %d: error = %v, want ErrOutOfRange", k, err)
		}
	}
}

func TestWHOWeightIncrement1MonthPercentiles_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sinkWHOIncrement, _ = WHOWeightIncrement1MonthPercentiles(true, 3)
	}); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkWHOWeightIncrement1MonthPercentiles(b *testing.B) {
	b.ReportAllocs()
	k := 3
	for b.Loop() {
		sinkWHOIncrement, _ = WHOWeightIncrement1MonthPercentiles(true, k)
	}
}
