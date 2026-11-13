package gastro

import "testing"

// HAND-DERIVED from Chong 2010 (abstract). Half-points are exact in float64, so comparisons are exact.
func TestRIPASAScore_example(t *testing.T) {
	all := RIPASAFindings{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true}
	tests := []struct {
		name string
		f    RIPASAFindings
		want float64
	}{
		{"minimum", RIPASAFindings{}, 1.5},
		{"maximum", all, 16},
		{"worked 7.5", RIPASAFindings{AgeUnder40: true, Male: true, SymptomsUnder48Hours: true, RightIliacFossaPain: true,
			RightIliacFossaTenderness: true, Guarding: true, RaisedWhiteCellCount: true}, 7.5},
		{"just below 7.5", RIPASAFindings{Male: true, SymptomsUnder48Hours: true, RightIliacFossaPain: true,
			RightIliacFossaTenderness: true, Guarding: true, RaisedWhiteCellCount: true}, 7},
		{"Rovsing", RIPASAFindings{RovsingSign: true}, 3.5},
		{"foreign ID", RIPASAFindings{ForeignNationalIDCard: true}, 2.5},
	}
	for _, tc := range tests {
		if got := RIPASAScore(tc.f); got != tc.want {
			t.Errorf("%s: got %v, want %v", tc.name, got, tc.want)
		}
	}
}

// parity: ledger ripasa, the app's point arrays (SEX [0.5, 1], AGE [1, 0.5], DURATION [1, 0.5], RIF pain 0.5,
// migration 0.5, guarding 2, Rovsing 2, others 1) agree with the abstract; all-"No" picker values → 1.5.
func TestRIPASAScore_parity(t *testing.T) {
	if got := RIPASAScore(RIPASAFindings{}); got != 1.5 {
		t.Errorf("all no: got %v, want 1.5", got)
	}
	f := RIPASAFindings{Male: true, AgeUnder40: true, NauseaAndVomiting: true, Anorexia: true, ReboundTenderness: true}
	if got := RIPASAScore(f); got != 5.5 {
		t.Errorf("got %v, want 5.5", got)
	}
}

func TestRIPASAScore_noAllocation(t *testing.T) {
	f := RIPASAFindings{Male: true, Guarding: true}
	if n := testing.AllocsPerRun(100, func() { sinkFloat = RIPASAScore(f) }); n != 0 {
		t.Errorf("RIPASAScore allocates %v times, want 0", n)
	}
}

func BenchmarkRIPASAScore(b *testing.B) {
	b.ReportAllocs()
	f := RIPASAFindings{Male: true, Guarding: true, Fever: true}
	for b.Loop() {
		sinkFloat = RIPASAScore(f)
	}
}
