package list

import (
	"testing"
)

// Worked example (Knuth Vol. 2 §3.1 Ex. 6 notation, hand-derived): 1→2→3→4→5→(3): μ = 2, λ = 3.
func TestCycleStart_workedExample(t *testing.T) {
	head, entry := cyclic([]int{1, 2, 3, 4, 5}, 2)
	start, n := CycleStart(head)
	if start != entry || start.Value != 3 || n != 3 {
		t.Errorf("CycleStart = (%v, %d), want (node 3, 3)", start, n)
	}
	if !DetectCycle(head) {
		t.Error("DetectCycle = false")
	}
}

func TestDetectCycle_acyclic(t *testing.T) {
	for n := range 12 {
		head := FromSlice(countUpTo(n))
		if DetectCycle(head) {
			t.Errorf("DetectCycle(%d nodes, acyclic) = true", n)
		}
		if s, l := CycleStart(head); s != nil || l != 0 {
			t.Errorf("CycleStart(%d nodes, acyclic) = (%v, %d)", n, s, l)
		}
	}
}

// Every tail length μ and cycle length λ up to 15: the entry and λ are exact.
func TestCycleStart_allShapes(t *testing.T) {
	for n := 1; n <= 15; n++ {
		for entry := range n {
			head, want := cyclic(countUpTo(n), entry)
			if !DetectCycle(head) {
				t.Fatalf("n=%d entry=%d: DetectCycle = false", n, entry)
			}
			start, l := CycleStart(head)
			if start != want || l != n-entry {
				t.Fatalf("n=%d entry=%d: CycleStart = (%v, %d), want (%v, %d)", n, entry, start, l, want, n-entry)
			}
		}
	}
}

// Problem-site pattern, in our words: a single node pointing at itself.
func TestCycleStart_selfLoop(t *testing.T) {
	n := &Node[int]{Value: 1}
	n.Next = n
	if start, l := CycleStart(n); !DetectCycle(n) || start != n || l != 1 {
		t.Errorf("self loop: (%v, %d)", start, l)
	}
}

func TestDetectCycle_noAlloc(t *testing.T) {
	head, _ := cyclic(countUpTo(100), 40)
	if a := testing.AllocsPerRun(100, func() { _ = DetectCycle(head); _, _ = CycleStart(head) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkDetectCycle(b *testing.B) {
	for _, s := range benchSizes {
		head, _ := cyclic(countUpTo(s.n), s.n/2)
		b.Run(s.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = DetectCycle(head)
			}
		})
	}
}

func BenchmarkCycleStart(b *testing.B) {
	for _, s := range benchSizes {
		head, _ := cyclic(countUpTo(s.n), s.n/2)
		b.Run(s.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = CycleStart(head)
			}
		})
	}
}
