package list

import (
	"errors"
	"fmt"
)

// Errors of LRUCache and LFUCache. ErrCapacity is wrapped with the value given; test with errors.Is.
var (
	// ErrCapacity: a cache's capacity is less than 1.
	ErrCapacity = errors.New("list: cache capacity must be at least 1")
	// ErrKeyNotComparable: a key is not equal to itself (a floating-point NaN, or a struct or array
	// holding one). A Go map can never find such a key again, so a cache cannot store or evict it.
	ErrKeyNotComparable = errors.New("list: key is not equal to itself (NaN)")
)

// Entry is a key and its value, as reported by a cache's Put when it evicts one.
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// LRUCache is a fixed-capacity key/value cache that, when full, evicts the least recently used key.
// A use is a Get or a Put of the key; Peek and Delete are not uses.
//
// Methods: NewLRUCache(capacity) returns an empty cache. Get(key) returns the value and marks the key most
// recently used. Peek(key) returns the value without marking it. Put(key, value) stores or replaces the
// value, marks the key most recently used, and, when a new key arrives at a full cache, evicts the least
// recently used entry and returns it with didEvict = true. Delete(key) removes the key and reports whether
// it was present. Len is the number of entries, Cap the capacity.
// Put refuses a key that is not equal to itself (ErrKeyNotComparable) and then changes nothing.
// The zero value and a nil *LRUCache have capacity 0: Get and Peek miss, Delete reports false, and Put
// stores nothing and returns the given pair as evicted. No method panics. No clock, no goroutines; not
// safe for concurrent use. In steady state (a full cache) Put reuses the evicted entry's node.
//
// ID:      ALG-LIST-015
// Time:    O(1) expected per operation (one map operation plus O(1) relinking); NewLRUCache O(1)
// Space:   O(c) for c = capacity entries: one map slot and one list element per entry
// Source:  Mattson, Gecsei, Slutz, Traiger, "Evaluation techniques for storage hierarchies", IBM Systems Journal, 1970, 9(2):78–117, doi:10.1147/sj.92.0078 (LRU replacement as a stack algorithm)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §10.2 (doubly linked list, O(1) delete) and ch. 11 (hash tables)
// Tags:    cache, lru, eviction, hash-map, linked-list, generic
type LRUCache[K comparable, V any] struct {
	capacity int
	items    map[K]*Element[Entry[K, V]]
	order    DoublyLinkedList[Entry[K, V]] // front = most recently used
}

// NewLRUCache returns an empty LRUCache that holds at most capacity entries.
// Errors: ErrCapacity when capacity < 1.
func NewLRUCache[K comparable, V any](capacity int) (*LRUCache[K, V], error) {
	if capacity < 1 {
		return nil, fmt.Errorf("%w: capacity = %d", ErrCapacity, capacity)
	}
	// The map is presized up to a modest bound so a huge capacity does not allocate up front.
	return &LRUCache[K, V]{capacity: capacity, items: make(map[K]*Element[Entry[K, V]], min(capacity, 1<<16))}, nil
}

// Len returns the number of entries.
func (c *LRUCache[K, V]) Len() int {
	if c == nil {
		return 0
	}
	return c.order.Len()
}

// Cap returns the capacity (0 for the zero value or nil).
func (c *LRUCache[K, V]) Cap() int {
	if c == nil {
		return 0
	}
	return c.capacity
}

// Get returns the value stored under key and marks key most recently used; ok is false on a miss.
func (c *LRUCache[K, V]) Get(key K) (value V, ok bool) {
	if c == nil {
		return value, false
	}
	e, ok := c.items[key]
	if !ok {
		return value, false
	}
	c.order.MoveToFront(e)
	return e.Value.Value, true
}

// Peek returns the value stored under key without marking it used; ok is false on a miss.
func (c *LRUCache[K, V]) Peek(key K) (value V, ok bool) {
	if c == nil {
		return value, false
	}
	e, ok := c.items[key]
	if !ok {
		return value, false
	}
	return e.Value.Value, true
}

// Put stores value under key and marks key most recently used. When key is new and the cache is full,
// the least recently used entry is removed and returned as evicted with didEvict = true. On a cache of
// capacity 0 nothing is stored and the given pair is returned as evicted.
// Errors: ErrKeyNotComparable when key != key; nothing changes.
func (c *LRUCache[K, V]) Put(key K, value V) (evicted Entry[K, V], didEvict bool, err error) {
	if key != key {
		return evicted, false, ErrKeyNotComparable
	}
	if c == nil || c.capacity < 1 {
		return Entry[K, V]{Key: key, Value: value}, true, nil
	}
	if e, ok := c.items[key]; ok {
		e.Value.Value = value
		c.order.MoveToFront(e)
		return evicted, false, nil
	}
	if c.order.Len() < c.capacity {
		c.items[key] = c.order.PushFront(Entry[K, V]{Key: key, Value: value})
		return evicted, false, nil
	}
	// Full: reuse the least recently used element for the new entry.
	e := c.order.Back()
	evicted = e.Value
	delete(c.items, evicted.Key)
	e.Value = Entry[K, V]{Key: key, Value: value}
	c.order.MoveToFront(e)
	c.items[key] = e
	return evicted, true, nil
}

// Delete removes key and reports whether it was present.
func (c *LRUCache[K, V]) Delete(key K) bool {
	if c == nil {
		return false
	}
	e, ok := c.items[key]
	if !ok {
		return false
	}
	delete(c.items, key)
	c.order.Remove(e)
	return true
}
