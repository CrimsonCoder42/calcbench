package dp

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// jumpsSlow is breadth-first search over the explicit edges i → i+1 … i+maxJump[i] (CLRS §22.2):
// O(n²) on these inputs. It returns −1 when the last index is unreachable.
func jumpsSlow(maxJump []int) int {
	n := len(maxJump)
	if n == 0 {
		return -1
	}
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	dist[0] = 0
	queue := []int{0}
	for len(queue) > 0 {
		i := queue[0]
		queue = queue[1:]
		for j := i + 1; j <= min(n-1, i+maxJump[i]); j++ {
			if dist[j] < 0 {
				dist[j] = dist[i] + 1
				queue = append(queue, j)
			}
		}
	}
	return dist[n-1]
}

// Worked by hand: from [2 0 1 0], index 0 reaches 2, and 2 reaches 3: two moves.
func TestMinJumpsToEnd_example(t *testing.T) {
	if got, ok, err := MinJumpsToEnd([]int{2, 0, 1, 0}); err != nil || !ok || got != 2 {
		t.Errorf("MinJumpsToEnd([2 0 1 0]) = (%d, %v, %v), want (2, true, nil)", got, ok, err)
	}
}

func TestJumps_cases(t *testing.T) {
	tests := []struct {
		name  string
		jumps []int
		min   int // −1: unreachable
	}{
		{"empty", nil, -1},
		{"one", []int{0}, 0},
		{"two, stuck", []int{0, 5}, -1},
		{"two", []int{1, 0}, 1},
		{"blocked by a zero", []int{1, 0, 2}, -1},
		{"one long jump", []int{5, 0, 0}, 1},
		{"huge jump, no overflow", []int{math.MaxInt, 0, 0, 0}, 1},
		{"all ones", []int{1, 1, 1, 1}, 3},
		{"two hops through the longer reach", []int{2, 4, 1, 1, 1, 3}, 2},
		{"zeros after the end are fine", []int{3, 0, 0, 0}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			can, err := CanJumpToEnd(tt.jumps)
			if err != nil || can != (tt.min >= 0) {
				t.Errorf("CanJumpToEnd(%v) = (%v, %v), want %v", tt.jumps, can, err, tt.min >= 0)
			}
			got, ok, err := MinJumpsToEnd(tt.jumps)
			if err != nil || ok != (tt.min >= 0) || (ok && got != tt.min) {
				t.Errorf("MinJumpsToEnd(%v) = (%d, %v, %v), want %d", tt.jumps, got, ok, err, tt.min)
			}
		})
	}
}

func TestJumps_errors(t *testing.T) {
	in := []int{2, -1, 0}
	if _, err := CanJumpToEnd(in); !errors.Is(err, ErrNegative) {
		t.Errorf("CanJumpToEnd: err = %v", err)
	}
	if _, _, err := MinJumpsToEnd(in); !errors.Is(err, ErrNegative) {
		t.Errorf("MinJumpsToEnd: err = %v", err)
	}
}

func TestJumps_allocations(t *testing.T) {
	in := []int{2, 4, 1, 1, 1, 3}
	if n := testing.AllocsPerRun(100, func() { _, _, _ = MinJumpsToEnd(in); _, _ = CanJumpToEnd(in) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func TestJumps_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		jumps := make([]int, rng.IntN(25))
		for i := range jumps {
			jumps[i] = rng.IntN(4)
		}
		want := jumpsSlow(jumps)
		if can, err := CanJumpToEnd(jumps); err != nil || can != (want >= 0) {
			t.Fatalf("trial %d: CanJumpToEnd(%v) = %v, want %v", trial, jumps, can, want >= 0)
		}
		got, ok, err := MinJumpsToEnd(jumps)
		if err != nil || ok != (want >= 0) || (ok && got != want) {
			t.Fatalf("trial %d: MinJumpsToEnd(%v) = (%d, %v), want %d", trial, jumps, got, ok, want)
		}
	}
}

func benchJumps(n int) []int {
	rng := rand.New(rand.NewPCG(1, 2))
	out := make([]int, n)
	for i := range out {
		out[i] = 1 + rng.IntN(3) // never 0: the whole array is scanned
	}
	return out
}

func BenchmarkCanJumpToEnd(b *testing.B) {
	for _, size := range linearSizes {
		in := benchJumps(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = CanJumpToEnd(in)
			}
		})
	}
}

func BenchmarkMinJumpsToEnd(b *testing.B) {
	for _, size := range linearSizes {
		in := benchJumps(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = MinJumpsToEnd(in)
			}
		})
	}
}
