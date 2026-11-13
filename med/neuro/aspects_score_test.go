package neuro

import "testing"

// Teekaput 2024 / Nukovic 2023: 10 minus one per affected region.
func TestASPECTSScore_example(t *testing.T) {
	var none, all [10]bool
	for i := range all {
		all[i] = true
	}
	if got := ASPECTSScore(none); got != 10 {
		t.Errorf("none = %d, want 10", got)
	}
	if got := ASPECTSScore(all); got != 0 {
		t.Errorf("all = %d, want 0", got)
	}
	two := [10]bool{0: true, 4: true}
	if got := ASPECTSScore(two); got != 8 {
		t.Errorf("two = %d, want 8", got)
	}
}

// parity: ledger aspects_stroke, source cases 1–2 (consumer: 10 − affected count).
func TestASPECTSScore_parity(t *testing.T) {
	if got := ASPECTSScore([10]bool{1: true, 2: true, 3: true}); got != 7 { // parity: ledger aspects_stroke, source case 1
		t.Errorf("case 1 = %d, want 7", got)
	}
	if got := ASPECTSScore([10]bool{}); got != 10 { // parity: ledger aspects_stroke, source case 2
		t.Errorf("case 2 = %d, want 10", got)
	}
}

func TestASPECTSScore_noAlloc(t *testing.T) {
	r := [10]bool{0: true}
	if n := testing.AllocsPerRun(100, func() { sinkInt = ASPECTSScore(r) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkASPECTSScore(b *testing.B) {
	b.ReportAllocs()
	r := [10]bool{0: true}
	for b.Loop() {
		sinkInt = ASPECTSScore(r)
	}
}
