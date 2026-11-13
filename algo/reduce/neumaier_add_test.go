package reduce

import (
	"math"
	"testing"
)

// Neumaier 1974: stepping 1, 1e100, 1, -1e100 leaves sum 0 and correction 2; the total is 2.
func TestNeumaierAdd_example(t *testing.T) {
	var s, c float64
	for _, x := range []float64{1, 1e100, 1, -1e100} {
		s, c = NeumaierAdd(s, c, x)
	}
	if s+c != 2 {
		t.Errorf("total = %v + %v = %v, want 2", s, c, s+c)
	}
}

func TestNeumaierAdd_cases(t *testing.T) {
	tests := []struct {
		name              string
		sum, comp, x      float64
		wantSum, wantComp float64
	}{
		{"zero", 0, 0, 0, 0, 0},
		{"exact", 1, 0, 2, 3, 0},
		{"small lost", 1, 0, 0x1p-60, 1, 0x1p-60},
		{"large x", 0x1p-60, 0, 1, 1, 0x1p-60},
	}
	for _, tc := range tests {
		s, c := NeumaierAdd(tc.sum, tc.comp, tc.x)
		if s != tc.wantSum || c != tc.wantComp {
			t.Errorf("%s: NeumaierAdd = (%v, %v), want (%v, %v)", tc.name, s, c, tc.wantSum, tc.wantComp)
		}
	}
}

func TestNeumaierAdd_nanPropagates(t *testing.T) {
	s, c := NeumaierAdd(0, 0, math.NaN())
	if AllFinite(s + c) {
		t.Errorf("NeumaierAdd(NaN) total = %v, want non-finite", s+c)
	}
}

func BenchmarkNeumaierAdd(b *testing.B) {
	b.ReportAllocs()
	var s, c float64
	for b.Loop() {
		s, c = NeumaierAdd(s, c, 0.1)
	}
	_ = s + c
}
