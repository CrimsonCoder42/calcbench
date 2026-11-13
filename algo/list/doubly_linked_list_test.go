package list

import (
	"math/rand/v2"
	"slices"
	"testing"
)

func dllValues[T any](l *DoublyLinkedList[T]) []T { return slices.Collect(l.All()) }

func dllBackward[T any](l *DoublyLinkedList[T]) []T { return slices.Collect(l.Backward()) }

// checkLinks walks both directions with Next and Prev and compares with want.
func checkLinks(t *testing.T, l *DoublyLinkedList[int], want []int) {
	t.Helper()
	if l.Len() != len(want) {
		t.Fatalf("Len = %d, want %d", l.Len(), len(want))
	}
	var fwd, back []int
	for e := l.Front(); e != nil; e = e.Next() {
		fwd = append(fwd, e.Value)
	}
	for e := l.Back(); e != nil; e = e.Prev() {
		back = append(back, e.Value)
	}
	rev := slices.Clone(want)
	slices.Reverse(rev)
	if !slices.Equal(fwd, want) || !slices.Equal(back, rev) {
		t.Fatalf("forward %v backward %v, want %v", fwd, back, want)
	}
	if got := dllValues(l); !slices.Equal(got, want) {
		t.Fatalf("All = %v, want %v", got, want)
	}
	if got := dllBackward(l); !slices.Equal(got, rev) {
		t.Fatalf("Backward = %v, want %v", got, rev)
	}
}

// Worked example (hand-derived, catalog/list.md ALG-LIST-005).
func TestDoublyLinkedList_workedExample(t *testing.T) {
	var l DoublyLinkedList[int]
	l.PushBack(1)
	l.PushBack(2)
	last := l.PushBack(3)
	front := l.PushFront(0)
	checkLinks(t, &l, []int{0, 1, 2, 3})
	if !l.MoveToBack(front) {
		t.Fatal("MoveToBack(front) = false")
	}
	checkLinks(t, &l, []int{1, 2, 3, 0})
	if !l.Remove(l.Back()) {
		t.Fatal("Remove(back) = false")
	}
	checkLinks(t, &l, []int{1, 2, 3})
	if !l.MoveToFront(last) {
		t.Fatal("MoveToFront = false")
	}
	checkLinks(t, &l, []int{3, 1, 2})
}

func TestDoublyLinkedList_zeroValue(t *testing.T) {
	var l DoublyLinkedList[string]
	if l.Len() != 0 || l.Front() != nil || l.Back() != nil || len(dllValues(&l)) != 0 || len(dllBackward(&l)) != 0 {
		t.Error("zero value is not an empty list")
	}
}

func TestDoublyLinkedList_nilList(t *testing.T) {
	var l *DoublyLinkedList[int]
	if l.Len() != 0 || l.Front() != nil || l.Back() != nil {
		t.Error("nil list reads not empty")
	}
	if l.PushFront(1) != nil || l.PushBack(1) != nil {
		t.Error("push on nil list returned an element")
	}
	e := &Element[int]{Value: 1}
	if l.Remove(e) || l.MoveToFront(e) || l.MoveToBack(e) || l.InsertBefore(1, e) != nil || l.InsertAfter(1, e) != nil {
		t.Error("nil list accepted an element")
	}
	if len(dllValues(l)) != 0 || len(dllBackward(l)) != 0 {
		t.Error("nil list iterated values")
	}
}

func TestDoublyLinkedList_insertBeforeAfter(t *testing.T) {
	var l DoublyLinkedList[int]
	mid := l.PushBack(2)
	l.InsertBefore(1, mid)
	l.InsertAfter(3, mid)
	l.InsertAfter(4, l.Back())
	l.InsertBefore(0, l.Front())
	checkLinks(t, &l, []int{0, 1, 2, 3, 4})
}

func TestDoublyLinkedList_foreignAndRemovedElements(t *testing.T) {
	var a, b DoublyLinkedList[int]
	ea := a.PushBack(1)
	b.PushBack(2)
	if b.Remove(ea) || b.MoveToFront(ea) || b.MoveToBack(ea) || b.InsertBefore(9, ea) != nil || b.InsertAfter(9, ea) != nil {
		t.Error("list b accepted list a's element")
	}
	if b.Remove(nil) || b.MoveToFront(nil) || b.InsertAfter(9, nil) != nil {
		t.Error("nil element accepted")
	}
	if !a.Remove(ea) || a.Remove(ea) {
		t.Error("second Remove of the same element succeeded")
	}
	if ea.Next() != nil || ea.Prev() != nil || ea.Value != 1 {
		t.Error("removed element still linked or lost its value")
	}
	var nilElem *Element[int]
	if nilElem.Next() != nil || nilElem.Prev() != nil {
		t.Error("nil element has neighbours")
	}
	checkLinks(t, &a, nil)
	checkLinks(t, &b, []int{2})
}

