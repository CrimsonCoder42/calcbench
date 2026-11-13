package hematology

import "testing"

var sinkBoolHemat bool

// Arber 2022, Table 11: the major criterion, or at least 3 of the 4 minor criteria; tryptase does not count
// in SM-AMN.
func TestICCSystemicMastocytosis_table(t *testing.T) {
	tests := []struct {
		name                          string
		major, morph, mark, kit, tryp bool
		amn                           bool
		want                          bool
	}{
		{"major alone", true, false, false, false, false, false, true},
		{"none", false, false, false, false, false, false, false},
		{"two minor", false, true, true, false, false, false, false},
		{"three minor", false, true, true, true, false, false, true},
		{"three minor with tryptase", false, true, false, true, true, false, true},
		{"four minor", false, true, true, true, true, false, true},
		{"SM-AMN: tryptase not counted", false, true, true, false, true, true, false},
		{"SM-AMN: three other minor", false, true, true, true, true, true, true},
		{"SM-AMN with major", true, false, false, false, true, true, true},
	}
	for _, tc := range tests {
		if got := ICCSystemicMastocytosis(tc.major, tc.morph, tc.mark, tc.kit, tc.tryp, tc.amn); got != tc.want {
			t.Errorf("%s: got %v, want %v", tc.name, got, tc.want)
		}
	}
}

// parity: ledger icc_systemic_mastocytosis. The source entry has no _test.go and its JSON no example; the
// consumer's rule is "major || minorCount >= 3" (it has no SM-AMN input), which agrees with Table 11.
func TestICCSystemicMastocytosis_parity(t *testing.T) {
	if !ICCSystemicMastocytosis(true, false, false, false, false, false) {
		t.Error("source case 1: major only, want true")
	}
	if ICCSystemicMastocytosis(false, true, false, false, true, false) {
		t.Error("source case 2: two minor, want false")
	}
}

func TestICCSystemicMastocytosis_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sinkBoolHemat = ICCSystemicMastocytosis(false, true, true, true, false, false)
	}); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkICCSystemicMastocytosis(b *testing.B) {
	b.ReportAllocs()
	x := true
	for b.Loop() {
		sinkBoolHemat = ICCSystemicMastocytosis(!x, x, x, x, !x, !x)
	}
}
