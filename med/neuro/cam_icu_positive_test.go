package neuro

import "testing"

// Pipanmekaporn 2014 / Kim 2026: features 1 and 2 plus 3 or 4; all sixteen combinations against the rule.
func TestCAMICUPositive_example(t *testing.T) {
	for m := 0; m < 16; m++ {
		f1, f2, f3, f4 := m&1 != 0, m&2 != 0, m&4 != 0, m&8 != 0
		want := f1 && f2 && (f3 || f4)
		if got := CAMICUPositive(f1, f2, f3, f4); got != want {
			t.Errorf("CAMICUPositive(%v, %v, %v, %v) = %v, want %v", f1, f2, f3, f4, got, want)
		}
	}
	if !CAMICUPositive(true, true, false, true) || CAMICUPositive(true, true, false, false) || CAMICUPositive(false, true, true, true) {
		t.Error("documented examples")
	}
}

// parity: ledger cam_icu, source cases 1–2 (consumer: f1 && f2 && (f3 || f4) once assessable; f2 from errors ≥ 3).
func TestCAMICUPositive_parity(t *testing.T) {
	if !CAMICUPositive(true, true, true, false) { // parity: ledger cam_icu, source case 1
		t.Error("case 1")
	}
	if CAMICUPositive(true, false, true, true) { // parity: ledger cam_icu, source case 2
		t.Error("case 2")
	}
}

func TestCAMICUPositive_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkBool = CAMICUPositive(true, true, false, true) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkCAMICUPositive(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkBool = CAMICUPositive(true, true, false, true)
	}
}
