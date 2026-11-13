package fitness

import (
	"math"
	"testing"
)

// The IPF GL document's practical examples, p. 4 (read 2026-11-12):
//
//  1. Men's Equipped Powerlifting, body weight 92.04 kg, total 1035.0 kg: 0.109039 × 1035.0 = 112.855365 points
//  2. Women's Classic Bench Press, body weight 70.50 kg, result 122.5 kg: 0.790069 × 122.5 = 96.783453 points
//
// Coefficients and points are 6-decimal values from the document's rounding; compared to 1e-12 absolute (the
// float64 spelling of a 6-decimal number).
func TestIPFGLPoints_example(t *testing.T) {
	cases := []struct {
		event        IPFEvent
		female       bool
		bw, result   float64
		coef, points float64
	}{
		{IPFEquippedPowerlifting, false, 92.04, 1035.0, 0.109039, 112.855365},
		{IPFClassicBenchPress, true, 70.50, 122.5, 0.790069, 96.783453},
	}
	for _, c := range cases {
		coef, err := IPFGLCoefficient(c.event, c.female, c.bw)
		if err != nil || math.Abs(coef-c.coef) > 1e-12 {
			t.Errorf("IPFGLCoefficient(%d, %v, %v) = (%v, %v), want %v", c.event, c.female, c.bw, coef, err, c.coef)
		}
		pts, err := IPFGLPoints(c.event, c.female, c.bw, c.result)
		if err != nil || math.Abs(pts-c.points) > 1e-12 {
			t.Errorf("IPFGLPoints(%d, %v, %v, %v) = (%v, %v), want %v", c.event, c.female, c.bw, c.result, pts, err, c.points)
		}
	}
}

// parity: ledger ipf_gl_points, source cases 1–8. Coefficient and points copied from the source test.
// Case 2's source value is 96.7834525 (unrounded product, compared within 1.5e-6); the document prints
// 96.783453 (rounded to 6 decimals, p. 4), and that is what is asserted here.
func TestIPFGLPoints_parity(t *testing.T) {
	cases := []struct {
		event        IPFEvent
		female       bool
		bw, result   float64
		coef, points float64
	}{
		{IPFEquippedPowerlifting, false, 92.04, 1035.0, 0.109039, 112.855365}, // parity: ledger ipf_gl_points, source case 1
		{IPFClassicBenchPress, true, 70.50, 122.5, 0.790069, 96.783453},       // parity: ledger ipf_gl_points, source case 2 (source 96.7834525)
		{IPFClassicPowerlifting, false, 83, 700, 0.138427, 96.8989},           // parity: ledger ipf_gl_points, source case 3
		{IPFClassicPowerlifting, true, 63, 450, 0.218783, 98.45235},           // parity: ledger ipf_gl_points, source case 4
		{IPFEquippedBenchPress, false, 105, 300, 0.310505, 93.1515},           // parity: ledger ipf_gl_points, source case 5
		{IPFEquippedBenchPress, true, 57, 150, 0.645613, 96.84195},            // parity: ledger ipf_gl_points, source case 6
		{IPFEquippedPowerlifting, true, 72, 550, 0.168294, 92.5617},           // parity: ledger ipf_gl_points, source case 7
		{IPFClassicBenchPress, false, 120, 250, 0.421882, 105.4705},           // parity: ledger ipf_gl_points, source case 8
	}
	for _, c := range cases {
		coef, err := IPFGLCoefficient(c.event, c.female, c.bw)
		if err != nil || math.Abs(coef-c.coef) > 1e-12 {
			t.Errorf("IPFGLCoefficient(%d, %v, %v) = (%v, %v), want %v", c.event, c.female, c.bw, coef, err, c.coef)
		}
		pts, err := IPFGLPoints(c.event, c.female, c.bw, c.result)
		if err != nil || math.Abs(pts-c.points) > 1e-12 {
			t.Errorf("IPFGLPoints(%d, %v, %v, %v) = (%v, %v), want %v", c.event, c.female, c.bw, c.result, pts, err, c.points)
		}
	}
}

