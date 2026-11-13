package seq

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// collectGroups copies each yielded group, since the iterator reuses one slice. Shared with the Batch
// and Window tests.
func collectGroups[T any](s func(yield func([]T) bool)) [][]T {
	var out [][]T
	for g := range s {
		out = append(out, slices.Clone(g))
	}
	return out
}

func equalGroups(a, b [][]int) bool {
	return slices.EqualFunc(a, b, func(x, y []int) bool { return slices.Equal(x, y) })
}

// slowChunks is the oracle: re-slicing the whole input.
func slowChunks(in []int, size int) [][]int {
	var out [][]int
	for i := 0; i < len(in); i += size {
		out = append(out, slices.Clone(in[i:min(i+size, len(in))]))
	}
	return out
}

// Hand-derived: 1 … 7 in groups of 3 is [1 2 3] [4 5 6] [7].
func TestChunk_example(t *testing.T) {
	s, err := Chunk(slices.Values([]int{1, 2, 3, 4, 5, 6, 7}), 3)
	if err != nil {
		t.Fatal(err)
	}
	want := [][]int{{1, 2, 3}, {4, 5, 6}, {7}}
	if got := collectGroups(s); !equalGroups(got, want) {
		t.Errorf("Chunk(1…7, 3) = %v, want %v", got, want)
	}
}

func TestChunk_cases(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		size int
		want [][]int
	}{
		{"empty", nil, 2, nil},
		{"one", []int{9}, 2, [][]int{{9}}},
		{"size 1", []int{1, 2}, 1, [][]int{{1}, {2}}},
		{"exact multiple", []int{1, 2, 3, 4}, 2, [][]int{{1, 2}, {3, 4}}},
		{"size larger than input", []int{1, 2}, 5, [][]int{{1, 2}}},
		{"largest size", []int{1, 2}, math.MaxInt, [][]int{{1, 2}}},
		{"duplicates and negatives", []int{-1, -1, 0}, 2, [][]int{{-1, -1}, {0}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := Chunk(slices.Values(tt.in), tt.size)
			if err != nil {
				t.Fatal(err)
			}
			if got := collectGroups(s); !equalGroups(got, tt.want) {
				t.Errorf("Chunk(%v, %d) = %v, want %v", tt.in, tt.size, got, tt.want)
			}
		})
	}
}

func TestChunk_errors(t *testing.T) {
	for _, size := range []int{0, -1, math.MinInt} {
		s, err := Chunk(slices.Values([]int{1}), size)
		if !errors.Is(err, ErrSize) || s != nil {
			t.Errorf("Chunk(size = %d) = (%v, %v), want (nil, ErrSize)", size, s, err)
		}
	}
}

func TestChunk_nilSeq(t *testing.T) {
	s, _ := Chunk[int](nil, 2)
	if got := collectGroups(s); len(got) != 0 {
		t.Errorf("Chunk(nil) = %v, want empty", got)
	}
}

func TestChunk_reusesOneSlice(t *testing.T) {
	s, _ := Chunk(slices.Values([]int{1, 2, 3, 4}), 2)
	var first []int
	for g := range s {
		if first == nil {
			first = g
			continue
		}
		if &first[0] != &g[0] {
			t.Errorf("Chunk yielded a new slice; the package documents one reused buffer")
		}
	}
}

func TestChunk_stopsUpstreamOnBreak(t *testing.T) {
	var p probe
	s, _ := Chunk(p.ints(-1), 3)
	groups := 0
	for range s {
		groups++
		if groups == 2 {
			break
		}
	}
	if p.produced != 6 || !p.stopped {
		t.Errorf("after break at the 2nd group: produced %d, stopped %v; want 6, true", p.produced, p.stopped)
	}
}

func TestChunk_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		in := randomInts(rng, rng.IntN(60))
		size := 1 + rng.IntN(10)
		s, _ := Chunk(slices.Values(in), size)
		if got, want := collectGroups(s), slowChunks(in, size); !equalGroups(got, want) {
			t.Fatalf("Chunk(%v, %d) = %v, want %v", in, size, got, want)
		}
	}
}

// The buffer grows to size once per range loop, then no group allocates.
func TestChunk_noAllocationPerGroup(t *testing.T) {
	run := func(n int) float64 {
		s, _ := Chunk(count(n), 8)
		return testing.AllocsPerRun(20, func() {
			for range s {
			}
		})
	}
	if small, large := run(16), run(16_000); small != large {
		t.Errorf("Chunk allocates %v times over 16 elements and %v over 16000; want equal", small, large)
	}
}

func BenchmarkChunk(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			s, _ := Chunk(count(size.n), 64)
			b.ReportAllocs()
			for b.Loop() {
				total := 0
				for g := range s {
					total += len(g)
				}
				_ = total
			}
		})
	}
}
