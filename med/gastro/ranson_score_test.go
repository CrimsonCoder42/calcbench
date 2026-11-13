package gastro

import "testing"

// Lee 2022 Table 4 / Abu-Eshy 2008 Table 1: one point per criterion met. Exhaustive over all 2^11 combinations
// against a popcount (2048 cases, well under a second).
func TestRansonScore_example(t *testing.T) {
	if got := RansonScore(RansonFindings{AgeOver55Years: true, WBCOver16000PerMm3: true, GlucoseOver200MgDL: true, CalciumBelow8MgDL: true}); got != 4 {
		t.Errorf("worked: got %v, want 4", got)
	}
	for mask := 0; mask < 1<<11; mask++ {
		bit := func(i int) bool { return mask&(1<<i) != 0 }
		f := RansonFindings{bit(0), bit(1), bit(2), bit(3), bit(4), bit(5), bit(6), bit(7), bit(8), bit(9), bit(10)}
		want := 0
		for i := 0; i < 11; i++ {
			if bit(i) {
				want++
			}
		}
		if got := RansonScore(f); got != want {
			t.Fatalf("mask %011b: got %v, want %v", mask, got, want)
		}
	}
}

// parity: ledger ranson_pancreatitis. The app sums eleven 0/1 items (admission 5 + 48 h 6): defaults all 0 → 0;
// all 1 → 11; the five admission items only → 5.
func TestRansonScore_parity(t *testing.T) {
	all := RansonFindings{true, true, true, true, true, true, true, true, true, true, true}
	adm := RansonFindings{AgeOver55Years: true, WBCOver16000PerMm3: true, GlucoseOver200MgDL: true, LDHOver350UL: true, ASTOver250UL: true}
	cases := []struct {
		f    RansonFindings
		want int
	}{{RansonFindings{}, 0}, {all, 11}, {adm, 5}}
	for i, c := range cases {
		if got := RansonScore(c.f); got != c.want { // parity: ledger ranson_pancreatitis, source case i+1
			t.Errorf("case %d: got %v, want %v", i+1, got, c.want)
		}
	}
}

func TestRansonScore_noAllocation(t *testing.T) {
	f := RansonFindings{AgeOver55Years: true, CalciumBelow8MgDL: true}
	if n := testing.AllocsPerRun(100, func() { sinkInt = RansonScore(f) }); n != 0 {
		t.Errorf("RansonScore allocates %v times, want 0", n)
	}
}

func BenchmarkRansonScore(b *testing.B) {
	b.ReportAllocs()
	f := RansonFindings{AgeOver55Years: true, CalciumBelow8MgDL: true}
	for b.Loop() {
		sinkInt = RansonScore(f)
	}
}
