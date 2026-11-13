package seq

import (
	"errors"
	"math/rand/v2"
	"slices"
	"testing"
)

// isZero marks a terminator element (0) that closes a batch.
func isZero(v int) bool { return v == 0 }

// slowBatches is the oracle: walk the input and cut after a terminator or at maxSize.
func slowBatches(in []int, maxSize int, ends func(int) bool) [][]int {
	var out [][]int
	start := 0
	for i, v := range in {
		if i+1-start == maxSize || (ends != nil && ends(v)) {
			out = append(out, slices.Clone(in[start:i+1]))
			start = i + 1
		}
	}
	if start < len(in) {
		out = append(out, slices.Clone(in[start:]))
	}
	return out
}

// Hand-derived: records 5 0 | 1 2 3 4 | 0 | 6 with terminator 0 and at most 4 per batch:
// [5 0] (terminator), [1 2 3 4] (size), [0] (terminator), [6] (the rest).
func TestBatch_example(t *testing.T) {
	s, err := Batch(slices.Values([]int{5, 0, 1, 2, 3, 4, 0, 6}), 4, isZero)
	if err != nil {
		t.Fatal(err)
	}
	want := [][]int{{5, 0}, {1, 2, 3, 4}, {0}, {6}}
	if got := collectGroups(s); !equalGroups(got, want) {
		t.Errorf("Batch = %v, want %v", got, want)
	}
}

func TestBatch_cases(t *testing.T) {
	tests := []struct {
		name    string
		in      []int
		maxSize int
		ends    func(int) bool
		want    [][]int
	}{
		{"empty", nil, 3, isZero, nil},
		{"one terminator", []int{0}, 3, isZero, [][]int{{0}}},
		{"one plain", []int{4}, 3, isZero, [][]int{{4}}},
		{"terminator last", []int{1, 0}, 3, isZero, [][]int{{1, 0}}},
		{"consecutive terminators", []int{0, 0, 0}, 3, isZero, [][]int{{0}, {0}, {0}}},
		{"size 1", []int{1, 2}, 1, isZero, [][]int{{1}, {2}}},
		{"nil predicate is Chunk", []int{1, 0, 2, 3}, 2, nil, [][]int{{1, 0}, {2, 3}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := Batch(slices.Values(tt.in), tt.maxSize, tt.ends)
			if err != nil {
				t.Fatal(err)
			}
			if got := collectGroups(s); !equalGroups(got, tt.want) {
				t.Errorf("Batch(%v, %d) = %v, want %v", tt.in, tt.maxSize, got, tt.want)
			}
		})
	}
}

func TestBatch_errors(t *testing.T) {
	for _, size := range []int{0, -5} {
		s, err := Batch(slices.Values([]int{1}), size, isZero)
		if !errors.Is(err, ErrSize) || s != nil {
			t.Errorf("Batch(maxSize = %d) = (%v, %v), want (nil, ErrSize)", size, s, err)
		}
	}
}

func TestBatch_nilSeq(t *testing.T) {
	s, _ := Batch[int](nil, 2, isZero)
	if got := collectGroups(s); len(got) != 0 {
		t.Errorf("Batch(nil) = %v, want empty", got)
	}
}

func TestBatch_stopsUpstreamOnBreak(t *testing.T) {
	var p probe
	s, _ := Batch(p.ints(-1), 100, func(v int) bool { return v%5 == 4 })
	for range s {
		break // the first batch is 0 1 2 3 4
	}
	if p.produced != 5 || !p.stopped {
		t.Errorf("after break at the 1st batch: produced %d, stopped %v; want 5, true", p.produced, p.stopped)
	}
}

func TestBatch_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		in := randomInts(rng, rng.IntN(60))
		size := 1 + rng.IntN(8)
		s, _ := Batch(slices.Values(in), size, isZero)
		if got, want := collectGroups(s), slowBatches(in, size, isZero); !equalGroups(got, want) {
			t.Fatalf("Batch(%v, %d) = %v, want %v", in, size, got, want)
		}
	}
}

func TestBatch_noAllocationPerGroup(t *testing.T) {
	run := func(n int) float64 {
		s, _ := Batch(count(n), 8, func(v int) bool { return v%3 == 0 })
		return testing.AllocsPerRun(20, func() {
			for range s {
			}
		})
	}
	if small, large := run(16), run(16_000); small != large {
		t.Errorf("Batch allocates %v times over 16 elements and %v over 16000; want equal", small, large)
	}
}

func BenchmarkBatch(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			s, _ := Batch(count(size.n), 64, func(v int) bool { return v%50 == 0 })
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