func TestDoublyLinkedList_moveOneElement(t *testing.T) {
	var l DoublyLinkedList[int]
	e := l.PushBack(1)
	if !l.MoveToFront(e) || !l.MoveToBack(e) {
		t.Error("move of the only element failed")
	}
	checkLinks(t, &l, []int{1})
}

func TestDoublyLinkedList_removeWhileIterating(t *testing.T) {
	var l DoublyLinkedList[int]
	for i := range 6 {
		l.PushBack(i)
	}
	var seen []int
	e := l.Front()
	for v := range l.All() {
		seen = append(seen, v)
		next := e.Next()
		if v%2 == 0 {
			l.Remove(e)
		}
		e = next
	}
	if !slices.Equal(seen, []int{0, 1, 2, 3, 4, 5}) {
		t.Errorf("seen = %v", seen)
	}
	checkLinks(t, &l, []int{1, 3, 5})
}

func TestDoublyLinkedList_earlyStop(t *testing.T) {
	var l DoublyLinkedList[int]
	for i := range 5 {
		l.PushBack(i)
	}
	var got []int
	for v := range l.Backward() {
		if v < 3 {
			break
		}
		got = append(got, v)
	}
	if !slices.Equal(got, []int{4, 3}) {
		t.Errorf("got %v", got)
	}
}

// A copied list must not loop: iterating the copy yields nothing (documented).
func TestDoublyLinkedList_copyDoesNotHang(t *testing.T) {
	var l DoublyLinkedList[int]
	l.PushBack(1)
	l.PushBack(2)
	cp := l // deliberate misuse
	if got := dllValues(&cp); len(got) != 0 {
		t.Errorf("copy iterated %v", got)
	}
	if got := dllBackward(&cp); len(got) != 0 {
		t.Errorf("copy iterated backward %v", got)
	}
}

// Oracle: random operations against a slice model.
func TestDoublyLinkedList_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		var l DoublyLinkedList[int]
		var model []int
		var elems []*Element[int] // elems[i] holds model[i]
		for op := range 60 {
			v := trial*1000 + op
			switch k := rng.IntN(6); {
			case k == 0 || len(model) == 0:
				elems = append([]*Element[int]{l.PushFront(v)}, elems...)
				model = append([]int{v}, model...)
			case k == 1:
				elems = append(elems, l.PushBack(v))
				model = append(model, v)
			case k == 2:
				i := rng.IntN(len(model))
				l.Remove(elems[i])
				elems = slices.Delete(elems, i, i+1)
				model = slices.Delete(model, i, i+1)
			case k == 3:
				i := rng.IntN(len(model))
				e := elems[i]
				l.MoveToFront(e)
				elems = append([]*Element[int]{e}, slices.Delete(elems, i, i+1)...)
				x := model[i]
				model = append([]int{x}, slices.Delete(model, i, i+1)...)
			case k == 4:
				i := rng.IntN(len(model))
				e := elems[i]
				l.MoveToBack(e)
				elems = append(slices.Delete(elems, i, i+1), e)
				x := model[i]
				model = append(slices.Delete(model, i, i+1), x)
			default:
				i := rng.IntN(len(model))
				ne := l.InsertAfter(v, elems[i])
				elems = slices.Insert(elems, i+1, ne)
				model = slices.Insert(model, i+1, v)
			}
		}
		checkLinks(t, &l, model)
	}
}

func TestDoublyLinkedList_noAllocMove(t *testing.T) {
	var l DoublyLinkedList[int]
	for i := range 10 {
		l.PushBack(i)
	}
	if a := testing.AllocsPerRun(100, func() { l.MoveToFront(l.Back()) }); a != 0 {
		t.Errorf("MoveToFront allocs = %v, want 0", a)
	}
}

func BenchmarkDoublyLinkedList(b *testing.B) {
	for _, s := range benchSizes {
		b.Run(s.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				var l DoublyLinkedList[int]
				for i := range s.n {
					l.PushBack(i)
				}
				for l.Len() > 0 {
					l.Remove(l.Front())
				}
			}
		})
	}
}

func BenchmarkDoublyLinkedListMoveToFront(b *testing.B) {
	for _, s := range benchSizes {
		var l DoublyLinkedList[int]
		for i := range s.n {
			l.PushBack(i)
		}
		b.Run(s.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				l.MoveToFront(l.Back())
			}
		})
	}
}
