package psychiatry

import "testing"

// Boudreaux 2013: positive for active ideation or an attempt within the past 6 months; all four combinations.
func TestPSS3Positive_example(t *testing.T) {
	cases := []struct{ ideation, attempt, want bool }{{false, false, false}, {true, false, true}, {false, true, true}, {true, true, true}}
	for _, c := range cases {
		if got := PSS3Positive(c.ideation, c.attempt); got != c.want {
			t.Errorf("PSS3Positive(%v, %v) = %v, want %v", c.ideation, c.attempt, got, c.want)
		}
	}
}

// parity: ledger ed_safe_pss3, source cases 1–2. The consumer counts any lifetime attempt as positive and adds
// an "Indeterminate" result for refusals; the source's rule uses an attempt within 6 months. Case 2 asserts the
// source's side (a lifetime attempt 2 years ago, no ideation → not positive; consumer: Positive).
func TestPSS3Positive_parity(t *testing.T) {
	if !PSS3Positive(true, false) { // parity: ledger ed_safe_pss3, source case 1
		t.Error("case 1")
	}
	if PSS3Positive(false, false) { // parity: ledger ed_safe_pss3, source case 2 (consumer value: Positive)
		t.Error("case 2")
	}
}

func TestPSS3Positive_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkBool = PSS3Positive(false, true) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkPSS3Positive(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkBool = PSS3Positive(false, true)
	}
}
