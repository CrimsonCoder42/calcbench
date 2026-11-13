package obstetric

import (
	"errors"
	"math"
	"testing"
)

// Cases from IOM 2009 Table S-1 (Summary): each category on both sides of its edges 18.5, 25 and 30.
func TestGestationalWeightGainIOM_tableS1(t *testing.T) {
	under := GWGRecommendation{GWGUnderweight, 28, 40, 1, 1, 1.3}
	normal := GWGRecommendation{GWGNormalWeight, 25, 35, 1, 0.8, 1}
	over := GWGRecommendation{GWGOverweight, 15, 25, 0.6, 0.5, 0.7}
	obese := GWGRecommendation{GWGObese, 11, 20, 0.5, 0.4, 0.6}
	tests := []struct {
		bmi  float64
		want GWGRecommendation
	}{
		{10, under}, {18.4999, under},
		{18.5, normal}, {22, normal}, {24.9999, normal},
		{25, over}, {29.9999, over},
		{30, obese}, {31, obese}, {60, obese},
	}
	for _, tc := range tests {
		if got, err := GestationalWeightGainIOM(tc.bmi); err != nil || got != tc.want {
			t.Errorf("GestationalWeightGainIOM(%v) = (%+v, %v), want %+v", tc.bmi, got, err, tc.want)
		}
	}
}

// Provisional twin ranges (IOM 2009 Summary): normal 37–54, overweight 31–50, obese 25–42 lb; none for underweight.
func TestGestationalWeightGainTwinsIOM_ranges(t *testing.T) {
	tests := []struct {
		bmi       float64
		cat       GWGCategory
		low, high float64
	}{
		{18.5, GWGNormalWeight, 37, 54}, {24.9999, GWGNormalWeight, 37, 54},
		{25, GWGOverweight, 31, 50}, {29.9999, GWGOverweight, 31, 50},
		{30, GWGObese, 25, 42},
	}
	for _, tc := range tests {
		c, lo, hi, err := GestationalWeightGainTwinsIOM(tc.bmi)
		if err != nil || c != tc.cat || lo != tc.low || hi != tc.high {
			t.Errorf("GestationalWeightGainTwinsIOM(%v) = (%q, %v, %v, %v), want (%q, %v, %v)", tc.bmi, c, lo, hi, err, tc.cat, tc.low, tc.high)
		}
	}
}

func TestGestationalWeightGainIOM_errors(t *testing.T) {
	for _, bmi := range []float64{math.NaN(), math.Inf(1), math.Inf(-1), 0, -5} {
		if _, err := GestationalWeightGainIOM(bmi); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("singleton %v: err = %v, want ErrOutOfRange", bmi, err)
		}
		if _, _, _, err := GestationalWeightGainTwinsIOM(bmi); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("twins %v: err = %v, want ErrOutOfRange", bmi, err)
		}
	}
	if _, _, _, err := GestationalWeightGainTwinsIOM(18.4999); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("twins underweight: err = %v, want ErrOutOfRange", err)
	}
}

func TestGestationalWeightGainIOM_noAllocation(t *testing.T) {
	var r GWGRecommendation
	var c GWGCategory
	if n := testing.AllocsPerRun(100, func() {
		r, _ = GestationalWeightGainIOM(22)
		c, _, _, _ = GestationalWeightGainTwinsIOM(27)
	}); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
	_, _ = r, c
}

func BenchmarkGestationalWeightGainIOM(b *testing.B) {
	b.ReportAllocs()
	var r GWGRecommendation
	for b.Loop() {
		r, _ = GestationalWeightGainIOM(22)
	}
	_ = r
}

func BenchmarkGestationalWeightGainTwinsIOM(b *testing.B) {
	b.ReportAllocs()
	var c GWGCategory
	for b.Loop() {
		c, _, _, _ = GestationalWeightGainTwinsIOM(27)
	}
	_ = c
}
