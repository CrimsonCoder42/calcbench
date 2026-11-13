package cardio

import (
	"errors"
	"math"
	"testing"
)

var sinkPPVPercent float64

// Michard 2000 (PMID 10903232) abstract: ΔPP = (PPmax − PPmin) / mean(PPmax, PPmin), as a percentage.
// The abstract prints no worked example; the cases are HAND-DERIVED from that definition.
// Tolerance 1e-12 relative: three float operations.
func TestPulsePressureVariationPercent_definition(t *testing.T) {
	tests := []struct {
		name     string
		max, min float64
		want     float64
	}{
		{"equal pressures", 40, 40, 0},
		{"50 and 43", 50, 43, 100 * 7 / 46.5},        // 15.0537…
		{"13 percent", 53.9, 47.5, 100 * 6.4 / 50.7}, // 12.6232…
		{"doubled", 60, 30, 100 * 30.0 / 45},         // 66.67
		{"tiny positive", 2e-300, 1e-300, 100 * 1e-300 / 1.5e-300},
	}
	for _, tc := range tests {
		got, err := PulsePressureVariationPercent(tc.max, tc.min)
		if err != nil {
			t.Fatalf("%s: err = %v", tc.name, err)
		}
		if math.Abs(got-tc.want) > 1e-12*math.Max(1, math.Abs(tc.want)) {
			t.Errorf("%s: PulsePressureVariationPercent(%v, %v) = %v, want %v", tc.name, tc.max, tc.min, got, tc.want)
		}
	}
}

// The ratio form keeps the result finite for any finite inputs whose sum is finite.
func TestPulsePressureVariationPercent_large(t *testing.T) {
	got, err := PulsePressureVariationPercent(1e307, 1)
	if err != nil || math.IsInf(got, 0) || math.Abs(got-200) > 1e-9 {
		t.Errorf("PulsePressureVariationPercent(1e307, 1) = (%v, %v), want ≈ 200", got, err)
	}
	if _, err := PulsePressureVariationPercent(math.MaxFloat64, math.MaxFloat64/2); !errors.Is(err, ErrDomain) {
		t.Errorf("sum overflow: err = %v, want ErrDomain", err)
	}
}

// Range edges: min > 0 (0 rejected, the smallest positive accepted); max ≥ min (equal accepted, below rejected).
func TestPulsePressureVariationPercent_rangeEdges(t *testing.T) {
	if _, err := PulsePressureVariationPercent(10, 0); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("min 0: err = %v, want ErrOutOfRange", err)
	}
	if _, err := PulsePressureVariationPercent(10, -1); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("min −1: err = %v, want ErrOutOfRange", err)
	}
	if _, err := PulsePressureVariationPercent(10, math.SmallestNonzeroFloat64); err != nil {
		t.Errorf("min smallest positive: err = %v, want nil", err)
	}
	if _, err := PulsePressureVariationPercent(40, 40); err != nil {
		t.Errorf("max = min: err = %v, want nil", err)
	}
	if _, err := PulsePressureVariationPercent(39.999, 40); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("max < min: err = %v, want ErrOutOfRange", err)
	}
}

func TestPulsePressureVariationPercent_notFinite(t *testing.T) {
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := PulsePressureVariationPercent(v, 40); !errors.Is(err, ErrNotFinite) {
			t.Errorf("max %v: err = %v, want ErrNotFinite", v, err)
		}
		if _, err := PulsePressureVariationPercent(50, v); !errors.Is(err, ErrNotFinite) {
			t.Errorf("min %v: err = %v, want ErrNotFinite", v, err)
		}
	}
}

func TestPulsePressureVariationPercent_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkPPVPercent, _ = PulsePressureVariationPercent(50, 43) }); n != 0 {
		t.Errorf("PulsePressureVariationPercent allocates %v times, want 0", n)
	}
}

func BenchmarkPulsePressureVariationPercent(b *testing.B) {
	b.ReportAllocs()
	ppMax, ppMin := 50.0, 43.0
	for b.Loop() {
		sinkPPVPercent, _ = PulsePressureVariationPercent(ppMax, ppMin)
	}
}
