package rheumatology

import (
	"errors"
	"math"
	"testing"
)

type goutIn struct {
	joint, chars, course int
	tophus               bool
	urate                float64
	sfNeg, usDeposit, xr bool
}

func (g goutIn) run() (int, error) {
	return GoutACREULAR2015Score(g.joint, g.chars, g.course, g.tophus, g.urate, g.sfNeg, g.usDeposit, g.xr)
}

// Neogi 2015 Table 2. Minimum −6 (urate < 4 and MSU-negative fluid); maximum 23 ("The maximum possible score in
// the final criteria is 23"); each urate band edge on both sides.
func TestGoutACREULAR2015Score_example(t *testing.T) {
	tests := []struct {
		name string
		in   goutIn
		want int
	}{
		{"minimum", goutIn{0, 0, 0, false, 3.9, true, false, false}, -6},
		{"maximum", goutIn{2, 3, 2, true, 10, false, true, true}, 23},
		{"typical", goutIn{2, 3, 2, false, 9, false, false, false}, 10},
		{"urate 3.99", goutIn{urate: 3.99}, -4},
		{"urate 4", goutIn{urate: 4}, 0},
		{"urate 5.99", goutIn{urate: 5.99}, 0},
		{"urate 6", goutIn{urate: 6}, 2},
		{"urate 7.99", goutIn{urate: 7.99}, 2},
		{"urate 8", goutIn{urate: 8}, 3},
		{"urate 9.99", goutIn{urate: 9.99}, 3},
		{"urate 10", goutIn{urate: 10}, 4},
		{"ankle/midfoot", goutIn{joint: 1, urate: 5}, 1},
		{"one typical episode", goutIn{course: 1, urate: 5}, 1},
		{"tophus", goutIn{tophus: true, urate: 5}, 4},
		{"erosion", goutIn{xr: true, urate: 5}, 4},
		{"fluid negative", goutIn{sfNeg: true, urate: 5}, -2},
	}
	for _, c := range tests {
		if got, err := c.in.run(); err != nil || got != c.want {
			t.Errorf("%s: = (%d, %v), want %d", c.name, got, err, c.want)
		}
	}
}

func TestGoutACREULAR2015Score_errors(t *testing.T) {
	cases := []struct {
		name string
		in   goutIn
		want error
	}{
		{"urate NaN", goutIn{urate: math.NaN()}, ErrNotFinite},
		{"urate +Inf", goutIn{urate: math.Inf(1)}, ErrNotFinite},
		{"urate -Inf", goutIn{urate: math.Inf(-1)}, ErrNotFinite},
		{"urate -0.1", goutIn{urate: -0.1}, ErrOutOfRange},
		{"urate 30.1", goutIn{urate: 30.1}, ErrOutOfRange},
		{"joint -1", goutIn{joint: -1, urate: 5}, ErrOutOfRange},
		{"joint 3", goutIn{joint: 3, urate: 5}, ErrOutOfRange},
		{"chars -1", goutIn{chars: -1, urate: 5}, ErrOutOfRange},
		{"chars 4", goutIn{chars: 4, urate: 5}, ErrOutOfRange},
		{"course -1", goutIn{course: -1, urate: 5}, ErrOutOfRange},
		{"course 3", goutIn{course: 3, urate: 5}, ErrOutOfRange},
	}
	for _, c := range cases {
		if _, err := c.in.run(); !errors.Is(err, c.want) {
			t.Errorf("%s: error = %v, want %v", c.name, err, c.want)
		}
	}
	for _, u := range []float64{0, 30} {
		if _, err := (goutIn{urate: u}).run(); err != nil {
			t.Errorf("urate %v: error = %v, want nil (range edge)", u, err)
		}
	}
}

func TestGoutACREULAR2015Score_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = GoutACREULAR2015Score(2, 3, 2, false, 9, false, false, false) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkGoutACREULAR2015Score(b *testing.B) {
	b.ReportAllocs()
	j, c, tc, u := 2, 3, 2, 9.0
	for b.Loop() {
		sinkInt, _ = GoutACREULAR2015Score(j, c, tc, false, u, false, false, false)
	}
}
