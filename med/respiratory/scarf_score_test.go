package respiratory

import (
	"errors"
	"math"
	"testing"
)

// HAND-DERIVED from Hardin 2019 (PMID 31425473), one point per item: below 50 % predicted + RR 24 = 2;
// all four = 4; none = 0.
func TestSCARFScore_handExample(t *testing.T) {
	if got, err := SCARFScore(true, 24, 4, false); err != nil || got != 2 {
		t.Errorf("got (%d, %v), want 2", got, err)
	}
	if got, err := SCARFScore(true, 30, 8, true); err != nil || got != 4 {
		t.Errorf("all: got (%d, %v), want 4", got, err)
	}
	if got, err := SCARFScore(false, 16, 0, false); err != nil || got != 0 {
		t.Errorf("none: got (%d, %v), want 0", got, err)
	}
}

// Edges: RR > 20 (20 does not score, 20.5 does); pain ≥ 5 (4 does not, 5 does).
func TestSCARFScore_edges(t *testing.T) {
	cases := []struct {
		rr   float64
		pain int
		want int
	}{{20, 0, 0}, {20.5, 0, 1}, {16, 4, 0}, {16, 5, 1}, {16, 10, 1}}
	for _, c := range cases {
		if got, err := SCARFScore(false, c.rr, c.pain, false); err != nil || got != c.want {
			t.Errorf("rr %v pain %d: got (%d, %v), want %d", c.rr, c.pain, got, err, c.want)
		}
	}
}

func TestSCARFScore_errors(t *testing.T) {
	tests := []struct {
		name    string
		rr      float64
		pain    int
		wantErr error
	}{
		{"negative RR", -1, 3, ErrOutOfRange},
		{"pain -1", 16, -1, ErrOutOfRange},
		{"pain 11", 16, 11, ErrOutOfRange},
		{"NaN RR", math.NaN(), 3, ErrNotFinite},
		{"+Inf RR", math.Inf(1), 3, ErrNotFinite},
		{"-Inf RR", math.Inf(-1), 3, ErrNotFinite},
	}
	for _, tc := range tests {
		if _, err := SCARFScore(false, tc.rr, tc.pain, false); !errors.Is(err, tc.wantErr) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.wantErr)
		}
	}
}

func TestSCARFScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = SCARFScore(true, 24, 4, false) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkSCARFScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = SCARFScore(true, 24, 4, false)
	}
}
