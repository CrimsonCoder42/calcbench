package seq

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// slowWindows is the oracle: every slice in[i : i+width].
func slowWindows(in []int, width int) [][]int {
	var out [][]int
	for i := 0; i+width <= len(in); i++ {
		out = append(out, slices.Clone(in[i:i+width]))
	}
	return out
}

// Hand-derived: windows of width 3 over 1 … 5 are [1 2 3] [2 3 4] [3 4 5].
func TestWindow_example(t *testing.T) {
	s, err := Window(slices.Values([]int{1, 2, 3, 4, 5}), 3)
	if err != nil {
		t.Fatal(err)
	}
	want := [][]int{{1, 2, 3}, {2, 3, 4}, {3, 4, 5}}
	if got := collectGroups(s); !equalGroups(got, want) {
		t.Errorf("Window(1…5, 3) = %v, want %v", got, want)
	}
}

func TestWindow_cases(t *testing.T) {
	tests := []struct {
		name  string
		in    []int
		width int
		want  [][]int
	}{
		{"empty", nil, 2, nil},
		{"shorter than width", []int{1, 2}, 3, nil},
		{"largest width", []int{1, 2}, math.MaxInt, nil},
		{"exactly width", []int{1, 2}, 2, [][]int{{1, 2}}},
		{"width 1", []int{3, 1}, 1, [][]int{{3}, {1}}},
		{"duplicates and negatives", []int{-1, -1, 0, -1}, 2, [][]int{{-1, -1}, {-1, 0}, {0, -1}}},
		{"wraps the ring twice", []int{1, 2, 3, 4, 5, 6, 7}, 2, [][]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 7}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := Window(slices.Values(tt.in), tt.width)
			if err != nil {
				t.Fatal(err)
			}
			if got := collectGroups(s); !equalGroups(got, tt.want) {
				t.Errorf("Window(%v, %d) = %v, want %v", tt.in, tt.width, got, tt.want)
			}
		})
	}
}

func TestWindow_errors(t *testing.T) {
	for _, w := range []int{0, -1, math.MinInt} {
		s, err := Window(slices.Values([]int{1}), w)
		if !errors.Is(err, ErrSize) || s != nil {
			t.Errorf("Window(width = %d) = (%v, %v), want (nil, ErrSize)", w, s, err)
		}
	}
}

func TestWindow_nilSeq(t *testing.T) {
	s, _ := Window[int](nil, 2)
	if got := collectGroups(s); len(got) != 0 {
		t.Errorf("Window(nil) = %v, want empty", got)
	}
}

func TestWindow_stopsUpstreamOnBreak(t *testing.T) {
	var p probe
	s, _ := Window(p.ints(-1), 4)
	n := 0
	for range s {
		n++
		if n == 3 {
			break
		}
	}
	if p.produced != 6 || !p.stopped {
		t.Errorf("after break at the 3rd window: produced %d, stopped %v; want 6, true", p.produced, p.stopped)
	}
}

func TestWindow_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		in := randomInts(rng, rng.IntN(60))
		width := 1 + rng.IntN(12)
		s, _ := Window(slices.Values(in), width)
		if got, want := collectGroups(s), slowWindows(in, width); !equalGroups(got, want) {
			t.Fatalf("Window(%v, %d) = %v, want %v", in, width, got, want)
		}
	}
}

// The ring is allocated once per range loop; no window allocates.
func TestWindow_noAllocationPerWindow(t *testing.T) {
	run := func(n int) float64 {
		s, _ := Window(count(n), 8)
		return testing.AllocsPerRun(20, func() {
			for range s {
			}
		})
	}
	if small, large := run(16), run(16_000); small != large {
		t.Errorf("Window allocates %v times over 16 elements and %v over 16000; want equal", small, large)
	}
}

func BenchmarkWindow(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			s, _ := Window(count(size.n), 50)
			b.ReportAllocs()
			for b.Loop() {
				total := 0
				for w := range s {
					total += w[0]
				}
				_ = total
			}
		})
	}
}
