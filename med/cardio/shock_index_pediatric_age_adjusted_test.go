package cardio

import (
	"math"
	"testing"
)

// Acker 2015 abstract: SI > 1.22 (4–6), > 1.0 (7–12), > 0.9 (13–16); both sides of every band edge.
func TestSIPACutoff_bands(t *testing.T) {
	tests := []struct{ age, want float64 }{
		{4, 1.22}, {6.99, 1.22}, {7, 1.0}, {12.99, 1.0}, {13, 0.9}, {16.99, 0.9},
	}
	for _, tc := range tests {
		if got, err := SIPACutoff(tc.age); err != nil || got != tc.want {
			t.Errorf("age %v: (%v, %v), want %v", tc.age, got, err, tc.want)
		}
	}
	for _, a := range []float64{3.99, 17} {
		_, err := SIPACutoff(a)
		checkErr(t, "range", err, ErrOutOfRange)
	}
	_, err := SIPACutoff(math.NaN())
	checkErr(t, "NaN", err, ErrNotFinite)
}

func TestSIPACutoff_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = SIPACutoff(10) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkSIPACutoff(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = SIPACutoff(10)
	}
}
