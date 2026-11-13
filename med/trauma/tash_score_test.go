package trauma

import "testing"

// Yücel 2006 abstract weights; every band edge on both sides. HAND-DERIVED.
func TestTASHScore_bands(t *testing.T) {
	base := func(hb, be, sbp, hr float64) int {
		s, err := TASHScore(false, hb, be, sbp, hr, false, 0)
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	tests := []struct {
		name            string
		hb, be, sbp, hr float64
		want            int
	}{
		{"all zero", 14, 0, 130, 80, 0},
		{"SBP 119.9", 14, 0, 119.9, 80, 1}, {"SBP 120", 14, 0, 120, 80, 0},
		{"SBP 99.9", 14, 0, 99.9, 80, 4}, {"SBP 100", 14, 0, 100, 80, 1},
		{"Hb 11.9", 11.9, 0, 130, 80, 2}, {"Hb 12", 12, 0, 130, 80, 0},
		{"Hb 10.9", 10.9, 0, 130, 80, 3}, {"Hb 11", 11, 0, 130, 80, 2},
		{"Hb 9.9", 9.9, 0, 130, 80, 4}, {"Hb 10", 10, 0, 130, 80, 3},
		{"Hb 8.9", 8.9, 0, 130, 80, 6}, {"Hb 9", 9, 0, 130, 80, 4},
		{"Hb 6.9", 6.9, 0, 130, 80, 8}, {"Hb 7", 7, 0, 130, 80, 6},
		{"HR 120", 14, 0, 130, 120, 0}, {"HR 120.1", 14, 0, 130, 120.1, 2},
		{"BE -2", 14, -2, 130, 80, 0}, {"BE -2.1", 14, -2.1, 130, 80, 1},
		{"BE -6", 14, -6, 130, 80, 1}, {"BE -6.1", 14, -6.1, 130, 80, 3},
		{"BE -10", 14, -10, 130, 80, 3}, {"BE -10.1", 14, -10.1, 130, 80, 4},
	}
	for _, tc := range tests {
		if got := base(tc.hb, tc.be, tc.sbp, tc.hr); got != tc.want {
			t.Errorf("%s: %d, want %d", tc.name, got, tc.want)
		}
	}
	for ais, want := range []int{0, 0, 0, 3, 3, 6} {
		if got, _ := TASHScore(false, 14, 0, 130, 80, false, ais); got != want {
			t.Errorf("AIS %d → %d, want %d", ais, got, want)
		}
	}
	if got, _ := TASHScore(true, 14, 0, 130, 80, true, 0); got != 4 {
		t.Errorf("male + fluid → %d, want 4", got)
	}
	// Maximum: 1 + 8 + 4 + 4 + 2 + 3 + 6 = 28.
	if got, _ := TASHScore(true, 6.5, -11, 95, 130, true, 5); got != 28 {
		t.Errorf("max → %d, want 28", got)
	}
}

// parity: ledger tash_score, source case 1. The app scores pelvic fracture 6 and femur fracture 3 as separate items
// (sum up to 9); the primary has one item, complex long bone and/or pelvic fractures, AIS 3/4 = 3, AIS 5 = 6. An
// app case with both fractures (9 points) is, by the primary, the worst fracture's AIS only: AIS 5 → 6.
func TestTASHScore_parity(t *testing.T) {
	if got, _ := TASHScore(false, 14, 0, 130, 80, false, 5); got != 6 {
		t.Errorf("got %d, want 6 (primary; app 9 for pelvic + femur)", got)
	}
}

func TestTASHScore_errors(t *testing.T) {
	for _, v := range nonFinite {
		for i := 0; i < 4; i++ {
			in := [4]float64{14, 0, 130, 80}
			in[i] = v
			_, err := TASHScore(false, in[0], in[1], in[2], in[3], false, 0)
			checkErr(t, "non-finite", err, ErrNotFinite)
		}
	}
	for _, in := range [][4]float64{{-1, 0, 130, 80}, {14, 0, -1, 80}, {14, 0, 130, -1}} {
		_, err := TASHScore(false, in[0], in[1], in[2], in[3], false, 0)
		checkErr(t, "range", err, ErrOutOfRange)
	}
	for _, a := range []int{-1, 6} {
		_, err := TASHScore(false, 14, 0, 130, 80, false, a)
		checkErr(t, "AIS", err, ErrOutOfRange)
	}
}

func TestTASHScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = TASHScore(true, 8, -7, 110, 125, true, 3) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkTASHScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = TASHScore(true, 8, -7, 110, 125, true, 3)
	}
}
