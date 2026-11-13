package list

import "iter"

// Element is one element of a DoublyLinkedList. Value is the caller's; the links are private so a list
// can never be made cyclic or crossed with another list.
type Element[T any] struct {
	Value      T
	next, prev *Element[T]
	list       *DoublyLinkedList[T] // nil once removed
}

// Next returns the element after e, or nil when e is the last, removed, or nil.
func (e *Element[T]) Next() *Element[T] {
	if e == nil || e.list == nil || e.next == &e.list.root {
		return nil
	}
	return e.next
}

// Prev returns the element before e, or nil when e is the first, removed, or nil.
func (e *Element[T]) Prev() *Element[T] {
	if e == nil || e.list == nil || e.prev == &e.list.root {
		return nil
	}
	return e.prev
}

// DoublyLinkedList is a generic doubly linked list: O(1) insert, remove and move given an element, with
// values stored unboxed (container/list stores `any`; see doc/algorithms_core_2026-11-12/catalog/list.md
// B5-3 for why this type exists beside it).
//
// The zero value is an empty list ready to use. A list must not be copied after first use: its sentinel's
// address is its identity (iterating a copy yields nothing rather than looping). Not safe for
// concurrent use.
//
// Methods: PushFront(v) and PushBack(v) add v and return its element. InsertBefore(v, mark) and
// InsertAfter(v, mark) add v beside mark and return its element, or nil when mark is not on this list.
// Remove(e) unlinks e and reports whether it was on this list. MoveToFront(e) and MoveToBack(e) move e and
// report whether it was on this list. Front and Back return the first and last element (nil when empty).
// Len is the element count. All and Backward iterate the values front to back and back to front; the
// next element is read before each yield, so removing the current element while iterating is safe.
// Every method given nil, a removed element or another list's element does nothing; none panics, on a
// nil *DoublyLinkedList either (reads report empty; pushes return nil).
//
// ID:      ALG-LIST-005
// Time:    O(1) worst case per operation; All and Backward O(n) for n elements
// Space:   O(n): one allocation per pushed element
// Source:  Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.5 (doubly linked lists with a list head: LLINK, RLINK; O(1) deletion and insertion)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §10.2 (LIST-INSERT, LIST-DELETE with a sentinel, Figure 10.4)
// Tags:    linked-list, doubly-linked, deque, sentinel, iterator, data-structure, generic
type DoublyLinkedList[T any] struct {
	root Element[T] // sentinel: root.next is the front, root.prev the back
	len  int
}

// lazyInit makes the zero value's sentinel point to itself.
func (l *DoublyLinkedList[T]) lazyInit() {
	if l.root.next == nil {
		l.root.next = &l.root
		l.root.prev = &l.root
	}
}

// Len returns the number of elements; 0 for a nil list.
func (l *DoublyLinkedList[T]) Len() int {
	if l == nil {
		return 0
	}
	return l.len
}

// Front returns the first element, or nil when the list is empty or nil.
func (l *DoublyLinkedList[T]) Front() *Element[T] {
	if l == nil || l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back returns the last element, or nil when the list is empty or nil.
func (l *DoublyLinkedList[T]) Back() *Element[T] {
	if l == nil || l.len == 0 {
		return nil
	}
	return l.root.prev
}

// PushFront adds v at the front and returns its element (nil on a nil list).
func (l *DoublyLinkedList[T]) PushFront(v T) *Element[T] {
	if l == nil {
		return nil
	}
	l.lazyInit()
	return l.link(&Element[T]{Value: v}, &l.root)
}

// PushBack adds v at the back and returns its element (nil on a nil list).
func (l *DoublyLinkedList[T]) PushBack(v T) *Element[T] {
	if l == nil {
		return nil
	}
	l.lazyInit()
	return l.link(&Element[T]{Value: v}, l.root.prev)
}

// InsertBefore adds v just before mark and returns its element, or nil when mark is not on l.
func (l *DoublyLinkedList[T]) InsertBefore(v T, mark *Element[T]) *Element[T] {
	if !l.owns(mark) {
		return nil
	}
	return l.link(&Element[T]{Value: v}, mark.prev)
}

// InsertAfter adds v just after mark and returns its element, or nil when mark is not on l.
func (l *DoublyLinkedList[T]) InsertAfter(v T, mark *Element[T]) *Element[T] {
	if !l.owns(mark) {
		return nil
	}
	return l.link(&Element[T]{Value: v}, mark)
}

// Remove unlinks e from l and reports whether it was on l. A removed element keeps its Value; its Next
// and Prev return nil from then on.
func (l *DoublyLinkedList[T]) Remove(e *Element[T]) bool {
	if !l.owns(e) {
		return false
	}
	l.unlink(e)
	return true
}

// MoveToFront moves e to the front of l and reports whether e was on l.
func (l *DoublyLinkedList[T]) MoveToFront(e *Element[T]) bool {
	if !l.owns(e) {
		return false
	}
	if l.root.next != e {
		l.unlink(e)
		l.link(e, &l.root)
	}
	return true
}

// MoveToBack moves e to the back of l and reports whether e was on l.
func (l *DoublyLinkedList[T]) MoveToBack(e *Element[T]) bool {
	if !l.owns(e) {
		return false
	}
	if l.root.prev != e {
		l.unlink(e)
		l.link(e, l.root.prev)
	}
	return true
}

// All returns an iterator over the values from front to back.
func (l *DoublyLinkedList[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		if l == nil || l.len == 0 {
			return
		}
		for e := l.root.next; e != nil && e.list == l; { // the sentinel's list is nil
			next := e.next
			if !yield(e.Value) {
				return
			}
			e = next
		}
	}
}

// Backward returns an iterator over the values from back to front.
func (l *DoublyLinkedList[T]) Backward() iter.Seq[T] {
	return func(yield func(T) bool) {
		if l == nil || l.len == 0 {
			return
		}
		for e := l.root.prev; e != nil && e.list == l; {
			prev := e.prev
			if !yield(e.Value) {
				return
			}
			e = prev
		}
	}
}

// owns reports whether e is an element currently on l.
func (l *DoublyLinkedList[T]) owns(e *Element[T]) bool {
	return l != nil && e != nil && e.list == l && e != &l.root
}

// link inserts e after at (an element of l or its sentinel) and returns e.
func (l *DoublyLinkedList[T]) link(e, at *Element[T]) *Element[T] {
	e.prev = at
	e.next = at.next
	at.next.prev = e
	at.next = e
	e.list = l
	l.len++
	return e
}

// unlink removes e (an element of l) and clears its links.
func (l *DoublyLinkedList[T]) unlink(e *Element[T]) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next, e.prev, e.list = nil, nil, nil
	l.len--
}

// moveToFrontOf moves e from its list to the front of dst without allocating (LFU buckets).
// e must be on a list; dst must not be nil.
func moveToFrontOf[T any](e *Element[T], dst *DoublyLinkedList[T]) {
	e.list.unlink(e)
	dst.lazyInit()
	dst.link(e, &dst.root)
}
