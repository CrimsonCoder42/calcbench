package seq

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Hand-derived: the first 3 elements of 0 1 2 … (an infinite iterator) are 0 1 2, and the upstream
// is asked for exactly 3 elements.
func TestTake_example(t *testing.T) {
	var p probe
	s, err := Take(p.ints(-1), 3)
	if err != nil {
		t.Fatal(err)
	}
	got := slices.Collect(s)
	if !slices.Equal(got, []int{0, 1, 2}) || p.produced != 3 || !p.stopped {
		t.Errorf("Take(0 1 2 …, 3) = %v, produced %d, stopped %v; want [0 1 2], 3, true", got, p.produced, p.stopped)
	}
}

func TestTake_cases(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		n    int
		want []int
	}{
		{"empty", nil, 3, nil},
		{"zero", []int{1, 2}, 0, nil},
		{"one", []int{5}, 1, []int{5}},
		{"fewer than n", []int{1, 2}, 5, []int{1, 2}},
		{"exactly n", []int{1, 2}, 2, []int{1, 2}},
		{"duplicates", []int{7, 7, 7}, 2, []int{7, 7}},
		{"largest n", []int{1, 2}, math.MaxInt, []int{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := Take(slices.Values(tt.in), tt.n)
			if err != nil {
				t.Fatal(err)
			}
			if got := slices.Collect(s); !slices.Equal(got, tt.want) {
				t.Errorf("Take(%v, %d) = %v, want %v", tt.in, tt.n, got, tt.want)
			}
		})
	}
}

func TestTake_zeroNeverStartsUpstream(t *testing.T) {
	var p probe
	s, _ := Take(p.ints(-1), 0)
	for range s {
	}
	if p.produced != 0 {
		t.Errorf("Take(n = 0) asked the upstream for %d elements, want 0", p.produced)
	}
}

func TestTake_doesNotReadPastN(t *testing.T) {
	var p probe
	s, _ := Take(p.ints(5), 2)
	for range s {
	}
	// Take stops the upstream at its n-th yield: it never asks for element n+1.
	if p.produced != 2 || !p.stopped || p.finished {
		t.Errorf("produced %d, stopped %v, finished %v; want 2, true, false", p.produced, p.stopped, p.finished)
	}
}

func TestTake_errors(t *testing.T) {
	for _, n := range []int{-1, math.MinInt} {
		s, err := Take(slices.Values([]int{1}), n)
		if !errors.Is(err, ErrCount) || s != nil {
			t.Errorf("Take(n = %d) = (%v, %v), want (nil, ErrCount)", n, s, err)
		}
	}
}

func TestTake_nilSeq(t *testing.T) {
	s, err := Take[int](nil, 3)
	if err != nil {
		t.Fatal(err)
	}
	if got := slices.Collect(s); len(got) != 0 {
		t.Errorf("Take(nil) = %v, want empty", got)
	}
}

func TestTake_stopsUpstreamOnBreak(t *testing.T) {
	var p probe
	s, _ := Take(p.ints(-1), 10)
	for v := range s {
		if v == 1 {
			break
		}
	}
	if p.produced != 2 || !p.stopped {
		t.Errorf("after break at 1: produced %d, stopped %v; want 2, true", p.produced, p.stopped)
	}
}

func TestTake_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		in := randomInts(rng, rng.IntN(30))
		n := rng.IntN(40)
		want := in[:min(n, len(in))]
		s, err := Take(slices.Values(in), n)
		if err != nil {
			t.Fatal(err)
		}
		if got := slices.Collect(s); !slices.Equal(got, want) {
			t.Fatalf("Take(%v, %d) = %v, want %v", in, n, got, want)
		}
	}
}

func TestTake_noAllocationPerElement(t *testing.T) {
	run := func(n int) float64 {
		s, _ := Take(count(n), n)
		return testing.AllocsPerRun(20, func() {
			for range s {
			}
		})
	}
	if small, large := run(10), run(10_000); small != large {
		t.Errorf("Take allocates %v times over 10 elements and %v over 10000; want equal", small, large)
	}
}

func BenchmarkTake(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			s, _ := Take(count(2*size.n), size.n)
			b.ReportAllocs()
			for b.Loop() {
				sum := 0
				for v := range s {
					sum += v
				}
				_ = sum
			}
		})
	}
}
