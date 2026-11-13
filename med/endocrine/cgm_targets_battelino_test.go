package endocrine

import (
	"errors"
	"math"
	"testing"
)

var allMet = CGMTargetsMet{true, true, true, true, true, true}

// Battelino 2019 Table 3 (type 1 and type 2) and Table 2 (%CV): just inside every target, and exactly on
// every printed number (strict targets fail at the number, "≤36%" passes at 36).
func TestCGMTargetsMetBattelino(t *testing.T) {
	cases := []struct {
		name string
		m    CGMMetricPercents
		want CGMTargetsMet
	}{
		{"inside", CGMMetricPercents{71, 3.9, 0.9, 24.9, 4.9, 36}, allMet},
		{"on the numbers", CGMMetricPercents{70, 4, 1, 25, 5, 36}, CGMTargetsMet{CV: true}},
		{"just outside", CGMMetricPercents{69.9, 4.1, 1.1, 25.1, 5.1, 36.1}, CGMTargetsMet{}},
		{"ideal", CGMMetricPercents{100, 0, 0, 0, 0, 0}, allMet},
		{"TIR 70.1", CGMMetricPercents{70.1, 0, 0, 0, 0, 20}, allMet},
	}
	for _, c := range cases {
		got, err := CGMTargetsMetBattelino(c.m)
		if err != nil {
			t.Errorf("%s: unexpected error %v", c.name, err)
			continue
		}
		if got != c.want {
			t.Errorf("%s: got %+v, want %+v", c.name, got, c.want)
		}
	}
}

// From MED-ENDO-011's output: the nested sums feed the targets.
func TestCGMTargetsMetBattelino_fromTimeInRanges(t *testing.T) {
	p, err := TimeInRangesMgDL([]float64{110, 180, 95, 150, 152, 90, 200, 120, 130, 85})
	if err != nil {
		t.Fatal(err)
	}
	got, err := CGMTargetsMetBattelino(CGMMetricPercents{
		InRangePercent:       p.InRangePercent,
		BelowRange70Percent:  p.BelowRangeLevel1Percent + p.BelowRangeLevel2Percent,
		BelowRange54Percent:  p.BelowRangeLevel2Percent,
		AboveRange180Percent: p.AboveRangeLevel1Percent + p.AboveRangeLevel2Percent,
		AboveRange250Percent: p.AboveRangeLevel2Percent,
		CVPercent:            25,
	})
	if err != nil || got != allMet {
		t.Errorf("got %+v, %v; want all met (TIR 90 %%, TAR>180 10 %%)", got, err)
	}
}

func TestCGMTargetsMetBattelino_errors(t *testing.T) {
	base := CGMMetricPercents{80, 2, 0.5, 15, 2, 30}
	set := func(i int, v float64) CGMMetricPercents {
		m := base
		ps := [...]*float64{&m.InRangePercent, &m.BelowRange70Percent, &m.BelowRange54Percent,
			&m.AboveRange180Percent, &m.AboveRange250Percent, &m.CVPercent}
		*ps[i] = v
		return m
	}
	for i := range 6 {
		for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
			if _, err := CGMTargetsMetBattelino(set(i, v)); !errors.Is(err, ErrNotFinite) {
				t.Errorf("field %d = %v: err = %v, want ErrNotFinite", i, v, err)
			}
		}
		if _, err := CGMTargetsMetBattelino(set(i, -0.1)); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("field %d = -0.1: err = %v, want ErrOutOfRange", i, err)
		}
	}
	// Percent edges 0 and 100 accepted, 100.1 refused; CV has no upper limit.
	for i := range 5 {
		if _, err := CGMTargetsMetBattelino(set(i, 100.1)); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("field %d = 100.1: err = %v, want ErrOutOfRange", i, err)
		}
	}
	if _, err := CGMTargetsMetBattelino(CGMMetricPercents{0, 100, 100, 100, 100, 0}); err != nil {
		t.Errorf("edges 0/100: %v", err)
	}
	if _, err := CGMTargetsMetBattelino(set(5, 250)); err != nil {
		t.Errorf("CV 250: %v", err)
	}
	// Nesting: TBR<54 ≤ TBR<70 and TAR>250 ≤ TAR>180; equal is accepted.
	if _, err := CGMTargetsMetBattelino(CGMMetricPercents{80, 2, 2.1, 15, 2, 30}); !errors.Is(err, ErrDomain) {
		t.Errorf("TBR54 > TBR70: err = %v, want ErrDomain", err)
	}
	if _, err := CGMTargetsMetBattelino(CGMMetricPercents{80, 2, 0.5, 15, 15.1, 30}); !errors.Is(err, ErrDomain) {
		t.Errorf("TAR250 > TAR180: err = %v, want ErrDomain", err)
	}
	if _, err := CGMTargetsMetBattelino(CGMMetricPercents{80, 2, 2, 15, 15, 30}); err != nil {
		t.Errorf("equal nested ranges: %v", err)
	}
}

