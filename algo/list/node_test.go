package list

import (
	"errors"
	"math/rand/v2"
	"slices"
	"testing"
)

// benchSizes are the collection sizes of .claude/rules/go-complexity.md. Shared by the benchmarks.
var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

// values returns the list's values; it fails the test on a cycle. Shared by every test file.
func values[T any](t *testing.T, head *Node[T]) []T {
	t.Helper()
	got, err := AppendValues(nil, head)
	if err != nil {
		t.Fatalf("AppendValues: %v", err)
	}
	return got
}

// cyclic returns a list of vals whose last node links back to the node at index entry (0-based), and
// that entry node. Shared by every test file.
func cyclic(vals []int, entry int) (head, entryNode *Node[int]) {
	head = FromSlice(vals)
	last := head
	for last.Next != nil {
		last = last.Next
	}
	entryNode = nodeAt(head, entry)
	last.Next = entryNode
	return head, entryNode
}

// randInts returns n fixed-seed values in [0, span).
func randInts(rng *rand.Rand, n, span int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = rng.IntN(span)
	}
	return s
}

// countUpTo returns 0, 1, …, n-1.
func countUpTo(n int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = i
	}
	return s
}

func TestFromSlice_cases(t *testing.T) {
	tests := []struct {
		name string
		in   []int
	}{
		{"empty", nil},
		{"one", []int{7}},
		{"two", []int{1, 2}},
		{"duplicates", []int{3, 3, 3}},
		{"negative and zero", []int{-1, 0, -5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := FromSlice(tt.in)
			if (head == nil) != (len(tt.in) == 0) {
				t.Fatalf("FromSlice(%v) nil = %v", tt.in, head == nil)
			}
			if got := values(t, head); !slices.Equal(got, tt.in) {
				t.Errorf("values = %v, want %v", got, tt.in)
			}
		})
	}
}

func TestFromSlice_doesNotAliasInput(t *testing.T) {
	in := []int{1, 2, 3}
	head := FromSlice(in)
	head.Value = 9
	if in[0] != 1 {
		t.Errorf("input modified: %v", in)
	}
}

func TestFromSlice_oneAllocation(t *testing.T) {
	in := countUpTo(100)
	var sink *Node[int]
	if a := testing.AllocsPerRun(100, func() { sink = FromSlice(in) }); a != 1 || sink == nil {
		t.Errorf("FromSlice allocs = %v, want 1", a)
	}
}

// Worked example (hand-derived, catalog/list.md ALG-LIST-001…004): 1→2→3.
func TestAppendValues_workedExample(t *testing.T) {
	got, err := AppendValues([]int{0}, FromSlice([]int{1, 2, 3}))
	if err != nil || !slices.Equal(got, []int{0, 1, 2, 3}) {
		t.Errorf("AppendValues = (%v, %v), want ([0 1 2 3], nil)", got, err)
	}
}

func TestAppendValues_nil(t *testing.T) {
	got, err := AppendValues[int](nil, nil)
	if err != nil || len(got) != 0 {
		t.Errorf("AppendValues(nil) = (%v, %v)", got, err)
	}
}

func TestAppendValues_cycle(t *testing.T) {
	head, _ := cyclic([]int{1, 2, 3}, 0)
	dst := []int{5}
	got, err := AppendValues(dst, head)
	if !errors.Is(err, ErrCycle) || !slices.Equal(got, dst) {
		t.Errorf("AppendValues(cycle) = (%v, %v), want ([5], ErrCycle)", got, err)
	}
}

func TestAppendValues_noAllocWithCapacity(t *testing.T) {
	head := FromSlice(countUpTo(100))
	dst := make([]int, 0, 100)
	if a := testing.AllocsPerRun(100, func() { _, _ = AppendValues(dst[:0], head) }); a != 0 {
		t.Errorf("AppendValues allocs = %v, want 0", a)
	}
}

func TestLength_cases(t *testing.T) {
	for n := range 10 {
		got, err := Length(FromSlice(countUpTo(n)))
		if err != nil || got != n {
			t.Errorf("Length(%d nodes) = (%d, %v)", n, got, err)
		}
	}
}

func TestLength_cycles(t *testing.T) {
	selfLoop := &Node[int]{Value: 1}
	selfLoop.Next = selfLoop
	if _, err := Length(selfLoop); !errors.Is(err, ErrCycle) {
		t.Errorf("Length(self loop) err = %v, want ErrCycle", err)
	}
	for n := 1; n <= 12; n++ {
		for entry := range n {
			head, _ := cyclic(countUpTo(n), entry)
			if got, err := Length(head); !errors.Is(err, ErrCycle) || got != 0 {
				t.Errorf("Length(n=%d, entry=%d) = (%d, %v), want (0, ErrCycle)", n, entry, got, err)
			}
		}
	}
}

func BenchmarkFromSlice(b *testing.B) {
	for _, s := range benchSizes {
		in := countUpTo(s.n)
		b.Run(s.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = FromSlice(in)
			}
		})
	}
}

func BenchmarkAppendValues(b *testing.B) {
	for _, s := range benchSizes {
		head := FromSlice(countUpTo(s.n))
		dst := make([]int, 0, s.n)
		b.Run(s.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = AppendValues(dst[:0], head)
			}
		})
	}
}

func BenchmarkLength(b *testing.B) {
	for _, s := range benchSizes {
		head := FromSlice(countUpTo(s.n))
		b.Run(s.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = Length(head)
			}
		})
	}
}
