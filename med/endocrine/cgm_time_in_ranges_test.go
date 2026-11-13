package endocrine

import (
	"errors"
	"math"
	"testing"
)

// pctTol: percents are 100·count/n, exact to a few ulps.
const pctTol = 1e-12

func checkRanges(t *testing.T, name string, got CGMRangePercents, want [5]float64, n int) {
	t.Helper()
	g := [5]float64{got.BelowRangeLevel2Percent, got.BelowRangeLevel1Percent, got.InRangePercent,
		got.AboveRangeLevel1Percent, got.AboveRangeLevel2Percent}
	for i := range want {
		if math.Abs(g[i]-want[i]) > pctTol {
			t.Errorf("%s: range %d = %v, want %v", name, i, g[i], want[i])
		}
	}
	if got.Readings != n {
		t.Errorf("%s: readings %d, want %d", name, got.Readings, n)
	}
}

// parity: ledger TimeInRange, source case 1 (9 of 10 readings in 70–180 mg/dL, one above; 180 counts as in range).
func TestTimeInRangesMgDL_parity(t *testing.T) {
	g := []float64{110, 180, 95, 150, 152, 90, 200, 120, 130, 85}
	got, err := TimeInRangesMgDL(g)
	if err != nil {
		t.Fatal(err)
	}
	checkRanges(t, "parity", got, [5]float64{0, 0, 90, 10, 0}, 10)
}

// Battelino 2019 Table 2 edges, both sides of each: 53.9|54, 69|69.5|70, 180|180.5|181, 250|250.5.
func TestTimeInRangesMgDL_edges(t *testing.T) {
	cases := []struct {
		g    float64
		band int
	}{
		{1, 0}, {53.9, 0}, {54, 1}, {69, 1}, {69.5, 1}, {70, 2}, {180, 2}, {180.5, 3}, {181, 3}, {250, 3},
		{250.5, 4}, {600, 4},
	}
	for _, c := range cases {
		got, err := TimeInRangesMgDL([]float64{c.g})
		if err != nil {
			t.Fatal(err)
		}
		var want [5]float64
		want[c.band] = 100
		checkRanges(t, "edge", got, want, 1)
	}
}

func TestTimeInRangesMmolL_edges(t *testing.T) {
	// Doc example, HAND-DERIVED: 2.9, 3.0, 3.9, 10.0, 10.1, 14.0.
	got, err := TimeInRangesMmolL([]float64{2.9, 3.0, 3.9, 10.0, 10.1, 14.0})
	if err != nil {
		t.Fatal(err)
	}
	checkRanges(t, "example", got, [5]float64{100.0 / 6, 100.0 / 6, 200.0 / 6, 100.0 / 6, 100.0 / 6}, 6)
	cases := []struct {
		g    float64
		band int
	}{
		{2.99, 0}, {3.0, 1}, {3.8, 1}, {3.85, 1}, {3.9, 2}, {10.0, 2}, {10.05, 3}, {13.9, 3}, {13.95, 4},
	}
	for _, c := range cases {
		got, _ := TimeInRangesMmolL([]float64{c.g})
		var want [5]float64
		want[c.band] = 100
		checkRanges(t, "edge", got, want, 1)
	}
}

func TestTimeInRanges_errors(t *testing.T) {
	for _, f := range []func([]float64) (CGMRangePercents, error){TimeInRangesMgDL, TimeInRangesMmolL} {
		cases := []struct {
			in   []float64
			want error
		}{
			{nil, ErrOutOfRange},
			{[]float64{100, 0}, ErrOutOfRange},
			{[]float64{-5}, ErrOutOfRange},
			{[]float64{math.NaN()}, ErrNotFinite},
			{[]float64{100, math.Inf(1)}, ErrNotFinite},
			{[]float64{math.Inf(-1)}, ErrNotFinite},
		}
		for _, c := range cases {
			if _, err := f(c.in); !errors.Is(err, c.want) {
				t.Errorf("%v: err = %v, want %v", c.in, err, c.want)
			}
		}
	}
}

func TestTimeInRanges_noAlloc(t *testing.T) {
	g := []float64{50, 60, 100, 200, 300}
	if a := testing.AllocsPerRun(100, func() { _, _ = TimeInRangesMgDL(g) }); a != 0 {
		t.Errorf("mg/dL allocs = %v, want 0", a)
	}
	if a := testing.AllocsPerRun(100, func() { _, _ = TimeInRangesMmolL(g) }); a != 0 {
		t.Errorf("mmol/L allocs = %v, want 0", a)
	}
}

func BenchmarkTimeInRangesMgDL(b *testing.B) {
	g := make([]float64, 1000)
	for i := range g {
		g[i] = float64(40 + i%300)
	}
	b.ReportAllocs()
	for b.Loop() {
		_, _ = TimeInRangesMgDL(g)
	}
}

func BenchmarkTimeInRangesMmolL(b *testing.B) {
	g := make([]float64, 1000)
	for i := range g {
		g[i] = 2 + float64(i%150)/10
	}
	b.ReportAllocs()
	for b.Loop() {
		_, _ = TimeInRangesMmolL(g)
	}
}
