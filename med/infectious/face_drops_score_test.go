package infectious

import "testing"

// McEntire 2025, Table 3: F +8, A +6, C +3, E +4, D +3, otalgia −1, S +3.
func TestFACEDROPSScore_table(t *testing.T) {
	tests := []struct {
		name                string
		f, a, c, e, d, o, s bool
		want                int
	}{
		{"none", false, false, false, false, false, false, false, 0},
		{"fever", true, false, false, false, false, false, false, 8},
		{"aches", false, true, false, false, false, false, false, 6},
		{"cephalalgia", false, false, true, false, false, false, false, 3},
		{"exhaustion", false, false, false, true, false, false, false, 4},
		{"dermatomal", false, false, false, false, true, false, false, 3},
		{"otalgia alone is the minimum", false, false, false, false, false, true, false, -1},
		{"stiff neck", false, false, false, false, false, false, true, 3},
		{"all positive items is the maximum", true, true, true, true, true, false, true, 27},
		{"edge ≤4: exhaustion alone is 4", false, false, false, true, false, false, false, 4},
		{"edge 5: cephalalgia + stiff neck − otalgia", false, false, true, false, false, true, true, 5},
		{"edge 6: cephalalgia + exhaustion − otalgia", false, false, true, true, false, true, false, 6},
		{"edge ≥7: cephalalgia + exhaustion", false, false, true, true, false, false, false, 7},
	}
	for _, tc := range tests {
		if got := FACEDROPSScore(tc.f, tc.a, tc.c, tc.e, tc.d, tc.o, tc.s); got != tc.want {
			t.Errorf("%s: got %d, want %d", tc.name, got, tc.want)
		}
	}
}

// parity: ledger face_drops. The source entry has no _test.go and its JSON no example; case 1 is the consumer's
// worked example: fever (8), aches (6), stiff neck (3), otalgia (−1) → 16.
func TestFACEDROPSScore_parity(t *testing.T) {
	if got := FACEDROPSScore(true, true, false, false, false, true, true); got != 16 {
		t.Errorf("source case 1: got %d, want 16", got)
	}
}

func TestFACEDROPSScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt = FACEDROPSScore(true, false, true, false, true, false, true) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkFACEDROPSScore(b *testing.B) {
	b.ReportAllocs()
	x := true
	for b.Loop() {
		sinkInt = FACEDROPSScore(x, !x, x, !x, x, !x, x)
	}
}
