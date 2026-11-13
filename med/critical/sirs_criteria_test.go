package critical

import "testing"

var sinkSIRS SIRSResult

// Singer 2016 Box 1 / Balk 2014: both sides of every criterion edge. HAND-DERIVED.
func TestSIRSCriteria_edges(t *testing.T) {
	type in struct {
		temp, hr, rr, paco2 float64
		paco2OK             bool
		wbc, bands          float64
	}
	base := in{37, 80, 16, 40, true, 8000, 0}
	tests := []struct {
		name string
		mod  func(*in)
		want int
	}{
		{"none", func(*in) {}, 0},
		{"temp 38.0", func(v *in) { v.temp = 38.0 }, 0},
		{"temp 38.01", func(v *in) { v.temp = 38.01 }, 1},
		{"temp 36.0", func(v *in) { v.temp = 36.0 }, 0},
		{"temp 35.99", func(v *in) { v.temp = 35.99 }, 1},
		{"hr 90", func(v *in) { v.hr = 90 }, 0},
		{"hr 91", func(v *in) { v.hr = 91 }, 1},
		{"rr 20", func(v *in) { v.rr = 20 }, 0},
		{"rr 21", func(v *in) { v.rr = 21 }, 1},
		{"paco2 32", func(v *in) { v.paco2 = 32 }, 0},
		{"paco2 31.9", func(v *in) { v.paco2 = 31.9 }, 1},
		{"paco2 31.9 not measured", func(v *in) { v.paco2 = 31.9; v.paco2OK = false }, 0},
		{"rr and paco2 count once", func(v *in) { v.rr = 25; v.paco2 = 28 }, 1},
		{"wbc 12000", func(v *in) { v.wbc = 12000 }, 0},
		{"wbc 12001", func(v *in) { v.wbc = 12001 }, 1},
		{"wbc 4000", func(v *in) { v.wbc = 4000 }, 0},
		{"wbc 3999", func(v *in) { v.wbc = 3999 }, 1},
		{"bands 10", func(v *in) { v.bands = 10 }, 0},
		{"bands 10.1", func(v *in) { v.bands = 10.1 }, 1},
		{"all four", func(v *in) { *v = in{39, 120, 28, 30, true, 15000, 15} }, 4},
	}
	for _, tc := range tests {
		v := base
		tc.mod(&v)
		r, err := SIRSCriteria(v.temp, v.hr, v.rr, v.paco2, v.paco2OK, v.wbc, v.bands)
		if err != nil || r.Count != tc.want || r.SIRS != (tc.want >= 2) {
			t.Errorf("%s: (%+v, %v), want count %d", tc.name, r, err, tc.want)
		}
	}
}

// Doc comment example: 38.5 °C and HR 95 → 2, SIRS (Box 1 "Two or more of").
func TestSIRSCriteria_example(t *testing.T) {
	r, err := SIRSCriteria(38.5, 95, 18, 0, false, 9000, 0)
	if err != nil || r.Count != 2 || !r.SIRS || !r.Temperature || !r.HeartRate || r.Respiratory || r.WhiteCells {
		t.Errorf("got (%+v, %v)", r, err)
	}
}

// parity: ledger sirs (source case 1) and sirs_calc (source case 1). sirs: the app scores temperature, HR, RR and
// WBC only (no PaCO2, no bands); 101 °F (38.33 °C), HR 110, RR 22, WBC 13 ×10³/µL → 4. sirs_calc: RR 18,
// PaCO2 30, 37 °C, HR 85, WBC 8000, bands 12 → 2 (PaCO2 and bands criteria).
func TestSIRSCriteria_parity(t *testing.T) {
	if r, _ := SIRSCriteria(38.33, 110, 22, 0, false, 13000, 0); r.Count != 4 {
		t.Errorf("sirs case: count %d, want 4", r.Count)
	}
	if r, _ := SIRSCriteria(37, 85, 18, 30, true, 8000, 12); r.Count != 2 || !r.SIRS {
		t.Errorf("sirs_calc case: %+v, want count 2", r)
	}
}

func TestSIRSCriteria_errors(t *testing.T) {
	for _, v := range nonFiniteB6 {
		_, err := SIRSCriteria(v, 80, 16, 40, true, 8000, 0)
		wantErrB6(t, "temp", err, ErrNotFinite)
		_, err = SIRSCriteria(37, v, 16, 40, true, 8000, 0)
		wantErrB6(t, "hr", err, ErrNotFinite)
		_, err = SIRSCriteria(37, 80, v, 40, true, 8000, 0)
		wantErrB6(t, "rr", err, ErrNotFinite)
		_, err = SIRSCriteria(37, 80, 16, v, false, 8000, 0)
		wantErrB6(t, "paco2", err, ErrNotFinite)
		_, err = SIRSCriteria(37, 80, 16, 40, true, v, 0)
		wantErrB6(t, "wbc", err, ErrNotFinite)
		_, err = SIRSCriteria(37, 80, 16, 40, true, 8000, v)
		wantErrB6(t, "bands", err, ErrNotFinite)
	}
	_, err := SIRSCriteria(9.9, 80, 16, 40, true, 8000, 0)
	wantErrB6(t, "temp 9.9", err, ErrOutOfRange)
	_, err = SIRSCriteria(50.1, 80, 16, 40, true, 8000, 0)
	wantErrB6(t, "temp 50.1", err, ErrOutOfRange)
	_, err = SIRSCriteria(37, -1, 16, 40, true, 8000, 0)
	wantErrB6(t, "hr -1", err, ErrOutOfRange)
	_, err = SIRSCriteria(37, 80, -1, 40, true, 8000, 0)
	wantErrB6(t, "rr -1", err, ErrOutOfRange)
	_, err = SIRSCriteria(37, 80, 16, -1, true, 8000, 0)
	wantErrB6(t, "paco2 -1", err, ErrOutOfRange)
	_, err = SIRSCriteria(37, 80, 16, 40, true, -1, 0)
	wantErrB6(t, "wbc -1", err, ErrOutOfRange)
	_, err = SIRSCriteria(37, 80, 16, 40, true, 8000, 100.1)
	wantErrB6(t, "bands 100.1", err, ErrOutOfRange)
}

func TestSIRSCriteria_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkSIRS, _ = SIRSCriteria(38.5, 95, 18, 30, true, 9000, 0) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkSIRSCriteria(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkSIRS, _ = SIRSCriteria(38.5, 95, 18, 30, true, 9000, 0)
	}
}
