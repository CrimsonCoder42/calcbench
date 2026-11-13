package cardio

import "testing"

// Nordestgaard 2013, Table 1: highest item per group; HAND-DERIVED totals.
func TestDutchLipidClinicFHScore_items(t *testing.T) {
	tests := []struct {
		name string
		ldl  DLCNLDLBand
		f    DLCNFindings
		want int
	}{
		{"nothing", DLCNLDLBelow4, DLCNFindings{}, 0},
		{"example", DLCNLDL65To84, DLCNFindings{TendonXanthoma: true, RelativePrematureCHD: true}, 12},
		{"group 1 highest", DLCNLDL4To49, DLCNFindings{RelativePrematureCHD: true, ChildLDLAbove95th: true}, 3},
		{"group 2 highest", DLCNLDL5To64, DLCNFindings{PrematureCHD: true, PrematureCerebralOrPeripheral: true}, 5},
		{"group 3 arcus", DLCNLDLBelow4, DLCNFindings{CornealArcusUnder45: true}, 4},
		{"maximum 26", DLCNLDLAbove85, DLCNFindings{true, true, true, true, true, true, true, true, true}, 26},
	}
	for _, tc := range tests {
		if got, err := DutchLipidClinicFHScore(tc.ldl, tc.f); err != nil || got != tc.want {
			t.Errorf("%s: (%d, %v), want %d", tc.name, got, err, tc.want)
		}
	}
	for _, b := range []DLCNLDLBand{-1, 5} {
		_, err := DutchLipidClinicFHScore(b, DLCNFindings{})
		checkErr(t, "band", err, ErrOutOfRange)
	}
}

// Table 1 footnote: > 8 definite, 6–8 probable, 3–5 possible, 0–2 unlikely.
func TestDutchLipidClinicFHGroup_edges(t *testing.T) {
	tests := []struct {
		s    int
		want DLCNDiagnosis
	}{{0, DLCNUnlikely}, {2, DLCNUnlikely}, {3, DLCNPossible}, {5, DLCNPossible}, {6, DLCNProbable}, {8, DLCNProbable}, {9, DLCNDefinite}, {26, DLCNDefinite}}
	for _, tc := range tests {
		if got, err := DutchLipidClinicFHGroup(tc.s); err != nil || got != tc.want {
			t.Errorf("%d: (%q, %v), want %q", tc.s, got, err, tc.want)
		}
	}
	for _, s := range []int{-1, 27} {
		_, err := DutchLipidClinicFHGroup(s)
		checkErr(t, "range", err, ErrOutOfRange)
	}
}

func TestDutchLipidClinicFHScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sinkInt, _ = DutchLipidClinicFHScore(DLCNLDL65To84, DLCNFindings{})
		g, _ := DutchLipidClinicFHGroup(7)
		sinkStr = string(g)
	}); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkDutchLipidClinicFHScore(b *testing.B) {
	b.ReportAllocs()
	f := DLCNFindings{TendonXanthoma: true}
	for b.Loop() {
		sinkInt, _ = DutchLipidClinicFHScore(DLCNLDL65To84, f)
	}
}

func BenchmarkDutchLipidClinicFHGroup(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		g, _ := DutchLipidClinicFHGroup(7)
		sinkStr = string(g)
	}
}
