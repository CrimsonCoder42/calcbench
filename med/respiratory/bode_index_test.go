package respiratory

import (
	"errors"
	"math"
	"testing"
)

// HAND-DERIVED from the BODE table (Funk 2009 Table 1, PMID 19134192; Kamath 2020, PMC7842444):
// FEV₁ 45 (2) + 300 m (1) + MMRC 2 (1) + BMI 20 (1) = 5; the extremes 0 and 10.
func TestBODEIndex_tableExample(t *testing.T) {
	cases := []struct {
		fev1, walk float64
		mmrc       int
		bmi        float64
		want       int
	}{{45, 300, 2, 20, 5}, {80, 500, 0, 25, 0}, {30, 100, 4, 18, 10}}
	for _, c := range cases {
		if got, err := BODEIndex(c.fev1, c.walk, c.mmrc, c.bmi); err != nil || got != c.want {
			t.Errorf("%v: got (%d, %v), want %d", c, got, err, c.want)
		}
	}
}

// Both sides of every band edge.
func TestBODEIndex_edges(t *testing.T) {
	score := func(fev1, walk float64, mmrc int, bmi float64) int {
		s, err := BODEIndex(fev1, walk, mmrc, bmi)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
		return s
	}
	for _, c := range []struct {
		v    float64
		want int
	}{{65, 0}, {64, 1}, {64.5, 1}, {50, 1}, {49, 2}, {36, 2}, {35.5, 3}, {35, 3}} {
		if got := score(c.v, 400, 0, 25); got != c.want {
			t.Errorf("FEV1 %v: got %d, want %d", c.v, got, c.want)
		}
	}
	for _, c := range []struct {
		v    float64
		want int
	}{{350, 0}, {349, 1}, {250, 1}, {249, 2}, {150, 2}, {149, 3}, {0, 3}} {
		if got := score(80, c.v, 0, 25); got != c.want {
			t.Errorf("walk %v: got %d, want %d", c.v, got, c.want)
		}
	}
	for _, c := range []struct{ m, want int }{{0, 0}, {1, 0}, {2, 1}, {3, 2}, {4, 3}} {
		if got := score(80, 400, c.m, 25); got != c.want {
			t.Errorf("MMRC %d: got %d, want %d", c.m, got, c.want)
		}
	}
	for _, c := range []struct {
		v    float64
		want int
	}{{21.01, 0}, {21, 1}} {
		if got := score(80, 400, 0, c.v); got != c.want {
			t.Errorf("BMI %v: got %d, want %d", c.v, got, c.want)
		}
	}
}

// parity: ledger bode_index, source case 1: FEV₁ 55, 260 m, MMRC 3, BMI 22 → 1 + 1 + 2 + 0 = 4.
func TestBODEIndex_parity(t *testing.T) {
	if got, err := BODEIndex(55, 260, 3, 22); err != nil || got != 4 {
		t.Errorf("got (%d, %v), want 4", got, err)
	}
}

func TestBODEIndex_errors(t *testing.T) {
	cases := []struct {
		fev1, walk float64
		mmrc       int
		bmi        float64
		want       error
	}{
		{0, 300, 1, 22, ErrOutOfRange}, {50, -1, 1, 22, ErrOutOfRange}, {50, 300, -1, 22, ErrOutOfRange},
		{50, 300, 5, 22, ErrOutOfRange}, {50, 300, 1, 0, ErrOutOfRange},
		{math.NaN(), 300, 1, 22, ErrNotFinite}, {50, math.Inf(1), 1, 22, ErrNotFinite}, {50, 300, 1, math.Inf(-1), ErrNotFinite},
	}
	for _, c := range cases {
		if _, err := BODEIndex(c.fev1, c.walk, c.mmrc, c.bmi); !errors.Is(err, c.want) {
			t.Errorf("%v: err = %v, want %v", c, err, c.want)
		}
	}
}

func TestBODEIndex_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = BODEIndex(45, 300, 2, 20) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkBODEIndex(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt, _ = BODEIndex(45, 300, 2, 20)
	}
}
