package infectious

import "testing"

var sinkInt int

// Centor 1981 abstract: four variables, each counted once; 0–4.
func TestCentorScore_example(t *testing.T) {
	for mask := 0; mask < 16; mask++ {
		e, n, c, f := mask&1 != 0, mask&2 != 0, mask&4 != 0, mask&8 != 0
		want := 0
		for m := mask; m != 0; m &= m - 1 {
			want++
		}
		if got := CentorScore(e, n, c, f); got != want {
			t.Errorf("CentorScore(%v, %v, %v, %v) = %d, want %d", e, n, c, f, got, want)
		}
	}
}

// parity: ledger centor_score, source cases 1–7. The source computes the McIsaac score (adds +1 for age 3–14 and
// −1 for age ≥ 45) while citing Centor 1981, which has no age term; each case asserts the Centor count, with the
// source's value in the comment. Inputs: exudates, nodes, no cough, fever.
func TestCentorScore_parity(t *testing.T) {
	tests := []struct {
		name          string
		ex, ln, nc, f bool
		want          int
	}{
		{"child all (age 10; source 5)", true, true, true, true, 4},
		{"age 14 (source 3)", true, false, false, true, 2},
		{"age 15 (source 2)", false, false, true, true, 2},
		{"age 44 (source 3)", true, false, true, true, 3},
		{"age 45 none (source -1)", false, false, false, false, 0},
		{"age 60 all (source 3)", true, true, true, true, 4},
		{"adult one (age 30; source 1)", false, false, false, true, 1},
	}
	for _, tc := range tests {
		if got := CentorScore(tc.ex, tc.ln, tc.nc, tc.f); got != tc.want {
			t.Errorf("%s: CentorScore = %d, want %d", tc.name, got, tc.want)
		}
	}
}

func TestCentorScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt = CentorScore(true, false, true, true) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkCentorScore(b *testing.B) {
	b.ReportAllocs()
	x := true
	for b.Loop() {
		sinkInt = CentorScore(x, !x, x, x)
	}
}
