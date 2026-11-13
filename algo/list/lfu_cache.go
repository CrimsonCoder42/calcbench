package list

import (
	"fmt"
	"math"
)

// lfuBucket holds every entry used exactly count times, most recently used at the front.
type lfuBucket[K comparable, V any] struct {
	count uint64
	items DoublyLinkedList[lfuEntry[K, V]]
}

// lfuEntry is one cached entry and the bucket (element of the frequency list) it is in.
type lfuEntry[K comparable, V any] struct {
	key    K
	value  V
	bucket *Element[lfuBucket[K, V]]
}

// LFUCache is a fixed-capacity key/value cache that, when full, evicts the least frequently used key;
// among keys used equally often it evicts the least recently used one. A use is a Get or a Put of the
// key (a Put of a new key is its first use); Delete forgets the key and its count. The use count
// saturates at 2^64−1.
//
// Methods: NewLFUCache(capacity) returns an empty cache. Get(key) returns the value and counts a use.
// Put(key, value) stores or replaces the value and counts a use; when a new key arrives at a full cache,
// it first evicts the entry chosen above and returns it with didEvict = true. Delete(key) removes the key
// and reports whether it was present. Len is the number of entries, Cap the capacity.
// Put refuses a key that is not equal to itself (ErrKeyNotComparable) and then changes nothing.
// The zero value and a nil *LFUCache have capacity 0: Get misses, Delete reports false, and Put stores
// nothing and returns the given pair as evicted. No method panics. No clock, no goroutines; not safe for
// concurrent use.
//
// The structure is Shah, Mitra and Matani's: a doubly linked frequency list of buckets in ascending use
// count, each bucket a linked list of its entries, every entry pointing to its bucket, and a hash map by
// key. The paper keeps each bucket's entries in a set; here the bucket is a recency-ordered list, which
// gives the least-recently-used tie-break at no cost (catalog/list.md B5-7).
//
// ID:      ALG-LIST-016
// Time:    O(1) expected per operation (one map operation plus O(1) relinking); NewLFUCache O(1)
// Space:   O(c) for c = capacity entries: one map slot and one list element per entry, at most c buckets
// Source:  Shah, Mitra, Matani, "An O(1) algorithm for implementing the LFU cache eviction scheme", technical report, 16 August 2010, arXiv:2110.11602, §5 (frequency list of node lists, O(1) insert, access and evict)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §10.2 (doubly linked list, O(1) delete) and ch. 11 (hash tables)
// Tags:    cache, lfu, eviction, frequency, hash-map, linked-list, generic
type LFUCache[K comparable, V any] struct {
	capacity int
	items    map[K]*Element[lfuEntry[K, V]]
	buckets  DoublyLinkedList[lfuBucket[K, V]] // ascending count; front = least frequently used
	len      int
}

// NewLFUCache returns an empty LFUCache that holds at most capacity entries.
// Errors: ErrCapacity when capacity < 1.
func NewLFUCache[K comparable, V any](capacity int) (*LFUCache[K, V], error) {
	if capacity < 1 {
		return nil, fmt.Errorf("%w: capacity = %d", ErrCapacity, capacity)
	}
	return &LFUCache[K, V]{capacity: capacity, items: make(map[K]*Element[lfuEntry[K, V]], min(capacity, 1<<16))}, nil
}

// Len returns the number of entries.
func (c *LFUCache[K, V]) Len() int {
	if c == nil {
		return 0
	}
	return c.len
}

// Cap returns the capacity (0 for the zero value or nil).
func (c *LFUCache[K, V]) Cap() int {
	if c == nil {
		return 0
	}
	return c.capacity
}

// Get returns the value stored under key and counts a use; ok is false on a miss.
func (c *LFUCache[K, V]) Get(key K) (value V, ok bool) {
	if c == nil {
		return value, false
	}
	e, ok := c.items[key]
	if !ok {
		return value, false
	}
	c.touch(e)
	return e.Value.value, true
}

// Put stores value under key and counts a use. When key is new and the cache is full, the least
// frequently used entry (least recently used among ties) is removed first and returned as evicted with
// didEvict = true. On a cache of capacity 0 nothing is stored and the given pair is returned as evicted.
// Errors: ErrKeyNotComparable when key != key; nothing changes.
func (c *LFUCache[K, V]) Put(key K, value V) (evicted Entry[K, V], didEvict bool, err error) {
	if key != key {
		return evicted, false, ErrKeyNotComparable
	}
	if c == nil || c.capacity < 1 {
		return Entry[K, V]{Key: key, Value: value}, true, nil
	}
	if e, ok := c.items[key]; ok {
		e.Value.value = value
		c.touch(e)
		return evicted, false, nil
	}
	if c.len < c.capacity {
		first := c.bucketOne()
		e := first.Value.items.PushFront(lfuEntry[K, V]{key: key, value: value, bucket: first})
		c.items[key] = e
		c.len++
		return evicted, false, nil
	}
	// Full: the victim is the least recently used entry of the lowest-count bucket. Its element is
	// reused for the new key, which enters with count 1.
	low := c.buckets.Front()
	e := low.Value.items.Back()
	evicted = Entry[K, V]{Key: e.Value.key, Value: e.Value.value}
	delete(c.items, e.Value.key)
	e.Value.key, e.Value.value = key, value
	switch {
	case low.Value.count == 1:
		low.Value.items.MoveToFront(e)
	case low.Value.items.Len() == 1:
		low.Value.count = 1 // the lowest bucket becomes the count-1 bucket; order is kept
	default:
		one := c.buckets.PushFront(lfuBucket[K, V]{count: 1})
		moveToFrontOf(e, &one.Value.items)
		e.Value.bucket = one
	}
	c.items[key] = e
	return evicted, true, nil
}

// Delete removes key and reports whether it was present.
func (c *LFUCache[K, V]) Delete(key K) bool {
	if c == nil {
		return false
	}
	e, ok := c.items[key]
	if !ok {
		return false
	}
	delete(c.items, key)
	b := e.Value.bucket
	b.Value.items.Remove(e)
	if b.Value.items.Len() == 0 {
		c.buckets.Remove(b)
	}
	c.len--
	return true
}

// bucketOne returns the count-1 bucket, creating it at the front when absent.
func (c *LFUCache[K, V]) bucketOne() *Element[lfuBucket[K, V]] {
	if f := c.buckets.Front(); f != nil && f.Value.count == 1 {
		return f
	}
	return c.buckets.PushFront(lfuBucket[K, V]{count: 1})
}

// touch counts one use of e: it moves to the front of the bucket for count+1.
func (c *LFUCache[K, V]) touch(e *Element[lfuEntry[K, V]]) {
	b := e.Value.bucket
	count := b.Value.count
	if count == math.MaxUint64 { // saturated: stays, but becomes the most recent of its bucket
		b.Value.items.MoveToFront(e)
		return
	}
	if next := b.Next(); next != nil && next.Value.count == count+1 {
		moveToFrontOf(e, &next.Value.items)
		e.Value.bucket = next
		if b.Value.items.Len() == 0 {
			c.buckets.Remove(b)
		}
		return
	}
	if b.Value.items.Len() == 1 {
		b.Value.count++ // e's bucket holds only e and no count+1 bucket exists: renumber in place
		return
	}
	nb := c.buckets.InsertAfter(lfuBucket[K, V]{count: count + 1}, b)
	moveToFrontOf(e, &nb.Value.items)
	e.Value.bucket = nb
}
