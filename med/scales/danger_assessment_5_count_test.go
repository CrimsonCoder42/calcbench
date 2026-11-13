package scales

import "testing"

// Snider 2009 abstract: five yes/no questions; the count of yes answers, 0–5. All 32 combinations.
func TestDangerAssessment5Count_example(t *testing.T) {
	for mask := 0; mask < 32; mask++ {
		want := 0
		for m := mask; m != 0; m &= m - 1 {
			want++
		}
		got := DangerAssessment5Count(mask&1 != 0, mask&2 != 0, mask&4 != 0, mask&8 != 0, mask&16 != 0)
		if got != want {
			t.Errorf("DangerAssessment5Count(mask %05b) = %d, want %d", mask, got, want)
		}
	}
}

// parity: ledger danger_assessment_domestic, source cases 1–5 (the first k answers yes).
func TestDangerAssessment5Count_parity(t *testing.T) {
	for i, k := range []int{0, 1, 2, 3, 5} { // source cases 1–5: zero, one, two, three edge, five
		var a [5]bool
		for j := 0; j < k; j++ {
			a[j] = true
		}
		if got := DangerAssessment5Count(a[0], a[1], a[2], a[3], a[4]); got != k {
			t.Errorf("source case %d: DangerAssessment5Count = %d, want %d", i+1, got, k)
		}
	}
}

func TestDangerAssessment5Count_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt = DangerAssessment5Count(true, false, true, false, true) }); n != 0 {
		t.Errorf("DangerAssessment5Count allocates %v times, want 0", n)
	}
}

func BenchmarkDangerAssessment5Count(b *testing.B) {
	b.ReportAllocs()
	x := true
	for b.Loop() {
		sinkInt = DangerAssessment5Count(x, !x, x, !x, x)
	}
}
