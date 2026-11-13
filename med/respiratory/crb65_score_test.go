package respiratory

import (
	"errors"
	"math"
	"testing"
)

// HAND-DERIVED from the CRB-65 criteria (Dwyer 2014, PMID 25478185; Chen 2016, PMID 27250351).
func TestCRB65Score_criteriaExample(t *testing.T) {
	cases := []struct {
		conf              bool
		rr, sbp, dbp, age float64
		want              int
	}{
		{true, 32, 110, 70, 70, 3},
		{true, 30, 89, 90, 65, 4},
		{false, 16, 120, 80, 40, 0},
	}
	for _, c := range cases {
		if got, err := CRB65Score(c.conf, c.rr, c.sbp, c.dbp, c.age); err != nil || got != c.want {
			t.Errorf("%v: got (%d, %v), want %d", c, got, err, c.want)
		}
	}
}

// Both sides of each criterion's edge: RR ≥ 30, SBP < 90 or DBP ≤ 60, age ≥ 65.
func TestCRB65Score_edges(t *testing.T) {
	cases := []struct {
		name              string
		rr, sbp, dbp, age float64
		want              int
	}{
		{"RR 29", 29, 120, 80, 50, 0}, {"RR 30", 30, 120, 80, 50, 1},
		{"SBP 90", 20, 90, 80, 50, 0}, {"SBP 89.9", 20, 89.9, 80, 50, 1},
		{"DBP 61", 20, 120, 61, 50, 0}, {"DBP 60", 20, 120, 60, 50, 1},
		{"both BP low count once", 20, 80, 50, 50, 1},
		{"age 64.9", 20, 120, 80, 64.9, 0}, {"age 65", 20, 120, 80, 65, 1},
	}
	for _, c := range cases {
		if got, err := CRB65Score(false, c.rr, c.sbp, c.dbp, c.age); err != nil || got != c.want {
			t.Errorf("%s: got (%d, %v), want %d", c.name, got, err, c.want)
		}
	}
}

// CRB-65 equals CURB-65 with urea at or below 7 mmol/L (the urea item off), over a capped grid.
func TestCRB65Score_matchesCURB65WithoutUrea(t *testing.T) {
	for _, conf := range []bool{false, true} {
		for rr := 10.0; rr <= 40; rr += 5 {
			for sbp := 70.0; sbp <= 130; sbp += 20 {
				for _, age := range []float64{30, 64, 65, 90} {
					crb, err1 := CRB65Score(conf, rr, sbp, 70, age)
					curb, err2 := CURB65Score(conf, 5, rr, sbp, 70, age)
					if err1 != nil || err2 != nil || crb != curb {
						t.Fatalf("conf %v rr %v sbp %v age %v: CRB-65 %d, CURB-65 %d", conf, rr, sbp, age, crb, curb)
					}
				}
			}
		}
	}
}

// parity: ledger crb65_calc, source case 1: the source scores confusion, RR ≥ 30, SBP < 90 or DBP ≤ 60, age ≥ 65,
// one point each, as the secondaries; RR 30, BP 100/60, age 66 → 3.
func TestCRB65Score_parity(t *testing.T) {
	if got, err := CRB65Score(false, 30, 100, 60, 66); err != nil || got != 3 {
		t.Errorf("got (%d, %v), want 3", got, err)
	}
}

func TestCRB65Score_errors(t *testing.T) {
	cases := []struct {
		rr, sbp, dbp, age float64
		want              error
	}{
		{-1, 120, 80, 50, ErrOutOfRange}, {20, 0, 80, 50, ErrOutOfRange}, {20, 120, 0, 50, ErrOutOfRange}, {20, 120, 80, -1, ErrOutOfRange},
		{math.NaN(), 120, 80, 50, ErrNotFinite}, {20, math.Inf(1), 80, 50, ErrNotFinite}, {20, 120, math.Inf(-1), 50, ErrNotFinite}, {20, 120, 80, math.NaN(), ErrNotFinite},
	}
	for _, c := range cases {
		if _, err := CRB65Score(false, c.rr, c.sbp, c.dbp, c.age); !errors.Is(err, c.want) {
			t.Errorf("%v: err = %v, want %v", c, err, c.want)
		}
	}
}

func TestCRB65Score_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = CRB65Score(true, 32, 110, 70, 70) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkCRB65Score(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt, _ = CRB65Score(true, 32, 110, 70, 70)
	}
}
