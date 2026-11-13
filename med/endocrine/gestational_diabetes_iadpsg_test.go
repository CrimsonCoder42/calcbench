package endocrine

import "testing"

// IADPSG 2010 Table 1: 92 / 180 / 153 mg/dL and 5.1 / 10.0 / 8.5 mmol/L, "equaled or exceeded". Each edge is
// tested on both sides with the other two values below their thresholds.
func TestGDMThresholdIADPSGMgDL_edges(t *testing.T) {
	tests := []struct {
		name        string
		f, one, two float64
		want        bool
	}{
		{"all below", 91, 179, 152, false},
		{"fasting equals", 92, 179, 152, true},
		{"fasting just below", 91.9, 100, 100, false},
		{"1 h equals", 80, 180, 100, true},
		{"1 h just below", 80, 179.9, 100, false},
		{"2 h equals", 80, 100, 153, true},
		{"2 h just below", 80, 100, 152.9, false},
		{"all above", 120, 250, 220, true},
	}
	for _, tc := range tests {
		if got, err := GDMThresholdIADPSGMgDL(tc.f, tc.one, tc.two); err != nil || got != tc.want {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger bs_gdm, the app's rule (value >= 92 / 180 / 153, any one met). The app also counts a 0 as
// "not entered"; here every value is required (ErrOutOfRange for 0), noted in the ledger.
func TestGDMThresholdIADPSGMgDL_parity(t *testing.T) {
	tests := []struct {
		f, one, two float64
		want        bool
	}{
		{85, 160, 140, false},
		{95, 160, 140, true},
		{85, 185, 140, true},
		{85, 160, 155, true},
	}
	for _, tc := range tests {
		if got, err := GDMThresholdIADPSGMgDL(tc.f, tc.one, tc.two); err != nil || got != tc.want {
			t.Errorf("(%v, %v, %v): got (%v, %v), want %v", tc.f, tc.one, tc.two, got, err, tc.want)
		}
	}
}

func TestGDMThresholdIADPSGMmolL_edges(t *testing.T) {
	tests := []struct {
		name        string
		f, one, two float64
		want        bool
	}{
		{"all below", 5.0, 9.9, 8.4, false},
		{"fasting equals", 5.1, 8, 7, true},
		{"1 h equals", 4.5, 10.0, 7, true},
		{"2 h equals", 4.5, 8, 8.5, true},
		{"2 h just below", 4.5, 8, 8.49, false},
	}
	for _, tc := range tests {
		if got, err := GDMThresholdIADPSGMmolL(tc.f, tc.one, tc.two); err != nil || got != tc.want {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

func TestGDMThresholdIADPSG_errors(t *testing.T) {
	mg := func(a, b, c float64) func() error {
		return func() error { _, err := GDMThresholdIADPSGMgDL(a, b, c); return err }
	}
	si := func(a, b, c float64) func() error {
		return func() error { _, err := GDMThresholdIADPSGMmolL(a, b, c); return err }
	}
	runErrCases(t, []errCase{
		{"mg fasting NaN", mg(nan, 100, 100), ErrNotFinite},
		{"mg 1h +Inf", mg(80, pInf, 100), ErrNotFinite},
		{"mg 2h -Inf", mg(80, 100, nInf), ErrNotFinite},
		{"mg NaN before range", mg(0, 100, nan), ErrNotFinite},
		{"mg fasting 0", mg(0, 100, 100), ErrOutOfRange},
		{"mg 1h negative", mg(80, -1, 100), ErrOutOfRange},
		{"mg 2h 0", mg(80, 100, 0), ErrOutOfRange},
		{"si fasting NaN", si(nan, 8, 7), ErrNotFinite},
		{"si 2h 0", si(4.5, 8, 0), ErrOutOfRange},
	})
}

func TestGDMThresholdIADPSG_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkBool, _ = GDMThresholdIADPSGMgDL(90, 170, 150) }); n != 0 {
		t.Errorf("GDMThresholdIADPSGMgDL allocates %v times, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkBool, _ = GDMThresholdIADPSGMmolL(5, 9, 8) }); n != 0 {
		t.Errorf("GDMThresholdIADPSGMmolL allocates %v times, want 0", n)
	}
}

func BenchmarkGDMThresholdIADPSGMgDL(b *testing.B) {
	b.ReportAllocs()
	f, o, w := 90.0, 170.0, 150.0
	for b.Loop() {
		sinkBool, _ = GDMThresholdIADPSGMgDL(f, o, w)
	}
}

func BenchmarkGDMThresholdIADPSGMmolL(b *testing.B) {
	b.ReportAllocs()
	f, o, w := 5.0, 9.0, 8.0
	for b.Loop() {
		sinkBool, _ = GDMThresholdIADPSGMmolL(f, o, w)
	}
}
