package respiratory

import (
	"errors"
	"math"
	"testing"
)

// HAND-DERIVED from Roca 2016's definition (PMID 27481760): (92 / 0.5) / 25 = 7.36. Tolerance 1e-12 (two divisions).
func TestROXIndex_handExample(t *testing.T) {
	if got, err := ROXIndex(92, 0.5, 25); err != nil || math.Abs(got-7.36) > 1e-12 {
		t.Errorf("got (%v, %v), want 7.36", got, err)
	}
	// Around the abstract's 4.88 cut-off: 97.6 / 0.5 / 40 = 4.88.
	if got, err := ROXIndex(97.6, 0.5, 40); err != nil || math.Abs(got-4.88) > 1e-12 {
		t.Errorf("got (%v, %v), want 4.88", got, err)
	}
}

func TestROXIndex_rangeEdges(t *testing.T) {
	for _, c := range [][3]float64{{100, 0.21, 20}, {50, 1.0, 20}} {
		if _, err := ROXIndex(c[0], c[1], c[2]); err != nil {
			t.Errorf("%v: err = %v, want nil", c, err)
		}
	}
	tests := []struct {
		name           string
		spo2, fio2, rr float64
		wantErr        error
	}{
		{"SpO2 0", 0, 0.5, 20, ErrOutOfRange},
		{"SpO2 100.1", 100.1, 0.5, 20, ErrOutOfRange},
		{"FiO2 0.2", 92, 0.2, 20, ErrOutOfRange},
		{"FiO2 1.01", 92, 1.01, 20, ErrOutOfRange},
		{"RR 0", 92, 0.5, 0, ErrOutOfRange},
		{"NaN SpO2", math.NaN(), 0.5, 20, ErrNotFinite},
		{"+Inf FiO2", 92, math.Inf(1), 20, ErrNotFinite},
		{"-Inf RR", 92, 0.5, math.Inf(-1), ErrNotFinite},
	}
	for _, tc := range tests {
		if _, err := ROXIndex(tc.spo2, tc.fio2, tc.rr); !errors.Is(err, tc.wantErr) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.wantErr)
		}
	}
}

func TestROXIndex_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = ROXIndex(92, 0.5, 25) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkROXIndex(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = ROXIndex(92, 0.5, 25)
	}
}

// A respiratory rate so small that the quotient leaves float64 is ErrDomain, never +Inf (review finding).
func TestROXIndex_overflow(t *testing.T) {
	if _, err := ROXIndex(100, 0.21, math.SmallestNonzeroFloat64); !errors.Is(err, ErrDomain) {
		t.Errorf("ROXIndex(100, 0.21, 5e-324) err = %v, want ErrDomain", err)
	}
}
