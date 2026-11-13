package neuro

import "testing"

// Kattah 2009 abstract: any of three findings; all eight combinations, HAND-DERIVED.
func TestHINTSAnyFinding_example(t *testing.T) {
	for m := 0; m < 8; m++ {
		a, b, c := m&1 != 0, m&2 != 0, m&4 != 0
		if got := HINTSAnyFinding(a, b, c); got != (m != 0) {
			t.Errorf("HINTSAnyFinding(%v, %v, %v) = %v", a, b, c, got)
		}
	}
}

// parity: ledger hints_stroke, source cases 1–2 (consumer: benign only when all three are peripheral).
func TestHINTSAnyFinding_parity(t *testing.T) {
	if HINTSAnyFinding(false, false, false) { // parity: ledger hints_stroke, source case 1
		t.Error("case 1")
	}
	if !HINTSAnyFinding(true, false, false) { // parity: ledger hints_stroke, source case 2
		t.Error("case 2")
	}
}

func TestHINTSAnyFinding_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkBool = HINTSAnyFinding(false, true, false) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkHINTSAnyFinding(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkBool = HINTSAnyFinding(false, true, false)
	}
}
