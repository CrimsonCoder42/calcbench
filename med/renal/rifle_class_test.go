package renal

import (
	"math"
	"testing"
)

var sinkRIFLE RIFLEClass

// Hoste 2006 Table 1 / Wlodzimirow 2012 / Bellomo 2004 legend: both sides of each edge. HAND-DERIVED.
func TestRIFLEClassify_edges(t *testing.T) {
	tests := []struct {
		name     string
		base, cr float64
		urine    RIFLEUrineOutput
		rrt, esd bool
		want     RIFLEClass
	}{
		{"none", 1.0, 1.0, RIFLEUrineNone, false, false, RIFLENone},
		{"1.49x", 1.0, 1.49, RIFLEUrineNone, false, false, RIFLENone},
		{"1.5x", 1.0, 1.5, RIFLEUrineNone, false, false, RIFLERisk},
		{"1.99x", 1.0, 1.99, RIFLEUrineNone, false, false, RIFLERisk},
		{"2x", 1.0, 2.0, RIFLEUrineNone, false, false, RIFLEInjury},
		{"2.99x", 1.0, 2.99, RIFLEUrineNone, false, false, RIFLEInjury},
		{"3x", 1.0, 3.0, RIFLEUrineNone, false, false, RIFLEFailure},
		{"4.0 rise 0.5 not above 4.0 (1.14x, below Risk)", 3.5, 4.0, RIFLEUrineNone, false, false, RIFLENone},
		{"4.1 rise 0.5", 3.6, 4.1, RIFLEUrineNone, false, false, RIFLEFailure},
		{"4.1 rise 0.4", 3.7, 4.1, RIFLEUrineNone, false, false, RIFLENone},
		{"urine 6h", 1.0, 1.0, RIFLEUrineBelow05For6h, false, false, RIFLERisk},
		{"urine 12h", 1.0, 1.0, RIFLEUrineBelow05For12h, false, false, RIFLEInjury},
		{"urine 24h", 1.0, 1.0, RIFLEUrineBelow03For24hOrAnuria12h, false, false, RIFLEFailure},
		{"worst of the two", 1.0, 2.0, RIFLEUrineBelow05For6h, false, false, RIFLEInjury},
		{"loss", 1.0, 1.0, RIFLEUrineNone, true, false, RIFLELoss},
		{"esrd over loss", 1.0, 1.0, RIFLEUrineNone, true, true, RIFLEESRD},
	}
	for _, tc := range tests {
		got, err := RIFLEClassify(tc.base, tc.cr, tc.urine, tc.rrt, tc.esd)
		if err != nil || got != tc.want {
			t.Errorf("%s: (%q, %v), want %q", tc.name, got, err, tc.want)
		}
		si, err := RIFLEClassifySI(umolOf(tc.base), umolOf(tc.cr), tc.urine, tc.rrt, tc.esd)
		if err != nil || si != tc.want {
			t.Errorf("%s SI: (%q, %v), want %q", tc.name, si, err, tc.want)
		}
	}
}

// parity: ledger rifle_criteria, source cases. The app takes the creatinine change as a picked band and the GFR
// decrease (> 25/50/75%), not built here; with creatinine 1.0 → 2.5 and urine < 0.5 × 6 h both give Injury, and
// with dialysis > 3 months both give ESRD (the app labels it "ESKD"-style "End Stage").
func TestRIFLEClassify_parity(t *testing.T) {
	if got, _ := RIFLEClassify(1.0, 2.5, RIFLEUrineBelow05For6h, false, false); got != RIFLEInjury {
		t.Errorf("got %q, want Injury", got)
	}
	if got, _ := RIFLEClassify(1.0, 1.0, RIFLEUrineNone, false, true); got != RIFLEESRD {
		t.Errorf("got %q, want ESRD", got)
	}
}

func TestRIFLEClassify_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"base NaN", func() error { _, e := RIFLEClassify(nan, 1, 0, false, false); return e }, ErrNotFinite},
		{"cr +Inf", func() error { _, e := RIFLEClassify(1, pInf, 0, false, false); return e }, ErrNotFinite},
		{"base 0", func() error { _, e := RIFLEClassify(0, 1, 0, false, false); return e }, ErrOutOfRange},
		{"cr -1", func() error { _, e := RIFLEClassify(1, -1, 0, false, false); return e }, ErrOutOfRange},
		{"urine 4", func() error { _, e := RIFLEClassify(1, 1, 4, false, false); return e }, ErrOutOfRange},
		{"urine -1", func() error { _, e := RIFLEClassify(1, 1, -1, false, false); return e }, ErrOutOfRange},
		{"overflow", func() error { _, e := RIFLEClassify(1, math.MaxFloat64, 0, false, false); return e }, ErrDomain},
		{"SI -Inf", func() error { _, e := RIFLEClassifySI(nInf, 1, 0, false, false); return e }, ErrNotFinite},
		{"SI 0", func() error { _, e := RIFLEClassifySI(88.4, 0, 0, false, false); return e }, ErrOutOfRange},
	})
}

func TestRIFLEClassify_noAlloc(t *testing.T) {
	n := testing.AllocsPerRun(100, func() {
		sinkRIFLE, _ = RIFLEClassify(1.0, 2.1, RIFLEUrineBelow05For6h, false, false)
		sinkRIFLE, _ = RIFLEClassifySI(88.4, 185, RIFLEUrineNone, false, false)
	})
	if n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkRIFLEClassify(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkRIFLE, _ = RIFLEClassify(1.0, 2.1, RIFLEUrineBelow05For6h, false, false)
	}
}

func BenchmarkRIFLEClassifySI(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkRIFLE, _ = RIFLEClassifySI(88.4, 185, RIFLEUrineNone, false, false)
	}
}
