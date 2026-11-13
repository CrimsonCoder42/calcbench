package renal

import "testing"

// Moore 2014 Table 3 points; duration edges on both sides.
func TestSTONEScore_table3(t *testing.T) {
	tests := []struct {
		male      bool
		hours     float64
		nonBlack  bool
		nv        STONENauseaVomiting
		hematuria bool
		want      int
	}{
		{false, 30, false, STONENoNauseaVomiting, false, 0},
		{true, 30, false, STONENoNauseaVomiting, false, 2},
		{false, 24.1, false, STONENoNauseaVomiting, false, 0},
		{false, 24, false, STONENoNauseaVomiting, false, 1},
		{false, 6, false, STONENoNauseaVomiting, false, 1},
		{false, 5.9, false, STONENoNauseaVomiting, false, 3},
		{false, 30, true, STONENoNauseaVomiting, false, 3},
		{false, 30, false, STONENauseaAlone, false, 1},
		{false, 30, false, STONEVomitingAlone, false, 2},
		{false, 30, false, STONENoNauseaVomiting, true, 3},
		{true, 3, true, STONEVomitingAlone, true, 13},
	}
	for i, tc := range tests {
		if got, err := STONEScore(tc.male, tc.hours, tc.nonBlack, tc.nv, tc.hematuria); err != nil || got != tc.want {
			t.Errorf("case %d: (%d, %v), want %d", i, got, err, tc.want)
		}
	}
}

// parity: ledger stone_score, source case 1. The app drops race ("race-free version"), gives hematuria 1 and an
// "erythrocytes" item 2; the primary's Table 3 has race 3 and hematuria on dipstick 3. Male, < 6 h, nausea alone,
// hematuria, non-black: primary 2+3+1+3+3 = 12 — asserted here (app 2+3+1+1+2 = 9 with erythrocytes).
func TestSTONEScore_parity(t *testing.T) {
	if got, _ := STONEScore(true, 2, true, STONENauseaAlone, true); got != 12 {
		t.Errorf("got %d, want 12 (primary)", got)
	}
}

func TestSTONEScore_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { _, e := STONEScore(true, nan, true, 0, true); return e }, ErrNotFinite},
		{"+Inf", func() error { _, e := STONEScore(true, pInf, true, 0, true); return e }, ErrNotFinite},
		{"-Inf", func() error { _, e := STONEScore(true, nInf, true, 0, true); return e }, ErrNotFinite},
		{"negative hours", func() error { _, e := STONEScore(true, -1, true, 0, true); return e }, ErrOutOfRange},
		{"nv -1", func() error { _, e := STONEScore(true, 1, true, -1, true); return e }, ErrOutOfRange},
		{"nv 3", func() error { _, e := STONEScore(true, 1, true, 3, true); return e }, ErrOutOfRange},
	})
}

func TestSTONEScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkB6Int, _ = STONEScore(true, 3, true, STONENauseaAlone, true) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkSTONEScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkB6Int, _ = STONEScore(true, 3, true, STONENauseaAlone, true)
	}
}
