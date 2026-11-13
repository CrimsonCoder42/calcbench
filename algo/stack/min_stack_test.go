package stack

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Hand-derived (catalog): push 5 3 7 3 → Min 3; pop, pop → Min 3; pop → Min 5.
func TestMinStack_example(t *testing.T) {
	var s MinStack[int]
	for _, v := range []int{5, 3, 7, 3} {
		s.Push(v)
	}
	check := func(wantMin, wantLen int) {
		t.Helper()
		if m, ok := s.Min(); !ok || m != wantMin || s.Len() != wantLen {
			t.Errorf("Min = (%d, %v), Len = %d; want (%d, true), %d", m, ok, s.Len(), wantMin, wantLen)
		}
	}
	check(3, 4)
	s.Pop()
	s.Pop()
	check(3, 2)
	s.Pop()
	check(5, 1)
}

func TestMinStack_empty(t *testing.T) {
	var s MinStack[string]
	if v, ok := s.Pop(); ok || v != "" {
		t.Errorf("Pop on empty = (%q, %v), want (\"\", false)", v, ok)
	}
	if v, ok := s.Top(); ok || v != "" {
		t.Errorf("Top on empty = (%q, %v), want (\"\", false)", v, ok)
	}
	if v, ok := s.Min(); ok || v != "" {
		t.Errorf("Min on empty = (%q, %v), want (\"\", false)", v, ok)
	}
	if s.Len() != 0 {
		t.Errorf("Len on empty = %d", s.Len())
	}
	s.Push("b")
	s.Pop()
	if _, ok := s.Pop(); ok {
		t.Error("Pop after emptying returned ok")
	}
}

func TestMinStack_cases(t *testing.T) {
	var s MinStack[int]
	s.Push(math.MaxInt)
	s.Push(math.MinInt)
	s.Push(0)
	s.Push(math.MinInt) // a duplicate minimum survives one pop
	if top, _ := s.Top(); top != math.MinInt {
		t.Errorf("Top = %d, want MinInt", top)
	}
	s.Pop()
	if m, _ := s.Min(); m != math.MinInt {
		t.Errorf("Min after popping one of two minima = %d, want MinInt", m)
	}
	s.Pop()
	s.Pop()
	if m, _ := s.Min(); m != math.MaxInt {
		t.Errorf("Min = %d, want MaxInt", m)
	}
}

func TestMinStack_nan(t *testing.T) {
	var s MinStack[float64]
	s.Push(1)
	s.Push(math.NaN())
	s.Push(-1)
	if m, _ := s.Min(); !math.IsNaN(m) {
		t.Errorf("Min = %v, want NaN (the least value in cmp.Compare order)", m)
	}
}

// Oracle: a plain slice, its minimum found by a scan, under random pushes and pops.
func TestMinStack_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 200 {
		var s MinStack[int]
		var model []int
		for step := range 300 {
			if rng.IntN(3) == 0 {
				v, ok := s.Pop()
				if ok != (len(model) > 0) || (ok && v != model[len(model)-1]) {
					t.Fatalf("trial %d step %d: Pop = (%d, %v), model %v", trial, step, v, ok, model)
				}
				if ok {
					model = model[:len(model)-1]
				}
			} else {
				v := rng.IntN(50) - 25
				s.Push(v)
				model = append(model, v)
			}
			m, ok := s.Min()
			if ok != (len(model) > 0) || (ok && m != slices.Min(model)) || s.Len() != len(model) {
				t.Fatalf("trial %d step %d: Min = (%d, %v), Len %d, model %v", trial, step, m, ok, s.Len(), model)
			}
		}
	}
}

func BenchmarkMinStack(b *testing.B) {
	for _, size := range benchSizes {
		values := randomInts(rand.New(rand.NewPCG(1, 2)), size.n, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				var s MinStack[int]
				for _, v := range values {
					s.Push(v)
				}
				for s.Len() > 0 {
					s.Min()
					s.Pop()
				}
			}
		})
	}
}
