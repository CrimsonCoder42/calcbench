package hepatic

import "testing"

// Grades per Benedict 2017 Table 1 (PMC5468341) and Dowman 2011 (PMC3080668); exhaustive over the 4 × 4 × 3
// valid grades against the plain sum, plus both sides of every range edge.
func TestNAFLDActivityScore_example(t *testing.T) {
	if got, err := NAFLDActivityScore(2, 2, 1); err != nil || got != 5 {
		t.Errorf("worked: got (%v, %v), want 5", got, err)
	}
	for s := 0; s <= 3; s++ {
		for l := 0; l <= 3; l++ {
			for b := 0; b <= 2; b++ {
				if got, err := NAFLDActivityScore(s, l, b); err != nil || got != s+l+b {
					t.Errorf("(%d,%d,%d): got (%v, %v), want %d", s, l, b, got, err, s+l+b)
				}
			}
		}
	}
	if got, _ := NAFLDActivityScore(3, 3, 2); got != 8 {
		t.Errorf("maximum: got %v, want 8", got)
	}
}

// parity: ledger nafld_activity_score. The app sums steatosis 0–3, inflammation 0–3, ballooning 0–2; its
// defaults (0, 0, 0) → 0, and (3, 3, 2) → 8. Inflammation 3 is accepted (the secondaries' 0–3, not the PubMed
// abstract's 0–2).
func TestNAFLDActivityScore_parity(t *testing.T) {
	cases := []struct{ s, l, b, want int }{{0, 0, 0, 0}, {3, 3, 2, 8}, {1, 3, 0, 4}}
	for i, c := range cases {
		if got, err := NAFLDActivityScore(c.s, c.l, c.b); err != nil || got != c.want { // parity: ledger nafld_activity_score, source case i+1
			t.Errorf("case %d: got (%v, %v), want %v", i+1, got, err, c.want)
		}
	}
}

func TestNAFLDActivityScore_errors(t *testing.T) {
	call := func(s, l, b int) func() error {
		return func() error { _, err := NAFLDActivityScore(s, l, b); return err }
	}
	runErrCases(t, []errCase{
		{"steatosis -1", call(-1, 0, 0), ErrOutOfRange},
		{"steatosis 4", call(4, 0, 0), ErrOutOfRange},
		{"inflammation -1", call(0, -1, 0), ErrOutOfRange},
		{"inflammation 4", call(0, 4, 0), ErrOutOfRange},
		{"ballooning -1", call(0, 0, -1), ErrOutOfRange},
		{"ballooning 3", call(0, 0, 3), ErrOutOfRange},
	})
}

func TestNAFLDActivityScore_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = NAFLDActivityScore(2, 2, 1) }); n != 0 {
		t.Errorf("NAFLDActivityScore allocates %v times, want 0", n)
	}
}

func BenchmarkNAFLDActivityScore(b *testing.B) {
	b.ReportAllocs()
	s, l, bal := 2, 2, 1
	for b.Loop() {
		sinkInt, _ = NAFLDActivityScore(s, l, bal)
	}
}
