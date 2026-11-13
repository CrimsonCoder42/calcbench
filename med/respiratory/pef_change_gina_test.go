package respiratory

import (
	"errors"
	"math"
	"testing"
)

var sinkPEFChange bool

// GINA 2023, "How much variation in expiratory airflow is consistent with asthma?": "a change in PEF of at
// least 20%, is accepted as being consistent with asthma" (adults). Both sides of +20 and −20.
func TestPEFChangeConsistentWithAsthmaGINA_edges(t *testing.T) {
	tests := []struct {
		change float64
		want   bool
	}{
		{20, true}, {19.999, false}, {20.001, true},
		{-20, true}, {-19.999, false}, {-20.001, true},
		{0, false}, {-100, true}, {1500, true},
	}
	for _, tc := range tests {
		got, err := PEFChangeConsistentWithAsthmaGINA(tc.change, 40)
		if err != nil || got != tc.want {
			t.Errorf("PEFChangeConsistentWithAsthmaGINA(%v, 40) = (%v, %v), want %v", tc.change, got, err, tc.want)
		}
	}
}

// Parity with the reference's reversibility bands: < 20 "No Significant Reversibility", ≥ 20 "Positive
// Reversibility — Consistent with Asthma" (signed change, increase only). GINA counts an increase or a
// decrease, so a fall of ≥ 20 % asserts GINA (reference: no significant reversibility).
func TestPEFChangeConsistentWithAsthmaGINA_parity(t *testing.T) {
	tests := []struct {
		change float64
		want   bool
	}{
		{19.99, false}, // parity: ledger pef_reversibility_category, source case 1 (< 20)
		{20, true},     // parity: ledger pef_reversibility_category, source case 2 (exactly 20: positive in both)
		{35, true},     // parity: ledger pef_reversibility_category, source case 3
		{-25, true},    // parity: ledger pef_reversibility_category, source case 4 (reference: no significant reversibility; GINA: a decrease counts)
		{-93.75, true}, // parity: ledger pef_reversibility_category, source case 5 (the reference's lower input bound)
	}
	for _, tc := range tests {
		got, err := PEFChangeConsistentWithAsthmaGINA(tc.change, 30)
		if err != nil || got != tc.want {
			t.Errorf("PEFChangeConsistentWithAsthmaGINA(%v, 30) = (%v, %v), want %v", tc.change, got, err, tc.want)
		}
	}
}

func TestPEFChangeConsistentWithAsthmaGINA_errors(t *testing.T) {
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := PEFChangeConsistentWithAsthmaGINA(v, 30); !errors.Is(err, ErrNotFinite) {
			t.Errorf("change %v: err = %v, want ErrNotFinite", v, err)
		}
		if _, err := PEFChangeConsistentWithAsthmaGINA(25, v); !errors.Is(err, ErrNotFinite) {
			t.Errorf("age %v: err = %v, want ErrNotFinite", v, err)
		}
	}
	if _, err := PEFChangeConsistentWithAsthmaGINA(-100.001, 30); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("change < −100: err = %v, want ErrOutOfRange", err)
	}
	if _, err := PEFChangeConsistentWithAsthmaGINA(25, 17.99); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("age 17.99: err = %v, want ErrOutOfRange", err)
	}
	if _, err := PEFChangeConsistentWithAsthmaGINA(25, 18); err != nil {
		t.Errorf("age 18: err = %v, want nil", err)
	}
}

func TestPEFChangeConsistentWithAsthmaGINA_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkPEFChange, _ = PEFChangeConsistentWithAsthmaGINA(22, 30) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkPEFChangeConsistentWithAsthmaGINA(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkPEFChange, _ = PEFChangeConsistentWithAsthmaGINA(22, 30)
	}
}
