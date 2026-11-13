package cardio

import "testing"

var sinkBool bool

// Urban 2019: "at least 1 major or 2 minor criteria"; both sides of each edge.
func TestARCHBRHighBleedingRisk(t *testing.T) {
	tests := []struct {
		major, minor int
		want         bool
	}{
		{0, 0, false}, {0, 1, false}, {0, 2, true}, {1, 0, true}, {1, 1, true}, {0, 6, true}, {3, 4, true},
	}
	for _, tc := range tests {
		if got, err := ARCHBRHighBleedingRisk(tc.major, tc.minor); err != nil || got != tc.want {
			t.Errorf("(%d, %d) = (%v, %v), want %v", tc.major, tc.minor, got, err, tc.want)
		}
	}
	for _, in := range [][2]int{{-1, 0}, {0, -1}} {
		_, err := ARCHBRHighBleedingRisk(in[0], in[1])
		checkErr(t, "negative", err, ErrOutOfRange)
	}
}

// parity: ledger arc_hbr — no source test file or JSON example; the app's rule (major ≥ 1 or minor ≥ 2) matches.
func TestARCHBRHighBleedingRisk_parity(t *testing.T) {
	if got, _ := ARCHBRHighBleedingRisk(0, 2); !got { // parity: ledger arc_hbr, source case 1
		t.Error("0 major, 2 minor: want true")
	}
}

func TestARCHBRHighBleedingRisk_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkBool, _ = ARCHBRHighBleedingRisk(0, 2) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkARCHBRHighBleedingRisk(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkBool, _ = ARCHBRHighBleedingRisk(0, 2)
	}
}