// Both sides of the body-weight condition ("Bwt ≥ 40kg for men and Bwt ≥ 35kg for women", p. 1), and a result
// of 0 (equation (1): 0 points). At the edges, HAND-DERIVED: men Classic Bench Press 40 kg →
// 100 / (320.98041 − 281.40258·e^(−0.4032)) = 0.752144; women Classic Bench Press 35 kg → 1.732945.
func TestIPFGLPoints_edges(t *testing.T) {
	if c, err := IPFGLCoefficient(IPFClassicBenchPress, false, 40); err != nil || math.Abs(c-0.752144) > 1e-12 {
		t.Errorf("men 40 kg: (%v, %v), want 0.752144", c, err)
	}
	if c, err := IPFGLCoefficient(IPFClassicBenchPress, true, 35); err != nil || math.Abs(c-1.732945) > 1e-12 {
		t.Errorf("women 35 kg: (%v, %v), want 1.732945", c, err)
	}
	if p, err := IPFGLPoints(IPFClassicPowerlifting, false, 83, 0); err != nil || p != 0 {
		t.Errorf("result 0: (%v, %v), want 0", p, err)
	}
	// Large body weight: the denominator tends to A, so the coefficient stays finite.
	if c, err := IPFGLCoefficient(IPFEquippedPowerlifting, false, 1e6); err != nil || math.Abs(c-math.Round(1e8/1236.25115)/1e6) > 1e-12 {
		t.Errorf("1e6 kg: (%v, %v)", c, err)
	}
}

// parity: ledger ipf_gl_points, source refusal cases: event 4, sex code 2 and 0.5 (not expressible: female is
// bool), body weight 30 kg, lb unit (not expressible: kg only), missing event (not expressible).
func TestIPFGLPoints_errors(t *testing.T) {
	coef := func(e IPFEvent, f bool, bw float64) func() error {
		return func() error { _, err := IPFGLCoefficient(e, f, bw); return err }
	}
	pts := func(e IPFEvent, f bool, bw, r float64) func() error {
		return func() error { _, err := IPFGLPoints(e, f, bw, r); return err }
	}
	runErrCases(t, []errCase{
		{"event 4 (parity)", pts(4, false, 90, 700), ErrOutOfRange},
		{"event -1", pts(-1, false, 90, 700), ErrOutOfRange},
		{"body weight 30 kg (parity)", pts(IPFEquippedPowerlifting, false, 30, 700), ErrOutOfRange},
		{"men just below 40 kg", coef(IPFClassicBenchPress, false, 39.99), ErrOutOfRange},
		{"women just below 35 kg", coef(IPFClassicBenchPress, true, 34.99), ErrOutOfRange},
		{"result negative", pts(IPFClassicPowerlifting, false, 83, -0.5), ErrOutOfRange},
		{"bw NaN", coef(IPFClassicPowerlifting, false, nan), ErrNotFinite},
		{"bw +Inf", coef(IPFClassicPowerlifting, false, pInf), ErrNotFinite},
		{"bw -Inf", coef(IPFClassicPowerlifting, false, nInf), ErrNotFinite},
		{"result NaN", pts(IPFClassicPowerlifting, false, 83, nan), ErrNotFinite},
		{"result +Inf", pts(IPFClassicPowerlifting, false, 83, pInf), ErrNotFinite},
		{"result -Inf", pts(IPFClassicPowerlifting, false, 83, nInf), ErrNotFinite},
		{"points overflow", pts(IPFClassicBenchPress, true, 35, math.MaxFloat64), ErrDomain},
	})
}

func TestIPFGLPoints_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = IPFGLPoints(IPFClassicBenchPress, true, 70.5, 122.5) }); n != 0 {
		t.Errorf("IPFGLPoints allocates %v times, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sink, _ = IPFGLCoefficient(IPFClassicBenchPress, true, 70.5) }); n != 0 {
		t.Errorf("IPFGLCoefficient allocates %v times, want 0", n)
	}
}

func BenchmarkIPFGLCoefficient(b *testing.B) {
	b.ReportAllocs()
	bw := 92.04
	for b.Loop() {
		sink, _ = IPFGLCoefficient(IPFEquippedPowerlifting, false, bw)
	}
}

func BenchmarkIPFGLPoints(b *testing.B) {
	b.ReportAllocs()
	bw, r := 92.04, 1035.0
	for b.Loop() {
		sink, _ = IPFGLPoints(IPFEquippedPowerlifting, false, bw, r)
	}
}