// parity: ledger bs_cgm_metrics_category, source cases 1–4 (time below 54 "< 1" on goal; time above 250
// "< 5" on goal; TIR "≥ 70" on goal; CV "≤ 36" stable, > 36 unstable). Source case 3 disagrees with the
// primary: the source passes exactly 70 % TIR (≥ 70), Table 3 prints ">70%"; the primary's "not met" is
// asserted.
func TestCGMTargetsMetBattelino_parity(t *testing.T) {
	cases := []struct {
		name string
		m    CGMMetricPercents
		want CGMTargetsMet
	}{
		// case 1: TBR<54 0.99 on goal, 1 above goal.
		{"case 1 below", CGMMetricPercents{80, 3, 0.99, 10, 2, 30}, allMet},
		{"case 1 on", CGMMetricPercents{80, 3, 1, 10, 2, 30}, CGMTargetsMet{true, true, false, true, true, true}},
		// case 2: TAR>250 4.99 on goal, 5 above goal.
		{"case 2 below", CGMMetricPercents{80, 3, 0.5, 10, 4.99, 30}, allMet},
		{"case 2 on", CGMMetricPercents{80, 3, 0.5, 10, 5, 30}, CGMTargetsMet{true, true, true, true, false, true}},
		// case 3: TIR 70 — source "Time in range at or above 70%" (on goal); primary ">70%": not met.
		{"case 3 on", CGMMetricPercents{70, 3, 0.5, 20, 2, 30}, CGMTargetsMet{false, true, true, true, true, true}},
		{"case 3 below", CGMMetricPercents{69.99, 3, 0.5, 20, 2, 30}, CGMTargetsMet{false, true, true, true, true, true}},
		// case 4: CV 36 stable (met), 36.01 unstable (not met).
		{"case 4 on", CGMMetricPercents{80, 3, 0.5, 10, 2, 36}, allMet},
		{"case 4 above", CGMMetricPercents{80, 3, 0.5, 10, 2, 36.01}, CGMTargetsMet{true, true, true, true, true, false}},
	}
	for _, c := range cases {
		got, err := CGMTargetsMetBattelino(c.m)
		if err != nil || got != c.want {
			t.Errorf("%s: got %+v, %v; want %+v", c.name, got, err, c.want)
		}
	}
}

func TestCGMTargetsMetBattelino_noAlloc(t *testing.T) {
	m := CGMMetricPercents{75, 3, 0.5, 20, 3, 33}
	if a := testing.AllocsPerRun(100, func() { _, _ = CGMTargetsMetBattelino(m) }); a != 0 {
		t.Errorf("%v allocs per run, want 0", a)
	}
}

func BenchmarkCGMTargetsMetBattelino(b *testing.B) {
	b.ReportAllocs()
	m := CGMMetricPercents{75, 3, 0.5, 20, 3, 33}
	for b.Loop() {
		_, _ = CGMTargetsMetBattelino(m)
	}
}
