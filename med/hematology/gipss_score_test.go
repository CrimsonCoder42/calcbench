package hematology

import (
	"errors"
	"testing"
)

// Tefferi 2018: VHR karyotype 2, unfavorable 1; no type 1/like CALR, ASXL1, SRSF2, U2AF1 Q157 one each.
func TestGIPSSScore_example(t *testing.T) {
	tests := []struct {
		k             GIPSSKaryotype
		calr, a, s, u bool
		want          int
	}{
		{GIPSSKaryotypeFavorable, false, false, false, false, 0},
		{GIPSSKaryotypeUnfavorable, false, false, false, false, 1},
		{GIPSSKaryotypeVeryHighRisk, false, false, false, false, 2},
		{GIPSSKaryotypeUnfavorable, true, true, false, false, 3},
		{GIPSSKaryotypeVeryHighRisk, true, true, true, true, 6},
	}
	for _, tc := range tests {
		if got, err := GIPSSScore(tc.k, tc.calr, tc.a, tc.s, tc.u); err != nil || got != tc.want {
			t.Errorf("GIPSSScore(%+v) = (%d, %v), want %d", tc, got, err, tc.want)
		}
	}
	for _, k := range []GIPSSKaryotype{-1, 3} {
		if _, err := GIPSSScore(k, false, false, false, false); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("GIPSSScore(karyotype %d) error = %v, want ErrOutOfRange", k, err)
		}
	}
}

func TestGIPSSScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = GIPSSScore(GIPSSKaryotypeUnfavorable, true, true, false, false) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkGIPSSScore(b *testing.B) {
	b.ReportAllocs()
	k := GIPSSKaryotypeUnfavorable
	for b.Loop() {
		sinkInt, _ = GIPSSScore(k, true, true, false, false)
	}
}
