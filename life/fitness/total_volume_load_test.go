package fitness

import (
	"errors"
	"fmt"
	"math"
	"testing"
)

// parity: ledger VolumeLoad (series), source case 1 (10 × 100 + 8 × 105 + 12 × 0 = 1840).
func TestTotalVolumeLoad_parity(t *testing.T) {
	got, err := TotalVolumeLoad([]LoadSet{{10, 100}, {8, 105}, {12, 0}})
	if err != nil || got != 1840 {
		t.Errorf("got %v, %v; want 1840", got, err)
	}
}

func TestTotalVolumeLoad_cases(t *testing.T) {
	if got, err := TotalVolumeLoad(nil); err != nil || got != 0 {
		t.Errorf("no sets: %v, %v", got, err)
	}
	if got, _ := TotalVolumeLoad([]LoadSet{{5, 60}}); got != 300 {
		t.Errorf("one set: %v", got)
	}
	// Compensated: a huge set does not swallow many small ones.
	sets := []LoadSet{{1, 1e16}}
	for range 10 {
		sets = append(sets, LoadSet{1, 1})
	}
	if got, _ := TotalVolumeLoad(sets); got != 1e16+10 {
		t.Errorf("compensated: %v, want 1e16+10", got)
	}
}

func TestTotalVolumeLoad_errors(t *testing.T) {
	cases := []struct {
		name string
		sets []LoadSet
		want error
	}{
		{"negative reps", []LoadSet{{5, 60}, {-1, 60}}, ErrOutOfRange},
		{"negative load", []LoadSet{{5, -60}}, ErrOutOfRange},
		{"NaN load", []LoadSet{{5, math.NaN()}}, ErrNotFinite},
		{"Inf load", []LoadSet{{5, math.Inf(1)}}, ErrNotFinite},
		{"set overflow", []LoadSet{{2, math.MaxFloat64}}, ErrDomain},
		{"total overflow", []LoadSet{{1, math.MaxFloat64}, {1, math.MaxFloat64}}, ErrDomain},
	}
	for _, c := range cases {
		if _, err := TotalVolumeLoad(c.sets); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestTotalVolumeLoad_noAlloc(t *testing.T) {
	sets := []LoadSet{{10, 100}, {8, 105}}
	if a := testing.AllocsPerRun(100, func() { _, _ = TotalVolumeLoad(sets) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkTotalVolumeLoad(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		sets := make([]LoadSet, n)
		for i := range sets {
			sets[i] = LoadSet{Reps: 1 + i%12, LoadKgOrLb: 20 + float64(i%100)}
		}
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := TotalVolumeLoad(sets); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
