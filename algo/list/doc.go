// Package list holds linked lists and the algorithms on them: a generic singly linked node with the
// classic relinking algorithms (reverse, reverse in blocks, merge, partition, reorder, remove from the
// end, middle, cycle detection, digit-list addition), a generic doubly linked list, and the two caches
// built on it (LRU and LFU, O(1) per operation).
//
// ID short form: LIST (IDs are ALG-LIST-NNN).
//
// Conventions: a singly linked list is named by its first *Node; nil is the empty list. Functions whose
// name ends in InPlace relink the caller's nodes and allocate none; every other function leaves its
// input unchanged. Every function that walks to the end of a list first checks for a cycle (Floyd,
// O(1) space) and returns ErrCycle with the list unchanged, so a cyclic list never hangs a call; nil
// never panics. Values are compared with cmp.Less (NaN first, doc/DECISIONS.md D-002). The caches take
// no clock and start no goroutine; none of the types is safe for concurrent use. The package name is the
// topic; DoublyLinkedList is a generic, unboxed alternative to container/list.
//
// Functions:
//
//	ALG-LIST-001  Node                     singly linked list node (Value, Next)
//	ALG-LIST-002  FromSlice                a list holding a slice's values, one allocation
//	ALG-LIST-003  AppendValues             append a list's values to a slice
//	ALG-LIST-004  Length                   the number of nodes, ErrCycle on a cycle
//	ALG-LIST-005  DoublyLinkedList         generic doubly linked list (Element, iterators), O(1) per operation
//	ALG-LIST-006  ReverseInPlace           reverse a list by relinking
//	ALG-LIST-007  ReverseKGroupInPlace     reverse each block of k nodes
//	ALG-LIST-008  DetectCycle              whether a list loops (Floyd)
//	ALG-LIST-009  CycleStart               the first node of the cycle and its length
//	ALG-LIST-010  Middle                   the node at index ⌊n/2⌋
//	ALG-LIST-011  MergeTwoSortedInPlace    stable merge of two ascending lists
//	ALG-LIST-012  RemoveNthFromEndInPlace  unlink the k-th node from the end
//	ALG-LIST-013  ReorderInPlace           relink to L0, Ln−1, L1, Ln−2, …
//	ALG-LIST-014  PartitionInPlace         stable partition around a pivot
//	ALG-LIST-015  LRUCache                 least-recently-used cache, O(1) per operation
//	ALG-LIST-016  LFUCache                 least-frequently-used cache (LRU among ties), O(1) per operation
//	ALG-LIST-017  AddTwoNumbers            sum of two digit lists, least significant digit first
package list
