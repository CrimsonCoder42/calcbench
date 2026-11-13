package float

import (
	"math"
	"slices"
	"testing"
)

func TestReplaceNonFinite_Cases(t *testing.T) {
	inf, nan := math.Inf(1), math.NaN()
	cases := []struct {
		name string
		in   []float64
		want []float64
	}{
		{"empty", nil, []float64{}},
		{"all finite", []float64{1, -2, 0}, []float64{1, -2, 0}},
		{"mixed", []float64{1, nan, inf, -inf}, []float64{1, 0, 9, -9}},
		{"extremes kept", []float64{math.MaxFloat64, -math.MaxFloat64, 5e-324}, []float64{math.MaxFloat64, -math.MaxFloat64, 5e-324}},
	}
	for _, c := range cases {
		got := ReplaceNonFinite(nil, c.in, 0, 9, -9)
		if !slices.Equal(got, c.want) {
			t.Errorf("%s: got %v; want %v", c.name, got, c.want)
		}
	}
}

func TestReplaceNonFinite_InPlace(t *testing.T) {
	v := []float64{math.NaN(), 2, math.Inf(-1)}
	got := ReplaceNonFinite(v, v, -1, 0, 0)
	if !slices.Equal(got, []float64{-1, 2, 0}) || &got[0] != &v[0] {
		t.Errorf("in place: got %v (same backing %v)", got, &got[0] == &v[0])
	}
	if n := testing.AllocsPerRun(100, func() { _ = ReplaceNonFinite(v, v, 0, 0, 0) }); n != 0 {
		t.Errorf("allocs = %v; want 0", n)
	}
}

func BenchmarkReplaceNonFinite(b *testing.B) {
	in := make([]float64, 1000)
	for i := range in {
		in[i] = float64(i)
		if i%10 == 0 {
			in[i] = math.NaN()
		}
	}
	dst := make([]float64, 0, len(in))
	b.ReportAllocs()
	for b.Loop() {
		dst = ReplaceNonFinite(dst, in, 0, math.MaxFloat64, -math.MaxFloat64)
	}
}
