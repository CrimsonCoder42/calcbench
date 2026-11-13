package neuro

import (
	"errors"
	"math"
	"testing"
)

type dragonIn struct {
	dense, early, mrs bool
	age, glu, ott     float64
	nihss             int
}

func (d dragonIn) run() (int, error) {
	return DRAGONScore(d.dense, d.early, d.mrs, d.age, d.glu, d.ott, d.nihss)
}

// HAND-DERIVED from Strbian 2012 abstract points; each edge on both sides.
func TestDRAGONScore_example(t *testing.T) {
	base := dragonIn{age: 64.9, glu: 8, ott: 90, nihss: 4}
	tests := []struct {
		name string
		mod  func(*dragonIn)
		want int
	}{
		{"none at edges", func(*dragonIn) {}, 0},
		{"dense only", func(d *dragonIn) { d.dense = true }, 1},
		{"early only", func(d *dragonIn) { d.early = true }, 1},
		{"both signs", func(d *dragonIn) { d.dense, d.early = true, true }, 2},
		{"mRS > 1", func(d *dragonIn) { d.mrs = true }, 1},
		{"age 65", func(d *dragonIn) { d.age = 65 }, 1},
		{"age 79.9", func(d *dragonIn) { d.age = 79.9 }, 1},
		{"age 80", func(d *dragonIn) { d.age = 80 }, 2},
		{"glucose 8.01", func(d *dragonIn) { d.glu = 8.01 }, 1},
		{"OTT 91", func(d *dragonIn) { d.ott = 91 }, 1},
		{"NIHSS 5", func(d *dragonIn) { d.nihss = 5 }, 1},
		{"NIHSS 9", func(d *dragonIn) { d.nihss = 9 }, 1},
		{"NIHSS 10", func(d *dragonIn) { d.nihss = 10 }, 2},
		{"NIHSS 15", func(d *dragonIn) { d.nihss = 15 }, 2},
		{"NIHSS 16", func(d *dragonIn) { d.nihss = 16 }, 3},
		{"max", func(d *dragonIn) { *d = dragonIn{true, true, true, 82, 9, 100, 18} }, 10},
	}
	for _, tc := range tests {
		in := base
		tc.mod(&in)
		if got, err := in.run(); err != nil || got != tc.want {
			t.Errorf("%s: DRAGONScore = (%d, %v), want %d", tc.name, got, err, tc.want)
		}
	}
}

func TestDRAGONScore_errors(t *testing.T) {
	ok := dragonIn{age: 70, glu: 6, ott: 60, nihss: 8}
	tests := []struct {
		mod  func(*dragonIn)
		want error
	}{
		{func(d *dragonIn) { d.age = math.NaN() }, ErrNotFinite},
		{func(d *dragonIn) { d.glu = math.Inf(1) }, ErrNotFinite},
		{func(d *dragonIn) { d.ott = math.Inf(-1) }, ErrNotFinite},
		{func(d *dragonIn) { d.age = -1 }, ErrOutOfRange},
		{func(d *dragonIn) { d.glu = -0.1 }, ErrOutOfRange},
		{func(d *dragonIn) { d.ott = -1 }, ErrOutOfRange},
		{func(d *dragonIn) { d.nihss = -1 }, ErrOutOfRange},
	}
	for i, tc := range tests {
		in := ok
		tc.mod(&in)
		if _, err := in.run(); !errors.Is(err, tc.want) {
			t.Errorf("case %d: error = %v, want %v", i, err, tc.want)
		}
	}
}

func TestDRAGONScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = DRAGONScore(true, false, false, 70, 9, 100, 12) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkDRAGONScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = DRAGONScore(true, false, false, 70, 9, 100, 12)
	}
}
