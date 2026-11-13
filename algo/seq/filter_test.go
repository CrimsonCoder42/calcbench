package seq

import (
	"math/rand/v2"
	"slices"
	"testing"
)

func isEven(v int) bool { return v%2 == 0 }

// Hand-derived: the even values of 1 … 6 are 2 4 6.
func TestFilter_example(t *testing.T) {
	got := slices.Collect(Filter(slices.Values([]int{1, 2, 3, 4, 5, 6}), isEven))
	if !slices.Equal(got, []int{2, 4, 6}) {
		t.Errorf("Filter(1…6, even) = %v, want [2 4 6]", got)
	}
}

func TestFilter_cases(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want []int
	}{
		{"empty", nil, nil},
		{"one kept", []int{2}, []int{2}},
		{"one dropped", []int{1}, nil},
		{"none kept", []int{1, 3, 5}, nil},
		{"all kept", []int{0, -2, 4}, []int{0, -2, 4}},
		{"duplicates kept in order", []int{2, 1, 2, 2}, []int{2, 2, 2}},
		{"negative", []int{-3, -4}, []int{-4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Collect(Filter(slices.Values(tt.in), isEven))
			if !slices.Equal(got, tt.want) {
				t.Errorf("Filter(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestFilter_nil(t *testing.T) {
	if got := slices.Collect(Filter(nil, isEven)); len(got) != 0 {
		t.Errorf("Filter(nil seq) = %v, want empty", got)
	}
	var p probe
	if got := slices.Collect(Filter(p.ints(4), nil)); len(got) != 0 || p.produced != 0 {
		t.Errorf("Filter(nil keep) = %v, produced %d, want empty and 0", got, p.produced)
	}
}

func TestFilter_stopsUpstreamOnBreak(t *testing.T) {
	var p probe
	for v := range Filter(p.ints(-1), isEven) {
		if v == 4 {
			break
		}
	}
	if p.produced != 5 || !p.stopped {
		t.Errorf("after break at 4: produced %d, stopped %v; want 5, true", p.produced, p.stopped)
	}
}

func TestFilter_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		in := randomInts(rng, rng.IntN(50))
		var want []int
		for _, v := range in {
			if v > 4 {
				want = append(want, v)
			}
		}
		got := slices.Collect(Filter(slices.Values(in), func(v int) bool { return v > 4 }))
		if !slices.Equal(got, want) {
			t.Fatalf("Filter(%v) = %v, want %v", in, got, want)
		}
	}
}

func TestFilter_noAllocationPerElement(t *testing.T) {
	run := func(n int) float64 {
		s := Filter(count(n), isEven)
		return testing.AllocsPerRun(20, func() {
			for range s {
			}
		})
	}
	if small, large := run(10), run(10_000); small != large {
		t.Errorf("Filter allocates %v times over 10 elements and %v over 10000; want equal", small, large)
	}
}

func BenchmarkFilter(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			s := Filter(count(size.n), isEven)
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
