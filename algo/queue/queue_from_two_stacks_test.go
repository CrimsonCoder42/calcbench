package queue

import (
	"math/rand/v2"
	"testing"
)

// benchSizes are the collection sizes of .claude/rules/go-complexity.md. Shared by the benchmarks.
var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

// Hand-derived (catalog): enqueue 1 2 3, dequeue → 1, enqueue 4, dequeue all → 2 3 4.
func TestQueueFromTwoStacks_example(t *testing.T) {
	var q QueueFromTwoStacks[int]
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	if v, ok := q.Dequeue(); !ok || v != 1 {
		t.Fatalf("Dequeue = (%d, %v), want (1, true)", v, ok)
	}
	q.Enqueue(4)
	for _, want := range []int{2, 3, 4} {
		if v, ok := q.Dequeue(); !ok || v != want {
			t.Fatalf("Dequeue = (%d, %v), want (%d, true)", v, ok, want)
		}
	}
	if q.Len() != 0 {
		t.Errorf("Len = %d, want 0", q.Len())
	}
}

func TestQueueFromTwoStacks_empty(t *testing.T) {
	var q QueueFromTwoStacks[string]
	if v, ok := q.Dequeue(); ok || v != "" {
		t.Errorf("Dequeue on empty = (%q, %v), want (\"\", false)", v, ok)
	}
	if v, ok := q.Front(); ok || v != "" {
		t.Errorf("Front on empty = (%q, %v), want (\"\", false)", v, ok)
	}
	q.Enqueue("a")
	if v, ok := q.Front(); !ok || v != "a" || q.Len() != 1 {
		t.Errorf("Front = (%q, %v), Len %d; want (\"a\", true), 1", v, ok, q.Len())
	}
	q.Dequeue()
	if _, ok := q.Dequeue(); ok {
		t.Error("Dequeue after emptying returned ok")
	}
}

func TestQueueFromTwoStacks_frontDoesNotRemove(t *testing.T) {
	var q QueueFromTwoStacks[int]
	q.Enqueue(7)
	q.Enqueue(8)
	for range 3 {
		if v, _ := q.Front(); v != 7 {
			t.Fatalf("Front = %d, want 7", v)
		}
	}
	if q.Len() != 2 {
		t.Errorf("Len = %d, want 2", q.Len())
	}
}

// Oracle: a plain slice used as a FIFO, under random operations.
func TestQueueFromTwoStacks_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 200 {
		var q QueueFromTwoStacks[int]
		var model []int
		for step := range 300 {
			switch rng.IntN(3) {
			case 0:
				v, ok := q.Dequeue()
				if ok != (len(model) > 0) || (ok && v != model[0]) {
					t.Fatalf("trial %d step %d: Dequeue = (%d, %v), model %v", trial, step, v, ok, model)
				}
				if ok {
					model = model[1:]
				}
			case 1:
				v, ok := q.Front()
				if ok != (len(model) > 0) || (ok && v != model[0]) {
					t.Fatalf("trial %d step %d: Front = (%d, %v), model %v", trial, step, v, ok, model)
				}
			default:
				v := rng.IntN(1000)
				q.Enqueue(v)
				model = append(model, v)
			}
			if q.Len() != len(model) {
				t.Fatalf("trial %d step %d: Len = %d, want %d", trial, step, q.Len(), len(model))
			}
		}
	}
}

// n operations: n/2 enqueues interleaved with n/2 dequeues. Amortized O(1) shows as flat ns per element.
func BenchmarkQueueFromTwoStacks(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				var q QueueFromTwoStacks[int]
				for i := range size.n / 2 {
					q.Enqueue(i)
					if i%2 == 1 {
						q.Dequeue()
					}
				}
				for q.Len() > 0 {
					q.Dequeue()
				}
			}
		})
	}
}
