package market

import (
	"math"
	"testing"
)

var sinkFloat2 float64

// ELO 2012 §2 gives the rule V·Z(ΔP/σ); hand-checked with NIST's Φ(1.96) = 0.97500 (float64
// 0.9750021048517795, NUM-PROB NormalCDF). Tolerance 1e-9: the CDF's stated accuracy times V = 1000.
func TestBulkVolumeClassification_handChecked(t *testing.T) {
	buy, sell, err := BulkVolumeClassification(1000, 1.96, 1)
	if err != nil || !near(buy, 975.0021048517795, 1e-9) || !near(sell, 24.9978951482205, 1e-9) {
		t.Errorf("got %v, %v, %v; want 975.0021…, 24.9978…", buy, sell, err)
	}
}

func TestBulkVolumeClassification_cases(t *testing.T) {
	buy, sell, err := BulkVolumeClassification(1000, 0, 2)
	if err != nil || buy != 500 || sell != 500 {
		t.Errorf("flat bar: %v, %v, %v; want 500, 500", buy, sell, err)
	}
	buy, sell, err = BulkVolumeClassification(0, 3, 1)
	if err != nil || buy != 0 || sell != 0 {
		t.Errorf("zero volume: %v, %v, %v", buy, sell, err)
	}
	// symmetry: a fall of the same size swaps buy and sell
	b1, s1, _ := BulkVolumeClassification(800, -0.7, 1.3)
	b2, s2, _ := BulkVolumeClassification(800, 0.7, 1.3)
	if b1 != s2 || s1 != b2 || !near(b1+s1, 800, 1e-12) {
		t.Errorf("symmetry: %v %v / %v %v", b1, s1, b2, s2)
	}
}

func TestBulkVolumeClassification_errors(t *testing.T) {
	_, _, err := BulkVolumeClassification(-1, 0, 1)
	wantErr(t, "negative volume", err, ErrNegative)
	_, _, err = BulkVolumeClassification(1, math.NaN(), 1)
	wantErr(t, "NaN change", err, ErrNotFinite)
	_, _, err = BulkVolumeClassification(1, 0, math.Inf(1))
	wantErr(t, "Inf sd", err, ErrNotFinite)
	_, _, err = BulkVolumeClassification(1, 0, 0)
	wantErr(t, "zero sd", err, ErrOutOfRange)
}

func BenchmarkBulkVolumeClassification(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkFloat2, sinkErr = BulkVolumeClassification(1000, 0.3, 1.1)
	}
}

// ν = 1 is the Cauchy case, t₁(1) = ½ + atan(1)/π = 0.75 (NUM-PROB-078's example). Tolerance 1e-12.
func TestBulkVolumeClassificationT_handChecked(t *testing.T) {
	buy, sell, err := BulkVolumeClassificationT(1000, 1, 1, 1)
	if err != nil || !near(buy, 750, 1e-12) || !near(sell, 250, 1e-12) {
		t.Errorf("got %v, %v, %v; want 750, 250", buy, sell, err)
	}
	// ν = 2, t = 2 → ½ + 1/√6 (NUM-PROB-078's example)
	buy, _, err = BulkVolumeClassificationT(1, 4, 2, 2)
	if err != nil || !near(buy, 0.5+1/math.Sqrt(6), 1e-14) {
		t.Errorf("ν = 2: %v, %v", buy, err)
	}
}

func TestBulkVolumeClassificationT_largeDFNearNormal(t *testing.T) {
	bt, _, err1 := BulkVolumeClassificationT(1000, 1.96, 1, 1e6)
	bn, _, err2 := BulkVolumeClassification(1000, 1.96, 1)
	if err1 != nil || err2 != nil || !near(bt, bn, 1e-3) {
		t.Errorf("ν = 1e6: %v vs normal %v (%v, %v)", bt, bn, err1, err2)
	}
}

func TestBulkVolumeClassificationT_errors(t *testing.T) {
	_, _, err := BulkVolumeClassificationT(1, 0, 1, 0)
	wantErr(t, "df 0", err, ErrOutOfRange)
	_, _, err = BulkVolumeClassificationT(1, 0, 1, math.Inf(1))
	wantErr(t, "df Inf", err, ErrNotFinite)
	_, _, err = BulkVolumeClassificationT(-1, 0, 1, 3)
	wantErr(t, "negative volume", err, ErrNegative)
	_, _, err = BulkVolumeClassificationT(1, 0, -1, 3)
	wantErr(t, "negative sd", err, ErrOutOfRange)
	_, _, err = BulkVolumeClassificationT(1, math.MaxFloat64, 1e-300, 3)
	wantErr(t, "z overflow", err, ErrOverflow)
}

func BenchmarkBulkVolumeClassificationT(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkFloat2, sinkErr = BulkVolumeClassificationT(1000, 0.3, 1.1, 4)
	}
}
