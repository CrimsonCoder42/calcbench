package hematology

import (
	"errors"
	"math"
	"testing"
)

// Tolerance 1e-12 absolute: one division and one addition on values ≤ 2.
const crpTol = 1e-12

// Wener 2000 abstract: CRP = age/50 (males), age/50 + 0.6 (females), ages 25–70.
func TestCRPUpperReferenceLimitWenerMgDL_formula(t *testing.T) {
	tests := []struct {
		age    float64
		female bool
		want   float64
	}{
		{50, false, 1.0}, {50, true, 1.6}, {25, false, 0.5}, {25, true, 1.1}, {70, false, 1.4}, {70, true, 2.0},
	}
	for _, tc := range tests {
		got, err := CRPUpperReferenceLimitWenerMgDL(tc.age, tc.female)
		if err != nil || math.Abs(got-tc.want) > crpTol {
			t.Errorf("CRPUpperReferenceLimitWenerMgDL(%v, %v) = (%v, %v), want %v", tc.age, tc.female, got, err, tc.want)
		}
	}
}

// parity: ledger age_adjusted_esr_crp. The source entry has no _test.go and its JSON no example; the consumer
// computes crpULN = female ? age/50 + 0.6 : age/50 (no age limit), so case 1 (woman 60) is 1.8 mg/dL. The
// consumer applies the formula from age 18 to 100; outside 25–70 this function refuses (ledger note).
func TestCRPUpperReferenceLimitWenerMgDL_parity(t *testing.T) {
	got, err := CRPUpperReferenceLimitWenerMgDL(60, true)
	if err != nil || math.Abs(got-1.8) > crpTol {
		t.Errorf("source case 1: got (%v, %v), want 1.8", got, err)
	}
}

func TestCRPUpperReferenceLimitWenerMgDL_errors(t *testing.T) {
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := CRPUpperReferenceLimitWenerMgDL(v, false); !errors.Is(err, ErrNotFinite) {
			t.Errorf("ageYears %v: error = %v, want ErrNotFinite", v, err)
		}
	}
	for _, v := range []float64{24.99, 70.01, 18, 100} {
		if _, err := CRPUpperReferenceLimitWenerMgDL(v, true); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("ageYears %v: error = %v, want ErrOutOfRange", v, err)
		}
	}
}

func TestCRPUpperReferenceLimitWenerMgDL_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = CRPUpperReferenceLimitWenerMgDL(50, true) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkCRPUpperReferenceLimitWenerMgDL(b *testing.B) {
	b.ReportAllocs()
	age := 50.0
	for b.Loop() {
		sinkFloat, _ = CRPUpperReferenceLimitWenerMgDL(age, true)
	}
}
