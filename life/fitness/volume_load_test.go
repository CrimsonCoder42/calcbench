package fitness

import (
	"math"
	"testing"
)

// HAND-DERIVED from McBride 2009's definition "volume load (VL) (repetitions [no.] x external load [kg])":
// a set of 10 repetitions (as in McBride's hypertrophy protocol) at an external load of 100 kg → 1000 kg·rep. Exact: integer-valued
// products well below 2^53.
func TestVolumeLoad_example(t *testing.T) {
	if got, err := VolumeLoad(10, 100); err != nil || got != 1000 {
		t.Errorf("VolumeLoad(10, 100) = (%v, %v), want 1000", got, err)
	}
}

// parity: ledger volume_load, source cases 1–5 (exact: products of small integers and halves).
func TestVolumeLoad_parity(t *testing.T) {
	cases := []struct {
		reps       int
		load, want float64
	}{
		{10, 135, 1350}, // parity: ledger volume_load, source case 1
		{5, 225, 1125},  // parity: ledger volume_load, source case 2
		{12, 0, 0},      // parity: ledger volume_load, source case 3 (no external load)
		{1, 405, 405},   // parity: ledger volume_load, source case 4
		{8, 62.5, 500},  // parity: ledger volume_load, source case 5
	}
	for _, c := range cases {
		if got, err := VolumeLoad(c.reps, c.load); err != nil || got != c.want {
			t.Errorf("VolumeLoad(%d, %v) = (%v, %v), want %v", c.reps, c.load, got, err, c.want)
		}
	}
}

// Edges: 0 reps and 0 load accepted (0); −1 refused on each; overflow reported.
func TestVolumeLoad_errors(t *testing.T) {
	if got, err := VolumeLoad(0, 100); err != nil || got != 0 {
		t.Errorf("VolumeLoad(0, 100) = (%v, %v), want 0", got, err)
	}
	call := func(r int, l float64) func() error {
		return func() error { _, err := VolumeLoad(r, l); return err }
	}
	runErrCases(t, []errCase{
		{"reps negative", call(-1, 100), ErrOutOfRange},
		{"load negative", call(10, -1), ErrOutOfRange},
		{"load NaN", call(10, nan), ErrNotFinite},
		{"load +Inf", call(10, pInf), ErrNotFinite},
		{"load -Inf", call(10, nInf), ErrNotFinite},
		{"product overflows", call(2, math.MaxFloat64), ErrDomain},
	})
}

func TestVolumeLoad_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = VolumeLoad(10, 135) }); n != 0 {
		t.Errorf("VolumeLoad allocates %v times, want 0", n)
	}
}

func BenchmarkVolumeLoad(b *testing.B) {
	b.ReportAllocs()
	r, l := 10, 135.0
	for b.Loop() {
		sink, _ = VolumeLoad(r, l)
	}
}
