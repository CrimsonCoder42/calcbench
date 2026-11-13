package list

import (
	"errors"
	"testing"
)

// Worked example (hand-derived, catalog/list.md ALG-LIST-010): 5 nodes → 3; 4 nodes → 3 (second middle).
// The same two shapes are the problem-site pattern, in our words.
func TestMiddle_workedExample(t *testing.T) {
	for _, tt := range []struct {
		in   []int
		want int
	}{{[]int{1, 2, 3, 4, 5}, 3}, {[]int{1, 2, 3, 4}, 3}} {
		m, err := Middle(FromSlice(tt.in))
		if err != nil || m == nil || m.Value != tt.want {
			t.Errorf("Middle(%v) = (%v, %v), want %d", tt.in, m, err, tt.want)
		}
	}
}

func TestMiddle_everyLength(t *testing.T) {
	if m, err := Middle[int](nil); m != nil || err != nil {
		t.Errorf("Middle(nil) = (%v, %v)", m, err)
	}
	for n := 1; n <= 40; n++ {
		m, err := Middle(FromSlice(countUpTo(n)))
		if err != nil || m.Value != n/2 {
			t.Errorf("Middle(%d nodes) = (%v, %v), want index %d", n, m, err, n/2)
		}
	}
}

func TestMiddle_cycle(t *testing.T) {
	for n := 1; n <= 12; n++ {
		for entry := range n {
			head, _ := cyclic(countUpTo(n), entry)
			if m, err := Middle(head); !errors.Is(err, ErrCycle) || m != nil {
				t.Fatalf("n=%d entry=%d: (%v, %v), want ErrCycle", n, entry, m, err)
			}
		}
	}
}

func TestMiddle_noAlloc(t *testing.T) {
	head := FromSlice(countUpTo(101))
	if a := testing.AllocsPerRun(100, func() { _, _ = Middle(head) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkMiddle(b *testing.B) {
	for _, s := range benchSizes {
		head := FromSlice(countUpTo(s.n))
		b.Run(s.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = Middle(head)
			}
		})
	}
}
