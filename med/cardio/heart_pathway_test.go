package cardio

import "testing"

// Mahler 2015: HEART 0–3 and negative troponin at 0 and 3 hours.
func TestHEARTPathwayLowRisk(t *testing.T) {
	tests := []struct {
		s      int
		t0, t3 bool
		want   bool
	}{
		{3, true, true, true}, {4, true, true, false}, {0, true, true, true},
		{3, false, true, false}, {3, true, false, false}, {10, true, true, false},
	}
	for _, tc := range tests {
		if got, err := HEARTPathwayLowRisk(tc.s, tc.t0, tc.t3); err != nil || got != tc.want {
			t.Errorf("%+v = (%v, %v), want %v", tc, got, err, tc.want)
		}
	}
	for _, s := range []int{-1, 11} {
		_, err := HEARTPathwayLowRisk(s, true, true)
		checkErr(t, "range", err, ErrOutOfRange)
	}
}

// parity: ledger heart_pathway — no source test file or JSON example. The app tests only one troponin (the
// score's troponin item = 0); here both serial troponins are inputs. Source case 1: score 2, negative → low risk.
func TestHEARTPathwayLowRisk_parity(t *testing.T) {
	if got, _ := HEARTPathwayLowRisk(2, true, true); !got { // parity: ledger heart_pathway, source case 1
		t.Error("want low risk")
	}
}

func TestHEARTPathwayLowRisk_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkBool, _ = HEARTPathwayLowRisk(3, true, true) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkHEARTPathwayLowRisk(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkBool, _ = HEARTPathwayLowRisk(3, true, true)
	}
}
